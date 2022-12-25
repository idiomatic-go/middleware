package actuator

import (
	"fmt"
	"time"
)

func Example_newRetry() {
	t := newRetry("test-route", newTable(true), NewRetryConfig([]int{504}, time.Millisecond*500, true))
	fmt.Printf("test: newRetry() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.current)

	t = newRetry("test-route2", newTable(true), NewRetryConfig([]int{503, 504}, time.Millisecond*2000, true))
	fmt.Printf("test: newRetry() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.current)

	t2 := cloneRetry(t)
	t2.enabled = false
	fmt.Printf("test: cloneRetry() -> [prev-enabled:%v] [prev-name:%v] [curr-enabled:%v] [curr-name:%v]\n", t.enabled, t.name, t2.enabled, t2.name)

	t = newRetry("test-route3", newTable(true), NewRetryConfig([]int{503, 504}, time.Millisecond*2000, false))
	fmt.Printf("test: newRetry() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.current)

	//Output:
	//test: newRetry() -> [enabled:true] [name:test-route] [config:{true 500000000 [504]}]
	//test: newRetry() -> [enabled:true] [name:test-route2] [config:{true 2000000000 [503 504]}]
	//test: cloneRetry() -> [prev-enabled:true] [prev-name:test-route2] [curr-enabled:false] [curr-name:test-route2]
	//test: newRetry() -> [enabled:false] [name:test-route3] [config:{false 2000000000 [503 504]}]

}

func Example_Retry_Status() {
	name := "test-route"
	config := NewRetryConfig([]int{504}, time.Millisecond*2000, true)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.Retry().IsEnabled())
	prevEnabled := act.Retry().IsEnabled()

	act.Retry().Disable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act1.Retry().IsEnabled())
	prevEnabled = act1.Retry().IsEnabled()

	act1.Retry().Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act.Retry().IsEnabled())
	prevEnabled = act.Retry().IsEnabled()

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func Example_Retry_IsRetryable() {
	name := "test-route"
	config := NewRetryConfig([]int{503, 504}, time.Millisecond, true)
	t := newTable(true)
	err := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", err, t.count())

	act := t.LookupByName(name)
	enabled, ok := act.Retry().IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(500)
	fmt.Printf("test: IsRetryable(500) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(502)
	fmt.Printf("test: IsRetryable(502) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(505)
	fmt.Printf("test: IsRetryable(505) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	act.Retry().Disable()
	act = t.LookupByName(name)
	enabled, ok = act.Retry().IsRetryable(200)
	fmt.Printf("test: IsRetryable(200) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(500)
	fmt.Printf("test: IsRetryable(500) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(502)
	fmt.Printf("test: IsRetryable(502) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(503)
	fmt.Printf("test: IsRetryable(503) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(504)
	fmt.Printf("test: IsRetryable(504) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	enabled, ok = act.Retry().IsRetryable(505)
	fmt.Printf("test: IsRetryable(505) -> [enabled:%v] [ok:%v]\n", enabled, ok)

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: IsRetryable(200) -> [enabled:true] [ok:false]
	//test: IsRetryable(500) -> [enabled:true] [ok:false]
	//test: IsRetryable(502) -> [enabled:true] [ok:false]
	//test: IsRetryable(503) -> [enabled:true] [ok:true]
	//test: IsRetryable(504) -> [enabled:true] [ok:true]
	//test: IsRetryable(505) -> [enabled:true] [ok:false]
	//test: IsRetryable(200) -> [enabled:false] [ok:false]
	//test: IsRetryable(500) -> [enabled:false] [ok:false]
	//test: IsRetryable(502) -> [enabled:false] [ok:false]
	//test: IsRetryable(503) -> [enabled:false] [ok:false]
	//test: IsRetryable(504) -> [enabled:false] [ok:false]
	//test: IsRetryable(505) -> [enabled:false] [ok:false]
}
