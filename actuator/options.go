package actuator

import (
	"github.com/idiomatic-go/middleware/accessdata"
	"net/http"
)

type Matcher func(req *http.Request) (routeName string)

type LogAccess func(entry *accessdata.Entry)

func SetLoggerAccess(fn LogAccess) {
	if fn != nil {
		defaultLogger.config.accessInvoke = fn
	}
}

// Extract - optionally allows extraction of log data
type Extract func(l *accessdata.Entry)

func EnableExtract(fn Extract) {
	if fn != nil {
		defaultExtract = fn
	}
}
