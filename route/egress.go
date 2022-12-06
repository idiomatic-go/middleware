package route

type EgressMatch func(url, method string) (name string)

var (
	egressDefault = &Route{Name: "/"}
	egressMatch   EgressMatch
	egressRoutes  routes
)

func init() {
	SetEgressMatchFn(nil)
	egressRoutes = NewTable()
}

func SetEgressDefault(r *Route) bool {
	if r == nil || r.Name == "" {
		return false
	}
	egressDefault = r
	return true
}

func SetEgressMatchFn(fn EgressMatch) {
	if fn != nil {
		egressMatch = fn
	} else {
		egressMatch = func(url, method string) (name string) {
			return egressDefault.Name
		}
	}
}

// LookupEgress - Find a Route based on the url and method, returning the Default if not found
func LookupEgress(url, method string) *Route {
	name := egressMatch(url, method)
	if name != "" {
		return egressRoutes.lookup(name)
	}
	return egressDefault
}

// AddEgress - Add a Route
func AddEgress(r *Route) bool {
	return egressRoutes.add(r)
}

// UpdateEgress - Update a Route
func UpdateEgress(r *Route) bool {
	return egressRoutes.update(r)
}

// RemoveEgress - Remove a route
func RemoveEgress(name string) bool {
	return egressRoutes.remove(name)
}
