package route

import "time"

type Route struct {
	Name      string
	Timeout   time.Duration
	RateLimit int
	AccessLog bool
	Ping      bool
}

func (route *Route) IsTimeout() bool {
	return route.Timeout != 0
}

type routeTable struct {
	m map[string]*Route
}

//func SetIngressMatchRequest
