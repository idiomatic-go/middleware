package egress

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/accesslog"
	"net/http"
	"time"
)

type wrapper struct {
	rt http.RoundTripper
}

func (w wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start time.Time
	var flags string
	var resp *http.Response
	var err error

	route := Lookup(req)
	allow := route.Allow()
	if !allow {
		flags = accesslog.RateLimitFlag
	}
	if route.IsLogging() {
		start = time.Now()
	}
	if allow {
		if route.IsTimeout() {
			ctx, cancel := context.WithTimeout(req.Context(), route.Duration())
			defer cancel()
			req = req.Clone(ctx)
		}
		resp, err = w.RoundTrip(req) //http.DefaultTransport.RoundTrip(req)
	} else {
		resp = &http.Response{Request: req, StatusCode: 503}
	}
	// TODO : check on timeout
	if err != nil && errors.As(err, context.DeadlineExceeded) {
		flags = accesslog.UpstreamTmeoutFlag
		resp.StatusCode = 504
	}
	if route.IsLogging() {
		accesslog.WriteEgress(start, time.Since(start), route, req, resp, flags)
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
