package actuator

import (
	"golang.org/x/time/rate"
	"net/http"
	"strings"
)

const (
	RateLimitName  = "rateLimit"
	RateBurstName  = "burst"
	StatusCodeName = "statusCode"
	InfValue       = "-1"
	DefaultBurst   = 1
)

type RateLimiterController interface {
	Attribute(name string) Attribute
	Allow() bool
	StatusCode() int
	SetLimit(limit rate.Limit)
	SetBurst(burst int)
	SetRateLimiter(limit rate.Limit, burst int)
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

/*
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
	r.table.enableRateLimiter(r.name, true)
}

func (r *rateLimiter) Reset() {
	r.table.setRateLimiter(r.name, r.defaultConfig)
}

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
	case RateBurstName:
		if val, ok := attr.Value().(int); ok {
			r.SetBurst(val)
		}
	default:
		errors.New(fmt.Sprintf("invalid attribute name: name not found [%v]", attr.Name()))
	}
	return nil
}


*/
func (r *rateLimiter) Attribute(name string) Attribute {
	if strings.Index(name, RateLimitName) != -1 {
		return NewAttribute(RateLimitName, r.config.limit)
	}
	if strings.Index(name, RateBurstName) != -1 {
		return NewAttribute(RateBurstName, r.config.burst)
	}
	if strings.Index(name, StatusCodeName) != -1 {
		return NewAttribute(StatusCodeName, r.config.statusCode)
	}
	return nilAttribute(name)
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
