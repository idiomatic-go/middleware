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
	Controller
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
		statusCode = NilValue
	}
	// TODO : validate status code
	return &TimeoutConfig{timeout: timeout, statusCode: statusCode}
}

type timeout struct {
	table    *table
	name     string
	defaultC TimeoutConfig
	current  TimeoutConfig
}

func cloneTimeout(act Actuator) *timeout {
	if act == nil {
		return nil
	}
	t := new(timeout)
	s := act.Timeout().(*timeout)
	*t = *s
	return t
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
	if tokens[0] != "timeout" {
		return errors.New("invalid event schema : timeout tag not found")
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
	if up {
		a.table.setTimeout(a.name, a.current.timeout+5)
	}
}

func (a *timeout) State() string { //(names []string, values []string) {
	//names = append(names, "timeout")
	//names = append(names, "statusCode")
	//values = append(values, strconv.Itoa(a.current.timeout))
	//values = append(values, strconv.Itoa(a.current.statusCode))
	return fmt.Sprintf("timeout: %v , statusCode:%v", a.current.timeout, a.current.statusCode)
}

func (a *timeout) Value(name string) string {
	if name == "" {
		return ""
	}
	if name == "timeout" {
		return strconv.Itoa(a.current.timeout)
	}
	return ""
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
