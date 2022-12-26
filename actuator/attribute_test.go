package actuator

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func Example_String() {
	a := NewAttribute("nil", nil)
	fmt.Printf("test: NewAttribute(nil) -> [name:%v] [string:%v]\n", a.Name(), a.String())

	a = NewAttribute("bool", true)
	fmt.Printf("test: NewAttribute(true) -> [name:%v] [string:%v]\n", a.Name(), a.String())

	a = NewAttribute("integer", 1234)
	fmt.Printf("test: NewAttribute(1234) -> [name:%v] [string:%v]\n", a.Name(), a.String())

	a = NewAttribute("duration", time.Duration(time.Millisecond*4500))
	fmt.Printf("test: NewAttribute(time.Duration(time.Millisecond*4500)) -> [name:%v] [string:%v]\n", a.Name(), a.String())

	var limit rate.Limit = 100
	a = NewAttribute("limit", limit)
	fmt.Printf("test: NewAttribute(100) -> [name:%v] [string:%v]\n", a.Name(), a.String())

	limit = rate.Inf
	a = NewAttribute("limit", limit)
	fmt.Printf("test: NewAttribute(rate.Inf) -> [name:%v] [string:%v]\n", a.Name(), a.String())

	a = NewAttribute("string", "test of a string attribute")
	fmt.Printf("test: NewAttribute(\"test of a string attribute\") -> [name:%v] [string:%v]\n", a.Name(), a.String())

	//Output:
	//test: NewAttribute(nil) -> [name:nil] [string:nil]
	//test: NewAttribute(true) -> [name:bool] [string:true]
	//test: NewAttribute(1234) -> [name:integer] [string:1234]
	//test: NewAttribute(time.Duration(time.Millisecond*4500)) -> [name:duration] [string:4500]
	//test: NewAttribute(100) -> [name:limit] [string:100]
	//test: NewAttribute(rate.Inf) -> [name:limit] [string:-1]
	//test: NewAttribute("test of a string attribute") -> [name:string] [string:test of a string attribute]

}

func Example_Tag() {
	a := NewAttribute("nil", nil)
	fmt.Printf("test: NewAttribute(nil) -> [%v]\n", a.Tag())

	a = NewAttribute("bool", true)
	fmt.Printf("test: NewAttribute(true) -> [%v]\n", a.Tag())

	a = NewAttribute("integer", 1234)
	fmt.Printf("test: NewAttribute(1234) -> [%v]\n", a.Tag())

	a = NewAttribute("duration", time.Duration(time.Millisecond*4500))
	fmt.Printf("test: NewAttribute(time.Duration(time.Millisecond*4500)) -> [%v]\n", a.Tag())

	var limit rate.Limit = 100
	a = NewAttribute("limit", limit)
	fmt.Printf("test: NewAttribute(100) -> [%v]\n", a.Tag())

	limit = rate.Inf
	a = NewAttribute("limit", limit)
	fmt.Printf("test: NewAttribute(rate.Inf) -> [%v]\n", a.Tag())

	a = NewAttribute("string", "test of a string attribute")
	fmt.Printf("test: NewAttribute(\"test of a string attribute\") -> [%v]\n", a.Tag())

	//Output:
	//test: NewAttribute(nil) -> [nil:nil]
	//test: NewAttribute(true) -> [bool:true]
	//test: NewAttribute(1234) -> [integer:1234]
	//test: NewAttribute(time.Duration(time.Millisecond*4500)) -> [duration:4500]
	//test: NewAttribute(100) -> [limit:100]
	//test: NewAttribute(rate.Inf) -> [limit:-1]
	//test: NewAttribute("test of a string attribute") -> [string:test of a string attribute]

}
