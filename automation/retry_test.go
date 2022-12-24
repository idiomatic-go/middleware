package automation

import (
	"fmt"
	"time"
)

func Example_newRetry() {
	t := newRetry("test-route", newTable(true), NewRetryConfig([]int{504}, time.Duration(time.Millisecond*500)))
	fmt.Printf("test: newRetry() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.current)

	t = newRetry("test-route2", newTable(true), NewRetryConfig([]int{503, 504}, time.Millisecond*2000))
	fmt.Printf("test: newRetry() -> [enabled:%v] [name:%v] [config:%v]\n", t.enabled, t.name, t.current)

	t2 := cloneRetry(t)
	t2.enabled = false
	fmt.Printf("test: cloneRetry() -> [prev-enabled:%v] [prev-name:%v] [curr-enabled:%v] [curr-name:%v]\n", t.enabled, t.name, t2.enabled, t2.name)

	//Output:
	//test: newRetry() -> [enabled:true] [name:test-route] [config:{500000000 [504]}]
	//test: newRetry() -> [enabled:true] [name:test-route2] [config:{2000000000 [503 504]}]
	//test: cloneRetry() -> [prev-enabled:true] [prev-name:test-route2] [curr-enabled:false] [curr-name:test-route2]

}

func Example_Retry_Status() {
	name := "test-route"
	config := NewRetryConfig([]int{504}, time.Millisecond*2000)
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
