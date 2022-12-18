package automation

import (
	"fmt"
	"net/http"
)

func ExampleTable_SetDefault() {
	t := newTable()

	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	a := t.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [default:%v]\n", a.(*actuator).name == DefaultName)
	//fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefault("not-default", nil, nil, nil, nil)
	a = t.Lookup(nil)
	fmt.Printf("test: Lookup(req) -> [default:%v]\n", a.(*actuator).name == DefaultName)

	//Output:
	//test: empty() -> [true]
	//test: Lookup(nil) -> [default:true]
	//test: Lookup(req) -> [default:false]

}

func ExampleTable_Add_Exists_LookupByName() {
	name := "test-route"
	t := newTable()
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	ok := t.Add("", nil, nil, nil, nil)
	fmt.Printf("test: Add(nil) -> [ok:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.exists(name), t.LookupByName(name))

	ok = t.Add(name, nil, nil, nil, nil)
	fmt.Printf("test: Add(actuator) -> [ok:%v] [count:%v] [exists:%v] [lookup:%v]\n", ok, t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove("")
	fmt.Printf("test: remove(\"\") -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove(name)
	fmt.Printf("test: remove(name) -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name))

	//Output:
	//test: empty() -> [true]
	//test: Add(nil) -> [ok:false] [count:0] [exists:false] [lookup:<nil>]
	//test: Add(actuator) -> [ok:true] [count:1] [exists:true] [lookup:true]
	//test: remove("") -> [count:1] [exists:true] [lookup:true]
	//test: remove(name) -> [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_Lookup() {
	name := "test-route"
	t := newTable()
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [actuator:%v]\n", r.Name())

	req, _ := http.NewRequest("", "http://localhost:8080/accesslog", nil)
	r = t.Lookup(req)
	fmt.Printf("test: Lookup(req) -> [actuator:%v]\n", r.Name())

	ok := t.Add(name, nil, NewTimeoutConfig(100, NilValue), nil, nil)
	fmt.Printf("test: Add(actuator) -> [actuator:%v] [count:%v] [exists:%v]\n", ok, t.count(), t.exists(name))

	t.SetMatcher(func(req *http.Request) string {
		return name
	},
	)
	r = t.Lookup(req)
	fmt.Printf("test: Lookup(req) ->  [actuator:%v]\n", r.Name())

	//Output:
	//test: empty() -> [true]
	//test: Lookup(nil) -> [actuator:*]
	//test: Lookup(req) -> [actuator:*]
	//test: Add(actuator) -> [actuator:true] [count:1] [exists:true]
	//test: Lookup(req) ->  [actuator:test-route]

}
