package actuator

import (
	"log"
	"net/http"
	"time"
)

type LogAccess func(traffic string, start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess))

func SetDefaultLogger(lc *LoggerConfig) {
	if lc != nil {
		defaultLogger = newLogger(lc)
	}
}

var defaultAccess LogAccess = func(traffic string, start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"statusFlags\":\"%v\"}\n", traffic, start, duration, req, resp, statusFlags)
}

type LoggingController interface {
	LogIngressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string)
	LogEgressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string)
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

func (l *logger) LogIngressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(IngressTraffic, start, duration, act, req, resp, statusFlags)
}

func (l *logger) LogEgressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(EgressTraffic, start, duration, act, req, resp, statusFlags)
}
