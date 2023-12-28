package parser

import (
	"fmt"
	"net/http"
)

type CookieParser struct {
}

func (p CookieParser) Parse(response *http.Response, params map[string]string, _ string) (string, error) {
	key, exists := params["key"]
	if !exists {
		return "", fmt.Errorf("cookie parser: requires key parameter")
	}

	cookie, err := response.Request.Cookie(key)
	if err != nil {
		return "", fmt.Errorf("cookie parser: %v", err)
	}
	return cookie.Value, nil
}
