package automation

import (
	"golang.org/x/time/rate"
	"time"
)

type RetryController interface {
	Controller
	//Allow() bool
	//SetLimit(limit rate.Limit)
	//SetBurst(burst int)
}

type RetryConfig struct {
	limit   rate.Limit
	burst   int
	steps   []int
	timeout int
	hold    time.Duration
}

func NewRetryConfig(limit rate.Limit, burst int, steps []int, timeout int, hold time.Duration) *RetryConfig {
	validateLimiter(&limit, &burst)
	c := new(RetryConfig)
	c.limit = limit
	c.burst = burst
	return c
}

type retry struct {
	name        string
	table       *table
	enabled     bool
	tripped     bool
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
	t.enabled = true
	if config != nil {
		t.current = *config
	}
	//t.rateLimiter = rate.NewLimiter(t.current.limit, t.current.burst)
	return t
}

func (r *retry) IsEnabled() bool { return r.enabled }

func (r *retry) Disable() { r.table.enableRetry(r.name, false) }

func (r *retry) Enable() { r.table.enableRetry(r.name, false) }

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
