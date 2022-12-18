package automation

import "fmt"

const (
	FailoverName = "failover"
)

type FailoverInvoke func(name string)

type FailoverController interface {
	Controller
	Value(name string) string
	Failover()
}

type FailoverConfig struct {
	invoke FailoverInvoke
}

func NewFailoverConfig(invoke FailoverInvoke) *FailoverConfig {
	return &FailoverConfig{invoke: invoke}
}

type failover struct {
	table     *table
	name      string
	isEnabled bool
	invoke    FailoverInvoke
}

func cloneFailover(act Actuator) *failover {
	if act == nil {
		return nil
	}
	t := new(failover)
	s := act.Failover().(*failover)
	*t = *s
	return t
}

func newFailover(name string, config *FailoverConfig, table *table) *failover {
	if config == nil {
		config = NewFailoverConfig(nil)
	}
	t := new(failover)
	t.table = table
	t.name = name
	t.invoke = config.invoke
	return t
}

func (f *failover) IsEnabled() bool { return f.isEnabled }

func (f *failover) Reset() {
	f.Disable()
}

func (f *failover) Disable() {
	// TODO : set f.isEnabled = false
}

func (f *failover) Configure(items ...attribute) error {
	return nil
}

func (f *failover) Adjust(up bool) {}

func (f *failover) Value(name string) string {
	return fmt.Sprintf("%v", f.isEnabled)
}

func (f *failover) Failover() {
	if f.invoke == nil {
		return
	}
	f.invoke(f.name)
}
