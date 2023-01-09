package vhost

import (
	"errors"
	"fmt"
	"sync"
)

type entry struct {
	uri string
	c   chan Message
}

type entryMap struct {
	m  map[string]*entry
	mu sync.RWMutex
}

//var directory = entryMap{m: make(map[string]*entry)}

//func init() {
//	directory
//}

func (e entryMap) add(uri string, c chan Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.m[uri] = &entry{
		uri: uri,
		c:   c,
	}
}

func (e entryMap) count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.m)
}

func (e entryMap) send(msg Message) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	entry := e.m[msg.To]
	if entry == nil {
		return errors.New(fmt.Sprintf("entry not found: [%v]", msg.To))
	}
	entry.c <- msg
	return nil
}

func (e entryMap) uri() []string {
	var uri []string
	e.mu.RLock()
	defer e.mu.RUnlock()
	for key, _ := range e.m {
		uri = append(uri, key)
	}
	return uri
}

func (e entryMap) shutdown() {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, entry := range e.m {
		entry.c <- Message{To: entry.uri, Event: ShutdownEvent}
	}
}

type entryResponse struct {
	m  map[string]Message
	mu sync.RWMutex
}

func newEntryResponse() entryResponse {
	return entryResponse{m: make(map[string]Message)}
}

func (e entryResponse) count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.m)
}

func (e entryResponse) compare(event string, status int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var uri []string
	for u, entry := range e.m {
		if entry.Status != status && entry.Event != event {
			uri = append(uri, u)
		}
	}
	return uri
}

func (e entryResponse) add(msg Message) error {
	if msg.From == "" {
		return errors.New("invalid argument: message from is empty")
	}
	//if msg.Event ==
	e.mu.Lock()
	defer e.mu.Unlock()
	e.m[msg.From] = msg
	return nil
}
