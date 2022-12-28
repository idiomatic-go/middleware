package actuator

import (
	"log"
	"net/http"
	"time"
)

type LogAccess func(traffic string, start time.Time, duration time.Duration, routeName string, timeout []string, rateLimiter []string, failover []string, retry []string, req *http.Request, resp *http.Response, statusFlags string)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess))

func SetDefaultLogger(lc *LoggerConfig) {
	if lc != nil {
		defaultLogger = newLogger(lc)
	}
}

var defaultAccess LogAccess = func(traffic string, start time.Time, duration time.Duration, routeName string, timeout []string, rateLimiter []string, failover []string, retry []string, req *http.Request, resp *http.Response, statusFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"statusFlags\":\"%v\"}\n", traffic, start, duration, req, resp, statusFlags)
}

type LoggingController interface {
	LogAccess(traffic string, start time.Time, duration time.Duration, act Actuator, retry string, req *http.Request, resp *http.Response, statusFlags string)
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

func (l *logger) LogAccess(traffic string, start time.Time, duration time.Duration, act Actuator, retry string, req *http.Request, resp *http.Response, statusFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(traffic, start, duration, act.Name(),
		timeoutAttributes(timeoutController(act)),
		rateLimiterAttributes(rateLimiterController(act)),
		failoverAttributes(failoverController(act)),
		retryAttributes(retryController(act), retry),
		req, resp, statusFlags)
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
