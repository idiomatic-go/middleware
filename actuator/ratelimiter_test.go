package actuator

import (
	"fmt"
	"golang.org/x/time/rate"
)

func Example_newRateLimiter() {
	t := newRateLimiter("test-route", newTable(true), NewRateLimiterConfig(1, 100, 503))
	limit, burst := t.LimitAndBurst()
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [statusCode:%v]\n", t.name, limit, burst, t.StatusCode())

	t = newRateLimiter("test-route2", newTable(true), NewRateLimiterConfig(rate.Inf, DefaultBurst, 429))
	limit, burst = t.LimitAndBurst()
	fmt.Printf("test: newRateLimiter() -> [name:%v] [limit:%v] [burst:%v] [statusCode:%v]\n", t.name, limit, burst, t.StatusCode())

	t2 := cloneRateLimiter(t)
	t2.config.limit = 123
	fmt.Printf("test: cloneRateLimiter() -> [prev-limit:%v] [prev-name:%v] [curr-limit:%v] [curr-name:%v]\n", t.config.limit, t.name, t2.config.limit, t2.name)

	//Output:
	//test: newRateLimiter() -> [name:test-route] [limit:1] [burst:100] [statusCode:503]
	//test: newRateLimiter() -> [name:test-route2] [limit:1.7976931348623157e+308] [burst:1] [statusCode:429]
	//test: cloneRateLimiter() -> [prev-limit:1.7976931348623157e+308] [prev-name:test-route2] [curr-limit:123] [curr-name:test-route2]

}

/*
func Example_RateLimiter_Status() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.RateLimiter().IsEnabled())
	prevEnabled := act.RateLimiter().IsEnabled()

	act.RateLimiter().Disable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act1.RateLimiter().IsEnabled())
	prevEnabled = act1.RateLimiter().IsEnabled()

	act1.RateLimiter().Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act.RateLimiter().IsEnabled())
	prevEnabled = act.RateLimiter().IsEnabled()

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}
*/

func Example_RateLimiter_Set() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503)
	t := newTable(true)
	err := t.Add(name, nil, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: rateLimiterState(map,t) -> %v\n", rateLimiterState(nil, act.t().rateLimiter))

	act.t().rateLimiter.SetLimit(rate.Inf)
	act1 := t.LookupByName(name)
	fmt.Printf("test: SetLimit(rate.Inf) -> %v\n", rateLimiterState(nil, act1.t().rateLimiter))

	act1.t().rateLimiter.SetBurst(1)
	act = t.LookupByName(name)
	fmt.Printf("test: SetBurst(1) -> %v\n", rateLimiterState(nil, act.t().rateLimiter))

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: rateLimiterState(map,t) -> map[burst:100 rateLimit:10]
	//test: SetLimit(rate.Inf) -> map[burst:100 rateLimit:99999]
	//test: SetBurst(1) -> map[burst:1 rateLimit:99999]

}

func Example_RateLimiter_Adjust() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503)
	t := newTable(true)
	err := t.Add(name, nil, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	act.t().rateLimiter.AdjustRateLimiter(10)
	act1 := t.LookupByName(name)
	fmt.Printf("test: rateLimiterState(map,t) -> %v\n", rateLimiterState(nil, act1.t().rateLimiter))

	//Output:
	//test: Add() -> [[]] [count:1]
	//test: rateLimiterState(map,t) -> map[burst:100 rateLimit:10]
	//test: SetLimit(rate.Inf) -> map[burst:100 rateLimit:99999]
	//test: SetBurst(1) -> map[burst:1 rateLimit:99999]

}
