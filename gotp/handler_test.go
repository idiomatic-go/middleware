package gotp

import (
	"errors"
	"fmt"
)

func ExampleNoOpHandler_Handle() {
	location := "/test"
	err := errors.New("test error")
	var h NoOpHandler

	fmt.Printf("test: Handle(location,nil) -> [%v]\n", h.Handle(location, nil))
	fmt.Printf("test: Handle(location,err) -> [%v]\n", h.Handle(location, err))

	s := NewStatus(StatusInternal, location, nil)
	fmt.Printf("test: HandleStatus(s) -> [%v]\n", h.HandleStatus(s))

	s = NewStatus(StatusInternal, location, err)
	fmt.Printf("test: HandleStatus(s) -> [prev:%v] [curr:%v]\n", s, h.HandleStatus(s))

	//Output:
	//test: Handle(location,nil) -> [0 Successful]
	//test: Handle(location,err) -> [13 Internal Error [test error]]
	//test: HandleStatus(s) -> [0 Successful]
	//test: HandleStatus(s) -> [prev:13 Internal Error [test error]] [curr:13 Internal Error [test error]]

}

func ExampleDebugHandler_Handle() {
	location := "/test"
	err := errors.New("test error")
	var h DebugHandler

	s := h.Handle(location, nil)
	fmt.Printf("test: Handle(location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())
	s = h.Handle(location, err)
	fmt.Printf("test: Handle(location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatus(StatusInternal, location, nil)
	fmt.Printf("test: HandleStatus(s) -> [%v] [errors:%v]\n", h.HandleStatus(s), s.IsErrors())

	s = NewStatus(StatusInternal, location, err)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s)
	fmt.Printf("test: HandleStatus(s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(location,nil) -> [0 Successful] [errors:false]
	//[/test [test error]]
	//test: Handle(location,err) -> [13 Internal Error] [errors:false]
	//test: HandleStatus(s) -> [0 Successful] [errors:false]
	//[/test [test error]]
	//test: HandleStatus(s) -> [prev:13 Internal Error] [prev-errors:true] [curr:13 Internal Error] [curr-errors:false]

}

func ExampleLogHandler_Handle() {
	location := "/test"
	err := errors.New("test error")
	var h LogHandler

	s := h.Handle(location, nil)
	fmt.Printf("test: Handle(location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())
	s = h.Handle(location, err)
	fmt.Printf("test: Handle(location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatus(StatusInternal, location, nil)
	fmt.Printf("test: HandleStatus(s) -> [%v] [errors:%v]\n", h.HandleStatus(s), s.IsErrors())

	s = NewStatus(StatusInternal, location, err)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s)
	fmt.Printf("test: HandleStatus(s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(location,nil) -> [0 Successful] [errors:false]
	//test: Handle(location,err) -> [13 Internal Error] [errors:false]
	//test: HandleStatus(s) -> [0 Successful] [errors:false]
	//test: HandleStatus(s) -> [prev:13 Internal Error] [prev-errors:true] [curr:13 Internal Error] [curr-errors:false]

}
