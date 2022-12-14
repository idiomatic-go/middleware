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

	//if req == nil {
	//	return nil, errors.New("invalid argument : http.Request is nil on RoundTrip call")
	//}
	//if w == nil || w.rt == nil {
	//	return nil, errors.New("invalid wrapper : http.RoundTripper is nil")
	//}
	route := Routes.Lookup(req)
	if route.Allow() {
		if route.IsTimeout() && req != nil {
			ctx, cancel := context.WithTimeout(req.Context(), route.Duration())
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
	if route.IsLogging() || accesslog.IsExtract() {
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
