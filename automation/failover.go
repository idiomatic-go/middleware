package automation

import "errors"

type FailoverInvoke func(name string, failover bool)

type FailoverController interface {
	Controller
}

type FailoverConfig struct {
	invoke FailoverInvoke
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
		return errors.New("invalid configuration: failover controller FailureInvoke function cannot be nil")
	}
	return nil
}

func (f *failover) IsEnabled() bool { return f.enabled }

func (f *failover) Disable() {
	if !f.IsEnabled() {
		return
	}
	if f.invoke != nil {
		f.invoke(f.name, false)
	}
	f.table.enableFailover(f.name, false)
}

func (f *failover) Enable() {
	if f.IsEnabled() {
		return
	}
	if f.invoke != nil {
		f.invoke(f.name, true)
	}
	f.table.enableFailover(f.name, true)
}

func (f *failover) Reset()                     {}
func (f *failover) Configure(Attribute) error  { return nil }
func (f *failover) Adjust(any)                 {}
func (f *failover) Attribute(string) Attribute { return nilAttribute("") }
