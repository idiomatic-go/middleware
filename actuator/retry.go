package actuator

import (
	"errors"
	"math/rand"
	"net/http"
	"time"
)

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
// https://github.com/keikoproj/inverse-exp-backoff

type RetryController interface {
	Controller
	IsRetryable(statusCode int) (bool, bool)
}

type RetryConfig struct {
	enabled bool
	wait    time.Duration
	codes   []int
}

func NewRetryConfig(validCodes []int, wait time.Duration, enabled bool) *RetryConfig {
	c := new(RetryConfig)
	c.wait = wait
	c.codes = validCodes
	c.enabled = enabled
	return c
}

type retry struct {
	name    string
	table   *table
	enabled bool
	rand    *rand.Rand
	current RetryConfig
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
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	t.enabled = false
	if config != nil {
		t.current = *config
		t.enabled = config.enabled
	}
	return t
}

func (r *retry) validate() error {
	if len(r.current.codes) == 0 {
		return errors.New("invalid configuration: retry controller status codes are empty")
	}
	return nil
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

func (r *retry) Reset() {}

func (r *retry) Adjust(change any) {
	return
}

func (r *retry) Configure(attr Attribute) error {
	return nil
}

func (r *retry) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (r *retry) IsRetryable(statusCode int) (bool, bool) {
	if !r.IsEnabled() {
		return false, false
	}
	if statusCode < http.StatusInternalServerError {
		return true, false
	}
	for _, code := range r.current.codes {
		if code == statusCode {
			jitter := time.Duration(r.rand.Int31n(1000))
			time.Sleep(r.current.wait + jitter)
			return true, true
		}
	}
	return true, false
}
