package actuator

import (
	"log"
	"net/http"
	"time"
)

type LogAccess func(traffic string, start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, respFlags string)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess))

func SetDefaultLogger(lc *LoggerConfig) {
	if lc != nil {
		defaultLogger = newLogger(lc)
	}
}

var defaultAccess LogAccess = func(traffic string, start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, respFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"responseFlags\":\"%v\"}\n", traffic, start, duration, req, resp, respFlags)
}

type LoggingController interface {
	LogIngressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, respFlags string)
	LogEgressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, respFlags string)
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

/*
func (l *logger) IsEnabled() bool {
	return l.enabled
}

func (l *logger) Reset() {
}

func (l *logger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = false
}

func (l *logger) Enable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = true
}

func (l *logger) Configure(Attribute) error {
	return nil
}

func (l *logger) Adjust(any) {
}


*/
func (l *logger) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (l *logger) LogIngressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, respFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(IngressTraffic, start, duration, act, req, resp, respFlags)
}

func (l *logger) LogEgressAccess(start time.Time, duration time.Duration, act Actuator, req *http.Request, resp *http.Response, respFlags string) {
	if l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(EgressTraffic, start, duration, act, req, resp, respFlags)
}
