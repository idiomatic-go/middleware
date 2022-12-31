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
	retryRoute     = "retry-route"
	googleUrl      = "https://www.google.com/search?q=test"
	twitterUrl     = "https://www.twitter.com"
	facebookUrl    = "https://www.facebook.com"
	instagramUrl   = "https://www.instagram.com"

	config = []accesslog.Reference{
		//{Operator: accesslog.StartTimeOperator},
		//{Operator: accesslog.DurationOperator},
		{Operator: accesslog.TrafficOperator},
		{Operator: accesslog.RouteNameOperator},

		{Operator: accesslog.RequestMethodOperator},
		{Operator: accesslog.RequestHostOperator},
		{Operator: accesslog.RequestPathOperator},
		{Operator: accesslog.RequestProtocolOperator},

		{Operator: accesslog.ResponseStatusCodeOperator},
		{Operator: accesslog.StatusFlagsOperator},
		{Operator: accesslog.ResponseBytesReceivedOperator},
		{Operator: accesslog.ResponseBytesSentOperator},

		{Operator: accesslog.TimeoutDurationOperator},
		{Operator: accesslog.RateLimitOperator},
		{Operator: accesslog.RateBurstOperator},
		{Operator: accesslog.RetryOperator},
		{Operator: accesslog.RetryRateLimitOperator},
		{Operator: accesslog.RetryRateBurstOperator},
		{Operator: accesslog.FailoverOperator},
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
		if req.URL.String() == facebookUrl {
			return retryRoute
		}
		return ""
	})
	actuator.EgressTable.Add(timeoutRoute, nil, actuator.NewTimeoutConfig(time.Millisecond, 504))
	actuator.EgressTable.Add(rateLimitRoute, nil, actuator.NewRateLimiterConfig(2000, 0, 503))
	actuator.EgressTable.Add(retryRoute, nil, actuator.NewTimeoutConfig(time.Millisecond, 504), actuator.NewRetryConfig([]int{503, 504}, 0, 0, 0))

	actuator.SetLoggerAccess(func(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) {
		accesslog.Log(traffic, start, duration, actState, req, resp, statusFlags)
	},
	)
}

func Example_Default_Actuator() {
	act := actuator.EgressTable.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [name:%v]\n", act.Name())

	//Output:
	//test: Lookup(nil) -> [name:*]
}

func Example_No_Wrapper() {
	req, _ := http.NewRequest("GET", googleUrl, nil)

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

func Example_Default() {
	req, _ := http.NewRequest("GET", instagramUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"*","method":"GET","host":"www.instagram.com","path":null,"protocol":"HTTP/1.1","status_code":"200","status_flags":null,"bytes_received":"-1","bytes_sent":"0","timeout_ms":-1,"rate_limit":-1,"rate_burst":-1,"retry":null,"retry_rate_limit":-1,"retry_rate_burst":-1,"failover":null}]
	//test: RoundTrip(egress:true) -> [status_code:200] [err:<nil>]

}

func Example_Default_Timeout() {
	req, _ := http.NewRequest("GET", googleUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"timeout-route","method":"GET","host":"www.google.com","path":"/search","protocol":"HTTP/1.1","status_code":"504","status_flags":"UT","bytes_received":"0","bytes_sent":"0","timeout_ms":1,"rate_limit":-1,"rate_burst":-1,"retry":null,"retry_rate_limit":-1,"retry_rate_burst":-1,"failover":null}]
	//test: RoundTrip(egress:true) -> [status_code:504] [err:<nil>]

}

func Example_Default_RateLimit() {
	req, _ := http.NewRequest("GET", twitterUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"rate-limit-route","method":"GET","host":"www.twitter.com","path":null,"protocol":"HTTP/1.1","status_code":"503","status_flags":"RL","bytes_received":"0","bytes_sent":"0","timeout_ms":-1,"rate_limit":2000,"rate_burst":0,"retry":null,"retry_rate_limit":-1,"retry_rate_burst":-1,"failover":null}]
	//test: RoundTrip(egress:true) -> [status_code:503] [err:<nil>]

}

func Example_Default_Retry_NotEnabled() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	act := actuator.EgressTable.LookupByName(retryRoute)
	if act != nil {
		if c, ok := act.Retry(); ok {
			c.Disable()
		}
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":null,"protocol":"HTTP/1.1","status_code":"504","status_flags":"NE","bytes_received":"0","bytes_sent":"0","timeout_ms":1,"rate_limit":-1,"rate_burst":-1,"retry":false,"retry_rate_limit":0,"retry_rate_burst":0,"failover":null}]
	//test: RoundTrip(egress:true) -> [status_code:504] [err:<nil>]

}

func Example_Default_Retry_RateLimited() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	act := actuator.EgressTable.LookupByName(retryRoute)
	if act != nil {
		if c, ok := act.Retry(); ok {
			c.Enable()
		}
	}
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":null,"protocol":"HTTP/1.1","status_code":"504","status_flags":"RL","bytes_received":"0","bytes_sent":"0","timeout_ms":1,"rate_limit":-1,"rate_burst":-1,"retry":false,"retry_rate_limit":0,"retry_rate_burst":0,"failover":null}]
	//test: RoundTrip(egress:true) -> [status_code:504] [err:<nil>]

}

func Example_Default_Retry() {
	req, _ := http.NewRequest("GET", facebookUrl, nil)

	if !isEnabled {
		isEnabled = true
		EnableDefaultHttpClient()
	}
	act := actuator.EgressTable.LookupByName(retryRoute)
	if act != nil {
		if c, ok := act.Retry(); ok {
			c.Enable()
		}
		if c, ok := act.Retry(); ok {
			c.SetRateLimiter(100, 10)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(egress:true) -> [status_code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":null,"protocol":"HTTP/1.1","status_code":"504","status_flags":"UT","bytes_received":"0","bytes_sent":"0","timeout_ms":1,"rate_limit":-1,"rate_burst":-1,"retry":false,"retry_rate_limit":100,"retry_rate_burst":10,"failover":null}]
	//test: WriteEgress() -> [{"traffic":"egress","route_name":"retry-route","method":"GET","host":"www.facebook.com","path":null,"protocol":"HTTP/1.1","status_code":"504","status_flags":"UT","bytes_received":"0","bytes_sent":"0","timeout_ms":1,"rate_limit":-1,"rate_burst":-1,"retry":true,"retry_rate_limit":100,"retry_rate_burst":10,"failover":null}]
	//test: RoundTrip(egress:true) -> [status_code:504] [err:<nil>]

}
