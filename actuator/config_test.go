package actuator

import (
	"encoding/json"
	"fmt"
)

func ExampleConfig_Marshal() {
	config := Route{Name: "test-route", Host: "google.com",
		Timeout: &TimeoutConfig{
			StatusCode: 504,
			Duration:   20000,
		},
		RateLimiter: &RateLimiterConfig{
			Limit:      100,
			Burst:      25,
			StatusCode: 503,
		},
		Retry: &RetryConfig{
			Limit: 100,
			Burst: 33,
			Wait:  500,
			Codes: []int{503, 504},
		},
		//Failover: &FailoverConfig{
		//	Enabled: false,
		//	invoke:  nil,
		//},
	}
	buf, err := json.Marshal(config)
	fmt.Printf("test: Config{} -> [error:%v] %v\n", err, string(buf))

	//list := []ConfigList{{Package: "package-one", Config: config}, {Package: "package-two", Config: config}}

	//buf, err = json.Marshal(list)
	//fmt.Printf("test: []ConfigList{} -> [error:%v] %v\n", err, string(buf))

	//Output:
	//test: Config{} -> [error:<nil>] {"Name":"test-route","Host":"google.com","Timeout":{"Duration":20000,"StatusCode":504},"RateLimiter":{"Limit":100,"Burst":25,"StatusCode":503},"Retry":{"Limit":100,"Burst":33,"Wait":500,"Codes":[503,504]},"Failover":null}

}

func _ExampleConfig_Unmarshal() {
	var config = Route{}
	s := "{\"Name\":\"test-route\",\"Timeout\":{\"StatusCode\":504,\"Timeout\":20000},\"RateLimiter\":{\"Limit\":100,\"Burst\":25,\"StatusCode\":503},\"Retry\":{\"Limit\":100,\"Burst\":33,\"Wait\":500,\"Codes\":[503,504]}}"

	err := json.Unmarshal([]byte(s), &config)

	//buf, err := json.Marshal(config)
	fmt.Printf("test: Config{} -> [error:%v] [%v]\n", err, config)

	//Output:
	//test: Config{} -> [error:<nil>] [{test-route {504 20µs} {100 25 503} {100 33 500ns [503 504]} {false <nil>}}]
}
