package ingress

import (
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
)

func TimeoutHandler(routeName string, h http.Handler) http.Handler {
	if h == nil {
		return h
	}
	act := actuator.IngressTable.LookupByName(routeName)
	if r, ok := act.Timeout(); ok {
		return http.TimeoutHandler(h, r.Duration(), "")
	}
	return h
}
