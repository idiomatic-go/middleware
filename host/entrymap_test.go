package host

import (
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"time"
)

/*
func createTestEntry(uri string, status int32) *entry {
	entry := createEntry(uri, nil)
	entry.msgs.add(CreateMessage(VirtualHost, VirtualHost, StartupEvent, status, nil))
	return entry
}
*/

func ExampleDirectory_Add() {
	uri := "urn:test"
	uri2 := "urn:test:two"

	directory.empty()

	fmt.Printf("test: count() -> : %v\n", directory.count())
	d2 := directory.get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	directory.add(uri, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri)
	fmt.Printf("test: count() -> : %v\n", directory.count())
	d2 = directory.get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	directory.add(uri2, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri2)
	fmt.Printf("test: count() -> : %v\n", directory.count())
	d2 = directory.get(uri2)
	fmt.Printf("test: get(%v) -> : %v\n", uri2, d2)

	fmt.Printf("test: uri() -> : %v\n", directory.uri())

	//Output:
	//test: count() -> : 0
	//test: get(urn:test) -> : <nil>
	//test: add(urn:test) -> : ok
	//test: count() -> : 1
	//test: get(urn:test) -> : &{urn:test <nil>}
	//test: add(urn:test:two) -> : ok
	//test: count() -> : 2
	//test: get(urn:test:two) -> : &{urn:test:two <nil>}
	//test: uri() -> : [urn:test urn:test:two]

}

func ExampleDirectory_SendError() {
	uri := "urn:test"
	directory.empty()

	fmt.Printf("test: send(%v) -> : %v\n", uri, directory.send(Message{To: uri}))

	directory.add(uri, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri)
	fmt.Printf("test: send(%v) -> : %v\n", uri, directory.send(Message{To: uri}))

	//Output:
	//test: send(urn:test) -> : entry not found: [urn:test]
	//test: add(urn:test) -> : ok
	//test: send(urn:test) -> : entry channel is nil: [urn:test]

}

func ExampleDirectory_Send() {
	uri1 := "urn:test-1"
	uri2 := "urn:test-2"
	uri3 := "urn:test-3"
	c := make(chan Message, 16)
	directory.empty()

	directory.add(uri1, c)
	directory.add(uri2, c)
	directory.add(uri3, c)

	directory.send(Message{To: uri1, From: Name, Event: StartupEvent})
	directory.send(Message{To: uri2, From: Name, Event: StartupEvent})
	directory.send(Message{To: uri3, From: Name, Event: StartupEvent})

	time.Sleep(time.Second * 1)
	resp1 := <-c
	resp2 := <-c
	resp3 := <-c
	fmt.Printf("test: <- c -> : [%v] [%v] [%v]\n", resp1.To, resp2.To, resp3.To)
	close(c)
	//Output:
	//test: <- c -> : [urn:test-1] [urn:test-2] [urn:test-3]

}

func _ExampleResponse_Add() {
	resp := newEntryResponse()

	resp.add(Message{To: "to-uri", From: "from-uri-0", Event: StartupEvent, Status: template.StatusNotProvided})
	resp.add(Message{To: "to-uri", From: "from-uri-1", Event: StartupEvent, Status: 100})
	resp.add(Message{To: "to-uri", From: "from-uri-2", Event: PingEvent, Status: template.StatusNotProvided})
	resp.add(Message{To: "to-uri", From: "from-uri-3", Event: PingEvent, Status: template.StatusNotProvided})
	resp.add(Message{To: "to-uri", From: "from-uri-4", Event: PingEvent, Status: 200})

	fmt.Printf("test: count() -> : %v\n", resp.count())

	fmt.Printf("test: include(%v,%v) -> : %v\n", ShutdownEvent, template.StatusNotProvided, resp.include(ShutdownEvent, template.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", ShutdownEvent, template.StatusNotProvided, resp.exclude(ShutdownEvent, template.StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", StartupEvent, template.StatusNotProvided, resp.include(StartupEvent, template.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", StartupEvent, template.StatusNotProvided, resp.exclude(StartupEvent, template.StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", PingEvent, template.StatusNotProvided, resp.include(PingEvent, template.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", PingEvent, template.StatusNotProvided, resp.exclude(PingEvent, template.StatusNotProvided))

	//Output:
	//test: count() -> : 5
	//test: include(event:shutdown,-101) -> : []
	//test: exclude(event:shutdown,-101) -> : [from-uri-0 from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:startup,-101) -> : [from-uri-0]
	//test: exclude(event:startup,-101) -> : [from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:ping,-101) -> : [from-uri-2 from-uri-3]
	//test: exclude(event:ping,-101) -> : [from-uri-0 from-uri-1 from-uri-4]

}
