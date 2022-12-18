package automation

type LoggingInvoke func(traffic string, act Actuator)

type LoggingController interface {
	//Controller
	//Timeout() int
	//StatusCode(defaultStatusCode int) int
	//Duration() time.Duration
	IsPingTraffic() bool
	WriteIngress() bool
	WriteEgress() bool
	Log(traffic string, act Actuator)
}

type LoggerConfig struct {
	enabled      bool
	isPing       bool
	writeIngress bool
	writeEgress  bool
	invoke       LoggingInvoke
}

func NewLoggerConfig(isPingTraffic, writeIngress, writeEgress bool, invoke LoggingInvoke) *LoggerConfig {
	return &LoggerConfig{enabled: true, isPing: isPingTraffic, writeIngress: writeIngress, writeEgress: writeEgress, invoke: invoke}
}

type logger struct {
	table    *table
	name     string
	defaultC LoggerConfig
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

func (l *logger) IsPingTraffic() bool {
	return l.defaultC.isPing
}

func (l *logger) WriteIngress() bool {
	return l.defaultC.writeIngress
}

func (l *logger) WriteEgress() bool {
	return l.defaultC.writeEgress
}

func (l *logger) Log(traffic string, act Actuator) {

}
