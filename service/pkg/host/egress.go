package host

import (
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
)

func initEgress() {
	actuator.EgressTable.SetHostActuator(actuate, actuator.NewRateLimiterConfig(100, 10, http.StatusTooManyRequests))
}

func actuate(act actuator.Actuator, events []actuator.Event) error {
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
