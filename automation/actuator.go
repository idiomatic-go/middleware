package automation

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type Invoke func(a Actuator)

type Controller interface {
	Name() string
	IsEnabled() bool
	Reset()
	Disable()
	Configure(event string)
	Adjust(up bool)
	State() []string
}

type Actuator interface {
	Name() string
	Timeout() TimeoutController
	//RateLimit() RateLimitAction
	Actuate(events string) error
}

type actuator struct {
	name    string
	timeout *timeout
	//limit   *rateLimitAction
}

func newActuator(t *timeout) *actuator {
	return &actuator{timeout: t}
}

func (a *actuator) Name() string {
	return a.name
}

func (a *actuator) Timeout() TimeoutController {
	return a.timeout
}

//func (a *actuator) RateLimit() RateLimitAction {
//	return a.limit
//}

func (a *actuator) Actuate(events string) error {
	return nil
}
