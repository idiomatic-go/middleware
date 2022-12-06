package egress

import (
	"context"
	"github.com/idiomatic-go/middleware/accesslog"
	"net/http"
	"time"
)

type wrapper struct {
	rt http.RoundTripper
}

func (w wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start time.Time

	route := Lookup(req)
	if route.IsLogging() {
		start = time.Now()
	}
	if route.IsTimeout() {
		ctx, cancel := context.WithTimeout(req.Context(), route.Duration())
		defer cancel()
		req = req.Clone(ctx)
	}
	resp, err := w.RoundTrip(req) //http.DefaultTransport.RoundTrip(req)
	if route.IsLogging() {
		accesslog.LogEgress(route, start, time.Since(start), req, resp, err)
	}
	return resp, err
}

func EnableDefaultHttpClient() {
	if http.DefaultClient.Transport == nil {
		http.DefaultClient.Transport = &wrapper{http.DefaultTransport}
	} else {
		http.DefaultClient.Transport = EnableRoundTrip(http.DefaultClient.Transport)
	}
}

func EnableRoundTrip(rt http.RoundTripper) http.RoundTripper {
	return &wrapper{rt}
}
