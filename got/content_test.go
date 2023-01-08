package got

import (
	"context"
	"errors"
	"fmt"
)

type location interface {
	Location() string
}

type address struct {
	Name  string
	Email string
	Cell  string
}

func (a *address) Location() string {
	return a.Name
}

type address2 struct {
	Name  string
	Email string
	Cell  string
	Zip   string
}

func ExampleProcessContext_ErrorHandlers() {
	err := errors.New("this is a test error")
	ctx := ContextWithContent(context.Background(), err)
	if ctx != nil {
	}

	_, s := ProcessContextContent[address, NoOpHandler](nil)
	fmt.Printf("test: ProcessContextContent() -> %v\n", s.Errors())

	_, s = ProcessContextContent[address, DebugHandler](nil)
	fmt.Printf("test: ProcessContextContent() -> %v\n", s.Errors())

	_, s = ProcessContextContent[address, LogHandler](nil)
	fmt.Printf("test: ProcessContextContent() -> %v\n", s.Errors())

	//Output:
	//test: ProcessContextContent() -> [invalid configuration: context is nil]
	//[github.com/idiomatic-go/middleware/got/ProcessContextContent [invalid configuration: context is nil]]
	//test: ProcessContextContent() -> []
	//test: ProcessContextContent() -> []

}

func ExampleProcessContextContentStatus() {
	status := NewStatusOk()
	ctx := ContextWithContent(context.Background(), status)
	t, s := ProcessContextContent[address, NoOpHandler](ctx)
	if t.Cell != "" {
	}
	fmt.Printf("Status : %v\n", s.Ok())

	//Output:
	//Status : true
}

func ExampleProcessContextContentError() {
	err := errors.New("this is a test error")
	ctx := ContextWithContent(context.Background(), err)
	_, s := ProcessContextContent[address, DebugHandler](ctx)
	fmt.Printf("Error : %v\n", s.Errors())

	//Output:
	//Error : [this is a test error]
}

func ExampleProcessContextContentType() {
	addr := address{Name: "Mark", Email: "mark@gmail.com", Cell: "123-456-7891"}
	ctx := ContextWithContent(context.Background(), addr)
	t, s := ProcessContextContent[address, NoOpHandler](ctx)
	if t.Cell != "" {
	}
	fmt.Printf("Address : %v\n", t)
	fmt.Printf("Status  : %v\n", s.Ok())

	//Output:
	//Address : {Mark mark@gmail.com 123-456-7891}
	//Status  : true
}

func ExampleProcessContentInterface() {
	addr := address{Name: "Mark", Email: "mark@gmail.com", Cell: "123-456-7891"}
	var loc location = &addr

	ctx := ContextWithContent(context.Background(), loc)
	l, s := ProcessContent[location, NoOpHandler](ContextContent(ctx))
	fmt.Printf("Address : %v\n", l.Location())
	fmt.Printf("Status  : %v\n", s.Ok())

	//Output:
	//Address : Mark
	//Status  : true
}

func ExampleProcessContentErrors() {
	addr := address2{Name: "Mark", Email: "mark@gmail.com", Cell: "123-456-7891", Zip: "50436"}

	ContextWithContent(context.Background(), addr)
	l, s := ProcessContent[address, NoOpHandler](nil)
	fmt.Printf("Address : %v\n", l)
	fmt.Printf("Ok      : %v\n", s.Ok())

	ctx := ContextWithContent(context.Background(), addr)
	l, s = ProcessContent[address, NoOpHandler](ctx)
	fmt.Printf("Address : %v\n", l)
	fmt.Printf("Ok      : %v\n", s.Ok())

	//Output:
	//Address : {  }
	//Ok      : false
	//Address : {  }
	//Ok      : false
}
