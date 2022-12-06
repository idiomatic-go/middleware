package route

type Route struct {
	Name      string
	Timeout   int // milliseconds
	RateLimit int
	AccessLog bool
	Ping      bool
}

func (route *Route) IsTimeout() bool {
	return route.Timeout != 0
}
