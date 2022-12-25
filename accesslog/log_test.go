package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"time"
)

func Example_WriteEgress_Error() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	a1 := actuator.NewActuatorWithLogger("egress-route", actuator.NewLoggerConfig(true, true, false, nil))
	a2 := actuator.NewActuatorWithLogger("egress-route", actuator.NewLoggerConfig(true, false, false, nil))

	WriteEgress(start, time.Since(start), nil, nil, nil, "")
	WriteEgress(start, time.Since(start), a1, nil, nil, "")
	WriteEgress(start, time.Since(start), a2, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"error": "egress route is nil"}]
	//test: WriteEgress() -> [{"error": "egress log entries are empty"}]

}

func Example_WriteIngress_Error() {
	SetTestIngressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	a1 := actuator.NewActuatorWithLogger("ingress-route", actuator.NewLoggerConfig(true, true, false, nil))
	a2 := actuator.NewActuatorWithLogger("ingress-route", actuator.NewLoggerConfig(false, true, false, nil))
	WriteIngress(start, time.Since(start), nil, nil, 0, 0, "")
	WriteIngress(start, time.Since(start), a1, nil, 0, 0, "")
	WriteIngress(start, time.Since(start), a2, nil, 0, 0, "")

	//Output:
	//test: WriteIngress() -> [{"error": "ingress route is nil"}]
	//test: WriteIngress() -> [{"error": "ingress log entries are empty"}]
}

func Example_WriteEgress_Origin_Route() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
		{Operator: RouteNameOperator}, {Operator: RouteTimeoutOperator}, {Operator: RouteLimitOperator}, {Operator: RouteBurstOperator}, {Operator: "static", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	r1, _ := route.NewRouteWithConfig("egress-route", 1000, 500, 100, true, false)
	var start1 time.Time
	WriteEgress(start1, time.Since(start), r1, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234","route_name":"egress-route","timeout":1000,"limit":"500","burst":100,"static":"value"}]

}

func Example_WriteEgress_Request() {
	SetTestEgressWrite()
	req, _ := http.NewRequest("", "www.google.com/search/documents", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")

	r1, _ := route.NewRouteWithConfig("egress-route", 1000, 500, 100, true, false)

	var start time.Time
	err := CreateEgressEntries([]Reference{{Operator: RequestProtocolOperator}, {Operator: RequestMethodOperator}, {Operator: RequestUrlOperator},
		{Operator: RequestPathOperator}, {Operator: RequestHostOperator}, {Operator: "%REQ(customer)%"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	WriteEgress(start, time.Since(start), r1, nil, nil, "")
	WriteEgress(start, time.Since(start), r1, req, nil, "")

	//Output:
	//test: WriteEgress() -> [{"protocol":null,"method":null,"url":null,"path":null,"host":null,"customer":null}]
	//test: WriteEgress() -> [{"protocol":"HTTP/1.1","method":"GET","url":"www.google.com/search/documents","path":"www.google.com/search/documents","host":null,"customer":"Ted's Bait & Tackle"}]

}

func Example_WriteEgress_Response() {
	SetTestEgressWrite()
	resp := &http.Response{StatusCode: 404, ContentLength: 1234}
	r1, _ := route.NewRouteWithConfig("egress-route", 1000, 500, 100, true, false)

	err := CreateEgressEntries([]Reference{{Operator: ResponseStatusCodeOperator}, {Operator: ResponseBytesReceivedOperator}, {Operator: ResponseFlagsOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start time.Time
	WriteEgress(start, time.Since(start), r1, nil, nil, "UT")
	WriteEgress(start, time.Since(start), r1, nil, resp, "UT")

	//Output:
	//test: WriteEgress() -> [{"status_code":"0","bytes_received":"0","response_flags":"UT"}]
	//test: WriteEgress() -> [{"status_code":"404","bytes_received":"1234","response_flags":"UT"}]

}
