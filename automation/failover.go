package automation

const ()

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

func newFailover(name string, config *FailoverConfig, table *table) *failover {
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
	f.table.enableFailover(f.name, false)
}

func (f *failover) Enable() {
	f.table.enableFailover(f.name, true)
}

func (f *failover) Reset() {
}

func (f *failover) Configure(Attribute) error {
	return nil
}

func (f *failover) Adjust(bool) {}

func (f *failover) Attribute(string) Attribute {
	return nilAttribute("")
}

func (f *failover) Invoke() {
	if !f.IsEnabled() || f.invoke == nil {
		return
	}
	f.invoke(f.name)
}
