package actuator

import (
	"fmt"
	"net/http"
	"time"
)

func ExampleTable_SetDefaultActuator() {
	t := newTable(true)

	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	a := t.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [default:%v]\n", a.(*actuator).name == DefaultActuatorName)
	//fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefaultActuator("not-default", nil)
	a = t.Lookup(nil)
	fmt.Printf("test: Lookup(req) -> [default:%v]\n", a.(*actuator).name == DefaultActuatorName)

	//Output:
	//test: empty() -> [true]
	//test: Lookup(nil) -> [default:true]
	//test: Lookup(req) -> [default:false]

}

func ExampleTable_SetHostActuator() {
	t := newTable(true)

	a := t.Host()
	fmt.Printf("test: Host() -> [name:%v] [timeout-controller:%v]\n", a.Name(), a.t().timeout)

	t.SetHostActuator(nil, NewTimeoutConfig(time.Millisecond*1500, 504))
	a = t.Host()
	fmt.Printf("test: SetHostActuator(NewTimeoutConfig()) -> [name:%v] [timeout-controller:%v] \n", a.Name(), a.t().timeout != nil)

	//Output:
	//test: Host() -> [name:host] [timeout-controller:<nil>]
	//test: SetHostActuator(NewTimeoutConfig()) -> [name:host] [timeout-controller:true]

}

func ExampleTable_Add_Exists_LookupByName() {
	name := "test-route"
	t := newTable(true)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	err := t.Add("", "/table", nil, nil, nil, nil)
	fmt.Printf("test: Add(nil) -> [err:%v] [count:%v] [exists:%v] [lookup:%v]\n", err, t.count(), t.exists(name), t.LookupByName(name))

	err = t.Add(name, "/table", nil, nil, nil, nil)
	fmt.Printf("test: Add(actuator) -> [err:%v] [count:%v] [exists:%v] [lookup:%v]\n", err, t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove("")
	fmt.Printf("test: remove(\"\") -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove(name)
	fmt.Printf("test: remove(name) -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name))

	//Output:
	//test: empty() -> [true]
	//test: Add(nil) -> [err:[invalid argument: name is empty]] [count:0] [exists:false] [lookup:<nil>]
	//test: Add(actuator) -> [err:[]] [count:1] [exists:true] [lookup:true]
	//test: remove("") -> [count:1] [exists:true] [lookup:true]
	//test: remove(name) -> [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_Lookup() {
	name := "test-route"
	t := newTable(true)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [actuator:%v]\n", r.Name())

	req, _ := http.NewRequest("", "http://localhost:8080/accesslog", nil)
	r = t.Lookup(req)
	fmt.Printf("test: Lookup(req) -> [actuator:%v]\n", r.Name())

	ok := t.Add(name, "http://localhost:8080/accesslog", nil, NewTimeoutConfig(100, 503), nil, nil, nil)
	fmt.Printf("test: Add(actuator) -> [actuator:%v] [count:%v] [exists:%v]\n", ok, t.count(), t.exists(name))

	r = t.Lookup(req)
	fmt.Printf("test: Lookup(req) ->  [actuator:%v]\n", r.Name())

	//Output:
	//test: empty() -> [true]
	//test: Lookup(nil) -> [actuator:*]
	//test: Lookup(req) -> [actuator:*]
	//test: Add(actuator) -> [actuator:[]] [count:1] [exists:true]
	//test: Lookup(req) ->  [actuator:test-route]

}
