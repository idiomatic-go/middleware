package actuator

import (
	"fmt"
	"net/http"
	"time"
)

func ExampleActuator_newActuator() {
	t := newTable(true)

	a := newActuator("test", t, NewTimeoutConfig(time.Millisecond*1500, 0), NewRateLimiterConfig(100, 10, 503))

	_, toOk := a.Timeout()
	_, rateOk := a.RateLimiter()
	_, retryOk := a.Retry()
	_, failOk := a.Failover()
	fmt.Printf("test: newActuator() -> [logger:%v] [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", a.Logger() != nil, toOk, rateOk, retryOk, failOk)

	d := a.timeout.Duration()
	a1 := cloneActuator[*timeout](a, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	d1 := a1.timeout.Duration()
	fmt.Printf("test: cloneActuator() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//Output:
	//test: newActuator() -> [logger:true] [timeout:true] [rateLimit:true] [retry:false] [failover:false]
	//test: cloneActuator() -> [prev-duration:1.5s] [curr-duration:500ms]

}
