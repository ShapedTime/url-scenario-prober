package task

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"url-scenario-prober/asserter"
	"url-scenario-prober/parser"
	"url-scenario-prober/vars"
)

type Tasks struct {
	Tasks map[string]*Task `yaml:"tasks"`
}

type Params struct {
	Headers   map[string]string `yaml:"headers,omitempty"`
	Cookies   map[string]string `yaml:"cookies,omitempty"`
	Method    string            `yaml:"method,omitempty"`
	Body      string            `yaml:"body,omitempty"`
	GetParams map[string]string `yaml:"get_params,omitempty"`
}

type Task struct {
	Url      string `yaml:"url"`
	Name     string `yaml:"name"`
	Register []struct {
		Var   string              `yaml:"var"`
		Parse []map[string]string `yaml:"parse"` // this will contain Parse["type"] to be checked
	} `yaml:"register,omitempty"`
	AssertResponse struct {
		Headers    map[string]string `yaml:"headers,omitempty"`
		Body       map[string]string `yaml:"body,omitempty"`
		Cookies    map[string]string `yaml:"cookies,omitempty"`
		StatusCode map[string]string `yaml:"status_code,omitempty"`
	} `yaml:"assert_response,omitempty"`
	DependsOn     []string `yaml:"depends_on,omitempty"`
	Params        Params   `yaml:"params,omitempty"`
	status        Status
	StatusMessage string
	*sync.RWMutex
}

func LoadTasks(fileName string) (Tasks, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return Tasks{}, err
	}

	var tasksSlice struct {
		Tasks []Task `yaml:"tasks"`
	}
	err = yaml.Unmarshal(file, &tasksSlice)
	if err != nil {
		return Tasks{}, err
	}

	tasksMap := make(map[string]*Task)
	for i := range tasksSlice.Tasks {
		task := tasksSlice.Tasks[i]
		tasksMap[task.Name] = &task
		task.RWMutex = &sync.RWMutex{}
	}

	return Tasks{
		Tasks: tasksMap,
	}, nil
}

func (t *Tasks) GetTask(name string) *Task {
	return t.Tasks[name]
}

func (t *Task) Run(vars *vars.Vars, client *http.Client) {
	t.fillTemplateValues(vars)

	req, err := t.composeRequest()
	if err != nil {
		t.SetStatus(STATUS_FAILED_UNEXPECTED)
		t.SetStatusMessage(fmt.Sprintf("failed to create request: %s", err.Error()))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		t.SetStatus(STATUS_FAILED_UNEXPECTED)
		t.SetStatusMessage(fmt.Sprintf("failed to send request: %s", err.Error()))
		return
	}

	defer resp.Body.Close()

	// assert response
	result, failReason, err := t.assertResponse(resp)
	if err != nil {
		t.SetStatus(STATUS_FAILED_UNEXPECTED)
		t.SetStatusMessage(fmt.Sprintf("failed to assert response for %s: %s", failReason, err.Error()))
		return
	}

	if !result {
		t.SetStatus(STATUS_FAILED)
		t.SetStatusMessage(fmt.Sprintf("failed to assert response for %s", failReason))
		return
	}

	err = t.runParsers(resp, vars)
	if err != nil {
		t.SetStatus(STATUS_FAILED_UNEXPECTED)
		t.SetStatusMessage(fmt.Sprintf("failed to run parsers: %s", err.Error()))
		return
	}

	t.SetStatus(STATUS_SUCCESS)
}

func (t *Task) composeRequest() (*http.Request, error) {
	u, err := url.Parse(t.Url)

	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s: %s", t.Url, err.Error())
	}

	params := t.GetParams()
	getParams := url.Values{}
	for key, val := range params.GetParams {
		getParams.Add(key, val)
	}

	u.RawQuery = getParams.Encode()

	req, err := http.NewRequest(params.Method, u.String(), strings.NewReader(params.Body))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err.Error())
	}

	for key, val := range params.Headers {
		req.Header.Add(key, val)
	}

	for key, val := range params.Cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: val})
	}

	return req, err
}

// fillTemplateValues replaces all {{.var}} with the value of var
func (t *Task) fillTemplateValues(vars *vars.Vars) {
	varPattern := regexp.MustCompile(`{{\.(.*?)}}`)

	params := t.GetParams()

	params.Body = varPattern.ReplaceAllStringFunc(params.Body, func(s string) string {
		key := s[3 : len(s)-2] // Extract variable name
		return vars.Get(key)
	})

	for k, v := range params.Headers {
		params.Headers[k] = varPattern.ReplaceAllStringFunc(v, func(s string) string {
			key := s[3 : len(s)-2] // Extract variable name
			return vars.Get(key)
		})
	}

	for k, v := range params.Cookies {
		params.Cookies[k] = varPattern.ReplaceAllStringFunc(v, func(s string) string {
			key := s[3 : len(s)-2] // Extract variable name
			return vars.Get(key)
		})
	}

	for k, v := range params.GetParams {
		params.GetParams[k] = varPattern.ReplaceAllStringFunc(v, func(s string) string {
			key := s[3 : len(s)-2] // Extract variable name
			return vars.Get(key)
		})
	}

	t.SetParams(params)

	t.SetUrl(varPattern.ReplaceAllStringFunc(t.GetUrl(), func(s string) string {
		key := s[3 : len(s)-2] // Extract variable name
		return vars.Get(key)
	}))
}

func (t *Task) runParsers(resp *http.Response, vars *vars.Vars) error {
	parserFactory := parser.NewFactory()

	register := t.GetRegister()

	for _, v := range register {
		for _, p := range v.Parse {
			parserType := p["type"]
			parser, err := parserFactory.GetParser(parserType)
			if err != nil {
				return fmt.Errorf("failed to get parser: %s", err.Error())
			}

			value, err := parser.Parse(resp, p, vars.Get(v.Var))
			if err != nil {
				return fmt.Errorf("failed to parse response: %s", err.Error())
			}

			vars.Set(v.Var, value)
		}
	}

	return nil
}

func (t *Task) assertResponse(resp *http.Response) (bool, string, error) {
	headerResult, err := asserter.HeaderAsserter{}.Assert(resp, t.GetAssertResponse().Headers)
	if err != nil || !headerResult {
		return false, "header result", err
	}

	statusResult, err := asserter.StatusCodeAsserter{}.Assert(resp, t.GetAssertResponse().StatusCode)
	if err != nil || !statusResult {
		return false, "status code", err
	}

	bodyResult, err := asserter.BodyAsserter{}.Assert(resp, t.GetAssertResponse().Body)
	if err != nil || !bodyResult {
		return false, "body result", err
	}

	cookieResult, err := asserter.CookieAsserter{}.Assert(resp, t.GetAssertResponse().Cookies)
	if err != nil || !cookieResult {
		return false, "cookie result", err
	}

	return true, "", nil
}
