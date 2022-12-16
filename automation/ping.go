package automation

const (
	PingName = "ping"
)

type pingAction struct {
	enabled bool
}

func NewPingAction(enabled bool) *pingAction {
	return &pingAction{enabled}
}

func (a *pingAction) Name() string {
	return PingName
}

func (a *pingAction) IsEnabled() bool {
	return a.enabled
}

func (a *pingAction) Reset() {
}

func (a *pingAction) Disable() {
}

func (a *pingAction) Configure(v ...any) {
}
