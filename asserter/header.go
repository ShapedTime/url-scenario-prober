package asserter

import (
	"net/http"
	"regexp"
)

type HeaderAsserter struct {
}

// Assert checks if the response headers match the params value regex
func (a HeaderAsserter) Assert(response *http.Response, params map[string]string) (bool, error) {
	for key, value := range params {
		headerValue := response.Header.Get(key)
		match, err := regexp.MatchString(value, headerValue)
		if err != nil || !match {
			return false, err
		}
	}

	return true, nil
}
