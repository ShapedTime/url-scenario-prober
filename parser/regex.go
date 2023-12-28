package parser

import (
	"fmt"
	"net/http"
	"regexp"
)

type RegexParser struct {
}

func (p RegexParser) Parse(_ *http.Response, params map[string]string, previousResponse string) (string, error) {
	regex, exists := params["regex"]
	if !exists {
		return "", fmt.Errorf("regex parser: requires regex parameter")
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return "", fmt.Errorf("regex parser: invalid regex: %v", err)
	}

	match := r.FindStringSubmatch(previousResponse)
	if len(match) == 0 {
		return "", fmt.Errorf("regex parser: no match found")
	}

	return match[1], nil
}
