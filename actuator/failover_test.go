package actuator

import "fmt"

var testFn FailoverInvoke = func(name string, failover bool) { fmt.Printf("test: Invoke(%v,%v)\n", name, failover) }

func Example_newFailover() {
	name := "failover-test"

	f := newFailover(name, nil, nil)
	fmt.Printf("test: newFailover(nil) -> [enabled:%v] [validate:%v]\n", f.enabled, f.validate())

	f = newFailover(name, nil, NewFailoverConfig(testFn))
	fmt.Printf("test: newFailover(testFn) -> [enabled:%v] [validate:%v]\n", f.enabled, f.validate())

	f2 := cloneFailover(f)
	f2.enabled = true
	fmt.Printf("test: cloneFailover(f1) -> [f2-enabled:%v] [f2-validate:%v]\n", f2.enabled, f2.validate())

	f.enabled = false
	fmt.Printf("test: failoverAttributes(nil) -> %v\n", failoverAttributes(nil))
	fmt.Printf("test: failoverAttributes(f1) -> %v\n", failoverAttributes(f))
	fmt.Printf("test: failoverAttributes(f2) -> %v\n", failoverAttributes(f2))

	//Output:
	//test: newFailover(nil) -> [enabled:false] [validate:invalid configuration: FailoverController FailureInvoke function cannot be nil]
	//test: newFailover(testFn) -> [enabled:false] [validate:<nil>]
	//test: cloneFailover(f1) -> [f2-enabled:true] [f2-validate:<nil>]
	//test: failoverAttributes(nil) -> [failover:null]
	//test: failoverAttributes(f1) -> [failover:false]
	//test: failoverAttributes(f2) -> [failover:true]

}

func Example_Failover_Status() {
	prevEnabled := false
	name := "failover-test"
	t := newTable(true)

	err := t.Add(name, NewFailoverConfig(testFn))
	fmt.Printf("test: Add() -> [error:%v] [count:%v]\n", err, t.count())

	f := t.LookupByName(name)
	fmt.Printf("test: IsEnabled() -> [%v]\n", f.t().failover.IsEnabled())
	prevEnabled = f.t().failover.IsEnabled()

	f.t().failover.Disable()
	f2 := t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f2.t().failover.IsEnabled())
	prevEnabled = f2.t().failover.IsEnabled()

	f2.t().failover.Enable()
	f = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f.t().failover.IsEnabled())
	prevEnabled = f.t().failover.IsEnabled()

	f.t().failover.Enable()
	f2 = t.LookupByName(name)
	fmt.Printf("test: Enable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f2.t().failover.IsEnabled())
	prevEnabled = f2.t().failover.IsEnabled()

	f2.t().failover.Disable()
	f = t.LookupByName(name)
	fmt.Printf("test: Disable() -> [prev-enabled:%v] [curr-enabled:%v]\n", prevEnabled, f.t().failover.IsEnabled())

	//Output:
	//test: Add() -> [error:<nil>] [count:1]
	//test: IsEnabled() -> [false]
	//test: Disable() -> [prev-enabled:false] [curr-enabled:false]
	//test: Enable() -> [prev-enabled:false] [curr-enabled:true]
	//test: Enable() -> [prev-enabled:true] [curr-enabled:true]
	//test: Disable() -> [prev-enabled:true] [curr-enabled:false]

}

func Example_Failover_Invoke() {
	name := "failover-test"
	t := newTable(true)
	err := t.Add(name, NewFailoverConfig(testFn))
	fmt.Printf("test: Add() -> [error:%v] [count:%v]\n", err, t.count())

	f := t.LookupByName(name)
	f.t().failover.Invoke(true)
	fmt.Printf("test: Invoke(true) -> []\n")

	f.t().failover.Invoke(false)
	fmt.Printf("test: Invoke(false) -> []\n")

	//Output:
	//test: Add() -> [error:<nil>] [count:1]
	//test: Invoke(failover-test,true)
	//test: Invoke(true) -> []
	//test: Invoke(failover-test,false)
	//test: Invoke(false) -> []
}
