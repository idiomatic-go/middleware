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
	var retry = false

	// No panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid egress round tripper configuration : http.RoundTripper is nil")
	}
	act := actuator.Egress.Lookup(req)
	if rlc, ok := act.RateLimiter(); ok && !rlc.Allow() {
		resp = &http.Response{Request: req, StatusCode: rlc.StatusCode()}
		act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, false, req, resp, actuator.RateLimitFlag)
		return resp, nil
	}
	resp, err = w.exchange(act, req)
	if rc, ok := act.Retry(); ok && err == nil {
		retry, statusFlags = rc.IsRetryable(resp.StatusCode)
		if retry {
			act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, false, req, resp, statusFlags)
			start = time.Now()
			resp, err = w.exchange(act, req)
		}
	}
	if w.deadlineExceeded(err) {
		err = nil
		statusFlags = actuator.UpstreamTimeoutFlag
		resp = &http.Response{Request: req, StatusCode: http.StatusGatewayTimeout}
	}
	act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, retry, req, resp, statusFlags)
	return resp, err
}

func (w *wrapper) exchange(act actuator.Actuator, req *http.Request) (resp *http.Response, err error) {
	if tc, ok := act.Timeout(); ok {
		ctx, cancel := context.WithTimeout(req.Context(), tc.Duration())
		defer cancel()
		req = req.Clone(ctx)
	}
	resp, err = w.rt.RoundTrip(req)
	if w.deadlineExceeded(err) {
		resp = &http.Response{Request: req, StatusCode: http.StatusGatewayTimeout}
		err = nil
	}
	return
}

func (w *wrapper) deadlineExceeded(err error) bool {
	return err != nil && errors.As(err, &context.DeadlineExceeded)
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
