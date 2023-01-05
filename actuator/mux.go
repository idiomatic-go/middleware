package actuator

import (
	"errors"
	"fmt"
	"net/http"
)

type mux struct {
	m     map[string]string
	hosts bool
}

func newMux() *mux {
	m := new(mux)
	m.m = make(map[string]string)
	return m
}

func (m *mux) add(pattern, name string) error {
	if pattern == "" || name == "" {
		return errors.New("invalid configuration: pattern or name is empty")
	}
	if _, exist := m.m[pattern]; exist {
		return errors.New(fmt.Sprintf("invalid configuration: pattern already exists [%v]", pattern))
	}
	m.m[pattern] = name
	if pattern[0] != '/' {
		m.hosts = true
	}
	return nil
}

func (m *mux) lookup(req *http.Request) (name string, ok bool) {
	if req == nil {
		return "", false
	}
	if m.hosts {
		name = m.match(req.Host + req.URL.Path)
	}
	if name == "" {
		name = m.match(req.URL.Path)
	}
	if name != "" {
		return name, true
	}
	return "", false
}

func (m *mux) match(path string) string {
	// Check for exact match first.
	v, ok := m.m[path]
	if ok {
		return v
	}

	// Check for longest valid match.  mux.es contains all patterns
	// that end in / sorted from longest to shortest.
	//for _, e := range mux.es {
	//	if strings.HasPrefix(path, e.pattern) {
	//		return e.h, e.pattern
	//	}
	//}
	return ""
}
