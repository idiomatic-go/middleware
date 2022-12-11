package route

import (
	"fmt"
	"golang.org/x/time/rate"
)

type Ok func(name string) bool

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

func setupCall(t Routes, name string, fn Ok) (Route, bool) {
	ok := fn(name)
	r := t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [ok:%v]", ok)
		return nil, false
	}
	return r, true
}

func Example_Timeout() {
	name := "timeout-route"
	t := NewTable()
	r := setup(t, name, 1000, NilValue, NilValue)
	if r == nil {
		return
	}
	prev := current(r)

	t.SetTimeout(name, 2000)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: SetTimeout(\"timeout-route\",2000) -> [prev:%v] [curr:%v]\n", prev.timeout, current(r).timeout)

	prev = current(r)
	t.ResetTimeout(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: ResetTimeout(\"timeout-route\") -> [prev:%v] [curr:%v]\n", prev.timeout, current(r).timeout)

	prev = current(r)
	t.DisableTimeout(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: DisableTimeout(\"timeout-route\") -> [prev:%v] [curr:%v]\n", prev.timeout, current(r).timeout)

	//Output:
	//test: SetTimeout("timeout-route",2000) -> [prev:1000] [curr:2000]
	//test: ResetTimeout("timeout-route") -> [prev:2000] [curr:1000]
	//test: DisableTimeout("timeout-route") -> [prev:1000] [curr:-1]
}
