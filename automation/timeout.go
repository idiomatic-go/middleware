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

func newTimeout(name string, config *TimeoutConfig, table *table) *timeout {
	if config == nil {
		config = NewTimeoutConfig(NilValue, NilValue)
	}
	t := new(timeout)
	if config.timeout <= 0 {
		config.timeout = NilValue
	}
	t.table = table
	t.name = name
	t.current.timeout = config.timeout
	t.current.statusCode = config.statusCode
	t.defaultC = t.current
	return t
}

func (t *timeout) IsEnabled() bool {
	return t.current.timeout != NilValue
}

func (t *timeout) Reset() {
	t.table.resetTimeout(t.name)
}

func (t *timeout) Disable() {
	t.table.disableTimeout(t.name)
}

func (t *timeout) Configure(event string) error {
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
	t.table.setTimeout(t.name, to)
	return nil
}

func (t *timeout) Adjust(up bool) {
	if up {
		t.table.setTimeout(t.name, t.current.timeout+5)
	}
}

func (t *timeout) State() string { //(names []string, values []string) {
	//names = append(names, "timeout")
	//names = append(names, "statusCode")
	//values = append(values, strconv.Itoa(a.current.timeout))
	//values = append(values, strconv.Itoa(a.current.statusCode))
	return fmt.Sprintf("timeout: %v , statusCode:%v", t.current.timeout, t.current.statusCode)
}

func (t *timeout) Value(name string) string {
	if name == "" {
		return ""
	}
	if strings.Index(name, TimeoutName) != -1 {
		return strconv.Itoa(t.current.timeout)
	}
	return ""
}

func (t *timeout) Timeout() int {
	return t.current.timeout
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
