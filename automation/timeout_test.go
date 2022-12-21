package automation

import "fmt"

// TODO : test nil attribute

func Example_newTimeout() {
	t := newTimeout("test-route", nil, newTable())
	fmt.Printf("test: newTimeout() -> [enabled:%v] [name:%v] [default:%v] [current:%v]\n", t.enabled, t.name, t.defaultC, t.current)

	t = newTimeout("test-route2", NewTimeoutConfig(2000), newTable())
	fmt.Printf("test: newTimeout() -> [enabled:%v] [name:%v] [default:%v] [current:%v]\n", t.enabled, t.name, t.defaultC, t.current)

	t2 := cloneTimeout(t)
	t2.enabled = false
	fmt.Printf("test: cloneTimeout() -> [prev-enabled:%v] [prev-name:%v] [curr-enabled:%v] [curr-name:%v]\n", t.enabled, t.name, t2.enabled, t2.name)

	//Output:
	//test: newTimeout() -> [enabled:false] [name:test-route] [default:{-1}] [current:{-1}]
	//test: newTimeout() -> [enabled:true] [name:test-route2] [default:{2000}] [current:{2000}]
	//test: cloneTimeout() -> [prev-enabled:true] [prev-name:test-route2] [curr-enabled:false] [curr-name:test-route2]
}

func Example_Controller_ReadOnly() {
	t := newTimeout("test-route", NewTimeoutConfig(2000), newTable())
	fmt.Printf("test: IsEnabled() -> [%v]\n", t.IsEnabled())

	fmt.Printf("test: Duration() -> [%v]\n", t.Duration())

	t = newTimeout("test-route", NewTimeoutConfig(2000), newTable())

	a := t.Attribute("")
	fmt.Printf("test: Attribute(\"\") -> [name:%v] [value:%v] [string:%v]\n", a.Name(), a.Value(), a)

	a = t.Attribute(TimeoutName)
	fmt.Printf("test: Attribute(\"Timeout\") -> [name:%v] [value:%v] [string:%v]\n", a.Name(), a.Value(), a)

	//Output:
	//test: IsEnabled() -> [true]
	//test: Duration() -> [2s]
	//test: Attribute("") -> [name:] [value:<nil>] [string:nil]
	//test: Attribute("Timeout") -> [name:timeout] [value:2000] [string:2000]
}

func Example_Timeout_Controller_Status() {
	name := "test-route"
	config := NewTimeoutConfig(2000)
	t := newTable()

	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())
	ok := t.Add(name, config, nil, nil, nil)
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	act := t.LookupByName(name)
	//fmt.Printf("test: LookupByName(%v) -> [%v]\n", name, act != nil)

	fmt.Printf("test: Duration() -> [%v]\n", act.Timeout().Duration())

	fmt.Printf("test: IsEnabled() -> [%v]\n", act.Timeout().IsEnabled())

	act.Timeout().Disable()
	act1 := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act1.Timeout().IsEnabled())

	act1.Timeout().Enable()
	act = t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", act.Timeout().IsEnabled())

	//Output:
	//test: empty() -> [true]
	//test: Add() -> [true] [count:1]
	//test: Duration() -> [2s]
	//test: IsEnabled() -> [true]
	//test: IsEnabled() -> [false]
	//test: IsEnabled() -> [true]
}
