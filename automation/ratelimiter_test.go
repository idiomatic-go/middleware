package automation

import (
	"fmt"
	"golang.org/x/time/rate"
)

func Example_newRateLimiter() {
	t := newRateLimiter("test-route", newTable(true), NewRateLimiterConfig(1, 100, 503, false))
	fmt.Printf("test: newRateLimiter() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.currentConfig)

	t = newRateLimiter("test-route2", newTable(true), NewRateLimiterConfig(rate.Inf, DefaultBurst, 429, true))
	fmt.Printf("test: newRateLimiter() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.currentConfig)

	t2 := cloneRateLimiter(t)
	t2.enabled = false
	fmt.Printf("test: cloneRateLimiter() -> [prev-enabled:%v] [prev-name:%v] [curr-enabled:%v] [curr-name:%v]\n", t.enabled, t.name, t2.enabled, t2.name)

	//Output:
	//test: newRateLimiter() -> [enabled:true] [name:test-route] [config:{1 100 503 false}]
	//test: newRateLimiter() -> [enabled:true] [name:test-route2] [config:{1.7976931348623157e+308 1 429 true}]
	//test: cloneRateLimiter() -> [prev-enabled:true] [prev-name:test-route2] [curr-enabled:false] [curr-name:test-route2]

}

func Example_RateLimiter_Status() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503, true)
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

func Example_RateLimiter_Mutate() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503, false)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: State() -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(BurstName), act.RateLimiter().Attribute(StatusCodeName), act.RateLimiter().Attribute(StaticName))

	act.RateLimiter().SetLimit(5000)
	act1 := t.LookupByName(name)
	fmt.Printf("test: SetLimit(5000) -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act1.RateLimiter().Attribute(RateLimitName), act1.RateLimiter().Attribute(BurstName), act1.RateLimiter().Attribute(StatusCodeName), act1.RateLimiter().Attribute(StaticName))

	act1.RateLimiter().SetBurst(1)
	act = t.LookupByName(name)
	fmt.Printf("test: SetBurst(1) -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(BurstName), act.RateLimiter().Attribute(StatusCodeName), act.RateLimiter().Attribute(StaticName))

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: State() -> [limit:10] [burst:100] [statusCode:503] [static:false]
	//test: SetLimit(5000) -> [limit:5000] [burst:100] [statusCode:503] [static:false]
	//test: SetBurst(1) -> [limit:5000] [burst:1] [statusCode:503] [static:false]

}

func Example_RateLimiter_Mutate_Static() {
	name := "test-route"
	config := NewRateLimiterConfig(10, 100, 503, true)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: State() -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(BurstName), act.RateLimiter().Attribute(StatusCodeName), act.RateLimiter().Attribute(StaticName))

	act.RateLimiter().SetLimit(5000)
	act1 := t.LookupByName(name)
	fmt.Printf("test: SetLimit(5000) -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act1.RateLimiter().Attribute(RateLimitName), act1.RateLimiter().Attribute(BurstName), act1.RateLimiter().Attribute(StatusCodeName), act1.RateLimiter().Attribute(StaticName))

	act1.RateLimiter().SetBurst(1)
	act = t.LookupByName(name)
	fmt.Printf("test: SetBurst(1) -> [limit:%v] [burst:%v] [statusCode:%v] [static:%v]\n", act.RateLimiter().Attribute(RateLimitName), act.RateLimiter().Attribute(BurstName), act.RateLimiter().Attribute(StatusCodeName), act.RateLimiter().Attribute(StaticName))

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: State() -> [limit:10] [burst:100] [statusCode:503] [static:true]
	//test: SetLimit(5000) -> [limit:10] [burst:100] [statusCode:503] [static:true]
	//test: SetBurst(1) -> [limit:10] [burst:100] [statusCode:503] [static:true]

}
