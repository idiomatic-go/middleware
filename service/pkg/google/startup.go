package google

import (
	"github.com/idiomatic-go/middleware/host"
)

const (
	Uri = "google"
)

var c = make(chan host.Message, 1)
var envMatcher host.EnvironmentMatcher
var isStarted bool

func IsStarted() bool { return isStarted }

func isDevEnv() bool {
	if envMatcher == nil {
		return true
	}
	return envMatcher(host.DevEnv)
}
func isTestEnv() bool {
	if envMatcher == nil {
		return false
	}
	return envMatcher(host.TestEnv)
}

func init() {
	host.RegisterResource(Uri, c)
	go receive()
}

var messageHandler host.MessageHandler = func(msg host.Message) {
	switch msg.Event {
	case host.StartupEvent:
		envMatcher = host.AccessEnvironmentMatcher(&msg)
	case host.ShutdownEvent:
	}
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			// Exit on a closed channel
			if !open {
				return
			}
			messageHandler(msg)
			if msg.ReplyTo != nil {
				if IsStarted() {
					msg.ReplyTo(host.NewStartupSuccessfulMessage(msg))
				} else {

				}
			}
		}
	}
}
