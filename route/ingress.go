package route

import "net/http"

type IngressMatch func(req *http.Request) (name string)

var (
	ingressDefault = &Route{Name: "/"}
	ingressMatch   IngressMatch
	ingressRoutes  routes
)

func init() {
	SetIngressMatchFn(nil)
	ingressRoutes = NewTable()
}

func SetIngressDefault(r *Route) bool {
	if r == nil || r.Name == "" {
		return false
	}
	ingressDefault = r
	return true
}

func SetIngressMatchFn(fn IngressMatch) {
	if fn != nil {
		ingressMatch = fn
	} else {
		ingressMatch = func(req *http.Request) (name string) {
			return ingressDefault.Name
		}
	}
}

// LookupIngress - Find a Route based on the url and method, returning the Default if not found
func LookupIngress(req *http.Request) *Route {
	name := ingressMatch(req)
	if name != "" {
		return ingressRoutes.lookup(name)
	}
	return ingressDefault
}

// AddIngress - Add a Route
func AddIngress(r *Route) bool {
	return ingressRoutes.add(r)
}

// UpdateIngress - Update a Route
func UpdateIngress(r *Route) bool {
	return ingressRoutes.update(r)
}

// RemoveIngress - Remove a route
func RemoveIngress(name string) bool {
	return ingressRoutes.remove(name)
}
