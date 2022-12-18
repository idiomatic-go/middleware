package automation

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type attribute struct {
	name  string
	value any
}

type Controller interface {
	IsEnabled() bool
	Disable()
	Reset()
	Configure(items ...attribute) error
	Adjust(up bool)
	//State() string
	//Value(name string) string
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
	name     string
	logger   *logger
	timeout  *timeout
	limiter  *rateLimiter
	failover *failover
}

func newActuator(l *logger, t *timeout, r *rateLimiter, f *failover) *actuator {
	return &actuator{logger: l, timeout: t, limiter: r, failover: f}
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

func (a *actuator) Failover() FailoverController {
	return a.failover
}

func (a *actuator) Actuate(events string) error {
	return nil
}
