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
	Duration() (bool, time.Duration)
	SetTimeout(timeout int, enable bool)
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
	if !t.IsEnabled() {
		return
	}
	t.table.enableTimeout(t.name, false)
}

func (t *timeout) Enable() {
	if t.IsEnabled() || t.current.timeout <= 0 {
		return
	}
	t.table.enableTimeout(t.name, true)
}

func (t *timeout) Reset()                    {}
func (t *timeout) Adjust(any)                {}
func (t *timeout) Configure(Attribute) error { return nil }

func (t *timeout) Attribute(name string) Attribute {
	if strings.Index(name, TimeoutName) != -1 {
		return NewAttribute(TimeoutName, t.current.timeout)
	}
	return nilAttribute(name)
}

func (t *timeout) Duration() (bool, time.Duration) {
	if !t.IsEnabled() {
		return false, 0
	}
	if t.current.timeout <= 0 {
		return true, 0
	}
	return true, time.Duration(t.current.timeout) * time.Millisecond
}

func (t *timeout) SetTimeout(timeout int, enable bool) {
	if t.current.timeout == timeout || timeout <= 0 {
		return
	}
	t.table.setTimeout(t.name, timeout, enable)
}
