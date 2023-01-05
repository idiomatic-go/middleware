package host

import (
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/egress"
	"github.com/idiomatic-go/middleware/service/pkg/google"
	"net/http"
	"time"
)

func Startup(r *http.ServeMux) (http.Handler, bool) {
	initLogging()
	initIngress()
	initEgress()
	initRoutes(r)

	return r, true
}

func Shutdown() {}

func initEgress() {
	egress.EnableDefaultHttpClient()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	google.Startup()
}

func initExtract() {
	//err := extract.Initialize(uri string, newClient *http.Client, fn ErrorHandler)

}
