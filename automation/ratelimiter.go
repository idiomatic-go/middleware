package automation

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"strings"
	"time"
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
	SetLimit(limit rate.Limit)
	SetBurst(burst int)
}

type RateLimiterConfig struct {
	enabled bool
	startup bool
	limit   rate.Limit
	burst   int
	steps   []int
	hold    time.Duration
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
		enabled:       true,
		canary:        false,
		defaultConfig: RateLimiterConfig{limit: rate.Inf, burst: DefaultBurst},
		currentConfig: RateLimiterConfig{limit: rate.Inf, burst: DefaultBurst},
		canaryConfig:  RateLimiterConfig{limit: rate.Inf, burst: DefaultBurst},
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
	r.table.setRateLimiter(r.name, r.defaultConfig, false)
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
	case BurstName:
		if val, ok := attr.Value().(int); ok {
			r.SetBurst(val)
		}
	case CanaryName:
		if val, ok := attr.Value().(bool); ok {
			if val == false && r.canary {
				r.Reset()
				return nil
			}
			if val == true && !r.canary {
				r.SetCanary()
				return nil
			}
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

func (r *rateLimiter) SetBurst(burst int) {
	r.table.setBurst(r.name, burst)
}

func (r *rateLimiter) SetCanary() {
	r.table.setRateLimiter(r.name, r.canaryConfig, true)
}

func (r *rateLimiter) SetRateLimiter(limit rate.Limit, burst int) {
	validateLimiter(&limit, &burst)
	r.table.setRateLimiter(r.name, RateLimiterConfig{limit: limit, burst: burst}, false)
}
