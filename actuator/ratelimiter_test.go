package actuator

import (
	"fmt"
	"golang.org/x/time/rate"
)

func Example_newRateLimiter() {
	t := newRateLimiter("test-route", newTable(true), NewRateLimiterConfig(1, 100, 503))
	fmt.Printf("test: newRateLimiter() -> [name:%v] [config:%v]\n", t.name, t.config)

	t = newRateLimiter("test-route2", newTable(true), NewRateLimiterConfig(rate.Inf, DefaultBurst, 429))
	fmt.Printf("test: newRateLimiter() -> [name:%v] [config:%v]\n", t.name, t.config)

	t2 := cloneRateLimiter(t)
	t2.config.limit = 123
	fmt.Printf("test: cloneRateLimiter() -> [prev-limit:%v] [prev-name:%v] [curr-limit:%v] [curr-name:%v]\n", t.config.limit, t.name, t2.config.limit, t2.name)

	//Output:
	//test: newRateLimiter() -> [name:test-route] [config:{1 100 503}]
	//test: newRateLimiter() -> [name:test-route2] [config:{1.7976931348623157e+308 1 429}]
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
func Example_RateLimiter_Mutate() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: rateLimiterAttributes() -> %v\n", rateLimiterAttributes(act.t().rateLimiter))

	act.t().rateLimiter.SetLimit(5000)
	act1 := t.LookupByName(name)
	//fmt.Printf("test: SetLimit(5000) -> [limit:%v] [burst:%v] [statusCode:%v]\n", act1.RateLimiter().Attribute(RateLimitName), act1.RateLimiter().Attribute(RateBurstName), act1.RateLimiter().Attribute(StatusCodeName))

	act1.t().rateLimiter.SetBurst(1)
	act = t.LookupByName(name)
	//fmt.Printf("test: SetBurst(1) -> [limit:%v] [burst:%v] [statusCode:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(RateBurstName), act.RateLimiter().Attribute(StatusCodeName))

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: rateLimiterAttributes() -> [rateLimit:10 burst:100 statusCode:503]

}

/*
func Example_RateLimiter_Mutate_Static() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: State() -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(RateBurstName), act.RateLimiter().Attribute(StatusCodeName), act.RateLimiter().Attribute(StaticName))

	act.RateLimiter().SetLimit(5000)
	act1 := t.LookupByName(name)
	fmt.Printf("test: SetLimit(5000) -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act1.RateLimiter().Attribute(RateLimitName), act1.RateLimiter().Attribute(RateBurstName), act1.RateLimiter().Attribute(StatusCodeName), act1.RateLimiter().Attribute(StaticName))

	act1.RateLimiter().SetBurst(1)
	act = t.LookupByName(name)
	fmt.Printf("test: SetBurst(1) -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(RateBurstName), act.RateLimiter().Attribute(StatusCodeName), act.RateLimiter().Attribute(StaticName))

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: State() -> [limit:10] [burst:100] [statusCode:503] [static:true]
	//test: SetLimit(5000) -> [limit:10] [burst:100] [statusCode:503] [static:true]
	//test: SetBurst(1) -> [limit:10] [burst:100] [statusCode:503] [static:true]

}

*/
