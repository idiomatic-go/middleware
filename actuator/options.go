package actuator

import (
	"net/http"
	"time"
)

type LogAccess func(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string)

func SetLoggerAccess(fn LogAccess) {
	if fn != nil {
		defaultLogger.config.accessInvoke = fn
	}
}
