package automation

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type RetryController interface {
	Controller
	IsRetryable(timeout TimeoutController, spent time.Duration, statusCode int) (bool, bool)
	//SetRateLimiter(limit rate.Limit, burst int)
}

type RetryConfig struct {
	limit rate.Limit
	burst int
	//wait  time.Duration
	codes []int
}

func NewRetryConfig(validCodes []int) *RetryConfig {
	//validateLimiter(&limit, &burst)
	c := new(RetryConfig)
	c.limit = rate.Inf
	c.burst = DefaultBurst
	//c.wait = wait
	c.codes = validCodes
	return c
}

type retry struct {
	name        string
	table       *table
	enabled     bool
	current     RetryConfig
	rateLimiter *rate.Limiter
}

func cloneRetry(curr *retry) *retry {
	t := new(retry)
	*t = *curr
	return t
}

func newRetry(name string, table *table, config *RetryConfig) *retry {
	t := new(retry)
	t.name = name
	t.table = table
	if config != nil {
		t.enabled = true
		t.current = *config
	}
	t.rateLimiter = rate.NewLimiter(t.current.limit, t.current.burst)
	return t
}

func (r *retry) IsEnabled() bool { return r.enabled }

func (r *retry) Disable() {
	if !r.IsEnabled() {
		return
	}
	r.table.enableRetry(r.name, false)
}

func (r *retry) Enable() {
	if r.IsEnabled() {
		return
	}
	r.table.enableRetry(r.name, false)
}

func (r *retry) Reset() {}

func (r *retry) Adjust(change any) {
	return
}

func (r *retry) Configure(attr Attribute) error {
	return nil
}

func (r *retry) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (r *retry) IsRetryable(timeout TimeoutController, spent time.Duration, statusCode int) (bool, bool) {
	if !r.IsEnabled() {
		return false, false
	}
	if timeout == nil {
		return true, false
	}
	if statusCode < http.StatusInternalServerError {
		return true, false
	}
	wait := time.Duration(time.Second * 1)
	for _, code := range r.current.codes {
		if code == statusCode {
			time.Sleep(wait)
			return true, true
		}
	}
	return true, false
}
