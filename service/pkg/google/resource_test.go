package google

import (
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/egress"
	"github.com/idiomatic-go/middleware/got"
	"net/http"
	"time"
)

func init() {
	egress.EnableDefaultHttpClient()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add("google:search", "https://www.google.com/search", nil)
}

func ExampleSearch_Success() {
	buff, status := Search[got.DebugHandler](nil)
	fmt.Printf("test: Search() -> [%v] [content:%v]\n", status, buff != nil)

	//Output:
	//test: Search() -> [0 Successful] [content:true]

}
