package template

import (
	"fmt"
	"log"
)

type ErrorHandler interface {
	Handle(location string, errs ...error) *Status
	HandleStatus(s *Status) *Status
}

type NoOpHandler struct{}

func (NoOpHandler) Handle(location string, errs ...error) *Status {
	if len(errs) == 0 || (len(errs) == 1 && errs[0] == nil) {
		return NewStatusOk()
	}
	return NewStatusError(location, errs...)
}

func (NoOpHandler) HandleStatus(s *Status) *Status {
	return s
}

type DebugHandler struct{}

func (h DebugHandler) Handle(location string, errs ...error) *Status {
	if len(errs) == 0 || (len(errs) == 1 && errs[0] == nil) {
		return NewStatusOk()
	}
	return h.HandleStatus(NewStatus(StatusInternal, location, errs...))
}

func (h DebugHandler) HandleStatus(s *Status) *Status {
	if s != nil && s.IsErrors() {
		if len(s.errs) == 0 || (len(s.errs) == 1 && s.errs[0] == nil) {
		} else {
			if s.location == "" {
				s.location = "[]"
			}
			fmt.Printf("[%v %v]\n", s.location, s.errs)
			s.errs = nil
		}
	}
	return s
}

type LogHandler struct{}

func (h LogHandler) Handle(location string, errs ...error) *Status {
	if len(errs) == 0 || (len(errs) == 1 && errs[0] == nil) {
		return NewStatusOk()
	}
	return h.HandleStatus(NewStatus(StatusInternal, location, errs...))
}

func (h LogHandler) HandleStatus(s *Status) *Status {
	if s != nil && s.IsErrors() {
		if len(s.errs) == 0 || (len(s.errs) == 1 && s.errs[0] == nil) {
		} else {
			if s.location == "" {
				s.location = "[]"
			}
			log.Println(s.location, s.errs)
			s.errs = nil
		}
	}
	return s
}
