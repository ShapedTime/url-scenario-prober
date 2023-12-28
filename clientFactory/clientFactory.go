package clientFactory

import (
	"net/http"
	"time"
)

type ClientFactory struct {
	timeoutDuration int
}

func NewClientFactory(timeoutDuration int) *ClientFactory {
	return &ClientFactory{
		timeoutDuration: timeoutDuration,
	}
}

func (c *ClientFactory) NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * time.Duration(c.timeoutDuration),
	}
}
