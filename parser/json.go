package parser

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
)

type JsonParser struct {
}

func (j JsonParser) Parse(_ *http.Response, params map[string]string, previousResponse string) (string, error) {
	key, exists := params["key"]
	if !exists {
		return "", fmt.Errorf("json parser: key not found in params")
	}

	return gjson.Get(previousResponse, key).String(), nil
}
