package automation

const (
	DefaultName = "*"
	NilValue    = -1
)

type Controller interface {
	IsEnabled() bool
	Disable()
	Enable()
	Reset()
	Configure(items ...Attribute) error
	Adjust(up bool)
	Attribute(name string) Attribute
}

type Actuator interface {
	Name() string
	Logger() LoggingController
	Timeout() TimeoutController
	RateLimiter() RateLimiterController
	Failover() FailoverController
	Actuate(events string) error
}

type actuator struct {
	name        string
	logger      *logger
	timeout     *timeout
	rateLimiter *rateLimiter
	failover    *failover
}

func cloneActuator[T *timeout | *rateLimiter | *failover](curr *actuator, controller T) *actuator {
	newAct := new(actuator)
	*newAct = *curr
	switch i := any(controller).(type) {
	case *timeout:
		newAct.timeout = i
	case *rateLimiter:
		newAct.rateLimiter = i
	case *failover:
		newAct.failover = i
	default:
	}
	return newAct
}

func newActuator(l *logger, t *timeout, r *rateLimiter, f *failover) *actuator {
	return &actuator{logger: l, timeout: t, rateLimiter: r, failover: f}
}

func (a *actuator) Name() string {
	return a.name
}

func (a *actuator) Logger() LoggingController {
	return a.logger
}

func (a *actuator) Timeout() TimeoutController {
	return a.timeout
}

func (a *actuator) RateLimiter() RateLimiterController {
	return a.rateLimiter
}

func (a *actuator) Failover() FailoverController {
	return a.failover
}

func (a *actuator) Actuate(events string) error {
	return nil
}
