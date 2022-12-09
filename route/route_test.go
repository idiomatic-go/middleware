package route

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func ExampleNewRouteWithConfig() {
	ri := NewRouteWithConfig("test", NilValue, 100, 10, true, false)
	fmt.Printf("Route(i)  : %v\n", ri)

	rt := ri.t()
	fmt.Printf("Route(t)  : %v\n", rt)
	fmt.Printf("Duration  : %v\n", ri.Duration())
	fmt.Printf("Allow     : %v\n", ri.Allow())

	ri = NewRouteWithConfig("test", 100, NilValue, 10, false, true)
	fmt.Printf("Route(i)  : %v\n", ri)
	fmt.Printf("Duration  : %v\n", ri.Duration())
	fmt.Printf("Allow     : %v\n", ri.Allow())

	//Output:
	//Route(i)  : &{test {-1 100 10} {-1 100 10} true false <nil>}
	//Route(t)  : {test {-1 100 10} {-1 100 10} true false <nil>}
	//Duration  : 0s
	//Allow     : true
	//Route(i)  : &{test {100 -1 10} {100 -1 10} false true <nil>}
	//Duration  : 100ms
	//Allow     : true
}

func ExampleRateLimiterDisallowAll() {
	ri := NewRouteWithConfig("test", NilValue, 0, 0, false, false)
	rt := ri.t()

	rt.newRateLimiter()
	fmt.Printf("Route  : %v\n", rt.current)
	fmt.Printf("Allow  : %v\n", rt.Allow())

	//Output:
	//Route  : {-1 -1 -1}
	//Allow  : false

}

func ExampleRateLimiterAllowAll() {
	ri := NewRouteWithConfig("test", NilValue, rate.Inf, 0, false, false)
	rt := ri.t()

	rt.newRateLimiter()
	fmt.Printf("Route  : %v\n", rt.current)
	fmt.Printf("Allow  : %v\n", rt.Allow())
	i := 0
	for ; i < 100; i++ {
		if !rt.Allow() {
			fmt.Printf("Allow  : fail\n")
		}
	}

	//Output:
	//Route  : {-1 1.7976931348623157e+308 -1}
	//Allow  : true

}

func ExampleRateLimiterAllowSome() {
	ri := NewRouteWithConfig("test", NilValue, 1, 1, false, false)
	rt := ri.t()
	fmt.Printf("Route  : %v\n", rt.current)
	fmt.Printf("Allow  : %v\n", rt.Allow())

	rt.newRateLimiter()
	fmt.Printf("Route  : %v\n", rt.current)
	fmt.Printf("Allow  : %v\n", rt.Allow())
	i := 0
	for ; i < 10; i++ {
		fmt.Printf("Allow  : %v\n", rt.Allow())
		if i == 1 || i == 4 || i == 7 {
			time.Sleep(time.Second)
		}
	}

	//Output:
	//Route  : {-1 1 1}
	//Allow  : true
	//Route  : {-1 1 1}
	//Allow  : true
	//Allow  : false
	//Allow  : false
	//Allow  : true
	//Allow  : false
	//Allow  : false
	//Allow  : true
	//Allow  : false
	//Allow  : false
	//Allow  : true
	//Allow  : false

}
