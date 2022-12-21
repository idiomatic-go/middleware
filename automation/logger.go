package automation

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type LoggingAccess func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)

var defaultLogger = newLogger(NewLoggerConfig(true, true, defaultAccess, nil))

func SetDefaultLogger(lc *LoggerConfig) {
	if lc != nil {
		defaultLogger = newLogger(lc)
	}
}

var defaultAccess LoggingAccess = func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"responseFlags\":\"%v\"}\n", traffic, start, duration, req, resp, respFlags)
}

type LoggingController interface {
	Controller
	IsPingTraffic(name string) bool
	WriteIngress() bool
	WriteEgress() bool
	LogAccess(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)
}

type LoggerConfig struct {
	writeIngress bool
	writeEgress  bool
	accessInvoke LoggingAccess
	exclude      []string
}

func NewLoggerConfig(writeIngress, writeEgress bool, accessInvoke LoggingAccess, exclude []string) *LoggerConfig {
	return &LoggerConfig{writeIngress: writeIngress, writeEgress: writeEgress, accessInvoke: accessInvoke, exclude: exclude}
}

type logger struct {
	enabled bool
	mu      sync.RWMutex
	config  LoggerConfig
}

func newLogger(config *LoggerConfig) *logger {
	if config == nil {
		config = NewLoggerConfig(true, true, defaultAccess, nil)
	}
	if config.accessInvoke == nil {
		config.accessInvoke = defaultAccess
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

func (l *logger) Adjust(bool) {
}

func (l *logger) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (l *logger) IsPingTraffic(name string) bool {
	for _, n := range l.config.exclude {
		if n == name {
			return true
		}
	}
	return false
}

func (l *logger) WriteIngress() bool {
	return l.config.writeIngress
}

func (l *logger) WriteEgress() bool {
	return l.config.writeEgress
}

func (l *logger) LogAccess(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	if !l.enabled || l.config.accessInvoke == nil {
		return
	}
	l.config.accessInvoke(act, traffic, start, duration, req, resp, respFlags)
}
