package actuator

import (
	"errors"
	"net/http"
)

const (
	HostActuatorName      = "host"
	DefaultActuatorName   = "*"
	NilValue              = -1
	RateLimitInfValue     = 99999
	EgressTraffic         = "egress"
	IngressTraffic        = "ingress"
	RateLimitFlag         = "RL"
	UpstreamTimeoutFlag   = "UT"
	HostTimeoutFlag       = "HT"
	NotEnabledFlag        = "NE"
	InvalidResponseFlag   = "IR"
	InvalidStatusCodeFlag = "SC"

	StateAttributeFmt  = "%v:%v"
	TimeoutName        = "timeout"
	FailoverName       = "failover"
	RetryName          = "retry"
	RetryRateLimitName = "retryRateLimit"
	RetryRateBurstName = "retryBurst"
	RateLimitName      = "rateLimit"
	RateBurstName      = "burst"
	StatusCodeName     = "statusCode"
	ActName            = "name"
)

type Matcher func(req *http.Request) (routeName string)

type Actuator interface {
	Name() string
	Logger() LoggingController
	Timeout() (TimeoutController, bool)
	RateLimiter() (RateLimiterController, bool)
	Retry() (RetryController, bool)
	Failover() (FailoverController, bool)
	Actuate(events string) error
	t() *actuator
}

type Configuration interface {
	SetMatcher(fn Matcher)
	SetDefaultActuator(name string, config ...any) error
	SetHostActuator(config ...any) error
	Add(name string, config ...any) error
}

type Actuators interface {
	Host() Actuator
	Lookup(req *http.Request) Actuator
	LookupByName(name string) Actuator
}

type Table interface {
	Configuration
	Actuators
}

var IngressTable = NewIngressTable()
var EgressTable = NewEgressTable()

func NewActuator(name string, config ...any) Actuator {
	return newActuator(name, newTable(true), config...)
}

//func NewActuatorWithLogger(name string, config *LoggerConfig) Actuator {
//	return &actuator{name: name, logger: newLogger(config)}
//}

type actuator struct {
	name        string
	logger      *logger
	timeout     *timeout
	rateLimiter *rateLimiter
	failover    *failover
	retry       *retry
}

func cloneActuator[T *timeout | *rateLimiter | *retry | *failover](curr *actuator, controller T) *actuator {
	newAct := new(actuator)
	*newAct = *curr
	switch i := any(controller).(type) {
	case *timeout:
		newAct.timeout = i
	case *rateLimiter:
		newAct.rateLimiter = i
	case *failover:
		newAct.failover = i
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
		case *FailoverConfig:
			act.failover = newFailover(name, t, c)
		case *RetryConfig:
			act.retry = newRetry(name, t, c)
		}
	}
	return act
}

func newDefaultActuator(name string) *actuator {
	return &actuator{name: name, logger: defaultLogger}
}

func (a *actuator) validate(egress bool) error {
	if egress {

	} else {
		if a.failover != nil {
			return errors.New("invalid configuration: FailoverController is not valid for ingress traffic")
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

func (a *actuator) Timeout() (TimeoutController, bool) {
	if a.timeout == nil {
		return nil, false
	}
	return a.timeout, true
}

func (a *actuator) RateLimiter() (RateLimiterController, bool) {
	if a.rateLimiter == nil {
		return nil, false
	}
	return a.rateLimiter, true
}

func (a *actuator) Retry() (RetryController, bool) {
	if a.retry == nil {
		return nil, false
	}
	return a.retry, true
}

func (a *actuator) Failover() (FailoverController, bool) {
	if a.failover == nil {
		return nil, false
	}
	return a.failover, true
}

func (a *actuator) Actuate(events string) error {
	return nil
}

func (a *actuator) t() *actuator {
	return a
}
