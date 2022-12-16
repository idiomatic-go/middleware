package automation

const (
	PingName = "ping"
)

type PingAction struct {
	enabled bool
}

type PingConfig struct {
	enabled bool
}

func NewPingConfig(enabled bool) *PingConfig {
	return &PingConfig{enabled: enabled}
}

type pingAction struct {
	enabled bool
}

func newPingAction(enabled bool) *pingAction {
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
