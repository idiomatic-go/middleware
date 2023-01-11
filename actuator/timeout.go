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
	Duration   time.Duration
	StatusCode int
}

func NewTimeoutConfig(duration time.Duration, statusCode int) *TimeoutConfig {
	if statusCode <= 0 {
		statusCode = http.StatusGatewayTimeout
	}
	return &TimeoutConfig{Duration: duration, StatusCode: statusCode}
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
	if t.config.Duration <= 0 {
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
	if t.config.Duration <= 0 {
		return 0
	}
	return t.config.Duration
}

func (t *timeout) SetTimeout(duration time.Duration) {
	if t.config.Duration == duration || duration <= 0 {
		return
	}
	t.table.setTimeout(t.name, duration)
}

func (t *timeout) StatusCode() int {
	return t.config.StatusCode
}
