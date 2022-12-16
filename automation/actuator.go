package automation

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type Actuator interface {
	Timeout() TimeoutAction
	RateLimit() RateLimitAction
	Actuate(events string) error
}

type actuator struct {
	timeout *timeoutAction
	limit   *rateLimitAction
}

func newActuator(t *timeoutAction, l *rateLimitAction) *actuator {
	return &actuator{timeout: t, limit: l}
}
