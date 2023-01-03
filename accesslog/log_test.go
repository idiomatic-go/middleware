package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"time"
)

func ExampleLog_Error() {
	SetTestEgressWrite()
	start := time.Now()
	accessdata.SetOrigin(accessdata.Origin{Region: "us-west", Zone: "dfw", SubZone: "", Service: "test-service", InstanceId: "123456-7890-1234"})

	Log(accessdata.EgressTraffic, start, time.Since(start), map[string]string{}, nil, nil, "")
	Log(accessdata.EgressTraffic, start, time.Since(start), map[string]string{accessdata.ActName: "egress-route"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"error": "egress route name is empty"}]
	//test: WriteEgress() -> [{"error": "egress log entries are empty"}]

}

func ExampleLog_Origin() {
	name := "ingress-origin-route"
	SetTestIngressWrite()
	start := time.Now()
	accessdata.SetOrigin(accessdata.Origin{Region: "us-west", Zone: "dfw", SubZone: "cluster", Service: "test-service", InstanceId: "123456-7890-1234"})
	err := CreateIngressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Value: accessdata.DurationOperator, Name: "duration_ms"},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}, {Value: accessdata.OriginRegionOperator}, {Value: accessdata.OriginZoneOperator}, {Value: accessdata.OriginSubZoneOperator}, {Value: accessdata.OriginServiceOperator}, {Value: accessdata.OriginInstanceIdOperator},
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(accessdata.IngressTraffic, start1, time.Since(start), map[string]string{accessdata.ActName: name}, nil, nil, "")

	//Output:
	//test: WriteIngress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ingress","route_name":"ingress-origin-route","region":"us-west","zone":"dfw","sub_zone":"cluster","service":"test-service","instance_id":"123456-7890-1234"}]

}

func ExampleLog_Ping() {
	name := "ingress-ping-route"
	SetTestIngressWrite()
	accessdata.SetPingRoutes([]string{name})
	start := time.Now()
	err := CreateIngressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Value: accessdata.DurationOperator, Name: "duration_ms"},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(accessdata.IngressTraffic, start1, time.Since(start), map[string]string{accessdata.ActName: name}, nil, nil, "")

	//Output:
	//test: WriteIngress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"ping","route_name":"ingress-ping-route"}]

}

func ExampleLog_Timeout() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressOperators([]accessdata.Operator{{Value: accessdata.StartTimeOperator}, {Name: "duration_ms", Value: accessdata.DurationOperator},
		{Value: accessdata.TrafficOperator}, {Value: accessdata.RouteNameOperator}, {Value: accessdata.TimeoutDurationOperator}, {Name: "static", Value: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(accessdata.EgressTraffic, start1, time.Since(start), map[string]string{accessdata.ActName: "egress-route", accessdata.TimeoutName: "5000"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","timeout_ms":5000,"static":"value"}]

}

/*
func ExampleLog_RateLimiter_500() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Value: StartTimeOperator}, {Value: DurationOperator, Name: "duration_ms"},
		{Value: TrafficOperator}, {Value: RouteNameOperator}, {Value: RateLimitOperator}, {Value: RateBurstOperator}, {Value: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", RateLimitName: "500", RateBurstName: "10"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","rate_limit":500,"rate_burst":10,"static2":"value"}]

}

func ExampleLog_RateLimiter_Inf() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Value: StartTimeOperator}, {Value: DurationOperator, Name: "duration_ms"},
		{Value: TrafficOperator}, {Value: RouteNameOperator}, {Value: RateLimitOperator}, {Value: RateBurstOperator}, {Value: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", RateLimitName: "1000", RateBurstName: "10"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","rate_limit":1000,"rate_burst":10,"static2":"value"}]

}

func ExampleLog_Failover() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Value: StartTimeOperator}, {Value: DurationOperator, Name: "duration_ms"},
		{Value: TrafficOperator}, {Value: RouteNameOperator}, {Value: FailoverOperator}, {Value: "static2", Name: "value"}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", FailoverName: "true"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","failover":true,"static2":"value"}]

}

func ExampleLog_Retry() {
	SetTestEgressWrite()
	start := time.Now()
	err := CreateEgressEntries([]Reference{{Value: StartTimeOperator}, {Value: DurationOperator, Name: "duration_ms"},
		{Value: TrafficOperator}, {Value: RouteNameOperator}, {Value: RetryOperator},
		{Value: RetryRateLimitOperator}, {Value: RetryRateBurstOperator}})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	var start1 time.Time
	Log(EgressTraffic, start1, time.Since(start), map[string]string{ActName: "egress-route", RetryName: "true", RetryRateLimitName: "123", RetryRateBurstName: "67"}, nil, nil, "")

	//Output:
	//test: WriteEgress() -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":"egress","route_name":"egress-route","retry":true,"retry_rate_limit":123,"retry_rate_burst":67}]

}

func ExampleLog_Request() {
	SetTestEgressWrite()
	req, _ := http.NewRequest("", "www.google.com/search/documents", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")

	var start time.Time
	err := CreateEgressEntries([]Reference{{Value: RequestProtocolOperator}, {Value: RequestMethodOperator}, {Value: RequestUrlOperator},
		{Value: RequestPathOperator}, {Value: RequestHostOperator}, {Value: "%REQ(customer)%"}})
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

func ExampleLog_Response() {
	SetTestEgressWrite()
	resp := &http.Response{StatusCode: 404, ContentLength: 1234}

	err := CreateEgressEntries([]Reference{{Value: ResponseStatusCodeOperator}, {Value: ResponseBytesReceivedOperator}, {Value: StatusFlagsOperator}})
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

*/
