package automation

import "fmt"

// TODO : test nil attribute

func Example_newTimeout() {
	t := newTimeout("test-route", nil, newTable())
	fmt.Printf("test: newTimeout() -> [enabled:%v] [name:%v] [default:%v] [current:%v]\n", t.enabled, t.name, t.defaultC, t.current)

	t = newTimeout("test-route2", NewTimeoutConfig(2000, 504), newTable())
	fmt.Printf("test: newTimeout() -> [enabled:%v] [name:%v] [default:%v] [current:%v]\n", t.enabled, t.name, t.defaultC, t.current)

	t2 := cloneTimeout(t)
	t2.enabled = false
	fmt.Printf("test: cloneTimeout() -> [prev-enabled:%v] [prev-name:%v] [curr-enabled:%v] [curr-name:%v]\n", t.enabled, t.name, t2.enabled, t2.name)

	//Output:
	//test: newTimeout() -> [enabled:false] [name:test-route] [default:{-1 -1}] [current:{-1 -1}]
	//test: newTimeout() -> [enabled:true] [name:test-route2] [default:{2000 504}] [current:{2000 504}]
	//test: cloneTimeout() -> [prev-enabled:true] [prev-name:test-route2] [curr-enabled:false] [curr-name:test-route2]
}

func Example_Controller_ReadOnly() {
	t := newTimeout("test-route", NewTimeoutConfig(2000, 504), newTable())
	fmt.Printf("test: IsEnabled() -> [%v]\n", t.IsEnabled())

	fmt.Printf("test: Duration() -> [%v]\n", t.Duration())

	fmt.Printf("test: StatusCode(200) -> [%v]\n", t.StatusCode(200))

	t = newTimeout("test-route", NewTimeoutConfig(2000, NilValue), newTable())
	fmt.Printf("test: StatusCode(503) -> [%v]\n", t.StatusCode(503))

	a := t.Attribute("")
	fmt.Printf("test: Attribute(\"\") -> [name:%v] [value:%v] [string:%v]\n", a.Name(), a.Value(), a)

	a = t.Attribute(TimeoutName)
	fmt.Printf("test: Attribute(\"Timeout\") -> [name:%v] [value:%v] [string:%v]\n", a.Name(), a.Value(), a)

	//Output:
	//test: IsEnabled() -> [true]
	//test: Duration() -> [2s]
	//test: StatusCode(200) -> [504]
	//test: StatusCode(503) -> [503]
	//test: Attribute("") -> [name:] [value:<nil>] [string:nil]
	//test: Attribute("Timeout") -> [name:timeout] [value:2000] [string:2000]
}
