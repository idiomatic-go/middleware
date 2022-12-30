package accesslog

import (
	"fmt"
	"net/http"
	"reflect"
	"time"
)

func Example_Log_Error() {
	SetTestEgressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	Log(EgressTraffic, start, time.Since(start), map[string]string{}, nil, nil, "")
	Log(EgressTraffic, start, time.Since(start), map[string]string{ActName: "egress-route"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"error": "egress route name is empty"}]
	//test: WriteEgress() -> [{"error": "egress log entries are empty"}]

}

func Example_Log_Origin() {
	name := "ingress-origin-route"
	SetTestIngressWrite()
	start := time.Now()
	SetOrigin(Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateIngressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}, {Operator: OriginRegionOperator}, {Operator: OriginZoneOperator}, {Operator: OriginSubZoneOperator}, {Operator: OriginServiceOperator}, {Operator: OriginInstanceIdOperator},
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(IngressTraffic, start1, time.Since(start), map[string]string{ActName: name}, nil, nil, "")

	//Output:
	//test: WriteIngress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ingress","route_name":"ingress-origin-route","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234"}]

}

func Example_Log_Ping() {
	name := "ingress-ping-route"
	SetTestIngressWrite()
	SetPingRoutes([]string{name})
	start := time.Now()
	err := CreateIngressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(IngressTraffic, start1, time.Since(start), map[string]string{ActName: name}, nil, nil, "")

	//Output:
	//test: WriteIngress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ping","route_name":"ingress-ping-route"}]

}

func Example_Log_Timeout() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}, {Operator: TimeoutDurationOperator}, {Operator: "static", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", TimeoutName: "5000"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","timeout_ms":5000,"static":"value"}]

}

func Example_Log_RateLimiter_500() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}, {Operator: RateLimitOperator}, {Operator: RateBurstOperator}, {Operator: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", RateLimitName: "500", RateBurstName: "10"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","rate_limit_s":500,"rate_burst":10,"static2":"value"}]

}

func Example_Log_RateLimiter_Inf() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}, {Operator: RateLimitOperator}, {Operator: RateBurstOperator}, {Operator: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", RateLimitName: "1000", RateBurstName: "10"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","rate_limit_s":1000,"rate_burst":10,"static2":"value"}]

}

func Example_Log_Failover() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}, {Operator: FailoverOperator}, {Operator: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", FailoverName: "true"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","failover":true,"static2":"value"}]

}

func Example_Log_Retry() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Operator: StartTimeOperator}, {Operator: DurationOperator, Name: "duration_ms"},
		{Operator: TrafficOperator}, {Operator: RouteNameOperator}, {Operator: RetryOperator},
		{Operator: RetryRateLimitOperator}, {Operator: RetryRateBurstOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", RetryName: "true", RetryRateLimitName: "123", RetryRateBurstName: "67"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","retry":true,"retry_rate_limit_s":123,"retry_rate_burst":67}]

}

func Example_Log_Request() {
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
	Log(EgressTraffic, start, time.Since(start), map[string]string{ActName: "egress-route"}, nil, nil, "")
	Log(EgressTraffic, start, time.Since(start), map[string]string{ActName: "egress-route"}, req, nil, "")

	//Output:
	//test: WriteEgress() -> [{"protocol":null,"method":null,"url":null,"path":null,"host":null,"customer":null}]
	//test: WriteEgress() -> [{"protocol":"HTTP/1.1","method":"GET","url":"www.google.com/search/documents","path":"www.google.com/search/documents","host":null,"customer":"Ted's Bait & Tackle"}]

}

func Example_Log_Response() {
	SetTestEgressWrite()
	resp := &http.Response{StatusCode: 404, ContentLength: 1234}

	err := CreateEgressEntries([]Reference{{Operator: ResponseStatusCodeOperator}, {Operator: ResponseBytesReceivedOperator}, {Operator: StatusFlagsOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start time.Time
	Log(EgressTraffic, start, time.Since(start), map[string]string{ActName: "egress-route"}, nil, nil, "UT")
	Log(EgressTraffic, start, time.Since(start), map[string]string{ActName: "egress-route"}, nil, resp, "UT")

	//Output:
	//test: WriteEgress() -> [{"status_code":"0","bytes_received":"0","status_flags":"UT"}]
	//test: WriteEgress() -> [{"status_code":"404","bytes_received":"1234","status_flags":"UT"}]

}

func __Example_Log_State() {
	t := time.Duration(time.Millisecond * 500)
	i := reflect.TypeOf(t)
	a := any(t)

	fmt.Printf("test 1 -> %v\n", a)

	fmt.Printf("test 2 -> %v\n", i)

	//Output:
	//fail
}
