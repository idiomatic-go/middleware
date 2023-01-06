package actuator

import (
	"github.com/idiomatic-go/middleware/accessdata"
	"log"
)

var defaultLogger = newLogger(NewLoggerConfig(defaultAccess))

var defaultAccess LogAccess = func(entry *accessdata.Entry) {
	log.Printf(accessdata.WriteJson(operators, entry))
}

type LoggingController interface {
	LogAccess(entry *accessdata.Entry)
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

func (l *logger) LogAccess(entry *accessdata.Entry) {
	if l.config.accessInvoke == nil || entry == nil {
		return
	}
	l.config.accessInvoke(entry)
}

var operators = []accessdata.Operator{
	{Name: "", Value: accessdata.StartTimeOperator},
	{Name: "", Value: accessdata.DurationOperator},
	{Name: "", Value: accessdata.TrafficOperator},
	{Name: "", Value: accessdata.RouteNameOperator},

	{Name: "", Value: accessdata.OriginRegionOperator},
	{Name: "", Value: accessdata.OriginZoneOperator},
	{Name: "", Value: accessdata.OriginSubZoneOperator},
	{Name: "", Value: accessdata.OriginServiceOperator},
	{Name: "", Value: accessdata.OriginInstanceIdOperator},

	{Name: "", Value: accessdata.RequestMethodOperator},
	{Name: "", Value: accessdata.RequestHostOperator},
	{Name: "", Value: accessdata.RequestPathOperator},

	{Name: "", Value: accessdata.ResponseStatusCodeOperator},
	{Name: "", Value: accessdata.StatusFlagsOperator},

	{Name: "", Value: accessdata.TimeoutDurationOperator},
	{Name: "", Value: accessdata.RateLimitOperator},
	{Name: "", Value: accessdata.RateBurstOperator},
	{Name: "", Value: accessdata.RetryOperator},
	{Name: "", Value: accessdata.RetryRateLimitOperator},
	{Name: "", Value: accessdata.RetryRateBurstOperator},
	{Name: "", Value: accessdata.FailoverOperator},
}
