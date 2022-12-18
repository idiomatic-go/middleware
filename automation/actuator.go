package automation

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type Invoke func(a Actuator)

type Controller interface {
	//Name() string
	IsEnabled() bool
	Reset()
	Disable()
	Configure(event string) error
	Adjust(up bool)
	//State() string
	Value(name string) string
}

type Actuator interface {
	Name() string
	Logger() LoggingController
	Timeout() TimeoutController
	RateLimiter() RateLimiterController
	Actuate(events string) error
}

type actuator struct {
	name    string
	logger  *logger
	timeout *timeout
	limiter *rateLimiter
}

func newActuator(l *logger, t *timeout, r *rateLimiter) *actuator {
	return &actuator{logger: l, timeout: t, limiter: r}
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
	return a.limiter
}

func (a *actuator) Actuate(events string) error {
	return nil
}
