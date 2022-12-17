package automation

import (
	"time"
)

const (
	TimeoutName = "timeout"
)

type TimeoutController interface {
	Timeout() int
	StatusCode(defaultStatusCode int) int
	Duration() time.Duration
}

type TimeoutConfig struct {
	timeout    int
	statusCode int
}

func NewTimeoutConfig(timeout int, statusCode int) *TimeoutConfig {
	if timeout <= 0 {
		timeout = NilValue
	}
	if statusCode <= 0 {
		timeout = NilValue
	}
	// TODO : validate status code
	return &TimeoutConfig{timeout: timeout, statusCode: statusCode}
}

func NewTimeout(duration int, statusCode int) TimeoutController {
	t := new(timeout)
	if duration <= 0 {
		duration = NilValue
	}
	t.current.timeout = duration
	t.current.statusCode = statusCode
	t.defaultC = t.current
	return t
}

type timeout struct {
	defaultC TimeoutConfig
	current  TimeoutConfig
}

func newTimeout(c *TimeoutConfig) *timeout {
	if c == nil {
		c = NewTimeoutConfig(NilValue, NilValue)
	}
	t := new(timeout)
	if c.timeout <= 0 {
		c.timeout = NilValue
	}
	t.current.timeout = c.timeout
	t.current.statusCode = c.statusCode
	t.defaultC = t.current
	return t
}

func (a *timeout) Name() string {
	return TimeoutName
}

func (a *timeout) IsEnabled() bool {
	return a.current.timeout != NilValue
}

func (a *timeout) Reset() {
	//configureTimeout(timeout,statusCode)a.configure(nil)
}

func (a *timeout) Disable() {
	//a.configure(NilValue)
}

func (a *timeout) Configure(s string) {
	//if len(v) != 0 {
	//		a.configure(v)
	//}
}

func (a *timeout) Adjust(up bool) {

}

func (a *timeout) State(s string) []string {

	return nil
}

func (a *timeout) Timeout() int {
	return a.current.timeout
}

func (a *timeout) StatusCode(defaultStatusCode int) int {
	if a.current.statusCode == NilValue {
		return defaultStatusCode
	}
	return a.current.statusCode
}

func (a *timeout) Duration() time.Duration {
	if a.current.timeout == NilValue {
		return 0
	}
	return time.Duration(a.current.timeout) * time.Millisecond
}

func (a *timeout) configure(v ...any) {
	//timeout := NilValue
	//statusCode := NilValue

	//configureTimeout(timeout,statusCode)

}
