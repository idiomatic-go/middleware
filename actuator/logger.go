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
	{Name: "start_time", Value: accessdata.StartTimeOperator},
	{Name: "duration_ms", Value: accessdata.DurationOperator},
	{Name: "traffic", Value: accessdata.TrafficOperator},
	{Name: "route_name", Value: accessdata.RouteNameOperator},

	{Name: "region", Value: accessdata.OriginRegionOperator},
	{Name: "zone", Value: accessdata.OriginZoneOperator},
	{Name: "sub_zone", Value: accessdata.OriginSubZoneOperator},
	{Name: "service", Value: accessdata.OriginServiceOperator},
	{Name: "instance_id", Value: accessdata.OriginInstanceIdOperator},

	{Name: "method", Value: accessdata.RequestMethodOperator},
	{Name: "host", Value: accessdata.RequestHostOperator},
	{Name: "path", Value: accessdata.RequestPathOperator},

	{Name: "status_code", Value: accessdata.ResponseStatusCodeOperator},
	{Name: "status_flags", Value: accessdata.StatusFlagsOperator},

	{Name: "timeout_ms", Value: accessdata.TimeoutDurationOperator},
	{Name: "rate_limit", Value: accessdata.RateLimitOperator},
	{Name: "rate_burst", Value: accessdata.RateBurstOperator},
	{Name: "retry", Value: accessdata.RetryOperator},
	{Name: "retry_rate_limit", Value: accessdata.RetryRateLimitOperator},
	{Name: "retry_rate_burst", Value: accessdata.RetryRateBurstOperator},
	{Name: "failover", Value: accessdata.FailoverOperator},
}
