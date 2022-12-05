package ingress

import (
	"github.com/idiomatic-go/middleware/route"
	"net/http"
)

var DefaultRoute = &route.Route{Name: "/"}

type Match func(r *http.Request) (name string)

var matchFn Match

func init() {
	SetMatchFn(nil)
}

func SetMatchFn(fn Match) {
	if fn != nil {
		matchFn = fn
	} else {
		matchFn = func(r *http.Request) (name string) {
			return DefaultRoute.Name
		}
	}
}

// LookupRoute - Find a Route based on the http.Request, returning the Default if not found
func LookupRoute(req *http.Request) *route.Route {
	if req != nil {
		//return DefaultRoute
	}
	return DefaultRoute
}
