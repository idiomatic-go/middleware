package host

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
	var retry = false

	// !panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid handler round tripper configuration : http.RoundTripper is nil")
	}
	act := actuator.EgressTable.Lookup(req)
	act.UpdateHeaders(req)
	if rlc, ok := act.RateLimiter(); ok && !rlc.Allow() {
		resp := &http.Response{Request: req, StatusCode: rlc.StatusCode()}
		act.LogEgress(start, time.Since(start), req, resp, actuator.RateLimitFlag, false)
		return resp, nil
	}
	tc, _ := act.Timeout()
	resp, err, statusFlags := w.exchange(tc, req)
	if err != nil {
		return resp, err
	}
	if rc, ok := act.Retry(); ok {
		prevFlags := statusFlags
		retry, statusFlags = rc.IsRetryable(resp.StatusCode)
		if retry {
			act.LogEgress(start, time.Since(start), req, resp, prevFlags, false)
			start = time.Now()
			resp, err, statusFlags = w.exchange(tc, req)
		}
	}
	act.LogEgress(start, time.Since(start), req, resp, statusFlags, retry)
	return resp, err
}

func (w *wrapper) exchange(tc actuator.TimeoutController, req *http.Request) (resp *http.Response, err error, statusFlags string) {
	if tc == nil {
		resp, err = w.rt.RoundTrip(req)
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), tc.Duration())
	defer cancel()
	req = req.Clone(ctx)
	resp, err = w.rt.RoundTrip(req)
	if w.deadlineExceeded(err) {
		resp = &http.Response{Request: req, StatusCode: tc.StatusCode()}
		err = nil
		statusFlags = actuator.UpstreamTimeoutFlag
	}
	return
}

func (w *wrapper) deadlineExceeded(err error) bool {
	return err != nil && errors.As(err, &context.DeadlineExceeded)
}

func WrapDefaultTransport() {
	if http.DefaultClient.Transport == nil {
		http.DefaultClient.Transport = &wrapper{http.DefaultTransport}
	} else {
		http.DefaultClient.Transport = WrapRoundTripper(http.DefaultClient.Transport)
	}
}

func WrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &wrapper{rt}
}
