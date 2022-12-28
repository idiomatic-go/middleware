package egress

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
	"strconv"
	"time"
)

type wrapper struct {
	rt http.RoundTripper
}

func (w *wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var start = time.Now()
	var statusFlags string
	var retry = false

	// No panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid egress round tripper configuration : http.RoundTripper is nil")
	}
	act := actuator.Egress.Lookup(req)

	if rlc, ok := act.RateLimiter(); ok && !rlc.Allow() {
		resp := &http.Response{Request: req, StatusCode: rlc.StatusCode()}
		act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, "null", req, resp, actuator.RateLimitFlag)
		return resp, nil
	}
	tc, _ := act.Timeout()
	resp, err := w.exchange(tc, req)
	if err != nil {
		return resp, err
	}
	if rc, ok := act.Retry(); ok {
		retry, statusFlags = rc.IsRetryable(resp.StatusCode)
		if retry {
			act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, strconv.FormatBool(retry), req, resp, "")
			start = time.Now()
			resp, err = w.exchange(tc, req)
		}
	}
	act.Logger().LogAccess(actuator.EgressTraffic, start, time.Since(start), act, strconv.FormatBool(retry), req, resp, statusFlags)
	return resp, err
}

func (w *wrapper) exchange(tc actuator.TimeoutController, req *http.Request) (resp *http.Response, err error) {
	if tc == nil {
		return w.rt.RoundTrip(req)
	}
	ctx, cancel := context.WithTimeout(req.Context(), tc.Duration())
	defer cancel()
	req = req.Clone(ctx)
	resp, err = w.rt.RoundTrip(req)
	if w.deadlineExceeded(err) {
		resp = &http.Response{Request: req, StatusCode: tc.StatusCode()}
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
