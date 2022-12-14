package egress

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accesslog"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
)

var (
	name = "egress-route"

	config = []accesslog.Reference{
		{Operator: accesslog.StartTimeOperator},
		{Operator: accesslog.TrafficOperator},
		{Operator: accesslog.RouteNameOperator},

		//{Operator: accesslog.OriginRegionOperator},
		//{Operator: accesslog.OriginZoneOperator},
		//{Operator: accesslog.OriginSubZoneOperator},
		//{Operator: accesslog.OriginServiceOperator},
		//{Operator: accesslog.OriginInstanceIdOperator},

		{Operator: accesslog.RequestMethodOperator},
		{Operator: accesslog.RequestHostOperator},
		{Operator: accesslog.RequestPathOperator},
		{Operator: accesslog.RequestProtocolOperator},
		//{Operator: accesslog.RequestIdOperator},
		//{Operator: accesslog.RequestForwardedForOperator},

		{Operator: accesslog.DurationOperator},
		{Operator: accesslog.ResponseStatusCodeOperator},
		{Operator: accesslog.ResponseFlagsOperator},
		{Operator: accesslog.ResponseBytesReceivedOperator},
		{Operator: accesslog.ResponseBytesSentOperator},

		{Operator: accesslog.RouteTimeoutOperator},
		{Operator: accesslog.RouteLimitOperator},
		{Operator: accesslog.RouteBurstOperator},
	}
)

func init() {
	err := accesslog.CreateEgressEntries(config)
	if err != nil {
		fmt.Printf("init() -> [:%v]\n", err)
	}
	accesslog.SetTestEgressWrite()
	Routes.SetMatcher(func(req *http.Request) string {
		return name
	})

}

func Example_Routes() {
	r := Routes.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [name:%v]\n", r.Name())

	//Output:
	//test: Lookup(nil) -> [name:*]
}

func Example_RoundTrip() {
	r1, _ := route.NewRouteWithConfig(name, 1000, 500, 100, true, false)
	Routes.Add(r1)
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)

	// Testing - check for a nil wrapper or round tripper
	w := wrapper{}
	resp, err := w.RoundTrip(req)
	fmt.Printf("test: RoundTrip(wrapper:nil) -> [resp:%v] [err:%v]\n", resp, err)

	// Testing - no wrapper, calling Google search
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(req) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	// Testing - enable egress, calling Google search
	EnableDefaultHttpClient()
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(req) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(wrapper:nil) -> [resp:<nil>] [err:invalid wrapper : http.RoundTripper is nil]
	//test: RoundTrip(req) -> [status_code:200] [err:<nil>]

}
