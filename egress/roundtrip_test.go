package egress

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accesslog"
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
	"time"
)

var (
	accessLogging  = false
	isEnabled      = false
	timeoutRoute   = "timeout-route"
	rateLimitRoute = "rate-limit-route"
	googleUrl      = "https://www.google.com/search?q=test"
	twitterUrl     = "https://www.twitter.com"
	facebookUrl    = "https://www.facebook.com"

	config = []accesslog.Reference{
		{Operator: accesslog.StartTimeOperator},
		{Operator: accesslog.TrafficOperator},
		{Operator: accesslog.RouteNameOperator},

		{Operator: accesslog.RequestMethodOperator},
		{Operator: accesslog.RequestHostOperator},
		{Operator: accesslog.RequestPathOperator},
		{Operator: accesslog.RequestProtocolOperator},

		{Operator: accesslog.DurationOperator},
		{Operator: accesslog.ResponseStatusCodeOperator},
		{Operator: accesslog.StatusFlagsOperator},
		{Operator: accesslog.ResponseBytesReceivedOperator},
		{Operator: accesslog.ResponseBytesSentOperator},

		{Operator: accesslog.TimeoutDurationOperator},
		{Operator: accesslog.RateLimitOperator},
		{Operator: accesslog.RateBurstOperator},
	}
)

func init() {
	err := accesslog.CreateEgressEntries(config)
	if err != nil {
		fmt.Printf("init() -> [:%v]\n", err)
	}
	accesslog.SetTestEgressWrite()
	actuator.EgressTable.SetMatcher(func(req *http.Request) string {
		if req == nil {
			return ""
		}
		if req.URL.String() == twitterUrl {
			return rateLimitRoute
		}
		if req.URL.String() == googleUrl {
			return timeoutRoute
		}
		return ""
	})
	actuator.EgressTable.Add(timeoutRoute, actuator.NewTimeoutConfig(time.Millisecond, 504))
	actuator.EgressTable.Add(rateLimitRoute, actuator.NewRateLimiterConfig(2000, 0, 503))
	//r, _ := route.NewRouteWithConfig(timeoutRoute, 1, 500, 100, accessLogging, false)
	//Routes.Add(r)
	//r, _ = route.NewRouteWithConfig(rateLimitRoute, 2000, 0, 0, accessLogging, false)
	//Routes.Add(r)
	actuator.SetAccessInvoke(actuator.NewLoggerConfig(func(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) {
		accesslog.Log(traffic, start, duration, actState, req, resp, statusFlags)
	},
	))
}

func Example_Routes() {
	act := actuator.EgressTable.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [name:%v]\n", act.Name())

	//Output:
	//test: Lookup(nil) -> [name:*]
}

func Example_RoundTrip_No_Wrapper() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	// Testing - check for a nil wrapper or round tripper
	w := wrapper{}
	resp, err := w.RoundTrip(req)
	fmt.Printf("test: RoundTrip(wrapper:nil) -> [resp:%v] [err:%v]\n", resp, err)

	// Testing - no wrapper, calling Google search
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:false) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(wrapper:nil) -> [resp:<nil>] [err:invalid egress round tripper configuration : http.RoundTripper is nil]
	//test: RoundTrip(egress:false) -> [status_code:200] [err:<nil>]

}

func Example_RoundTrip_Default() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(egress:true) -> [status_code:200] [err:<nil>]

}

func Example_RoundTrip_Timeout() {
	req, _ := http.NewRequest("GET", googleUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(egress:true) -> [status_code:504] [err:<nil>]

}

func Example_RoundTrip_RateLimit() {
	req, _ := http.NewRequest("GET", twitterUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: RoundTrip(egress:true) -> [status_code:503] [err:<nil>]

}
