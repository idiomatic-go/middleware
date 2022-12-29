package actuator

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
// https://github.com/keikoproj/inverse-exp-backoff

type RetryController interface {
	IsEnabled() bool
	Enable()
	Disable()
	SetRateLimiter(limit rate.Limit, burst int)
	IsRetryable(statusCode int) (ok bool, status string)
}

type RetryConfig struct {
	limit rate.Limit
	burst int
	//wait  time.Duration
	codes []int
}

func NewRetryConfig(validCodes []int, limit rate.Limit, burst int) *RetryConfig {
	c := new(RetryConfig)
	//c.wait = wait
	c.limit = limit
	c.burst = burst
	c.codes = validCodes
	return c
}

type retry struct {
	name        string
	table       *table
	enabled     bool
	rand        *rand.Rand
	config      RetryConfig
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
	t.enabled = false
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if config != nil {
		t.config = *config
	}
	t.rateLimiter = rate.NewLimiter(t.config.limit, t.config.burst)
	return t
}

func (r *retry) validate() error {
	if len(r.config.codes) == 0 {
		return errors.New("invalid configuration: retry controller status codes are empty")
	}
	if r.config.limit <= 0 || r.config.limit == rate.Inf {
		return errors.New("invalid configuration: retry controller limit is <= 0 or == rate.Inf")
	}
	return nil
}

func retryPut(r RetryController, retried bool, m map[string]string) {
	var limit rate.Limit = -1
	var burst = -1
	var name = ""
	if r != nil {
		name = strconv.FormatBool(retried)
		limit = r.(*retry).config.limit
		if limit == rate.Inf {
			limit = RateLimitInfValue
		}
		burst = r.(*retry).config.burst
	}
	if m != nil {
		m[RetryName] = name
		m[RetryRateLimitName] = fmt.Sprintf("%v", limit)
		m[RetryRateBurstName] = strconv.Itoa(burst)
	}

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
	r.table.enableRetry(r.name, true)
}

func (r *retry) SetRateLimiter(limit rate.Limit, burst int) {
	if r.config.limit == limit {
		return
	}
	r.table.setRetryRateLimit(r.name, limit, burst)
}

func (r *retry) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (r *retry) IsRetryable(statusCode int) (bool, string) {
	if !r.IsEnabled() {
		return false, NotEnabledFlag
	}
	if statusCode < http.StatusInternalServerError {
		return false, InvalidStatusCodeFlag
	}
	if !r.rateLimiter.Allow() {
		return false, RateLimitFlag
	}
	for _, code := range r.config.codes {
		if code == statusCode {
			//jitter := time.Duration(r.rand.Int31n(1000))
			//time.Sleep(r.current.wait + jitter)
			return true, ""
		}
	}
	return false, InvalidStatusCodeFlag
}
