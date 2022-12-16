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

type timeoutAction struct {
	t       *table
	Default int
	current int
}

func newTimeoutAction(timeout int, t *table) *timeoutAction {
	if timeout <= 0 {
		timeout = NilValue
	}
	return &timeoutAction{Default: timeout, current: timeout, t: t}
}

func (a *timeoutAction) Name() string {
	return TimeoutName
}

func (a *timeoutAction) IsEnabled() bool {
	return a.current != NilValue
}

func (a *timeoutAction) Reset() {
	a.configure(nil)
}

func (a *timeoutAction) Disable() {
	a.configure(NilValue)
}

func (a *timeoutAction) Configure(v ...any) {
	if len(v) == 0 {
		return
	}
	a.configure(v)
}

func (a *timeoutAction) Duration() time.Duration {
	if a.current == NilValue {
		return 0
	}
	return time.Duration(a.current) * time.Millisecond
}

func (a *timeoutAction) configure(v ...any) {
	a.t.configureTimeout(v)

}
