package actuator

import (
	"fmt"
	"time"
)

func Example_newTimeout() {
	t := newTimeout("test-route", newTable(true), NewTimeoutConfig(NilValue))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.timeout)

	t = newTimeout("test-route2", newTable(true), NewTimeoutConfig(time.Millisecond*2000))
	fmt.Printf("test: newTimeout() -> [name:%v] [current:%v]\n", t.name, t.config.timeout)

	t2 := cloneTimeout(t)
	t2.config.timeout = time.Millisecond * 1000
	fmt.Printf("test: cloneTimeout() -> [prev-config:%v] [prev-name:%v] [curr-config:%v] [curr-name:%v]\n", t.config, t.name, t2.config, t2.name)

	//Output:
	//test: newTimeout() -> [name:test-route] [current:-1ns]
	//test: newTimeout() -> [name:test-route2] [current:2s]
	//test: cloneTimeout() -> [prev-config:{2000000000}] [prev-name:test-route2] [curr-config:{1000000000}] [curr-name:test-route2]

}

func Example_Timeout_State() {
	t := newTimeout("test-route", newTable(true), NewTimeoutConfig(time.Millisecond*2000))

	d := t.Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)

	t = newTimeout("test-route", newTable(true), NewTimeoutConfig(time.Millisecond*2000))

	st := timeoutState(nil)
	fmt.Printf("test: state(nil) -> %v\n", st)

	st = timeoutState(t)
	fmt.Printf("test: state(t) -> %v\n", st)

	//Output:
	//test: Duration() -> [2s]
	//test: state(nil) -> [timeout:-1]
	//test: state(t) -> [timeout:2000]

}

func Example_Timeout_SetTimeout() {
	name := "test-route"
	config := NewTimeoutConfig(time.Millisecond * 1500)
	t := newTable(true)

	ok := t.Add(name, config)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	act := t.LookupByName(name)
	d := act.Timeout().Duration()
	fmt.Printf("test: Duration() -> [%v]\n", d)
	prevDuration := act.Timeout().Duration()

	act.Timeout().SetTimeout(time.Second * 2)
	act1 := t.LookupByName(name)
	d = act1.Timeout().Duration()
	fmt.Printf("test: SetTimeout(2s) -> [prev-duration:%v] [curr-duration:%v]\n", prevDuration, d)
	prevDuration = act1.Timeout().Duration()

	st := timeoutState(act1.Timeout().(*timeout))
	fmt.Printf("test: state(t) -> %v\n", st)

	//Output:
	//test: Add() -> [<nil>] [count:1]
	//test: Duration() -> [1.5s]
	//test: SetTimeout(2s) -> [prev-duration:1.5s] [curr-duration:2s]
	//test: state(t) -> [timeout:2000]

}
