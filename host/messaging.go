package host

import "github.com/idiomatic-go/middleware/template"

const (
	StartupEvent  = "event:startup"
	ShutdownEvent = "event:shutdown"
	PingEvent     = "event:ping"
	Name          = "host"
)

type MessageHandler func(msg Message)

type Message struct {
	To      string
	From    string
	Event   string
	Status  int
	Content []any
	ReplyTo MessageHandler
}

func SendMessage(msg Message) error {
	return directory.send(msg)
}

func NewStartupSuccessfulMessage(from Message) Message {
	return Message{To: from.From, From: from.To, Event: StartupEvent, Status: template.StatusOk}
}
