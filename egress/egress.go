package egress

import (
	"github.com/idiomatic-go/middleware/route"
)

var DefaultRoute = &route.Route{Name: "/"}

type Match func(url, method string) (name string)

var matchFn Match

func init() {
	SetMatchFn(nil)
}

func SetMatchFn(fn Match) {
	if fn != nil {
		matchFn = fn
	} else {
		matchFn = func(url, method string) (name string) {
			return DefaultRoute.Name
		}
	}
}

// LookupRoute - Find a Route based on the url and method, returning the Default if not found
func LookupRoute(url, method string) *route.Route {
	name := matchFn(url, method)
	if name != "" {
		//return DefaultRoute
	}
	return DefaultRoute
}

//func With
