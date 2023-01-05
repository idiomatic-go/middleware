package host

import (
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
	"time"
)

func initIngress() {
	actuator.IngressTable.SetHostActuator(hostActuate, actuator.NewRateLimiterConfig(100, 10, http.StatusTooManyRequests))
	actuator.IngressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))
}

func hostActuate(act actuator.Actuator, events []actuator.Event) error {
	if len(events) == 0 {
		return nil
	}
	if events[0].IsWatch() {
		actuator.AdjustRateLimiter(act, -10)
		return nil
	}
	if events[0].IsCancel() {
		actuator.AdjustRateLimiter(act, 10)
		return nil
	}
	return nil
}
