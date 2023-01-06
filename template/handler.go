package template

import (
	"fmt"
	"log"
)

type ErrorHandler interface {
	New(location string, errs ...error) *Status
}

var PassThrough ErrorHandler
var Log ErrorHandler
var Debug ErrorHandler

func init() {
	PassThrough = newPassThrough()
	Log = newLogger()
	Debug = newLogger()
}

type passThrough struct{}

func newPassThrough() ErrorHandler {
	return new(passThrough)
}

func (p *passThrough) New(location string, errs ...error) *Status {
	return NewStatusError(location, errs...)
}

type logger struct{}

func newLogger() ErrorHandler {
	return new(logger)
}

func (p *logger) New(location string, errs ...error) *Status {
	log.Println(errs)
	return NewStatus(StatusInternal, location, nil)
}

type debugger struct{}

func newDebugger() ErrorHandler {
	return new(debugger)
}

func (p *debugger) New(location string, errs ...error) *Status {
	fmt.Printf("%v\n", errs)
	return NewStatus(StatusInternal, location, nil)
}
