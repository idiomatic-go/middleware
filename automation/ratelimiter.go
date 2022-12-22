package automation

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
)

const (
	RateLimitName = "rateLimit"
	BurstName     = "burst"
	InfValue      = "INF"
	DefaultBurst  = 1
)

type RateLimiterController interface {
	Controller
	Allow() (bool, bool)
	StatusCode() (bool, int)
	SetLimit(limit rate.Limit)
	SetBurst(burst int)
}

type RateLimiterConfig struct {
	limit      rate.Limit
	burst      int
	step       int
	statusCode int
}

func NewRateLimiterConfig(limit rate.Limit, burst int, statusCode int, step int) *RateLimiterConfig {
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
	name          string
	table         *table
	enabled       bool
	defaultConfig RateLimiterConfig
	currentConfig RateLimiterConfig
	rateLimiter   *rate.Limiter
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
	t.enabled = true
	t.currentConfig = RateLimiterConfig{limit: rate.Inf, burst: DefaultBurst}
	if config != nil {
		t.currentConfig = *config
	}
	t.defaultConfig = t.currentConfig
	t.rateLimiter = rate.NewLimiter(t.currentConfig.limit, t.currentConfig.burst)
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

func (r *rateLimiter) IsEnabled() bool { return r.enabled }

func (r *rateLimiter) Disable() {
	if !r.IsEnabled() {
		return
	}
	r.table.enableRateLimiter(r.name, false)
}

func (r *rateLimiter) Enable() {
	if r.IsEnabled() {
		return
	}
	r.table.enableRateLimiter(r.name, false)
}

func (r *rateLimiter) Reset() { r.table.setRateLimiter(r.name, r.defaultConfig) }

func (r *rateLimiter) Adjust(change any) {
	return
}

func (r *rateLimiter) Configure(attr Attribute) error {
	err := attr.Validate()
	if err != nil {
		return err
	}
	switch attr.Name() {
	case RateLimitName:
		if val, ok := attr.Value().(rate.Limit); ok {
			r.SetLimit(val)
		}
	case BurstName:
		if val, ok := attr.Value().(int); ok {
			r.SetBurst(val)
		}
	default:
		errors.New(fmt.Sprintf("invalid attribute name: name not found [%v]", attr.Name()))
	}
	return nil
}

func (r *rateLimiter) Attribute(name string) Attribute {
	if strings.Index(name, RateLimitName) != -1 {
		return NewAttribute(RateLimitName, r.currentConfig.limit)
	}
	if strings.Index(name, BurstName) != -1 {
		return NewAttribute(BurstName, r.currentConfig.burst)
	}
	return nilAttribute(name)
}

func (r *rateLimiter) Allow() (bool, bool) {
	if !r.IsEnabled() {
		return false, true
	}
	if r.currentConfig.limit == rate.Inf {
		return true, true
	}
	return true, r.rateLimiter.Allow()
}

func (r *rateLimiter) StatusCode() (bool, int) {
	return r.enabled, r.currentConfig.statusCode
}

func (r *rateLimiter) SetLimit(limit rate.Limit) {
	if r.currentConfig.limit == limit {
		return
	}
	r.table.setLimit(r.name, limit)
}

func (r *rateLimiter) SetBurst(burst int) {
	if r.currentConfig.burst == burst {
		return
	}
	r.table.setBurst(r.name, burst)
}

func (r *rateLimiter) SetRateLimiter(limit rate.Limit, burst int) {
	validateLimiter(&limit, &burst)
	if r.currentConfig.limit == limit && r.currentConfig.burst == burst {
		return
	}
	r.table.setRateLimiter(r.name, RateLimiterConfig{limit: limit, burst: burst})
}
