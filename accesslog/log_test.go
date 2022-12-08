package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"time"
)

func _ExampleLogEgressError() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	WriteEgress(nil, start, time.Since(start), nil, nil, "", nil)
	WriteEgress(&route.Route{Name: "egress-route", WriteAccessLog: true}, start, time.Since(start), nil, nil, "", nil)
	WriteEgress(&route.Route{Name: "egress-route", WriteAccessLog: false}, start, time.Since(start), nil, nil, "", nil)
	//Output:
	//fail
}

func _ExampleLogIngressError() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	WriteIngress(nil, start, time.Since(start), nil, 0, 0, "", nil)
	WriteIngress(&route.Route{Name: "ingress-route", WriteAccessLog: true}, start, time.Since(start), nil, 0, 0, "", nil)
	WriteIngress(&route.Route{Name: "ingress-route", WriteAccessLog: false}, start, time.Since(start), nil, 0, 0, "", nil)
	//Output:
	//fail
}

func _ExampleLogEgress() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")
	CreateEgressEntries([]Reference{{Operator: "%START_TIME%"}, {Operator: "%TRAFFIC%"}, {Operator: "%REGION%"}, {Operator: "%SUB_ZONE%"}, {Operator: "%INSTANCE_ID%"}, {Operator: "%ROUTE_NAME%"}, {Operator: HttpMethodOperator}, {Operator: "%REQ(customer)%"}, {Operator: ResponseCodeOperator}, {Operator: "%DURATION%", Name: "duration_ms_start"}, {Operator: "static", Name: "value"}})
	WriteEgress(&route.Route{Name: "egress-route", WriteAccessLog: true}, start, time.Since(start), req, nil, "", nil)
	//Output:
	//fail

}
