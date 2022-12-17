package automation

import (
	"fmt"
	"golang.org/x/time/rate"
)

type Ok func(name string) bool

/*
func current(r Route) config {
	if r == nil {
		return config{NilValue, NilValue, NilValue}
	}
	return r.(*route).current
}


*/
func setup(t *table, name string, tcfg *TimeoutConfig, limit rate.Limit, burst int) Actuator {
	//r, err := NewRouteWithConfig(name, timeout, limit, burst, false, false)
	//if err != nil {
	//	fmt.Printf("test: New(2000,_,_) -> [err:%v]\n", err)
	//		return nil
	//}
	ok := t.Add(name, tcfg)
	if !ok {
		fmt.Printf("test: Add(actuator) -> [ok:%v]", ok)
		return nil
	}
	return t.LookupByName(name)
}

/*
func setupCall(t Routes, name string, fn Ok) (Route, bool) {
	ok := fn(name)
	r := t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [ok:%v]", ok)
		return nil, false
	}
	return r, true
}


*/
func Example_Timeout() {
	name := "timeout-route"
	t := newTable()
	a0 := setup(t, name, NewTimeoutConfig(1000, NilValue), NilValue, NilValue)
	if a0 == nil {
		return
	}
	prev := a0.Timeout().(Controller).State()

	t.setTimeout(name, 2000)
	a := t.LookupByName(name)
	if a == nil {
		fmt.Printf("test: LookupByName(%v) -> [route:%v]", name, a)
		return
	}
	curr := a.Timeout().(Controller).State()
	fmt.Printf("test: setTimeout(%v) -> [prev:%v] [curr:%v]\n", name, prev, curr)

	prev = curr
	t.resetTimeout(name)
	a = t.LookupByName(name)
	if a == nil {
		fmt.Printf("test: LookupByName(%v) -> [actuator:%v]", name, a)
		return
	}
	curr = a.Timeout().(Controller).State()
	fmt.Printf("test: resetTimeout(%v) -> [prev:%v] [curr:%v]\n", name, prev, curr)

	prev = curr
	t.disableTimeout(name)
	a = t.LookupByName(name)
	if a == nil {
		fmt.Printf("test: LookupByName(%v) -> [route:%v]", name, a)
		return
	}
	curr = a.Timeout().(Controller).State()
	fmt.Printf("test: disableTimeout(%v) -> [prev:%v] [curr:%v]\n", name, prev, curr)

	//Output:
	//test: SetTimeout("timeout-route",2000) -> [prev:1000] [curr:2000]
	//test: ResetTimeout("timeout-route") -> [prev:2000] [curr:1000]
	//test: DisableTimeout("timeout-route") -> [prev:1000] [curr:-1]
}

/*
func Example_RateLimit_Limit() {
	name := "limit-route"
	t := NewTable()
	r := setup(t, name, NilValue, 100, 25)
	if r == nil {
		return
	}
	prev := current(r)

	t.SetLimit(name, 50)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: SetLimit(\"limit-route\",50) -> [prev:%v] [curr:%v]\n", prev.limit, current(r).limit)

	prev = current(r)
	t.ResetLimit(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: ResetLimit(\"limit-route\") -> [prev:%v] [curr:%v]\n", prev.limit, current(r).limit)

	prev = current(r)
	t.DisableLimiter(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: DisableLimiter(\"limit-route\") -> [prev:%v] [curr:%v]\n", prev.limit, current(r).limit)

	//Output:
	//test: SetLimit("limit-route",50) -> [prev:100] [curr:50]
	//test: ResetLimit("limit-route") -> [prev:50] [curr:100]
	//test: DisableLimiter("limit-route") -> [prev:100] [curr:1.7976931348623157e+308]

}

func Example_RateLimit_Burst() {
	name := "limit-route"
	t := NewTable()
	r := setup(t, name, NilValue, 100, 25)
	if r == nil {
		return
	}
	prev := current(r)

	t.SetBurst(name, 10)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: SetBurst(\"limit-route\",10) -> [prev:%v] [curr:%v]\n", prev.burst, current(r).burst)

	prev = current(r)
	t.ResetBurst(name)
	r = t.LookupByName(name)
	if r == nil {
		fmt.Printf("test: Lookup(route) -> [route:%v]", r)
		return
	}
	fmt.Printf("test: ResetBurst(\"limit-route\") -> [prev:%v] [curr:%v]\n", prev.burst, current(r).burst)

	//Output:
	//test: SetBurst("limit-route",10) -> [prev:25] [curr:10]
	//test: ResetBurst("limit-route") -> [prev:10] [curr:25]

}

*/
