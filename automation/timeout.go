package automation

import "time"

const (
	TimeoutName = "timeout"
)

type TimeoutAction interface {
	Duration() time.Duration
}

type timeoutAction struct {
	Default int
	current int
}

func (a *timeoutAction) Name() string {
	return TimeoutName
}

func (a *timeoutAction) IsEnabled() bool {
	return a.current != NilValue
}

func (a *timeoutAction) Reset() {

}

func (a *timeoutAction) Disable() {
}

func (a *timeoutAction) Configure(v ...any) {
}

func (a *timeoutAction) Duration() time.Duration {
	if a.current == NilValue {
		return 0
	}
	return time.Duration(a.current) * time.Millisecond
}
