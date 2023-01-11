package host

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

type entry struct {
	uri string
	c   chan Message
}

type entryDirectory struct {
	m  map[string]*entry
	mu sync.RWMutex
}

func newEntryDirectory() *entryDirectory {
	return &entryDirectory{m: make(map[string]*entry)}
}

func (e *entryDirectory) get(uri string) *entry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.m[uri]
}

func (e *entryDirectory) add(uri string, c chan Message) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.m[uri] = &entry{
		uri: uri,
		c:   c,
	}
}

func (e *entryDirectory) count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.m)
}

func (e *entryDirectory) uri() []string {
	var uri []string
	e.mu.RLock()
	defer e.mu.RUnlock()
	for key, _ := range e.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}

func (e *entryDirectory) send(msg Message) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if entry, ok := e.m[msg.To]; ok {
		if entry.c == nil {
			return errors.New(fmt.Sprintf("entry channel is nil: [%v]", msg.To))
		}
		entry.c <- msg
		return nil
	}
	return errors.New(fmt.Sprintf("entry not found: [%v]", msg.To))
}

func (e *entryDirectory) shutdown() {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, entry := range e.m {
		if entry.c != nil {
			entry.c <- Message{To: entry.uri, Event: ShutdownEvent}
		}
	}
}

func (e *entryDirectory) empty() {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for key, entry := range e.m {
		if entry.c != nil {
			close(entry.c)
		}
		delete(e.m, key)
	}
}

type entryResponse struct {
	m  map[string]Message
	mu sync.RWMutex
}

func newEntryResponse() *entryResponse {
	return &entryResponse{m: make(map[string]Message)}
}

func (e *entryResponse) count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.m)
}

func (e *entryResponse) exclude(event string, status int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var uri []string
	for u, entry := range e.m {
		if entry.Status != status || entry.Event != event {
			uri = append(uri, u)
		}
	}
	sort.Strings(uri)
	return uri
}

func (e *entryResponse) include(event string, status int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var uri []string
	for u, entry := range e.m {
		if entry.Status == status && entry.Event == event {
			uri = append(uri, u)
		}
	}
	sort.Strings(uri)
	return uri
}

func (e *entryResponse) add(msg Message) error {
	if msg.From == "" {
		return errors.New("invalid argument: message from is empty")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.m[msg.From] = msg
	return nil
}
