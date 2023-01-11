package facebook

import (
	"github.com/idiomatic-go/middleware/host"
	"reflect"
	"sync/atomic"
)

type pkg struct{}

var (
	Uri        = pkgPath
	pkgPath    = reflect.TypeOf(any(pkg{})).PkgPath()
	c          = make(chan host.Message, 1)
	envMatcher host.EnvironmentMatcher
	started    int64
)

func IsStarted() bool { return atomic.LoadInt64(&started) != 0 }

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
		if envMatcher != nil {
			atomic.StoreInt64(&started, 1)
			host.StartupReplyTo(msg, true)
		} else {
			host.StartupReplyTo(msg, false)
		}
	case host.ShutdownEvent:
	}
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			messageHandler(msg)
		}
	}
}
