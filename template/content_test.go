package template

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

func ExampleProcessContextContentStatus() {
	status := NewStatusOk()
	ctx := ContextWithContent(context.Background(), status)
	t, s := ProcessContextContent[address](ctx)
	if t.Cell != "" {
	}
	fmt.Printf("Status : %v\n", s.Ok())

	//Output:
	//Status : true
}

func ExampleProcessContextContentError() {
	err := errors.New("this is a test error")
	ctx := ContextWithContent(context.Background(), err)
	_, s := ProcessContextContent[address](ctx)
	fmt.Printf("Error : %v\n", s.Errors())

	//Output:
	//Error : [this is a test error]
}

func ExampleProcessContextContentType() {
	addr := address{Name: "Mark", Email: "mark@gmail.com", Cell: "123-456-7891"}
	ctx := ContextWithContent(context.Background(), addr)
	t, s := ProcessContextContent[address](ctx)
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
	l, s := ProcessContent[location](ContextContent(ctx))
	fmt.Printf("Address : %v\n", l.Location())
	fmt.Printf("Status  : %v\n", s.Ok())

	//Output:
	//Address : Mark
	//Status  : true
}

func ExampleProcessContentErrors() {
	addr := address2{Name: "Mark", Email: "mark@gmail.com", Cell: "123-456-7891", Zip: "50436"}

	ContextWithContent(context.Background(), addr)
	l, s := ProcessContent[address](nil)
	fmt.Printf("Address : %v\n", l)
	fmt.Printf("Ok      : %v\n", s.Ok())

	ctx := ContextWithContent(context.Background(), addr)
	l, s = ProcessContent[address](ctx)
	fmt.Printf("Address : %v\n", l)
	fmt.Printf("Ok      : %v\n", s.Ok())

	//Output:
	//Address : {  }
	//Ok      : false
	//Address : {  }
	//Ok      : false
}
