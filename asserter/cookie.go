package asserter

import (
	"net/http"
	"regexp"
)

type CookieAsserter struct {
}

func (a CookieAsserter) Assert(response *http.Response, params map[string]string) (bool, error) {
	for key, value := range params {
		cookies := response.Cookies()
		matched := false
		for _, cookie := range cookies {
			if cookie.Name == key {
				match, err := regexp.MatchString(value, cookie.Value)
				if err != nil || !match {
					return false, err
				}
				matched = true
				break
			}
		}
		if !matched {
			return false, nil
		}
	}

	return true, nil
}
