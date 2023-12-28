package parser

import (
	"fmt"
	"net/http"
)

type HeaderParser struct {
}

func (p HeaderParser) Parse(response *http.Response, params map[string]string, _ string) (string, error) {
	key, exists := params["key"]
	if !exists {
		return "", fmt.Errorf("header parser: requires key parameter")
	}

	return response.Header.Get(key), nil
}
