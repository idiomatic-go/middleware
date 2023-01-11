package google

import (
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/host"
	"github.com/idiomatic-go/middleware/template"
	"net/http"
	"time"
)

func init() {
	name := "google:search"
	host.WrapDefaultTransport()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
	actuator.EgressTable.Add(name, nil)
	actuator.EgressTable.SetMatcher(func(req *http.Request) string {
		return name
	},
	)
}

func ExampleSearch_Success() {
	buff, status := Search[template.DebugHandler](nil)
	fmt.Printf("test: Search() -> [%v] [content:%v]\n", status, buff != nil)

	//Output:
	//test: Search() -> [0 Successful] [content:true]

}
