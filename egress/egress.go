package egress

import (
	"github.com/idiomatic-go/middleware/route"
)

var Routes = route.NewTable(&route.Route{Name: "/"})

/*

//defaultRoute =
//match  route.MatchFn

func init() {
	SetMatchFn(nil)
	//Routes =
}

func SetDefaultRoute(r *route.Route) bool {
	if r == nil || r.Name == "" {
		return false
	}
	//defaultRoute = r
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

// LookupRoute - Find a Route based on the url and method, returning the Default if not found
func LookupRoute(req *http.Request) (route.Route, bool) {
	name := match(req)
	if name != "" {
		return routes.Lookup(name)
	}
	return *defaultRoute, true
}


 // AddRoute - Add a Route
func AddRoute(r *route.Route) bool {
	return routes.Add(r)
}

func AddRouteWithLimiter(r *route.Route, max rate.Limit, b int) bool {
	return routes.AddWithLimiter(r, max, b)
}

// UpdateRouteTimeout - Update a Route
func UpdateRouteTimeout(name string, timeout int) bool {
	return routes.UpdateTimeout(name, timeout)
}

func UpdateRouteLimiter(name string, max rate.Limit, b int) bool {
	return routes.UpdateLimiter(name, max, b)
}

// RemoveRoute - Remove a route
func RemoveRoute(name string) bool {
	return routes.Remove(name)
}

func RemoveRouteLimiter(name string) bool {
	return routes.RemoveLimiter(name)
}
*/
