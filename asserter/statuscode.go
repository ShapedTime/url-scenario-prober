package asserter

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

type StatusCodeAsserter struct {
}

func (a StatusCodeAsserter) Assert(response *http.Response, params map[string]string) (bool, error) {
	var err error
	for _, value := range params {
		var match bool
		match, err = regexp.MatchString(value, strconv.Itoa(response.StatusCode))

		// delete this
		body, err := io.ReadAll(response.Body)
		err = response.Body.Close()
		if err != nil {
			return false, err
		}

		response.Body = io.NopCloser(bytes.NewReader(body))
		// delete this

		if err == nil && match {
			return true, nil
		}
	}

	return false, err
}
