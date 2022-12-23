package automation

import "fmt"

var testFn FailoverInvoke = func(name string, failover bool) { fmt.Printf("test: Invoke(%v,%v)\n", name, failover) }

func Example_newFailover() {
	name := "failover-test"

	f := newFailover(name, nil, nil)
	fmt.Printf("test: newFailover(nil) -> [enabled:%v]\n", f.enabled)
	fmt.Printf("test: newFailover(nil) -> [validate:%v]\n", f.validate())

	f = newFailover(name, nil, NewFailoverConfig(testFn))
	fmt.Printf("test: newFailover(testFn) -> [enabled:%v]\n", f.enabled)
	fmt.Printf("test: newFailover(testFn) -> [validate:%v]\n", f.validate())

	//f2 := cloneFailover(f)
	//fmt.Printf("test: cloneFailover(f1) -> [f2-enabled:%v]\n", f2.enabled)

	//f.enabled = false
	//fmt.Printf("test: Attribute(f1) -> [enabled:%v]\n", f.enabled)
	//fmt.Printf("test: Attribute(f2) -> [enabled:%v]\n", f2.enabled)

	//Output:
	//test: newFailover(nil) -> [enabled:false]
	//test: newFailover(nil) -> [validate:invalid configuration: failover controller FailureInvoke function cannot be nil]
	//test: newFailover(testFn) -> [enabled:false]
	//test: newFailover(testFn) -> [validate:<nil>]

}

func Example_Failover_Controller_Status() {
	prevEnabled := false
	name := "failover-test"
	t := newTable(true)

	err := t.Add(name, NewFailoverConfig(testFn))
	fmt.Printf("test: Add() -> [error:%v] [count:%v]\n", err, t.count())

	f := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", f.Failover().IsEnabled())
	prevEnabled = f.Failover().IsEnabled()

	f.Failover().Disable()
	f2 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f2.Failover().IsEnabled())
	prevEnabled = f2.Failover().IsEnabled()

	f2.Failover().Enable()
	f = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f.Failover().IsEnabled())
	prevEnabled = f.Failover().IsEnabled()

	f.Failover().Enable()
	f2 = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f2.Failover().IsEnabled())
	prevEnabled = f2.Failover().IsEnabled()

	f2.Failover().Disable()
	f = t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f.Failover().IsEnabled())

	//Output:
	//test: Add() -> [error:<nil>] [count:1]
	//test: IsEnabled() -> [false]
	//test: Disable() -> [prev-enabled:false] [curr-enabled:false]
	//test: Invoke(failover-test,true)
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]
	//test: Enable() -> [prev-enabled:true] [curr-enabled:true]
	//test: Invoke(failover-test,false)
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]

}
