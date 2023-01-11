package host

import (
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/host"
	"net/http"
	"time"
)

func initEgress() {
	host.WrapDefaultTransport()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add("google:search", nil)

}
