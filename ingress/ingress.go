package ingress

import (
	"github.com/idiomatic-go/middleware/route"
)

var Routes = route.NewTable()

/*


//defaultRoute =
//match        route.MatchFn
//func init() {
//SetMatchFn(nil)
//}

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

// LookupRoute - Find a Route based on the url and method, returning the Default if not found
func LookupRoute(req *http.Request) (route.Route,bool) {
	name := match(req)
	if name != "" {
		return routes.Lookup(name)
	}
	return *defaultRoute,true
}

// AddRoute - Add a Route
func AddRoute(r *route.Route) bool {
	return routes.Add(r)
}

// UpdateRoute - Update a Route
func UpdateTimeout(name string,timeout int) bool {
	return routes.UpdateTimeout(name,timeout)
}

// Remove - Remove a route
func RemoveRoute(name string) bool {
	return routes.Remove(name)
}


*/
