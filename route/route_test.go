package route

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func Example_NewRouteWithConfig_Error() {
	r, err := newRouteWithConfig("", NilValue, 100, 10, true, false)
	fmt.Printf("Route  : %v [%v]\n", r, err)

	r, err = newRouteWithConfig("test", NilValue, 100, 0, true, false)
	fmt.Printf("Route  : %v [%v]\n", r, err)

	r, err = newRouteWithConfig("test", NilValue, 100, NilValue, true, false)
	fmt.Printf("Route  : %v [%v]\n", r, err)

	r, err = newRouteWithConfig("test", NilValue, 0, 10, true, false)
	fmt.Printf("Route  : %v [%v]\n", r, err)

	r, err = newRouteWithConfig("test", NilValue, NilValue, 10, true, false)
	fmt.Printf("Route  : %v [%v]\n", r, err)

	//r,err = newRouteWithConfig("test", NilValue, 100, 10, true, false)
	//fmt.Printf("Route  : %v\n", r)
	//fmt.Printf("Error  : %v\n", err)
	//fmt.Printf("Duration  : %v\n", r.Duration())
	//fmt.Printf("Allow     : %v\n", r.Allow())

	//r,err = newRouteWithConfig("test", 100, NilValue, 10, false, true)
	//fmt.Printf("Route  : %v\n", r)
	//fmt.Printf("Error  : %v\n", err)
	//fmt.Printf("Duration  : %v\n", r.Duration())
	//fmt.Printf("Allow     : %v\n", r.Allow())

	//Output:
	//Route  : <nil> [invalid argument : route name is empty]
	//Route  : <nil> [invalid argument : limit is configured but burst is not]
	//Route  : <nil> [invalid argument : limit is configured but burst is not]
	//Route  : <nil> [invalid argument : burst is configured but limit is not]
	//Route  : <nil> [invalid argument : burst is configured but limit is not]
}

func Example_NewRouteWithConfig() {
	r, err := newRouteWithConfig("test", NilValue, 0, 0, true, false)
	fmt.Printf("Route       : %v [%v]\n", r.original, err)
	fmt.Printf("RateLimiter : %v\n", r.IsRateLimiter())
	fmt.Printf("Duration    : %v\n", r.Duration())
	fmt.Printf("Allow       : %v\n", r.Allow())

	r, err = newRouteWithConfig("test", NilValue, 100, 10, true, false)
	fmt.Printf("Route       : %v [%v]\n", r.original, err)
	fmt.Printf("RateLimiter : %v\n", r.IsRateLimiter())
	fmt.Printf("Allow       : %v\n", r.Allow())

	r, err = newRouteWithConfig("test", NilValue, rate.Inf, NilValue, true, false)
	fmt.Printf("Route       : %v [%v]\n", r.original, err)
	fmt.Printf("RateLimiter : %v\n", r.IsRateLimiter())
	fmt.Printf("Allow       : %v\n", r.Allow())

	r, err = newRouteWithConfig("test", NilValue, rate.Inf, 0, true, false)
	fmt.Printf("Route       : %v [%v]\n", r.original, err)
	fmt.Printf("RateLimiter : %v\n", r.IsRateLimiter())
	fmt.Printf("Allow       : %v\n", r.Allow())

	//r,err = newRouteWithConfig("test", 100, NilValue, 10, false, true)
	//fmt.Printf("Route  : %v\n", r)
	//fmt.Printf("Error  : %v\n", err)
	//fmt.Printf("Duration  : %v\n", r.Duration())
	//fmt.Printf("Allow     : %v\n", r.Allow())

	//Output:
	//
}

/*
func Example_RateLimiterDisallowAll() {
	r,err := newRouteWithConfig("test", NilValue, 0, 0, false, false)

	fmt.Printf("Route  : %v\n", r.current)
	fmt.Printf("Allow  : %v\n", r.Allow())

	//Output:
	//Route  : {-1 -1 -1}
	//Allow  : false
}

*/

func _Example_RateLimiterAllowAll() {
	r, err := newRouteWithConfig("test", NilValue, rate.Inf, 0, false, false)
	fmt.Printf("Error  : %v\n", err)

	fmt.Printf("Route  : %v\n", r.current)
	fmt.Printf("Allow  : %v\n", r.Allow())
	i := 0
	for ; i < 100; i++ {
		if !r.Allow() {
			fmt.Printf("Allow  : fail\n")
		}
	}

	//Output:
	//Route  : {-1 1.7976931348623157e+308 -1}
	//Allow  : true

}

func _Example_RateLimiterAllowSome() {
	r, err := newRouteWithConfig("test", NilValue, 1, 1, false, false)
	fmt.Printf("Error  : %v\n", err)

	fmt.Printf("Route  : %v\n", r.current)
	fmt.Printf("Allow  : %v\n", r.Allow())

	for i := 0; i < 10; i++ {
		fmt.Printf("Allow  : %v\n", r.Allow())
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
