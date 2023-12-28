package parser

import (
	"fmt"
	"net/http"
)

type Parser interface {
	Parse(response *http.Response, params map[string]string, previousResponse string) (string, error)
}

type Factory struct {
	parsers map[string]Parser
}

func NewFactory() *Factory {
	parsers := map[string]Parser{
		"cookie": CookieParser{},
		"header": HeaderParser{},
		"regex":  RegexParser{},
		"body":   BodyParser{},
		"json":   JsonParser{},
	}

	return &Factory{
		parsers: parsers,
	}
}

func (f *Factory) GetParser(parserType string) (Parser, error) {
	parser, ok := f.parsers[parserType]
	if !ok {
		return nil, fmt.Errorf("parser type %s not found", parserType)
	}

	return parser, nil
}
