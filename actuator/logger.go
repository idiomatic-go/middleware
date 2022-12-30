package actuator

import (
	"log"
	"net/http"
	"time"
)

type LogAccess func(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess))

func SetLoggerAccessInvoke(lc *LoggerConfig) {
	if lc != nil && lc.accessInvoke != nil {
		defaultLogger.config.accessInvoke = lc.accessInvoke
	}
}

var defaultAccess LogAccess = func(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"statusFlags\":\"%v\"}\n", traffic, start, duration, req, resp, statusFlags)
}

type LoggingController interface {
	LogAccess(traffic string, start time.Time, duration time.Duration, act Actuator, retry bool, req *http.Request, resp *http.Response, statusFlags string)
}

type LoggerConfig struct {
	accessInvoke LogAccess
}

func NewLoggerConfig(accessInvoke LogAccess) *LoggerConfig {
	return &LoggerConfig{accessInvoke: accessInvoke}
}

type logger struct {
	config LoggerConfig
}

func newLogger(config *LoggerConfig) *logger {
	if config == nil {
		config = NewLoggerConfig(defaultAccess)
	}
	if config.accessInvoke == nil {
		config.accessInvoke = defaultAccess
	}
	return &logger{config: *config}
}

func (l *logger) LogAccess(traffic string, start time.Time, duration time.Duration, act Actuator, retry bool, req *http.Request, resp *http.Response, statusFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	state := make(map[string]string, 12)
	state[ActName] = act.Name()
	timeoutState(state, timeoutController(act))
	rateLimiterState(state, rateLimiterController(act))
	failoverState(state, failoverController(act))
	retryState(state, retryController(act), retry)
	l.config.accessInvoke(traffic, start, duration, state, req, resp, statusFlags)
}

func timeoutController(act Actuator) TimeoutController {
	c, _ := act.Timeout()
	return c
}

func rateLimiterController(act Actuator) RateLimiterController {
	c, _ := act.RateLimiter()
	return c
}

func failoverController(act Actuator) FailoverController {
	c, _ := act.Failover()
	return c
}

func retryController(act Actuator) RetryController {
	c, _ := act.Retry()
	return c
}
