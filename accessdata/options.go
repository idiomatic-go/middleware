package accessdata

// SetOrigin - required to track service identification
func SetOrigin(o Origin) {
	opt.origin = o
}

// SetPingRoutes - initialize the ping routes
func SetPingRoutes(routes []string) {
	opt.pingRoutes = routes
}

func IsPingRoute(name string) bool {
	for _, n := range opt.pingRoutes {
		if n == name {
			return true
		}
	}
	return false
}

type options struct {
	origin     Origin
	pingRoutes []string
}

var opt options
