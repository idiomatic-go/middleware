package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"time"
)

func _ExampleLogEgressError() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	WriteEgress(start, time.Since(start), nil, nil, nil, "")
	WriteEgress(start, time.Since(start), route.NewRouteWithLogging("egress-route", true), nil, nil, "")
	WriteEgress(start, time.Since(start), route.NewRouteWithLogging("egress-route", false), nil, nil, "")
	//Output:
	//fail
}

func _ExampleLogIngressError() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	WriteIngress(start, time.Since(start), nil, nil, 0, 0, "")
	WriteIngress(start, time.Since(start), route.NewRouteWithLogging("ingress-route", true), nil, 0, 0, "")
	WriteIngress(start, time.Since(start), route.NewRouteWithLogging("ingress-route", false), nil, 0, 0, "")
	//Output:
	//fail
}

func _ExampleLogEgress() {
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")
	CreateEgressEntries([]Reference{{Operator: "%START_TIME%"}, {Operator: "%TRAFFIC%"}, {Operator: "%REGION%"}, {Operator: "%SUB_ZONE%"}, {Operator: "%INSTANCE_ID%"}, {Operator: "%ROUTE_NAME%"}, {Operator: RequestMethodOperator}, {Operator: "%REQ(customer)%"}, {Operator: ResponseCodeOperator}, {Operator: "%DURATION%", Name: "duration_ms_start"}, {Operator: "static", Name: "value"}})
	WriteEgress(start, time.Since(start), route.NewRouteWithLogging("egress-route", true), req, nil, "")
	//Output:
	//fail

}
