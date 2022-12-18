package automation

import (
	"fmt"
	"golang.org/x/time/rate"
	"strconv"
	"strings"
)

const (
	RateLimitName = "rateLimit"
	BurstName     = "burst"
	InfValue      = "INF"
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
	name        string
	table       *table
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

func (a *rateLimiter) IsEnabled() bool {
	return a.current.limit != rate.Inf
}

func (a *rateLimiter) Reset() {

}

func (a *rateLimiter) Disable() {
}

func (a *rateLimiter) Configure(events string) error {
	return nil
}

func (a *rateLimiter) Adjust(up bool) {
}

func (a *rateLimiter) Value(name string) string {
	if name == "" {
		return ""
	}
	if strings.Index(name, RateLimitName) != -1 {
		if a.current.limit == rate.Inf {
			return InfValue
		}
		return fmt.Sprintf("%v", a.current.limit)
	}
	if strings.Index(name, BurstName) != -1 {
		return strconv.Itoa(a.current.burst)
	}
	return ""
}

func (a *rateLimiter) Allow() bool {
	return false
}
