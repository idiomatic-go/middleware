package route

import (
	"net/http"
	"time"
)

type MatchFn func(req *http.Request) (name string)

type Route struct {
	Name           string
	Timeout        int // milliseconds
	RateLimit      int
	WriteAccessLog bool
	Ping           bool
}

func (r *Route) IsTimeout() bool {
	return r != nil && r.Timeout != 0
}

func (r *Route) Duration() time.Duration {
	if r == nil {
		return 0
	}
	return time.Duration(r.Timeout) * time.Millisecond
}

func (r *Route) IsLogging() bool {
	return r != nil && r.WriteAccessLog
}
