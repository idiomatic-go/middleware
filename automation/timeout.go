package automation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	name     string
	defaultC TimeoutConfig
	current  TimeoutConfig
	table    *table
}

func newTimeout(name string, c *TimeoutConfig, table *table) *timeout {
	if c == nil {
		c = NewTimeoutConfig(NilValue, NilValue)
	}
	t := new(timeout)
	if c.timeout <= 0 {
		c.timeout = NilValue
	}
	t.table = table
	t.name = name
	t.current.timeout = c.timeout
	t.current.statusCode = c.statusCode
	t.defaultC = t.current
	return t
}

//func (a *timeout) Name() string {
//	return a.name
//}

func (a *timeout) IsEnabled() bool {
	return a.current.timeout != NilValue
}

func (a *timeout) Reset() {
	a.table.resetTimeout(a.name)
}

func (a *timeout) Disable() {
	a.table.disableTimeout(a.name)
}

func (a *timeout) Configure(event string) error {
	if event == "" {
		return errors.New("invalid event : event is empty")
	}
	tokens := strings.Split(event, ":")
	if len(tokens) <= 1 {
		return errors.New("invalid event schema : no value found")
	}
	to, err := strconv.Atoi(tokens[1])
	if err != nil {
		return err
	}
	if to <= 0 {
		to = NilValue
	}
	a.table.setTimeout(a.name, to)
	return nil
}

func (a *timeout) Adjust(up bool) {
	a.table.setTimeout(a.name, a.current.timeout+5)
}

func (a *timeout) State() (tags []string) {
	tags = append(tags, fmt.Sprintf("name:%v", a.name))
	tags = append(tags, fmt.Sprintf("timeout:%v", a.current.timeout))
	tags = append(tags, fmt.Sprintf("statusCode:%v", a.current.statusCode))
	return
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
