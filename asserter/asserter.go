package asserter

import "net/http"

type Asserter interface {
	Assert(response *http.Response, params map[string]string) (bool, error)
}
