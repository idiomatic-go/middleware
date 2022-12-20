package automation

import "fmt"

var testFn FailoverInvoke = func(name string) { fmt.Printf("test: Invoke() -> []\n") }

func Example_newFailover() {
	name := "failover-test"

	f := newFailover(name, nil, nil)
	fmt.Printf("test: newFailover(nil) -> [enabled:%v]\n", f.enabled)

	f = newFailover(name, NewFailoverConfig(testFn), nil)
	fmt.Printf("test: newFailover(testFn) -> [enabled:%v]\n", f.enabled)

	f2 := cloneFailover(f)
	fmt.Printf("test: cloneFailover(f1) -> [f2-enabled:%v]\n", f2.enabled)

	f.enabled = false
	fmt.Printf("test: Attribute(f1) -> [enabled:%v]\n", f.enabled)
	fmt.Printf("test: Attribute(f2) -> [enabled:%v]\n", f2.enabled)

	//Output:
	//test: newFailover(nil) -> [enabled:false]
	//test: newFailover(testFn) -> [enabled:true]
	//test: cloneFailover(f1) -> [f2-enabled:true]
	//test: Attribute(f1) -> [enabled:false]
	//test: Attribute(f2) -> [enabled:true]
}

func Example_Failover_Controller_Status() {
	name := "failover-test"
	t := newTable()

	ok := t.Add(name, nil, nil, NewFailoverConfig(testFn))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	f := t.LookupByName(name)
	fmt.Printf("test: IsEnabled(f1) -> [%v]\n", f.Failover().IsEnabled())

	f.Failover().Disable()
	f2 := t.LookupByName(name)
	fmt.Printf("test: IsEnabled(f2) -> [%v]\n", f2.Failover().IsEnabled())

	f2.Failover().Enable()
	f = t.LookupByName(name)
	fmt.Printf("test: IsEnabled(f1) -> [%v]\n", f.Failover().IsEnabled())

	//Output:
	//test: Add() -> [true] [count:1]
	//test: IsEnabled(f1) -> [true]
	//test: IsEnabled(f2) -> [false]
	//test: IsEnabled(f1) -> [true]
}

func Example_Failover_Controller_State() {
	name := "failover-test"
	t := newTable()

	ok := t.Add(name, nil, nil, NewFailoverConfig(testFn))
	fmt.Printf("test: Add() -> [%v] [count:%v]\n", ok, t.count())

	f1 := t.LookupByName(name)
	f1.Failover().Invoke()

	//Output:
	//test: Add() -> [true] [count:1]
	//test: Invoke() -> []
}
