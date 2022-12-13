package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"time"
)

func Example_WriteEgress_Error() {
	egressWrite = func(s string) {
		fmt.Printf("test: WriteEgress() -> [%v]\n", s)
	}
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	r1, _ := route.NewRouteWithLogging("egress-route", true)
	r2, _ := route.NewRouteWithLogging("egress-route", false)

	WriteEgress(start, time.Since(start), nil, nil, nil, "")
	WriteEgress(start, time.Since(start), r1, nil, nil, "")
	WriteEgress(start, time.Since(start), r2, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"error": "egress route is nil"}]
	//test: WriteEgress() -> [{"error": "egress log entries are empty"}]

}

func Example_WriteIngress_Error() {
	ingressWrite = func(s string) {
		fmt.Printf("test: WriteIngress() -> [%v]\n", s)
	}
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	r1, _ := route.NewRouteWithLogging("ingress-route", true)
	r2, _ := route.NewRouteWithLogging("ingress-route", false)
	WriteIngress(start, time.Since(start), nil, nil, 0, 0, "")
	WriteIngress(start, time.Since(start), r1, nil, 0, 0, "")
	WriteIngress(start, time.Since(start), r2, nil, 0, 0, "")

	//Output:
	//test: WriteIngress() -> [{"error": "ingress route is nil"}]
	//test: WriteIngress() -> [{"error": "ingress log entries are empty"}]
}

func Example_WriteEgress() {
	egressWrite = func(s string) {
		fmt.Printf("test: WriteEgress() -> [%v]\n", s)
	}
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")

	r1, _ := route.NewRouteWithLogging("egress-route", true)

	var start1 time.Time
	CreateEgressEntries([]Reference{{Operator: "%START_TIME%"}, {Operator: "%TRAFFIC%"}, {Operator: "%REGION%"}, {Operator: "%SUB_ZONE%"}, {Operator: "%INSTANCE_ID%"}, {Operator: "%ROUTE_NAME%"}, {Operator: RequestMethodOperator}, {Operator: "%REQ(customer)%"}, {Operator: ResponseCodeOperator}, {Operator: "%DURATION%", Name: "duration_ms_start"}, {Operator: "static", Name: "value"}})
	WriteEgress(start1, time.Since(start), r1, req, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","traffic":"ingress","region":"us-west","sub_zone":null,"instance_id":"123456-7890-1234","route_name":"egress-route","method":"GET","customer":"Ted's Bait & Tackle","status_code":"0","duration_ms_start":0,"static":"value"}]

}
