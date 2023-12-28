package vars

import "sync"

type Vars struct {
	v   map[string]string
	mux sync.Mutex
}

func NewVars() *Vars {
	return &Vars{
		v: make(map[string]string),
	}
}

func (v *Vars) Set(key, value string) {
	v.mux.Lock()
	v.v[key] = value
	v.mux.Unlock()
}

func (v *Vars) Get(key string) string {
	v.mux.Lock()
	value := v.v[key]
	v.mux.Unlock()
	return value
}
