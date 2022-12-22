package automation

type FailoverInvoke func(name string)

type FailoverController interface {
	Controller
	Invoke()
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
	if config == nil {
		config = NewFailoverConfig(nil)
	}
	t := new(failover)
	t.table = table
	t.name = name
	t.invoke = config.invoke
	if t.invoke != nil {
		t.enabled = true
	}
	return t
}

func (f *failover) IsEnabled() bool { return f.enabled }

func (f *failover) Disable() {
	if f.IsEnabled() {
		f.table.enableFailover(f.name, false)
	}
}

func (f *failover) Enable() {
	if !f.enabled && f.invoke != nil {
		f.table.enableFailover(f.name, true)
	}
}

func (f *failover) Reset()                     {}
func (f *failover) Configure(Attribute) error  { return nil }
func (f *failover) Adjust(any)                 {}
func (f *failover) Attribute(string) Attribute { return nilAttribute("") }

func (f *failover) Invoke() {
	if !f.IsEnabled() {
		return
	}
	f.invoke(f.name)
}

func (f *failover) SetInvoke(fn FailoverInvoke) {
	if fn == nil {
		return
	}
	f.table.setFailoverInvoke(f.name, fn)
}
