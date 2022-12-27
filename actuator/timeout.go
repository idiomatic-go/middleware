package actuator

import (
	"errors"
	"fmt"
	"time"
)

const (
	TimeoutName = "timeout"
)

type TimeoutController interface {
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
	table  *table
	name   string
	config TimeoutConfig
}

func cloneTimeout(curr *timeout) *timeout {
	t := new(timeout)
	*t = *curr
	return t
}

func newTimeout(name string, table *table, config *TimeoutConfig) *timeout {
	t := new(timeout)
	t.table = table
	t.name = name
	if config != nil {
		t.config = *config
	}
	return t
}

func (t *timeout) validate() error {
	if t.config.timeout <= 0 {
		return errors.New("invalid configuration: TimeoutController duration cannot be <= 0")
	}
	return nil
}

/*
func (t *timeout) Attribute(name string) Attribute {
	if strings.Index(name, TimeoutName) != -1 {
		return NewAttribute(TimeoutName, t.config.timeout)
	}
	return nilAttribute(name)
}
*/

func timeoutState(t *timeout) []string {
	var val int64 = -1
	if t != nil {
		val = int64(t.Duration() / time.Millisecond)
	}
	return []string{fmt.Sprintf(StateAttributeFmt, TimeoutName, val)}
}

func (t *timeout) Duration() time.Duration {
	if t.config.timeout <= 0 {
		return 0
	}
	return t.config.timeout
}

func (t *timeout) SetTimeout(timeout time.Duration) {
	if t.config.timeout == timeout || timeout <= 0 {
		return
	}
	t.table.setTimeout(t.name, timeout)
}
