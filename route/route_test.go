package route

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func Example_NewRouteWithConfig_Error() {
	_, err := newRouteWithConfig("", NilValue, 100, 10, true, false)
	fmt.Printf("test: new(\"\",_,_) -> [err:%v]\n", err)

	_, err = newRouteWithConfig("test", NilValue, 100, 0, true, false)
	fmt.Printf("test: new(_,100,0) -> [err:%v]\n", err)

	_, err = newRouteWithConfig("test", NilValue, 100, NilValue, true, false)
	fmt.Printf("test: new(_,100,NilValue) -> [err:%v]\n", err)

	_, err = newRouteWithConfig("test", NilValue, 0, 10, true, false)
	fmt.Printf("test: new(_,0,10) -> [err:%v]\n", err)

	_, err = newRouteWithConfig("test", NilValue, NilValue, 10, true, false)
	fmt.Printf("test: new(_,NilValue,NilValue) -> [err:%v]\n", err)

	//Output:
	//test: new("",_,_) -> [err:invalid argument : route name is empty]
	//test: new(_,100,0) -> [err:invalid argument : limit is configured but burst is not]
	//test: new(_,100,NilValue) -> [err:invalid argument : limit is configured but burst is not]
	//test: new(_,0,10) -> [err:invalid argument : burst is configured but limit is not]
	//test: new(_,NilValue,NilValue) -> [err:invalid argument : burst is configured but limit is not]
}

func Example_NewRouteWithConfig() {
	r, err := newRouteWithConfig("test", NilValue, 0, 0, true, false)
	fmt.Printf("test: new(0,0) -> [config:%v] [err:%v] [limiter:%v] [duration:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Duration(), r.Allow())

	r, err = newRouteWithConfig("test", NilValue, 100, 10, true, false)
	fmt.Printf("test: new(100,0) -> [config:%v] [err:%v] [limiter:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Allow())

	r, err = newRouteWithConfig("test", NilValue, rate.Inf, NilValue, true, false)
	fmt.Printf("test: new(Inf,NilValue) -> [config:%v] [err:%v] [limiter:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Allow())

	r, err = newRouteWithConfig("test", NilValue, rate.Inf, 0, true, false)
	fmt.Printf("test: new(Inf,0) -> [config:%v] [err:%v] [limiter:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Allow())

	r, err = newRouteWithConfig("test", NilValue, rate.Inf, 1234, true, false)
	fmt.Printf("test: new(Inf,1234) -> [config:%v] [err:%v] [limiter:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Allow())

	//Output:
	//test: new(0,0) -> [config:{-1 -1 -1}] [err:<nil>] [limiter:false] [duration:0s] [allow:true]
	//test: new(100,0) -> [config:{-1 100 10}] [err:<nil>] [limiter:true] [allow:true]
	//test: new(Inf,NilValue) -> [config:{-1 1.7976931348623157e+308 1}] [err:<nil>] [limiter:true] [allow:true]
	//test: new(Inf,0) -> [config:{-1 1.7976931348623157e+308 1}] [err:<nil>] [limiter:true] [allow:true]
	//test: new(Inf,1234) -> [config:{-1 1.7976931348623157e+308 1234}] [err:<nil>] [limiter:true] [allow:true]

}

/*
func Example_RateLimiter_DisallowAll() {
	r,err := newRouteWithConfig("test", NilValue, 0, 0, false, false)

	fmt.Printf("Route  : %v\n", r.current)
	fmt.Printf("Allow  : %v\n", r.Allow())

	//Output:
	//Route  : {-1 -1 -1}
	//Allow  : false
}

*/

func Example_RateLimiter_AllowAll() {
	r, err := newRouteWithConfig("test", NilValue, rate.Inf, 0, false, false)
	fmt.Printf("test: new(Inf,0) -> [config:%v] [err:%v] [limiter:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Allow())

	i := 0
	for ; i < 100; i++ {
		if !r.Allow() {
			fmt.Printf("Allow    : fail\n")
		}
	}

	//Output:
	//test: new(Inf,0) -> [config:{-1 1.7976931348623157e+308 1}] [err:<nil>] [limiter:true] [allow:true]

}

func Example_RateLimiter_AllowLimit() {
	r, err := newRouteWithConfig("test", NilValue, 1, 1, false, false)
	fmt.Printf("test: new(1,1) -> [config:%v] [err:%v] [limiter:%v] [allow:%v]\n", r.default_, err, r.IsRateLimiter(), r.Allow())

	allow := 0
	disallow := 0
	for i := 0; i < 10; i++ {
		if r.Allow() {
			allow++
		} else {
			disallow++
		}
		//fmt.Printf("Allow   : %v\n", r.Allow())
		if i == 1 || i == 4 || i == 7 {
			time.Sleep(time.Second)
		}
	}
	fmt.Printf("test: () [count:10] [allow:%v] [disallow:%v]\n", allow, disallow)
	//Output:
	//test: new(1,1) -> [config:{-1 1 1}] [err:<nil>] [limiter:true] [allow:true]
	//test: () [count:10] [allow:3] [disallow:7]

}
