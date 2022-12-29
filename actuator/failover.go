package actuator

import (
	"errors"
	"strconv"
)

type FailoverInvoke func(name string, failover bool)

type FailoverController interface {
	IsEnabled() bool
	Enable()
	Disable()
	Invoke(failover bool)
}

type FailoverConfig struct {
	enabled bool
	invoke  FailoverInvoke
}

func NewFailoverConfig(invoke FailoverInvoke) *FailoverConfig {
	return &FailoverConfig{invoke: invoke}
}

type failover struct {
	table   *table
	name    string
	enabled bool
	invoke  FailoverInvoke
}

func cloneFailover(curr *failover) *failover {
	t := new(failover)
	*t = *curr
	return t
}

func newFailover(name string, table *table, config *FailoverConfig) *failover {
	t := new(failover)
	t.table = table
	t.name = name
	if config != nil {
		t.invoke = config.invoke
	}
	t.enabled = false
	return t
}

func (f *failover) validate() error {
	if f.invoke == nil {
		return errors.New("invalid configuration: FailoverController FailureInvoke function cannot be nil")
	}
	return nil
}

func failoverState(m map[string]string, f FailoverController) {
	if f == nil {
		m[FailoverName] = ""
	} else {
		m[FailoverName] = strconv.FormatBool(f.IsEnabled())
	}
}

func (f *failover) IsEnabled() bool { return f.enabled }

func (f *failover) Disable() {
	if !f.IsEnabled() {
		return
	}
	f.table.enableFailover(f.name, false)
}

func (f *failover) Enable() {
	if f.IsEnabled() {
		return
	}
	f.table.enableFailover(f.name, true)
}

func (f *failover) Invoke(failover bool) {
	if f.invoke == nil {
		return
	}
	f.invoke(f.name, failover)
}
