package google

import (
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/egress"
	"github.com/idiomatic-go/middleware/template"
	"net/http"
	"time"
)

func init() {
	egress.EnableDefaultHttpClient()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add("google:search", "https://www.google.com/search", nil)
}

func ExampleSearch_Success() {
	_, status := Search[template.DebugHandler](nil)
	fmt.Printf("test: Search() -> [%v]\n", status)

	//Output:
	//test: Search() -> [0 The operation was successful []]
}
