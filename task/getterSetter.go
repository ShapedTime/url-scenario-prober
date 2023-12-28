package task

func (t *Task) GetUrl() string {
	t.RLock()
	url := t.Url
	t.RUnlock()
	return url
}

func (t *Task) GetName() string {
	t.RLock()
	name := t.Name
	t.RUnlock()
	return name
}

func (t *Task) GetRegister() []struct {
	Var   string              `yaml:"var"`
	Parse []map[string]string `yaml:"parse"`
} {
	t.RLock()
	register := t.Register
	t.RUnlock()
	return register
}

func (t *Task) GetAssertResponse() struct {
	Headers    map[string]string `yaml:"headers,omitempty"`
	Body       map[string]string `yaml:"body,omitempty"`
	Cookies    map[string]string `yaml:"cookies,omitempty"`
	StatusCode map[string]string `yaml:"status_code,omitempty"`
} {
	t.RLock()
	assertResponse := t.AssertResponse
	t.RUnlock()
	return assertResponse
}

func (t *Task) GetDependsOn() []string {
	t.RLock()
	dependsOn := t.DependsOn
	t.RUnlock()
	return dependsOn
}

func (t *Task) GetParams() Params {
	t.RLock()
	params := t.Params
	t.RUnlock()
	return params
}

func (t *Task) GetStatusMessage() string {
	t.RLock()
	statusMessage := t.StatusMessage
	t.RUnlock()
	return statusMessage
}

func (t *Task) SetStatusMessage(statusMessage string) {
	t.Lock()
	t.StatusMessage = statusMessage
	t.Unlock()
}

func (t *Task) SetParams(params Params) {
	t.Lock()
	t.Params = params
	t.Unlock()
}

func (t *Task) SetUrl(url string) {
	t.Lock()
	t.Url = url
	t.Unlock()
}
