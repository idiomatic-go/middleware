package route

import (
	"fmt"
	"net/http"
)

func ExampleTable_SetDefault() {
	t := newTable()

	fmt.Printf("IsEmpty   : %v\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("Route     : %v\n", r.(*route).default_)
	fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefault(NewRoute("not-default"))
	r = t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	//Output:
	//IsEmpty   : true
	//Route     : {-1 -1 -1}
	//IsDefault : true
	//IsDefault : false
}

func ExampleTable_Interact() {
	name := "test-route"
	t := newTable()
	fmt.Printf("result : [empty:%v]\n", t.isEmpty())

	ok := t.Add(nil)
	fmt.Printf("result : [add:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Add(NewRoute(name))
	fmt.Printf("result : [add:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Remove("")
	fmt.Printf("result : [remove:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Remove(name)
	fmt.Printf("result : [remove:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	//Output:
	//result : [empty:true]
	//result : [add:false] [count:0] [exists:false] [lookup:<nil>]
	//result : [add:true] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [remove:false] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [remove:true] [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_Lookup() {
	name := "test-route"
	t := newTable()
	fmt.Printf("result : [empty:%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("result : [lookup:%v]\n", r)

	req, _ := http.NewRequest("", "http://localhost:8080/accesslog", nil)
	r = t.Lookup(req)
	fmt.Printf("result : [lookup:%v]\n", r)

	ok := t.Add(NewRoute(name))
	fmt.Printf("result : [add:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	t.SetMatcher(func(req *http.Request) string {
		return name
	},
	)
	r = t.Lookup(req)
	fmt.Printf("result : [lookup:%v]\n", r)

	//Output:
	//result : [empty:true]
	//result : [lookup:&{* {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [lookup:&{* {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [add:true] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//result : [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]

}
