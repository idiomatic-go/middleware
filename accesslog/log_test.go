package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"time"
)

func ExampleLogEgress() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})
	LogEgress(&route.Route{Name: "egress-route", WriteAccessLog: true}, start, time.Since(start), nil, nil, nil)

	//Output:
	//fail

}
