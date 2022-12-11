package route

import (
	"fmt"
	"net/http"
)

func ExampleTable_SetDefault() {
	t := newTable()

	fmt.Printf("result : empty() -> [%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("result : lookup(nil) -> [default:%v]\n", r.(*route).name == DefaultName)
	//fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefault(newRoute("not-default"))
	r = t.Lookup(nil)
	fmt.Printf("result : lookup(req) -> [default:%v]\n", r.(*route).name == DefaultName)

	//Output:
	//result : empty() -> [true]
	//result : lookup(nil) -> [default:true]
	//result : lookup(req) -> [default:false]

}

func ExampleTable_Interact() {
	name := "test-route"
	t := newTable()
	fmt.Printf("result : empty() -> [%v]\n", t.isEmpty())

	ok := t.Add(nil)
	fmt.Printf("result : add(nil) -> [%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Add(newRoute(name))
	fmt.Printf("result : add(route) -> [%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	t.Remove("")
	fmt.Printf("result : remove(\"\") -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.Exists(name), t.LookupByName(name))

	t.Remove(name)
	fmt.Printf("result : remove(name) -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.Exists(name), t.LookupByName(name))

	//Output:
	//result : empty() -> [true]
	//result : add(nil) -> [false] [count:0] [exists:false] [lookup:<nil>]
	//result : add(route) -> [true] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : remove("") -> [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : remove(name) -> [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_Lookup() {
	name := "test-route"
	t := newTable()
	fmt.Printf("result : [empty:%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("result : [lookup(nil):%v]\n", r)

	req, _ := http.NewRequest("", "http://localhost:8080/accesslog", nil)
	r = t.Lookup(req)
	fmt.Printf("result : [lookup(req):%v]\n", r)

	ok := t.Add(newRoute(name))
	fmt.Printf("result : [add(route):%v] [count:%v] [exists:%v]\n", ok, t.count(), t.Exists(name))

	t.SetMatcher(func(req *http.Request) string {
		return name
	},
	)
	r = t.Lookup(req)
	fmt.Printf("result : [lookup(req):%v]\n", r)

	//Output:
	//result : [empty:true]
	//result : [lookup(nil):&{* {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [lookup(req):&{* {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [add(route):true] [count:1] [exists:true]
	//result : [lookup(req):&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]

}
