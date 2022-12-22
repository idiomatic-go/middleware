package automation

import (
	"net/http"
)

type Matcher func(req *http.Request) (routeName string)

type Configuration interface {
	SetMatcher(fn Matcher)
	SetDefault(name string, config ...any) error
	Add(name string, config ...any) error
}

type Actuators interface {
	Lookup(req *http.Request) Actuator
	LookupByName(name string) Actuator
}

type Automation interface {
	Configuration
	Actuators
}

var Ingress = NewIngressTable()
var Egress = NewEgressTable()
