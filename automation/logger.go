package automation

import (
	"net/http"
	"time"
)

type LoggingAccess func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)

type LoggingController interface {
	Controller
	IsPingTraffic() bool
	WriteIngress() bool
	WriteEgress() bool
	LogAccess(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string)
}

type LoggerConfig struct {
	enabled      bool
	isPing       bool
	writeIngress bool
	writeEgress  bool
	accessInvoke LoggingAccess
}

func NewLoggerConfig(isPingTraffic, writeIngress, writeEgress bool, accessInvoke LoggingAccess) *LoggerConfig {
	return &LoggerConfig{enabled: true, isPing: isPingTraffic, writeIngress: writeIngress, writeEgress: writeEgress, accessInvoke: accessInvoke}
}

type logger struct {
	table     *table
	name      string
	isEnabled bool
	defaultC  LoggerConfig
}

func cloneLogging(act Actuator) *logger {
	if act == nil {
		return nil
	}
	t := new(logger)
	s := act.Logger().(*logger)
	*t = *s
	return t
}

func newLogger(name string, config *LoggerConfig, table *table) *logger {
	if config == nil {
		config = NewLoggerConfig(false, false, false, nil)
	}
	return &logger{name: name, defaultC: *config, table: table}
}

func (l *logger) IsEnabled() bool {
	return l.isEnabled
}

func (l *logger) Reset() {

}
func (l *logger) Disable() {

}
func (l *logger) Configure(items ...Attribute) error {
	return nil
}

func (l *logger) Adjust(up bool) {
}

func (l *logger) Attribute(name string) Attribute {
	return nilAttribute(name)
}

func (l *logger) IsPingTraffic() bool {
	return l.defaultC.isPing
}

func (l *logger) WriteIngress() bool {
	return l.defaultC.writeIngress
}

func (l *logger) WriteEgress() bool {
	return l.defaultC.writeEgress
}

func (l *logger) LogAccess(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
	if act == nil || l.defaultC.accessInvoke == nil {
		return
	}
	l.defaultC.accessInvoke(act, traffic, start, duration, req, resp, respFlags)
}
