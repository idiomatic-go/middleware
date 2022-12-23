package automation

import "errors"

const (
	HostActuatorName    = "host"
	DefaultActuatorName = "*"
	NilValue            = -1
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
	Retry() RetryController
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
	retry          *retry
	failoverBackup *failover
}

func cloneActuator[T *timeout | *rateLimiter | *circuitBreaker | *retry | *failover](curr *actuator, controller T) *actuator {
	newAct := new(actuator)
	*newAct = *curr
	switch i := any(controller).(type) {
	case *timeout:
		newAct.timeout = i
	case *rateLimiter:
		newAct.rateLimiter = i
	case *failover:
		newAct.failover = i
	case *circuitBreaker:
		newAct.circuitBreaker = i
	case *retry:
		newAct.retry = i
	default:
	}
	return newAct
}

func newActuator(name string, t *table, config ...any) *actuator {
	act := new(actuator)
	act.name = name
	act.logger = defaultLogger
	for _, cfg := range config {
		switch c := cfg.(type) {
		case *TimeoutConfig:
			act.timeout = newTimeout(name, t, c)
		case *RateLimiterConfig:
			act.rateLimiter = newRateLimiter(name, t, c)
		case *CircuitBreakerConfig:
			act.circuitBreaker = newCircuitBreaker(name, t, c)
		case *FailoverConfig:
			act.failover = newFailover(name, t, c)
		case *RetryConfig:
			act.retry = newRetry(name, t, c)
		}
	}
	return act
}

func newDefaultActuator(name string, t *table) *actuator {
	return newActuator(name, t, newTimeout(name, t, nil),
		newRateLimiter(name, t, nil),
		newCircuitBreaker(name, t, nil),
		newRetry(name, t, nil),
		newFailover(name, t, nil))
}

func (a *actuator) validate(egress bool) error {
	if egress {
		if a.rateLimiter != nil {
			return errors.New("invalid configuration: rate limiter controller is not valid for egress traffic")
		}
	} else {
		if a.circuitBreaker != nil || a.failover != nil || a.retry != nil {
			return errors.New("invalid configuration: circuit breaker, failover, and retry controllers are not valid for ingress traffic")
		}
	}
	return nil
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

func (a *actuator) Retry() RetryController {
	return a.retry
}

func (a *actuator) Failover() FailoverController {
	return a.failover
}

func (a *actuator) Actuate(events string) error {
	return nil
}
