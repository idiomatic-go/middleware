package automation

import (
	"golang.org/x/time/rate"
	"time"
)

type CircuitBreakerController interface {
	Controller
	Allow() (bool, bool)
	//SetLimit(limit rate.Limit)
	//SetBurst(burst int)
}

type CircuitBreakerConfig struct {
	limit   rate.Limit
	burst   int
	steps   []int
	timeout int
	hold    time.Duration
}

func NewCircuitBreakerConfig(limit rate.Limit, burst int, steps []int, timeout int, hold time.Duration) *CircuitBreakerConfig {
	validateLimiter(&limit, &burst)
	c := new(CircuitBreakerConfig)
	c.limit = limit
	c.burst = burst
	return c
}

func cloneCircuitBreaker(curr *circuitBreaker) *circuitBreaker {
	t := new(circuitBreaker)
	*t = *curr
	return t
}

type circuitBreaker struct {
	name        string
	table       *table
	enabled     bool
	tripped     bool
	current     CircuitBreakerConfig
	rateLimiter *rate.Limiter
}

func newCircuitBreaker(name string, table *table, config *CircuitBreakerConfig) *circuitBreaker {
	t := new(circuitBreaker)
	t.name = name
	t.table = table
	t.enabled = true
	if config != nil {
		t.current = *config
	}
	t.rateLimiter = rate.NewLimiter(t.current.limit, t.current.burst)
	return t
}

func (c *circuitBreaker) IsEnabled() bool { return c.enabled }

func (c *circuitBreaker) Disable() {
	if !c.IsEnabled() {
		return
	}
	c.table.enableCircuitBreaker(c.name, false)
}

func (c *circuitBreaker) Enable() {
	if c.IsEnabled() {
		return
	}
	c.table.enableCircuitBreaker(c.name, false)
}

func (c *circuitBreaker) Reset() {}

func (c *circuitBreaker) Adjust(change any) {
	return
}

func (c *circuitBreaker) Configure(attr Attribute) error {
	return nil
}

func (c *circuitBreaker) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (c *circuitBreaker) Allow() (bool, bool) {
	return c.enabled, true
}
