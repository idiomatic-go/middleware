package actuator

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
)

const (
	InfValue     = "-1"
	DefaultBurst = 1
)

type RateLimiterController interface {
	Allow() bool
	StatusCode() int
	SetLimit(limit rate.Limit)
	SetBurst(burst int)
	SetRateLimiter(limit rate.Limit, burst int)
	LimitAndBurst() (rate.Limit, int)
}

type RateLimiterConfig struct {
	limit      rate.Limit
	burst      int
	statusCode int
}

func NewRateLimiterConfig(limit rate.Limit, burst int, statusCode int) *RateLimiterConfig {
	validateLimiter(&limit, &burst)
	c := new(RateLimiterConfig)
	c.limit = limit
	c.burst = burst
	if statusCode <= 0 {
		statusCode = http.StatusTooManyRequests
	}
	c.statusCode = statusCode
	return c
}

type rateLimiter struct {
	name        string
	table       *table
	config      RateLimiterConfig
	rateLimiter *rate.Limiter
}

func cloneRateLimiter(curr *rateLimiter) *rateLimiter {
	newLimiter := new(rateLimiter)
	*newLimiter = *curr
	return newLimiter
}

func newRateLimiter(name string, table *table, config *RateLimiterConfig) *rateLimiter {
	t := new(rateLimiter)
	t.name = name
	t.table = table
	t.config = RateLimiterConfig{limit: rate.Inf, burst: DefaultBurst}
	if config != nil {
		t.config = *config
	}
	t.rateLimiter = rate.NewLimiter(t.config.limit, t.config.burst)
	return t
}

func validateLimiter(max *rate.Limit, burst *int) {
	if max != nil && *max < 0 {
		*max = rate.Inf
	}
	if burst != nil && *burst < 0 {
		*burst = DefaultBurst
	}
}

func (r *rateLimiter) validate() error {
	return nil
}

func rateLimiterState(m map[string]string, r RateLimiterController) {
	var limit rate.Limit = -1
	var burst = -1

	if r != nil {
		limit = r.(*rateLimiter).config.limit
		if limit == rate.Inf {
			limit = RateLimitInfValue
		}
		burst = r.(*rateLimiter).config.burst
	}
	m[RateLimitName] = fmt.Sprintf("%v", limit)
	m[RateBurstName] = strconv.Itoa(burst)
}

func (r *rateLimiter) Allow() bool {
	if r.config.limit == rate.Inf {
		return true
	}
	return r.rateLimiter.Allow()
}

func (r *rateLimiter) StatusCode() int {
	return r.config.statusCode
}

func (r *rateLimiter) SetLimit(limit rate.Limit) {
	if r.config.limit == limit {
		return
	}
	r.table.setRateLimit(r.name, limit)
}

func (r *rateLimiter) SetBurst(burst int) {
	if r.config.burst == burst {
		return
	}
	r.table.setRateBurst(r.name, burst)
}

func (r *rateLimiter) SetRateLimiter(limit rate.Limit, burst int) {
	validateLimiter(&limit, &burst)
	if r.config.limit == limit && r.config.burst == burst {
		return
	}
	r.table.setRateLimiter(r.name, RateLimiterConfig{limit: limit, burst: burst})
}

func (r *rateLimiter) LimitAndBurst() (rate.Limit, int) {
	return r.config.limit, r.config.burst
}
