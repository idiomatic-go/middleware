package google

import (
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/egress"
	"net/http"
	"time"
)

func init() {
	egress.EnableDefaultHttpClient()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add("google:search", "https://www.google.com/search", nil)
}

func ExampleSearch() {
	bytes := Search()
	fmt.Printf("test: Search() -> [%v]\n", len(bytes))

	//Output:
	//test: Search() -> [165067]
}
