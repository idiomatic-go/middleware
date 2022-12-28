package actuator

import (
	"log"
	"net/http"
	"time"
)

type LogAccess func(traffic string, start time.Time, duration time.Duration, timeout []string, rateLimiter []string, failover []string, retried any, req *http.Request, resp *http.Response, statusFlags string)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess))

func SetDefaultLogger(lc *LoggerConfig) {
	if lc != nil {
		defaultLogger = newLogger(lc)
	}
}

var defaultAccess LogAccess = func(traffic string, start time.Time, duration time.Duration, timeout []string, rateLimiter []string, failover []string, retried any, req *http.Request, resp *http.Response, statusFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"statusFlags\":\"%v\"}\n", traffic, start, duration, req, resp, statusFlags)
}

type LoggingController interface {
	LogAccess(traffic string, start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string)
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

func (l *logger) LogAccess(traffic string, start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, statusFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(traffic, start, duration, nil, nil, nil, nil, req, resp, statusFlags)
}
