package route

import "fmt"

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
	fmt.Printf("Result : [empty:%v]\n", t.isEmpty())

	ok := t.Add(nil)
	fmt.Printf("Result : [add:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Add(NewRoute(name))
	fmt.Printf("Result : [add:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Remove("")
	fmt.Printf("Result : [remove:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	ok = t.Remove(name)
	fmt.Printf("Result : [remove:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.Exists(name), t.LookupByName(name))

	//Output:
	//Result : [empty:true]
	//Result : [add:false] [count:0] [exists:false] [lookup:<nil>]
	//Result : [add:true] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//Result : [remove:false] [count:1] [exists:true] [lookup:&{test-route {-1 -1 -1} {-1 -1 -1} false false <nil>}]
	//Result : [remove:true] [count:0] [exists:false] [lookup:<nil>]

}
