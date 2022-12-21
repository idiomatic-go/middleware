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
	name          string
	table         *table
	enabled       bool
	canary        bool
	defaultConfig RateLimiterConfig
	currentConfig RateLimiterConfig
	canaryConfig  RateLimiterConfig
	rateLimiter   *rate.Limiter
}

func cloneRateLimiter(curr *rateLimiter) *rateLimiter {
	newLimiter := new(rateLimiter)
	*newLimiter = *curr
	return newLimiter
}

func newRateLimiter(name string, config []*RateLimiterConfig, table *table) *rateLimiter {
	var t = rateLimiter{
		name:          name,
		table:         table,
		enabled:       false,
		canary:        false,
		defaultConfig: RateLimiterConfig{rate.Inf, DefaultBurst},
		currentConfig: RateLimiterConfig{rate.Inf, DefaultBurst},
		canaryConfig:  RateLimiterConfig{rate.Inf, DefaultBurst},
		rateLimiter:   nil,
	}
	if len(config) > 0 {
		t.currentConfig = *config[0]
		if len(config) == 2 {
			t.canaryConfig = *config[1]
		}
		t.defaultConfig = t.currentConfig
	}
	t.rateLimiter = rate.NewLimiter(t.currentConfig.limit, t.currentConfig.burst)
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

func (r *rateLimiter) IsEnabled() bool { return r.enabled }

func (r *rateLimiter) Disable() { r.table.enableRateLimiter(r.name, false) }

func (r *rateLimiter) Enable() { r.table.enableRateLimiter(r.name, false) }

func (r *rateLimiter) Reset() {
	// TODO : set isCanary to false
}

func (r *rateLimiter) Configure(items ...Attribute) error {
	// TODO : how to set canary
	return nil
}

func (r *rateLimiter) Adjust(up bool) {
}

func (r *rateLimiter) Attribute(name string) Attribute {
	if strings.Index(name, RateLimitName) != -1 {
		return NewAttribute(RateLimitName, r.currentConfig.limit)
	}
	if strings.Index(name, BurstName) != -1 {
		return NewAttribute(BurstName, r.currentConfig.burst)
	}
	if strings.Index(name, CanaryName) != -1 {
		return NewAttribute(CanaryName, r.canary)
	}
	return nilAttribute(name)
}

func (r *rateLimiter) Allow() bool {
	return r.rateLimiter.Allow()
}

func (r *rateLimiter) SetLimit(limit rate.Limit) {
	r.table.setLimit(r.name, limit)
}
