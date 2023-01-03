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
	IsRetryable(statusCode int) (ok bool, status string)
	SetRateLimiter(limit rate.Limit, burst int)
	AdjustRateLimiter(percentage int) bool
	LimitAndBurst() (rate.Limit, int)
}

type RetryConfig struct {
	limit rate.Limit
	burst int
	wait  time.Duration
	codes []int
}

func NewRetryConfig(validCodes []int, limit rate.Limit, burst int, wait time.Duration) *RetryConfig {
	c := new(RetryConfig)
	c.wait = wait
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
	t.enabled = true
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if config != nil {
		t.config = *config
	}
	t.rateLimiter = rate.NewLimiter(t.config.limit, t.config.burst)
	return t
}

func (r *retry) validate() error {
	if len(r.config.codes) == 0 {
		return errors.New("invalid configuration: RetryController status codes are empty")
	}
	if r.config.limit < 0 {
		return errors.New("invalid configuration: RetryController limit is < 0")
	}
	if r.config.burst < 0 {
		return errors.New("invalid configuration: RetryController burst is < 0")
	}
	return nil
}

func retryState(m map[string]string, r *retry, retried bool) map[string]string {
	var limit rate.Limit = -1
	var burst = -1
	var name = ""
	if r != nil {
		name = strconv.FormatBool(retried)
		limit = r.config.limit
		if limit == rate.Inf {
			limit = RateLimitInfValue
		}
		burst = r.config.burst
	}
	if m == nil {
		m = make(map[string]string, 16)
	}
	m[RetryName] = name
	m[RetryRateLimitName] = fmt.Sprintf("%v", limit)
	m[RetryRateBurstName] = strconv.Itoa(burst)
	return m

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
	if r.config.limit == limit && r.config.burst == burst {
		return
	}
	r.table.setRetryRateLimit(r.name, limit, burst)
}

func (r *retry) IsRetryable(statusCode int) (bool, string) {
	if !r.IsEnabled() {
		return false, NotEnabledFlag
	}
	if statusCode < http.StatusInternalServerError {
		return false, ""
	}
	if !r.rateLimiter.Allow() {
		return false, RateLimitFlag
	}
	for _, code := range r.config.codes {
		if code == statusCode {
			jitter := time.Duration(r.rand.Int31n(1000))
			time.Sleep(r.config.wait + jitter)
			return true, ""
		}
	}
	return false, ""
}

func (r *retry) AdjustRateLimiter(percentage int) bool {
	newLimit, ok := limitAdjust(float64(r.config.limit), percentage)
	if !ok {
		return false
	}
	newBurst, ok1 := burstAdjust(r.config.burst, percentage)
	if !ok1 {
		return false
	}
	r.table.setRateLimiter(r.name, RateLimiterConfig{limit: rate.Limit(newLimit), burst: newBurst})
	return true
}

func (r *retry) LimitAndBurst() (rate.Limit, int) {
	return r.config.limit, r.config.burst
}
