package accessdata

type options struct {
	origin     Origin
	pingRoutes []string
}

var opt options

// SetOrigin - required to track service identification
func SetOrigin(o Origin) {
	opt.origin = o
}

func getOrigin() *Origin {
	return &opt.origin
}

// SetPingRoutes - initialize the ping routes
func SetPingRoutes(routes []string) {
	opt.pingRoutes = routes
}

func IsPingTraffic(name string) bool {
	for _, n := range opt.pingRoutes {
		if n == name {
			return true
		}
	}
	return false
}
