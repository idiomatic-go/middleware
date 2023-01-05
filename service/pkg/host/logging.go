package host

import "github.com/idiomatic-go/middleware/accessdata"

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
	accessdata.SetPingRoutes(nil)
}
