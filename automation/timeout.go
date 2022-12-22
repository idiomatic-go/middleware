package automation

import (
	"strings"
	"time"
)

const (
	TimeoutName = "timeout"
)

type TimeoutController interface {
	Controller
	Duration() time.Duration
	SetTimeout(timeout int)
}

type TimeoutConfig struct {
	timeout int
}

func NewTimeoutConfig(timeout int) *TimeoutConfig {
	if timeout < 0 {
		timeout = NilValue
	}
	return &TimeoutConfig{timeout: timeout}
}

type timeout struct {
	table   *table
	name    string
	enabled bool
	current TimeoutConfig
}

func cloneTimeout(curr *timeout) *timeout {
	t := new(timeout)
	*t = *curr
	return t
}

func newTimeout(name string, table *table, config *TimeoutConfig) *timeout {
	if config == nil {
		config = NewTimeoutConfig(NilValue)
	}
	t := new(timeout)
	t.table = table
	t.name = name
	t.current.timeout = config.timeout
	t.enabled = t.current.timeout > 0
	return t
}

func (t *timeout) IsEnabled() bool { return t.enabled }

func (t *timeout) Disable() {
	if t.IsEnabled() {
		t.table.enableTimeout(t.name, false)
	}
}

func (t *timeout) Enable() {
	if !t.enabled && t.current.timeout > 0 {
		t.table.enableTimeout(t.name, true)
	}
}

func (t *timeout) Reset() {}

func (t *timeout) Adjust(any) {}

func (t *timeout) Configure(Attribute) error { return nil }

func (t *timeout) Attribute(name string) Attribute {
	if strings.Index(name, TimeoutName) != -1 {
		return NewAttribute(TimeoutName, t.current.timeout)
	}
	return nilAttribute(name)
}

func (t *timeout) Duration() time.Duration {
	if t.current.timeout == NilValue {
		return 0
	}
	return time.Duration(t.current.timeout) * time.Millisecond
}

func (t *timeout) SetTimeout(timeout int) {
	if timeout <= 0 {
		return
	}
	t.table.setTimeout(t.name, timeout)
}
