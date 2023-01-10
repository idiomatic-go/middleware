package resource

import (
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"time"
)

func ExampleCreateToSend() {
	none := "/resource/none"
	one := "/resource/one"

	registerResourceUnchecked(none, nil)
	registerResourceUnchecked(one, nil)

	m := createToSend(nil, nil)
	fmt.Printf("test: registerResourceUnchecked(nil,nil) -> %v\n", m)

	cm := ContentMap{one: []any{"test content", "additional content"}}
	m = createToSend(cm, nil)
	fmt.Printf("test: registerResourceUnchecked(map,nil) -> %v\n", m)

	//Output:
	//test: registerResourceUnchecked(nil,nil) -> map[/resource/none:{/resource/none host event:startup -4 [] <nil>} /resource/one:{/resource/one host event:startup -4 [] <nil>}]
	//test: registerResourceUnchecked(map,nil) -> map[/resource/none:{/resource/none host event:startup -4 [] <nil>} /resource/one:{/resource/one host event:startup -4 [test content additional content] <nil>}]

}

func Example_Startup() {
	uri1 := "urn:good"
	uri2 := "urn:bad"
	uri3 := "urn:ugly"

	directory.empty()

	c := make(chan Message, 16)
	RegisterResource(uri1, c)
	go good(c)

	c = make(chan Message, 16)
	RegisterResource(uri2, c)
	go bad(c)

	c = make(chan Message, 16)
	RegisterResource(uri3, c)
	go ugly(c)

	status := Startup[template.DebugHandler](time.Second*3, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//fail
}

func good(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			// Exit on a closed channel
			if !open {
				return
			}
			if msg.ReplyTo != nil {
				msg.ReplyTo(Message{To: msg.From, From: msg.To, Event: StartupEvent, Status: template.StatusOk})
			}
		default:
		}
	}
}

func bad(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			// Exit on a closed channel
			if !open {
				return
			}
			if msg.ReplyTo != nil {
				time.Sleep(time.Second)
				msg.ReplyTo(Message{To: msg.From, From: msg.To, Event: StartupEvent, Status: template.StatusOk})
			}
		default:
		}
	}
}

func ugly(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			// Exit on a closed channel
			if !open {
				return
			}
			if msg.ReplyTo != nil {
				time.Sleep(time.Second)
				msg.ReplyTo(Message{To: msg.From, From: msg.To, Event: StartupEvent, Status: template.StatusInternal})
			}
		default:
		}
	}
}
