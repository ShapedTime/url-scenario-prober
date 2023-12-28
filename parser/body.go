package parser

import (
	"bytes"
	"io"
	"net/http"
)

type BodyParser struct {
}

func (p BodyParser) Parse(response *http.Response, _ map[string]string, _ string) (string, error) {
	body, err := io.ReadAll(response.Body)
	err = response.Body.Close()
	if err != nil {
		return "", err
	}

	response.Body = io.NopCloser(bytes.NewReader(body))

	return string(body), nil
}
