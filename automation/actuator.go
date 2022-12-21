package automation

const (
	DefaultName = "*"
	NilValue    = -1
)

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
// https://github.com/keikoproj/inverse-exp-backoff

type Controller interface {
	IsEnabled() bool
	Disable()
	Enable()
	Reset()
	Adjust(change any)
	Configure(attr Attribute) error
	Attribute(name string) Attribute
}

type Actuator interface {
	Name() string
	Logger() LoggingController
	Timeout() TimeoutController
	RateLimiter() RateLimiterController
	CircuitBreaker() CircuitBreakerController
	Failover() FailoverController
	Actuate(events string) error
}

type actuator struct {
	name           string
	logger         *logger
	timeout        *timeout
	rateLimiter    *rateLimiter
	failover       *failover
	circuitBreaker *circuitBreaker
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

func newActuator(l *logger, t *timeout, r *rateLimiter, c *circuitBreaker, f *failover) *actuator {
	return &actuator{logger: l, timeout: t, rateLimiter: r, circuitBreaker: c, failover: f}
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

func (a *actuator) CircuitBreaker() CircuitBreakerController {
	return a.rateLimiter
}

func (a *actuator) Failover() FailoverController {
	return a.failover
}

func (a *actuator) Actuate(events string) error {
	return nil
}
