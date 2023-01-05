package google

import "github.com/idiomatic-go/middleware/actuator"

const (
	ActuatorName = "google:search"
	Pattern      = "https://www.google.com/search?q=test"
)

func Startup() {
	actuator.EgressTable.Add(ActuatorName, Pattern, nil)

}
