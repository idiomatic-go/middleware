package host

import (
	"encoding/json"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/host"
	"github.com/idiomatic-go/middleware/service/pkg/resource"
	"net/http"
	"time"
)

func initEgress() {
	host.WrapDefaultTransport()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add("google:search", nil)

}

func readRoutes(name string) ([]actuator.Route, error) {
	var routes []actuator.Route

	buf, err := resource.ReadFile(name)
	if err != nil {
		return routes, err
	}
	err1 := json.Unmarshal(buf, &routes)
	return routes, err1
}
