package host

import (
	"github.com/idiomatic-go/middleware/accessdata"
	"github.com/idiomatic-go/middleware/accesslog"
)

func initLogging() {
	// accessdata options
	//   SetOrigin() - part of the access log data, and will show on each log entry
	//   SetPingRoutes() - determine which routes/actuator are health liveness check routes
	accessdata.SetOrigin(accessdata.Origin{
		Region:     "region",
		Zone:       "zone",
		SubZone:    "sub-zone",
		Service:    "example-middleware",
		InstanceId: "1234-567-8901",
	})
	// Used to determine which routes are health liveness so that the traffic can be labeled as "ping"
	accessdata.SetPingRoutes(nil)

	// Options that are defaulted to true for the statuses, and log.Printf() for the log functions
	//accesslog.SetIngressLogStatus(true)
	//accesslog.SetEgressLogStatus(true)
	//accesslog.SetIngressLogFn(nil)
	//accesslog.SetEgressLogFn(nil)

	// Enable logging for actuator
	//actuator.SetLoggerAccess(accesslog.Log)

	// Access log attributes for ingress and handler
	accesslog.CreateIngressOperators([]accessdata.Operator{
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
		{Name: "", Value: accessdata.RequestProtocolOperator},
		{Name: "", Value: accessdata.RequestIdOperator},

		{Name: "", Value: accessdata.ResponseStatusCodeOperator},
		{Name: "", Value: accessdata.StatusFlagsOperator},
		{Name: "", Value: accessdata.ResponseBytesSentOperator},

		{Name: "", Value: accessdata.TimeoutDurationOperator},
		{Name: "", Value: accessdata.RateLimitOperator},
		{Name: "", Value: accessdata.RateBurstOperator},
	})

	accesslog.CreateEgressOperators([]accessdata.Operator{
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
		{Name: "", Value: accessdata.RequestProtocolOperator},
		{Name: "", Value: accessdata.RequestIdOperator},

		{Name: "", Value: accessdata.ResponseStatusCodeOperator},
		{Name: "", Value: accessdata.StatusFlagsOperator},
		{Name: "", Value: accessdata.ResponseBytesSentOperator},

		{Name: "", Value: accessdata.TimeoutDurationOperator},
		{Name: "", Value: accessdata.RateLimitOperator},
		{Name: "", Value: accessdata.RateBurstOperator},
		{Name: "", Value: accessdata.RetryOperator},
		{Name: "", Value: accessdata.RetryRateLimitOperator},
		{Name: "", Value: accessdata.RetryRateBurstOperator},
		{Name: "", Value: accessdata.FailoverOperator},
	})
}
