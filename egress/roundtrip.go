package egress

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/accesslog"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"time"
)

var Routes = route.NewTable()

type wrapper struct {
	rt http.RoundTripper
}

func (w *wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start = time.Now()
	var flags string
	var resp *http.Response
	var err error

	// No panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid egress round tripper configuration : http.RoundTripper is nil")
	}
	rt := Routes.Lookup(req)
	if rt.Allow() {
		if rt.IsTimeout() {
			ctx, cancel := context.WithTimeout(req.Context(), rt.Duration())
			defer cancel()
			req = req.Clone(ctx)
		}
		resp, err = w.rt.RoundTrip(req)
		if err != nil && errors.As(err, &context.DeadlineExceeded) {
			err = nil
			flags = accesslog.UpstreamTimeoutFlag
			resp = &http.Response{Request: req, StatusCode: http.StatusGatewayTimeout}
		}
	} else {
		flags = accesslog.RateLimitFlag
		resp = &http.Response{Request: req, StatusCode: http.StatusServiceUnavailable}
	}
	if rt.IsLogging() || accesslog.IsExtract() {
		accesslog.WriteEgress(start, time.Since(start), rt, req, resp, flags)
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
