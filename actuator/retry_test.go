package actuator

import (
	"fmt"
)

func Example_newRetry() {
	t := newRetry("test-route", newTable(true), NewRetryConfig([]int{504}, 5, 10))
	fmt.Printf("test: newRetry() -> [name:%v] [config:%v]\n", t.name, t.config)

	t = newRetry("test-route2", newTable(true), NewRetryConfig([]int{503, 504}, 2, 20))
	fmt.Printf("test: newRetry() -> [name:%v] [config:%v]\n", t.name, t.config)

	t2 := cloneRetry(t)
	t2.enabled = true
	fmt.Printf("test: cloneRetry() -> [prev-enabled:%v] [curr-enabled:%v]\n", t.enabled, t2.enabled)

	//t = newRetry("test-route3", newTable(true), NewRetryConfig([]int{503, 504}, time.Millisecond*2000, false))
	fmt.Printf("test: retryAttributes(nil) -> %v\n", retryAttributes(nil, false))
	fmt.Printf("test: retryAttributes(t,false) -> %v\n", retryAttributes(t, false))
	fmt.Printf("test: retryAttributes(t,true) -> %v\n", retryAttributes(t, true))

	//Output:
	//test: newRetry() -> [name:test-route] [config:{5 10 [504]}]
	//test: newRetry() -> [name:test-route2] [config:{2 20 [503 504]}]
	//test: cloneRetry() -> [prev-enabled:false] [curr-enabled:true]
	//test: retryAttributes(nil) -> [retry:null]
	//test: retryAttributes(t,false) -> [retry:false retryRateLimit:2 retryBurst:20]
	//test: retryAttributes(t,true) -> [retry:true retryRateLimit:2 retryBurst:20]

}

func Example_Retry_Status() {
	name := "test-route"
	config := NewRetryConfig([]int{504}, 5, 10)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.Retry().IsEnabled())
	prevEnabled := act.Retry().IsEnabled()

	act.Retry().Enable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act1.Retry().IsEnabled())
	prevEnabled = act1.Retry().IsEnabled()

	act1.Retry().Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act.Retry().IsEnabled())
	prevEnabled = act.Retry().IsEnabled()

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsEnabled() -> [false]
	//test: Disable() -> [prev-enabled:false] [curr-enabled:true]
	//test: Enable() -> [prev-enabled:true] [curr-enabled:true]

}

func Example_Retry_IsRetryable_Disabled() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, 100, 10)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	//act.Retry().Enable()
	//act = t.LookupByName(name)
	ok, status := act.Retry().IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsRetryable(200) -> [ok:false] [status:NE]
	//test: IsRetryable(503) -> [ok:false] [status:NE]
	//test: IsRetryable(504) -> [ok:false] [status:NE]

}

func Example_Retry_IsRetryable_StatusCode() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, 100, 10)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	act.Retry().Enable()
	act = t.LookupByName(name)
	ok, status := act.Retry().IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(500)
	fmt.Printf("test: IsRetryable(500) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(502)
	fmt.Printf("test: IsRetryable(502) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(505)
	fmt.Printf("test: IsRetryable(505) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsRetryable(200) -> [ok:false] [status:SC]
	//test: IsRetryable(500) -> [ok:false] [status:SC]
	//test: IsRetryable(502) -> [ok:false] [status:SC]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:true] [status:]
	//test: IsRetryable(505) -> [ok:false] [status:SC]

}

func Example_Retry_IsRetryable_RateLimit() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, 1, 1)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	act.Retry().Enable()
	act = t.LookupByName(name)
	ok, status := act.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	act.Retry().SetRateLimiter(100, 10)
	act = t.LookupByName(name)
	ok, status = act.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [ok:%v] [status:%v]\n", ok, status)

	ok, status = act.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [ok:%v] [status:%v]\n", ok, status)

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:false] [status:RL]
	//test: IsRetryable(503) -> [ok:true] [status:]
	//test: IsRetryable(504) -> [ok:true] [status:]

}
