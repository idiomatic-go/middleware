package route

import (
	"fmt"
	"golang.org/x/time/rate"
)

func safeLookup(t Routes, name string) Route {
	r := t.LookupByName(name)
	if r == nil {
		return newRoute("empty")
	}
	return r
}

func current(r Route) config {
	if r == nil {
		return config{NilValue, NilValue, NilValue}
	}
	return r.(*route).current
}

func setup(t Routes, name string, timeout int, limit rate.Limit, burst int) Route {
	r, err := NewRouteWithConfig(name, timeout, limit, burst, false, false)
	if err != nil {
		fmt.Printf("test: New(2000,_,_) -> [err:%v]\n", err)
		return nil
	}
	ok := t.Add(r)
	if !ok {
		fmt.Printf("test: Add(route) -> [ok:%v]", ok)
		return nil
	}
	return r
}

func Example_Timeout() {
	name := "timeout-route"
	t := NewTable()
	r := setup(t, name, 1000, NilValue, NilValue)
	if r == nil {
		return
	}
	prev := current(r)

	ok := t.SetTimeout(name, 2000)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [ok:%v]", ok)
		return
	}
	fmt.Printf("test: SetTimeout(\"timeout-route\",2000) -> [ok:%v] [prev:%v] [curr:%v]\n", ok, prev.timeout, current(r).timeout)

	prev = current(r)
	ok = t.ResetTimeout(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [ok:%v]", ok)
		return
	}
	fmt.Printf("test: ResetTimeout(\"timeout-route\") -> [ok:%v] [prev:%v] [curr:%v]\n", ok, prev.timeout, current(r).timeout)

	prev = current(r)
	ok = t.DisableTimeout(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [ok:%v]", ok)
		return
	}
	fmt.Printf("test: DisableTimeout(\"timeout-route\") -> [ok:%v] [prev:%v] [curr:%v]\n", ok, prev.timeout, current(r).timeout)

	//Output:
	//test: SetTimeout("timeout-route",2000) -> [ok:true] [prev:1000] [curr:2000]
	//test: ResetTimeout("timeout-route") -> [ok:true] [prev:2000] [curr:1000]
	//test: DisableTimeout("timeout-route") -> [ok:true] [prev:1000] [curr:-1]
}
