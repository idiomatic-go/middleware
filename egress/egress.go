package egress

import (
	"github.com/idiomatic-go/middleware/route"
	"net/http"
)

var (
	defaultRoute = &route.Route{Name: "/"}
	match        route.MatchFn
	routes       route.Routes
)

func init() {
	SetMatchFn(nil)
	routes = route.NewTable()
}

func SetDefaultRoute(r *route.Route) bool {
	if r == nil || r.Name == "" {
		return false
	}
	defaultRoute = r
	return true
}

func SetMatchFn(fn route.MatchFn) {
	if fn != nil {
		match = fn
	} else {
		match = func(req *http.Request) (name string) {
			return ""
		}
	}
}

// Lookup - Find a Route based on the url and method, returning the Default if not found
func Lookup(req *http.Request) *route.Route {
	name := match(req)
	if name != "" {
		return routes.Lookup(name)
	}
	return defaultRoute
}

// Add - Add a Route
func Add(r *route.Route) bool {
	return routes.Add(r)
}

// Update - Update a Route
func Update(r *route.Route) bool {
	return routes.Update(r)
}

// Remove - Remove a route
func Remove(name string) bool {
	return routes.Remove(name)
}
