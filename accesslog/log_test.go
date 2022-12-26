package accesslog

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func Example_WriteEgress_Error() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	a1 := ActuatorState{Name: "egress-route", WriteIngress: true, WriteEgress: true}
	a2 := ActuatorState{Name: "egress-route", WriteIngress: true, WriteEgress: false}

	WriteEgress(start, time.Since(start), ActuatorState{}, nil, nil, "")
	WriteEgress(start, time.Since(start), a1, nil, nil, "")
	WriteEgress(start, time.Since(start), a2, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"error": "egress route name is empty"}]
	//test: WriteEgress() -> [{"error": "egress log entries are empty"}]

}

func Example_WriteIngress_Error() {
	SetTestIngressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	a1 := ActuatorState{Name: "ingress-route", WriteIngress: true, WriteEgress: true}
	a2 := ActuatorState{Name: "ingress-route", WriteIngress: false, WriteEgress: true}

	WriteIngress(start, time.Since(start), ActuatorState{}, nil, nil, "")
	WriteIngress(start, time.Since(start), a1, nil, nil, "")
	WriteIngress(start, time.Since(start), a2, nil, nil, "")

	//Output:
	//test: WriteIngress() -> [{"error": "ingress route name is empty"}]
	//test: WriteIngress() -> [{"error": "ingress log entries are empty"}]
}

func Example_WriteIngress_Ping() {
	name := "ingress-ping-route"
	SetTestIngressWrite()
	SetPingRoutes([]string{name})
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateIngressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
		{Operator: RouteNameOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	WriteIngress(start1, time.Since(start), ActuatorState{Name: name, WriteIngress: true}, nil, nil, "")

	//Output:
	//test: WriteIngress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ping","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234","route_name":"ingress-ping-route"}]

}

func Example_WriteEgress_Origin_Timeout() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
		{Operator: RouteNameOperator}, {Operator: TimeoutDurationOperator}, {Operator: "static", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	WriteEgress(start1, time.Since(start), NewActuatorStateWithTimeout("egress-route", 5000), nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234","route_name":"egress-route","timeout_ms":5000,"static":"value"}]

}

func Example_WriteEgress_Origin_RateLimiter_500() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
		{Operator: RouteNameOperator}, {Operator: RateLimitOperator}, {Operator: RateBurstOperator}, {Operator: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	WriteEgress(start1, time.Since(start), NewActuatorStateWithRateLimiter("egress-route", 500, 10), nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234","route_name":"egress-route","rate_limit_s":500,"rate_burst":10,"static2":"value"}]

}

func Example_WriteEgress_Origin_RateLimiter_Inf() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
		{Operator: RouteNameOperator}, {Operator: RateLimitOperator}, {Operator: RateBurstOperator}, {Operator: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	WriteEgress(start1, time.Since(start), NewActuatorStateWithRateLimiter("egress-route", rate.Inf, 10), nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234","route_name":"egress-route","rate_limit_s":-1,"rate_burst":10,"static2":"value"}]

}

func Example_WriteEgress_Origin_Failover() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
		{Operator: RouteNameOperator}, {Operator: FailoverOperator}, {Operator: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	WriteEgress(start1, time.Since(start), NewActuatorStateWithFailover("egress-route", true), nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234","route_name":"egress-route","failover":true,"static2":"value"}]

}

func Example_WriteEgress_Request() {
	SetTestEgressWrite()
	req, _ := http.NewRequest("", "www.google.com/search/documents", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")

	var start time.Time
	err := CreateEgressEntries([]Reference{{Operator: RequestProtocolOperator}, {Operator: RequestMethodOperator}, {Operator: RequestUrlOperator},
		{Operator: RequestPathOperator}, {Operator: RequestHostOperator}, {Operator: "%REQ(customer)%"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	WriteEgress(start, time.Since(start), ActuatorState{Name: "egress-route", WriteEgress: true}, nil, nil, "")
	WriteEgress(start, time.Since(start), ActuatorState{Name: "egress-route", WriteEgress: true}, req, nil, "")

	//Output:
	//test: WriteEgress() -> [{"protocol":null,"method":null,"url":null,"path":null,"host":null,"customer":null}]
	//test: WriteEgress() -> [{"protocol":"HTTP/1.1","method":"GET","url":"www.google.com/search/documents","path":"www.google.com/search/documents","host":null,"customer":"Ted's Bait & Tackle"}]

}

func Example_WriteEgress_Response() {
	SetTestEgressWrite()
	resp := &http.Response{StatusCode: 404, ContentLength: 1234}

	err := CreateEgressEntries([]Reference{{Operator: ResponseStatusCodeOperator}, {Operator: ResponseBytesReceivedOperator}, {Operator: ResponseFlagsOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start time.Time
	WriteEgress(start, time.Since(start), ActuatorState{Name: "egress-route", WriteEgress: true}, nil, nil, "UT")
	WriteEgress(start, time.Since(start), ActuatorState{Name: "egress-route", WriteEgress: true}, nil, resp, "UT")

	//Output:
	//test: WriteEgress() -> [{"status_code":"0","bytes_received":"0","response_flags":"UT"}]
	//test: WriteEgress() -> [{"status_code":"404","bytes_received":"1234","response_flags":"UT"}]

}
