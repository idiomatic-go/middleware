package template

import (
	"fmt"
	"log"
)

type ErrorHandler interface {
	Handle(location string, errs ...error) *Status
}

/*
type ErrorHandlerFn func(location string, errs ...error) *Status

func (ErrorHandlerFn)Handle(string,...error) *Status {
    return nil
}


*/
type NoOpHandler struct{}

func (NoOpHandler) Handle(location string, errs ...error) *Status {
	return NewStatusError(location, errs...)
}

type DebugHandler struct{}

func (DebugHandler) Handle(location string, errs ...error) *Status {
	if location == "" {
		location = "[]"
	}
	fmt.Printf("[%v %v]\n", location, errs)
	return NewStatus(StatusInternal, location, nil)
}

type LogHandler struct{}

func (LogHandler) Handle(location string, errs ...error) *Status {
	if location == "" {
		location = "[]"
	}
	log.Println(location, errs)
	return NewStatus(StatusInternal, location, nil)
}
