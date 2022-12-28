package egress

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
	"time"
)

type wrapper struct {
	rt http.RoundTripper
}

func (w *wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start = time.Now()
	var statusFlags string
	var resp *http.Response
	var err error

	// No panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid egress round tripper configuration : http.RoundTripper is nil")
	}
	act := actuator.Egress.Lookup(req)
	if act.RateLimiter() != nil && act.RateLimiter().Allow() {
		if act.Timeout() != nil {
			ctx, cancel := context.WithTimeout(req.Context(), act.Timeout().Duration())
			defer cancel()
			req = req.Clone(ctx)
		}
		resp, err = w.rt.RoundTrip(req)
		if err != nil && errors.As(err, &context.DeadlineExceeded) {
			err = nil
			statusFlags = actuator.UpstreamTimeoutFlag
			resp = &http.Response{Request: req, StatusCode: http.StatusGatewayTimeout}
		}
	} else {
		statusFlags = actuator.RateLimitFlag
		resp = &http.Response{Request: req, StatusCode: act.RateLimiter().StatusCode()}
	}
	act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, req, resp, statusFlags)
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
