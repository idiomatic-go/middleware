package actuator

import (
	"errors"
	"strings"
	"time"
)

const (
	TimeoutName = "timeout"
)

type TimeoutController interface {
	Controller
	Duration() time.Duration
	SetTimeout(timeout time.Duration)
}

type TimeoutConfig struct {
	timeout time.Duration
}

func NewTimeoutConfig(timeout time.Duration) *TimeoutConfig {
	return &TimeoutConfig{timeout: timeout}
}

type timeout struct {
	table         *table
	name          string
	enabled       bool
	defaultConfig TimeoutConfig
	currentConfig TimeoutConfig
}

func cloneTimeout(curr *timeout) *timeout {
	t := new(timeout)
	*t = *curr
	return t
}

func newTimeout(name string, table *table, config *TimeoutConfig) *timeout {
	//if config == nil {
	//	config = NewTimeoutConfig(NilValue)
	//}
	t := new(timeout)
	t.table = table
	t.name = name
	if config != nil {
		t.currentConfig = *config
		t.defaultConfig = *config
	}
	t.enabled = true
	return t
}

func (t *timeout) validate() error {
	if t.currentConfig.timeout <= 0 {
		return errors.New("invalid configuration: TimeoutController duration cannot be <= 0")
	}
	return nil
}

func (t *timeout) IsEnabled() bool { return t.enabled }

func (t *timeout) Disable() {
	if !t.IsEnabled() {
		return
	}
	t.table.enableTimeout(t.name, false)
}

func (t *timeout) Enable() {
	if t.IsEnabled() || t.currentConfig.timeout <= 0 {
		return
	}
	t.table.enableTimeout(t.name, true)
}

func (t *timeout) Reset() {
	t.SetTimeout(t.defaultConfig.timeout)
}
func (t *timeout) Adjust(any)                {}
func (t *timeout) Configure(Attribute) error { return nil }

func (t *timeout) Attribute(name string) Attribute {
	if strings.Index(name, TimeoutName) != -1 {
		return NewAttribute(TimeoutName, t.currentConfig.timeout)
	}
	return nilAttribute(name)
}

func (t *timeout) Duration() time.Duration {
	if !t.IsEnabled() {
		return 0
	}
	if t.currentConfig.timeout <= 0 {
		return 0
	}
	return t.currentConfig.timeout
}

func (t *timeout) SetTimeout(timeout time.Duration) {
	if t.currentConfig.timeout == timeout || timeout <= 0 {
		return
	}
	t.table.setTimeout(t.name, timeout)
}
