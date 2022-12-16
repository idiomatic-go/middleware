package automation

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type Actuator interface {
	Name() string
	Timeout() TimeoutAction
	RateLimit() RateLimitAction
	Actuate(events string) error
}

type actuator struct {
	name    string
	timeout *timeout
	limit   *rateLimitAction
}

func newActuator(t *timeout, l *rateLimitAction) *actuator {
	return &actuator{timeout: t, limit: l}
}

func (a *actuator) Name() string {
	return a.name
}

func (a *actuator) Timeout() TimeoutAction {
	return a.timeout
}

func (a *actuator) RateLimit() RateLimitAction {
	return a.limit
}

func (a *actuator) Actuate(events string) error {
	return nil
}
