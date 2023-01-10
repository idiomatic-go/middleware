package resource

const (
	StartupEvent   = "event:startup"
	ShutdownEvent  = "event:shutdown"
	StartupPause   = "event:pause"
	ShutdownResume = "event:resume"
	ErrorEvent     = "event:error"
	RestartEvent   = "event:restart"
	FailoverEvent  = "event:failover"
	FailbackEvent  = "event:failback"
	PingEvent      = "event:ping"
	ProfileEvent   = "event:profile"
	VirtualHost    = "host"

	StatusInProgress  = -100
	StatusNotProvided = -101
)

type ContentMap map[string]any

type MessageHandler func(msg Message)

type Message struct {
	To      string
	From    string
	Event   string
	Status  int
	Content any
	ReplyTo MessageHandler
}

func SendMessage(msg Message) error {
	return directory.send(msg)
}
