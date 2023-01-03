package actuator

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"net/http"
	"time"
)

const (
	HostActuatorName    = "host"
	DefaultActuatorName = "*"
	RateLimitInfValue   = 99999
	EgressTraffic       = "egress"
	IngressTraffic      = "ingress"
	RateLimitFlag       = "RL"
	UpstreamTimeoutFlag = "UT"
	HostTimeoutFlag     = "HT"
	NotEnabledFlag      = "NE"

	TimeoutName        = "timeout"
	FailoverName       = "failover"
	RetryName          = "retry"
	RetryRateLimitName = "retryRateLimit"
	RetryRateBurstName = "retryBurst"
	RateLimitName      = "rateLimit"
	RateBurstName      = "burst"
	ActName            = "name"
)

type Actuate func(act Actuator, events []Event) error

type Actuator interface {
	Name() string
	Logger() LoggingController
	Extract() ExtractController
	Timeout() (TimeoutController, bool)
	RateLimiter() (RateLimiterController, bool)
	Retry() (RetryController, bool)
	Failover() (FailoverController, bool)
	LogIngress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string)
	LogEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, retry bool)
	Actuate(events []Event) error
	t() *actuator
}

type Matcher func(req *http.Request) (routeName string)

type Configuration interface {
	SetMatcher(fn Matcher)
	SetDefaultActuator(name string, fn Actuate, config ...any) []error
	SetHostActuator(fn Actuate, config ...any) []error
	Add(name string, fn Actuate, config ...any) []error
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

type actuator struct {
	name        string
	logger      *logger
	timeout     *timeout
	rateLimiter *rateLimiter
	failover    *failover
	retry       *retry
	extract     *extract
	actuate     Actuate
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

func newActuator(name string, t *table, fn Actuate, config ...any) (*actuator, []error) {
	var errs []error
	var err error
	act := new(actuator)
	act.name = name
	act.actuate = fn
	act.logger = defaultLogger
	act.extract = newExtract()
	for _, cfg := range config {
		err = nil
		switch c := cfg.(type) {
		case *TimeoutConfig:
			act.timeout = newTimeout(name, t, c)
			err = act.timeout.validate()
		case *RateLimiterConfig:
			act.rateLimiter = newRateLimiter(name, t, c)
			err = act.rateLimiter.validate()
		case *FailoverConfig:
			act.failover = newFailover(name, t, c)
			err = act.failover.validate()
		case *RetryConfig:
			act.retry = newRetry(name, t, c)
			err = act.retry.validate()
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return act, errs
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
		if a.retry != nil {
			return errors.New("invalid configuration: RetryController is not valid for ingress traffic")
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

func (a *actuator) Extract() ExtractController {
	return a.extract
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

func (a *actuator) Actuate(events []Event) error {
	if a.actuate == nil {
		return errors.New(fmt.Sprintf("invalid configuration: Actuate function is nil for : %v", a.name))
	}
	return a.actuate(a, events)
}

func (a *actuator) t() *actuator {
	return a
}

func (a *actuator) state() map[string]string {
	state := make(map[string]string, 12)
	state[ActName] = a.Name()
	timeoutState(state, a.timeout)
	rateLimiterState(state, a.rateLimiter)
	return state
}

func (a *actuator) LogIngress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string) {
	entry := accessdata.NewIngressEntry(start, duration, a.state(), req, resp, statusFlags)
	a.Extract().Extract(entry)
	a.Logger().LogAccess(entry)
}

func (a *actuator) LogEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, retry bool) {
	state := a.state()
	failoverState(state, a.failover)
	retryState(state, a.retry, retry)
	entry := accessdata.NewEgressEntry(start, duration, state, req, resp, statusFlags)
	a.Extract().Extract(entry)
	a.Logger().LogAccess(entry)
}
