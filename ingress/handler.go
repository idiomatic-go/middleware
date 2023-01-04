package ingress

import (
	"github.com/felixge/httpsnoop"
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
	"time"
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

func HttpMetricsHandler(appHandler http.Handler) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusFlags := ""
		start := time.Now()
		act := actuator.IngressTable.Host()

		if rlc, ok := act.RateLimiter(); ok && !rlc.Allow() {
			act.LogIngress(start, time.Since(start), r, rlc.StatusCode(), 0, actuator.RateLimitFlag)
			return
		}
		m := httpsnoop.CaptureMetrics(appHandler, w, r)
		/*
			log.Printf("%s %s (code=%d dt=%s written=%d)",r.Method,r.URL,m.Code,m.Duration,m.Written,)
		*/
		act = actuator.IngressTable.Lookup(r)
		if toc, ok := act.Timeout(); ok && m.Code == http.StatusServiceUnavailable {
			m.Code = toc.StatusCode()
			statusFlags = actuator.HostTimeoutFlag
		}
		act.LogIngress(start, time.Since(start), r, m.Code, m.Written, statusFlags)
	})
	return wrappedH
}
