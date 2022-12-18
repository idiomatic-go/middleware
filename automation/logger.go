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
	defaultLogger = newLogger(lc)
}

var defaultAccess LoggingAccess = func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	log.Printf("traffic: %v start_time: %v duration_ms: %v request: %v response: %v responseFlags: %v\n", traffic, start, duration, req, resp, respFlags)
}

type LoggingController interface {
	Controller
	IsPingTraffic(name string) bool
	WriteIngress() bool
	WriteEgress() bool
	LogAccess(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)
}

type LoggerConfig struct {
	enabled      bool
	writeIngress bool
	writeEgress  bool
	accessInvoke LoggingAccess
	exclude      []string
}

func NewLoggerConfig(writeIngress, writeEgress bool, accessInvoke LoggingAccess, exclude []string) *LoggerConfig {
	return &LoggerConfig{enabled: true, writeIngress: writeIngress, writeEgress: writeEgress, accessInvoke: accessInvoke, exclude: exclude}
}

type logger struct {
	isEnabled bool
	mu        sync.RWMutex
	defaultC  LoggerConfig
}

func newLogger(config *LoggerConfig) *logger {
	if config == nil {
		config = NewLoggerConfig(true, true, defaultAccess, nil)
	}
	if config.accessInvoke == nil {
		config.accessInvoke = defaultAccess
	}
	config.enabled = true
	return &logger{defaultC: *config}
}

func (l *logger) IsEnabled() bool {
	return l.isEnabled
}

func (l *logger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isEnabled = true
}

func (l *logger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isEnabled = false
}

func (l *logger) Enable() {

}

func (l *logger) Configure(items ...Attribute) error {
	return nil
}

func (l *logger) Adjust(_ bool) {
}

func (l *logger) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (l *logger) IsPingTraffic(name string) bool {
	for _, n := range l.defaultC.exclude {
		if n == name {
			return true
		}
	}
	return false
}

func (l *logger) WriteIngress() bool {
	return l.defaultC.writeIngress
}

func (l *logger) WriteEgress() bool {
	return l.defaultC.writeEgress
}

func (l *logger) LogAccess(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	if !l.isEnabled || act == nil || l.defaultC.accessInvoke == nil {
		return
	}
	l.defaultC.accessInvoke(act, traffic, start, duration, req, resp, respFlags)
}
