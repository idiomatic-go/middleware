package actuator

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type TimeoutController interface {
	Duration() time.Duration
	SetTimeout(timeout time.Duration)
	StatusCode() int
}

type TimeoutConfig struct {
	StatusCode int
	Timeout    time.Duration
}

func NewTimeoutConfig(timeout time.Duration, statusCode int) *TimeoutConfig {
	if statusCode <= 0 {
		statusCode = http.StatusGatewayTimeout
	}
	return &TimeoutConfig{Timeout: timeout, StatusCode: statusCode}
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
	if t.config.Timeout <= 0 {
		return errors.New("invalid configuration: TimeoutController duration is <= 0")
	}
	return nil
}

func timeoutState(m map[string]string, t *timeout) {
	var val int64 = -1
	//var statusCode = -1
	if t != nil {
		val = int64(t.Duration() / time.Millisecond)
		//	statusCode = t.StatusCode()
	}
	m[TimeoutName] = strconv.Itoa(int(val))
}

func (t *timeout) Duration() time.Duration {
	if t.config.Timeout <= 0 {
		return 0
	}
	return t.config.Timeout
}

func (t *timeout) SetTimeout(timeout time.Duration) {
	if t.config.Timeout == timeout || timeout <= 0 {
		return
	}
	t.table.setTimeout(t.name, timeout)
}

func (t *timeout) StatusCode() int {
	return t.config.StatusCode
}
