package automation

import (
	"fmt"
	"time"
)

// TODO : test nil attribute

func Example_newTimeout() {
	t := newTimeout("test-route", newTable(true), NewTimeoutConfig(NilValue))
	fmt.Printf("test: newTimeout() -> [enabled:%v] [name:%v] [current:%v] [default:%v]\n", t.enabled, t.name, t.currentConfig.timeout, t.defaultConfig.timeout)

	t = newTimeout("test-route2", newTable(true), NewTimeoutConfig(time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [enabled:%v] [name:%v] [current:%v] [default:%v]\n", t.enabled, t.name, t.currentConfig.timeout, t.defaultConfig.timeout)

	t2 := cloneTimeout(t)
	t2.enabled = false
	fmt.Printf("test: cloneTimeout() -> [prev-enabled:%v] [prev-name:%v] [curr-enabled:%v] [curr-name:%v]\n", t.enabled, t.name, t2.enabled, t2.name)

	//Output:
	//test: newTimeout() -> [enabled:true] [name:test-route] [current:-1ns] [default:-1ns]
	//test: newTimeout() -> [enabled:true] [name:test-route2] [current:2s] [default:2s]
	//test: cloneTimeout() -> [prev-enabled:true] [prev-name:test-route2] [curr-enabled:false] [curr-name:test-route2]

}

func Example_Timeout_Attribute() {
	t := newTimeout("test-route", newTable(true), NewTimeoutConfig(time.Millisecond*2000))
	fmt.Printf("test: IsEnabled() -> [%v]\n", t.IsEnabled())

	_, d := t.Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)

	t = newTimeout("test-route", newTable(true), NewTimeoutConfig(time.Millisecond*2000))

	a := t.Attribute("")
	fmt.Printf("test: Attribute(\"\") -> [name:%v] [value:%v] [string:%v]\n", a.Name(), a.Value(), a)

	a = t.Attribute(TimeoutName)
	fmt.Printf("test: Attribute(\"Timeout\") -> [name:%v] [value:%v] [string:%v]\n", a.Name(), a.Value(), a)

	//Output:
	//test: IsEnabled() -> [true]
	//test: Duration() -> [2s]
	//test: Attribute("") -> [name:] [value:<nil>] [string:nil]
	//test: Attribute("Timeout") -> [name:timeout] [value:2s] [string:2s]

}

func Example_Timeout_Status() {
	name := "test-route"
	config := NewTimeoutConfig(time.Millisecond * 2000)
	t := newTable(true)

	//fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())
	ok := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	act := t.LookupByName(name)
	_, d := act.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.Timeout().IsEnabled())
	prevEnabled := act.Timeout().IsEnabled()

	act.Timeout().Disable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act1.Timeout().IsEnabled())
	prevEnabled = act.Timeout().IsEnabled()

	act1.Timeout().Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, act.Timeout().IsEnabled())
	prevEnabled = act.Timeout().IsEnabled()

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: Duration() -> [2s]
	//test: IsEnabled() -> [true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]

}

func Example_Timeout_SetTimeout() {
	name := "test-route"
	config := NewTimeoutConfig(time.Millisecond * 1500)
	t := newTable(true)

	ok := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	act := t.LookupByName(name)
	_, d := act.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)
	_, prevDuration := act.Timeout().Duration()

	act.Timeout().SetTimeout(time.Second * 2)
	act1 := t.LookupByName(name)
	_, d = act1.Timeout().Duration()
	fmt.Printf("test: SetTimeout(2s) -> [prev-duration:%v] [curr-duration:%v]\n", prevDuration, d)
	_, prevDuration = act1.Timeout().Duration()

	act1.Timeout().Reset()
	act = t.LookupByName(name)
	_, d = act.Timeout().Duration()
	fmt.Printf("test: Reset() -> [prev-duration:%v] [curr-duration:%v]\n", prevDuration, d)

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: Duration() -> [1.5s]
	//test: SetTimeout(2s) -> [prev-duration:1.5s] [curr-duration:2s]
	//test: Reset() -> [prev-duration:2s] [curr-duration:1.5s]

}
