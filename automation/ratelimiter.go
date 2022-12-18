package automation

import (
	"golang.org/x/time/rate"
	"strings"
)

const (
	RateLimitName = "rateLimit"
	BurstName     = "burst"
	CanaryName    = "canary"
	InfValue      = "INF"
	DefaultBurst  = 1
)

type RateLimiterController interface {
	Controller
	Value(name string) string
	Allow() bool
}

type RateLimiterConfig struct {
	limit rate.Limit
	burst int
}

func NewRateLimiterConfig(limit rate.Limit, burst int, canaryLimit rate.Limit, canaryBurst int) (configs []*RateLimiterConfig) {
	validateLimiter(&limit, &burst)
	c := new(RateLimiterConfig)
	c.limit = limit
	c.burst = burst
	configs = append(configs, c)

	validateLimiter(&canaryLimit, &canaryBurst)
	c = new(RateLimiterConfig)
	c.limit = canaryLimit
	c.burst = canaryBurst
	configs = append(configs, c)
	return
}

type rateLimiter struct {
	name        string
	table       *table
	isCanary    bool
	defaultC    RateLimiterConfig
	current     RateLimiterConfig
	canary      RateLimiterConfig
	rateLimiter *rate.Limiter
}

func cloneRateLimiter(act Actuator) *rateLimiter {
	if act == nil {
		return nil
	}
	t := new(rateLimiter)
	s := act.RateLimiter().(*rateLimiter)
	*t = *s
	return t
}

func newRateLimiter(name string, config []*RateLimiterConfig, table *table) *rateLimiter {
	var t = rateLimiter{
		name:        name,
		table:       table,
		defaultC:    RateLimiterConfig{rate.Inf, DefaultBurst},
		current:     RateLimiterConfig{rate.Inf, DefaultBurst},
		canary:      RateLimiterConfig{rate.Inf, DefaultBurst},
		rateLimiter: nil,
	}
	if len(config) > 0 {
		t.current = *config[0]
		if len(config) == 2 {
			t.canary = *config[1]
		}
		t.defaultC = t.current
	}
	t.rateLimiter = rate.NewLimiter(t.current.limit, t.current.burst)
	return &t
}

func validateLimiter(max *rate.Limit, burst *int) {
	if max != nil && *max < 0 {
		*max = rate.Inf
	}
	if burst != nil && *burst < 0 {
		*burst = DefaultBurst
	}
}

func (r *rateLimiter) IsEnabled() bool {
	return r.current.limit != rate.Inf
}

func (r *rateLimiter) Reset() {
	// TODO : set isCanary to false
}

func (r *rateLimiter) Disable() {
}

func (r *rateLimiter) Configure(items ...Attribute) error {
	// TODO : how to set canary
	return nil
}

func (r *rateLimiter) Adjust(up bool) {
}

func (r *rateLimiter) Attribute(name string) Attribute {
	if strings.Index(name, RateLimitName) != -1 {
		return NewAttribute(RateLimitName, r.current.limit)
	}
	if strings.Index(name, BurstName) != -1 {
		return NewAttribute(BurstName, r.current.burst)
	}
	if strings.Index(name, CanaryName) != -1 {
		return NewAttribute(CanaryName, r.isCanary)
	}
	return nilAttribute(name)
}

func (r *rateLimiter) Value(name string) string {
	if name == "" {
		return ""
	}

	return ""
}

func (r *rateLimiter) Allow() bool {
	return r.rateLimiter.Allow()
}
