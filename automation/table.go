package automation

import (
	"net/http"
	"sync"
)

type Routes interface {
	SetDefault(a Actuator)
	SetMatcher(fn Matcher)

	Lookup(req *http.Request) Actuator
	LookupByName(name string) Actuator

	Exists(name string) bool
	Add(a Actuator) bool
	Remove(name string)

	IsEnabled(name string) bool
	Reset(name string)
	Disable(name string) bool
	Configure(name string, v ...any)

	Timeout() TimeoutAction
}

type table struct {
	mu        sync.RWMutex
	actuators map[string]Actuator
	Default   Actuator
	match     Matcher
}
