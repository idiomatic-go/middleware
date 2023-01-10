package resource

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"reflect"
	"time"
)

type messageMap map[string]Message

type pkg struct{}

var pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()
var startupLocation = pkgPath + "/startup"

var directory = newEntryDirectory()

// RegisterResource - function to register a resource uri
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

func Startup[E template.ErrorHandler](duration time.Duration, content ContentMap) (status *template.Status) {
	var e E
	var failures []string

	if directory.count() == 0 {
		return nil
	}
	resp := newEntryResponse()
	toSend := createToSend(content, func(msg Message) {
		resp.add(msg)
	})
	sendMessages(toSend)
	for wait := time.Duration(float64(duration) * 0.6); duration > 0; duration -= wait {
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

func createToSend(cm ContentMap, fn MessageHandler) messageMap {
	m := make(messageMap)
	for _, k := range directory.uri() {
		msg := Message{To: k, From: VirtualHost, Event: StartupEvent, Status: StatusNotProvided, ReplyTo: fn}
		if cm != nil {
			if content, ok := cm[k]; ok {
				msg.Content = append(msg.Content, content...)
			}
		}
		m[k] = msg
	}
	return m
}

func sendMessages(msgs messageMap) {
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
