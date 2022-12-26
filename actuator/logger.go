package actuator

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	EgressTraffic  = "egress"
	IngressTraffic = "ingress"
)

type LogAccess func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess, defaultAccess))

func SetDefaultLogger(lc *LoggerConfig) {
	if lc != nil {
		defaultLogger = newLogger(lc)
	}
}

var defaultAccess LogAccess = func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"responseFlags\":\"%v\"}\n", traffic, start, duration, req, resp, respFlags)
}

type LoggingController interface {
	Controller
	LogIngressAccess(act Actuator, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)
	LogEgressAccess(act Actuator, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)
}

type LoggerConfig struct {
	ingressInvoke LogAccess
	egressInvoke  LogAccess
}

func NewLoggerConfig(ingressInvoke LogAccess, egressInvoke LogAccess) *LoggerConfig {
	return &LoggerConfig{ingressInvoke: ingressInvoke, egressInvoke: egressInvoke}
}

type logger struct {
	enabled bool
	mu      sync.RWMutex
	config  LoggerConfig
}

func newLogger(config *LoggerConfig) *logger {
	if config == nil {
		config = NewLoggerConfig(defaultAccess, defaultAccess)
	}
	if config.ingressInvoke == nil {
		config.ingressInvoke = defaultAccess
	}
	if config.egressInvoke == nil {
		config.egressInvoke = defaultAccess
	}
	return &logger{enabled: true, config: *config}
}

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

func (l *logger) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (l *logger) LogIngressAccess(act Actuator, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	if !l.enabled || l.config.ingressInvoke == nil {
		return
	}
	l.config.ingressInvoke(act, IngressTraffic, start, duration, req, resp, respFlags)
}

func (l *logger) LogEgressAccess(act Actuator, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	if !l.enabled || l.config.egressInvoke == nil {
		return
	}
	l.config.egressInvoke(act, EgressTraffic, start, duration, req, resp, respFlags)
}
