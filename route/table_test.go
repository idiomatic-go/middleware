package route

import (
	"fmt"
	"net/http"
)

func ExampleTable_SetDefault() {
	t := newTable()

	fmt.Printf("test: Empty() -> [%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [default:%v]\n", r.(*route).name == DefaultName)
	//fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefault(newRoute("not-default"))
	r = t.Lookup(nil)
	fmt.Printf("test: Lookup(req) -> [default:%v]\n", r.(*route).name == DefaultName)

	//Output:
	//test: Empty() -> [true]
	//test: Lookup(nil) -> [default:true]
	//test: Lookup(req) -> [default:false]

}

func ExampleTable_Interact() {
	name := "test-route"
	t := newTable()
	fmt.Printf("test: Empty() -> [%v]\n", t.isEmpty())

	ok := t.Add(nil)
	fmt.Printf("test: Add(nil) -> [ok:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Add(newRoute(name))
	fmt.Printf("test: Add(route) -> [ok:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	t.Remove("")
	fmt.Printf("test: Remove(\"\") -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.Exists(name), t.LookupByName(name))

	t.Remove(name)
	fmt.Printf("test: Remove(name) -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.Exists(name), t.LookupByName(name))

	//Output:
	//test: Empty() -> [true]
	//test: Add(nil) -> [ok:false] [count:0] [exists:false] [lookup:<nil>]
	//test: Add(route) -> [ok:true] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//test: Remove("") -> [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//test: Remove(name) -> [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_Lookup() {
	name := "test-route"
	t := newTable()
	fmt.Printf("test: Empty() -> [%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [route:%v]\n", r)

	req, _ := http.NewRequest("", "http://localhost:8080/accesslog", nil)
	r = t.Lookup(req)
	fmt.Printf("test: Lookup(req) -> [route:%v]\n", r)

	ok := t.Add(newRoute(name))
	fmt.Printf("test: Add(route) -> [route:%v] [count:%v] [exists:%v]\n", ok, t.count(), t.Exists(name))

	t.SetMatcher(func(req *http.Request) string {
		return name
	},
	)
	r = t.Lookup(req)
	fmt.Printf("test: Lookup(req) ->  [route:%v]\n", r)

	//Output:
	//test: Empty() -> [true]
	//test: Lookup(nil) -> [route:&{* {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//test: Lookup(req) -> [route:&{* {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//test: Add(route) -> [route:true] [count:1] [exists:true]
	//test: Lookup(req) ->  [route:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]

}
