package vhost

import (
	"errors"
	"fmt"
	"time"
)

var directory = entryMap{m: make(map[string]*entry)}

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

func Startup(ticks, iterations int, init MessageMap) error {
	if directory.count() == 0 {
		return nil
	}
	resp := newEntryResponse()
	count := 0
	toSend := createToSend(init, nil)
	sendMessages(toSend)
	for {
		if count > iterations {
			Shutdown()
			return errors.New(fmt.Sprintf("vhost startup failure %v, max iterations exceeded: %v", "", count))
		}
		time.Sleep(time.Second * time.Duration(ticks))
		count++
		failures := resp.compare(StartupEvent, 0)
		if len(failures) == 0 && resp.count() == directory.count() {
			return nil
		}
	}
	return nil
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
		//eventing.Directory.Add(k, eventing.CreateMessage(eventing.VirtualHost, eventing.VirtualHost, eventing.StartupEvent, eventing.StatusInProgress, nil))
	}
}
