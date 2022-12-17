package automation

import (
	"net/http"
)

type Matcher func(req *http.Request) (routeName string)

type Configuration interface {
	SetDefault(name string, t *TimeoutConfig)
	SetMatcher(fn Matcher)
	// IsPingEnabled(name string) bool
	Add(name string, t *TimeoutConfig) bool
	//Exists(name string) bool
	//Remove(name string)
}

type Actuators interface {
	Lookup(req *http.Request) Actuator
	LookupByName(name string) Actuator
}

type Automation interface {
	Configuration
	Actuators
}

var Ingress = NewTable()
var Egress = NewTable()
