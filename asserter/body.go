package asserter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type BodyAsserter struct {
}

func (a BodyAsserter) Assert(response *http.Response, params map[string]string) (bool, error) {
	for _, value := range params {
		body, err := io.ReadAll(response.Body)
		err = response.Body.Close()
		if err != nil {
			return false, err
		}

		response.Body = io.NopCloser(bytes.NewReader(body))
		if err != nil {
			return false, fmt.Errorf("body asserter: error reading response body: %w", err)
		}

		match, err := regexp.MatchString(value, string(body))
		if err != nil || !match {
			return false, err
		}
	}

	return true, nil
}
