package automation

import (
	"time"
)

const (
	TimeoutName = "timeout"
)

type TimeoutAction interface {
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
	// TODO : validate status code
	return &TimeoutConfig{timeout: timeout, statusCode: statusCode}
}

type timeout struct {
	defaultC TimeoutConfig
	current  TimeoutConfig
}

func newTimeout(c *TimeoutConfig) *timeout {
	if c == nil {
		c = NewTimeoutConfig(NilValue, NilValue)
	}
	action := new(timeout)
	if c.timeout <= 0 {
		c.timeout = NilValue
	}
	action.current.timeout = c.timeout
	action.current.statusCode = c.statusCode
	action.defaultC = action.current
	return action
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

func (a *timeout) Configure(v ...any) {
	//if len(v) != 0 {
	//		a.configure(v)
	//}
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
