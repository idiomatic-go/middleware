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
	StatusCode(defaultStatusCode int) int
	Duration() time.Duration
}

type TimeoutConfig struct {
	timeout    int
	statusCode int
}

func NewTimeoutConfig(timeout int, statusCode int) *TimeoutConfig {
	if timeout < 0 {
		timeout = NilValue
	}
	if statusCode <= 0 {
		statusCode = NilValue
	}
	// TODO : validate status code
	return &TimeoutConfig{timeout: timeout, statusCode: statusCode}
}

type timeout struct {
	table    *table
	name     string
	enabled  bool
	defaultC TimeoutConfig
	current  TimeoutConfig
}

func cloneTimeout(curr *timeout) *timeout {
	t := new(timeout)
	*t = *curr
	return t
}

func newTimeout(name string, config *TimeoutConfig, table *table) *timeout {
	if config == nil {
		config = NewTimeoutConfig(NilValue, NilValue)
	}
	t := new(timeout)
	t.table = table
	t.name = name
	t.current.timeout = config.timeout
	t.current.statusCode = config.statusCode
	t.defaultC = t.current
	t.enabled = t.current.timeout > 0
	return t
}

func (t *timeout) IsEnabled() bool {
	return t.enabled && t.current.timeout > 0
}

func (t *timeout) Disable() {
	t.table.enableTimeout(t.name, false)
}

func (t *timeout) Enable() {
	t.table.enableTimeout(t.name, true)
}

func (t *timeout) Reset() {
	t.table.setTimeout(t.name, t.defaultC.timeout)
}

func (t *timeout) Configure(items ...Attribute) error {
	if len(items) == 0 {
		return nil //errors.New("invalid event : event is empty")
	}
	if items[0].Name() == TimeoutName {
		if val, ok := items[0].Value().(int); ok {
			t.table.setTimeout(t.name, val)
		}
	}
	return nil
}

func (t *timeout) Adjust(up bool) {
	if up {
		t.table.setTimeout(t.name, t.current.timeout+5)
	}
}

func (t *timeout) Attribute(name string) Attribute {
	if strings.Index(name, TimeoutName) != -1 {
		return NewAttribute(TimeoutName, t.current.timeout)
	}
	return nilAttribute(name)
}

func (t *timeout) StatusCode(defaultStatusCode int) int {
	if t.current.statusCode == NilValue {
		return defaultStatusCode
	}
	return t.current.statusCode
}

func (t *timeout) Duration() time.Duration {
	if t.current.timeout == NilValue {
		return 0
	}
	return time.Duration(t.current.timeout) * time.Millisecond
}
