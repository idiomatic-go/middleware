package actuator

import (
	"fmt"
)

var eventingFn Actuate = func(act Actuator, events []Event) error {
	if len(events) == 0 {
		return nil
	}
	if events[0].SLOAlertStatus == "watch" {
		if c, ok := act.RateLimiter(); ok {
			c.AdjustRateLimiter(-10)
		}
		return nil
	}
	if events[0].SLOAlertStatus == "cancel" {
		if c, ok := act.RateLimiter(); ok {
			c.AdjustRateLimiter(10)
		}
		return nil
	}
	return nil
}

func Example_Actuate() {
	table := newTable(true)
	table.Add("test", "/eventing", eventingFn, NewRateLimiterConfig(100, 25, 503))
	act := table.LookupByName("test")

	r, _ := act.RateLimiter()
	limit, burst := r.LimitAndBurst()
	fmt.Printf("test: RateLimiter() -> [limit:%v] [burst:%v]\n", limit, burst)

	act.Actuate([]Event{{SLOAlertStatus: "watch"}})
	act = table.LookupByName("test")
	r, _ = act.RateLimiter()
	limit, burst = r.LimitAndBurst()
	fmt.Printf("test: RateLimiter() -> [limit:%v] [burst:%v]\n", limit, burst)

	act.Actuate([]Event{{SLOAlertStatus: "cancel"}})
	act = table.LookupByName("test")
	r, _ = act.RateLimiter()
	limit, burst = r.LimitAndBurst()
	fmt.Printf("test: RateLimiter() -> [limit:%v] [burst:%v]\n", limit, burst)

	//Output:
	//test: RateLimiter() -> [limit:100] [burst:25]
	//test: RateLimiter() -> [limit:90] [burst:22]
	//test: RateLimiter() -> [limit:99] [burst:24]
}
