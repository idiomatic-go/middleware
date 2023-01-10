package resource

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"reflect"
	"time"
)

type pkg struct{}

var pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()
var startupLocation = pkgPath + "/startup"

var directory = newEntryMap()

// RegisterResource - function to register a package uri
func RegisterResource(uri string, c chan Message) error {
	if uri == "" {
		return errors.New("invalid argument: uri is empty")
	}
	if c == nil {
		return errors.New(fmt.Sprintf("invalid argument: channel is nil for [%v]", uri))
	}
	registerResourceUnchecked(uri, c)
	return nil
}

func registerResourceUnchecked(uri string, c chan Message) error {
	directory.add(uri, c)
	return nil
}

// Shutdown - virtual host shutdown
func Shutdown() {
	directory.shutdown()
}

func Startup[E template.ErrorHandler](duration time.Duration, init MessageMap) (status *template.Status) {
	var e E
	var failures []string

	if directory.count() == 0 {
		return nil
	}
	resp := newEntryResponse()
	toSend := createToSend(init, func(msg Message) {
		resp.add(msg)
	})
	sendMessages(toSend)
	wait := time.Duration(float64(duration) * 0.6)
	for ; duration > 0; duration -= wait {
		time.Sleep(wait)
		// Check for completion
		if resp.count() < directory.count() {
			continue
		}
		// Check for failed resources
		failures = resp.compare(StartupEvent, 0)
		if len(failures) == 0 {
			return template.NewStatusOk()
		}
		break
	}
	Shutdown()
	if len(failures) > 0 {
		return e.Handle(startupLocation, errors.New(fmt.Sprintf("status failures %v", failures))).SetCode(template.StatusInternal)
	}
	return e.Handle(startupLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", resp.count(), directory.count()))).SetCode(template.StatusDeadlineExceeded)
}

func createToSend(msgs MessageMap, fn MessageHandler) MessageMap {
	m := make(MessageMap)
	for _, k := range directory.uri() {
		if msgs != nil {
			message, ok := msgs[k]
			if ok {
				message.Event = StartupEvent
				message.From = VirtualHost
				message.Status = StatusNotProvided
				message.ReplyTo = fn
				m[k] = message
				continue
			}
		}
		m[k] = Message{To: k, From: VirtualHost, Event: StartupEvent, Status: StatusNotProvided, ReplyTo: fn}
	}
	return m
}

func sendMessages(msgs MessageMap) {
	for k := range msgs {
		directory.send(msgs[k])
	}
}

/*
func Startup(ticks, iterations int, init MessageMap) error {
	if directory.count() == 0 {
		return nil
	}
	resp := newEntryResponse()
	count := 1
	toSend := createToSend(init, func(msg Message) {
		resp.add(msg)
	})
	sendMessages(toSend)
	for {
		if count > iterations {
			Shutdown()
			return errors.New(fmt.Sprintf("startup failure: response counts < directory entries [%v] [%v]", resp.count(), directory.count()))
		}
		time.Sleep(time.Second * time.Duration(ticks))
		count++
		if resp.count() < directory.count() {
			continue
		}
		failures := resp.compare(StartupEvent, 0)
		if len(failures) == 0 {
			return nil
		}
		Shutdown()
		return errors.New(fmt.Sprintf("startup failure: status failures %v", failures))
	}
	return nil
}

*/
