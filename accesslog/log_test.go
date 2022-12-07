package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"time"
)

func _ExampleLogEgressError() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	WriteEgress(nil, start, time.Since(start), nil, nil, nil)
	WriteEgress(&route.Route{Name: "egress-route", WriteAccessLog: true}, start, time.Since(start), nil, nil, nil)
	WriteEgress(&route.Route{Name: "egress-route", WriteAccessLog: false}, start, time.Since(start), nil, nil, nil)
	//Output:
	//fail
}

func _ExampleLogEgress() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	AddEgressAttributes([]Entry{{Operator: "%START_TIME%"}, {Operator: "%DURATION%", Name: "duration_ms"}})
	WriteEgress(&route.Route{Name: "egress-route", WriteAccessLog: true}, start, time.Since(start), nil, nil, nil)
	//Output:
	//fail

}
