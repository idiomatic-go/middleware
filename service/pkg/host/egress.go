package host

import (
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/handler"
	"net/http"
	"time"
)

func initEgress() {
	handler.EnableDefaultHttpClient()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add("google:search", "https://www.google.com/search", nil)

}
