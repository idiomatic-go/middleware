package automation

import (
	"fmt"
	"time"
)

func ExampleActuator_newActuator() {
	t := newTable(true)

	a := newActuator("test", t, NewTimeoutConfig(time.Millisecond*1500), NewRateLimiterConfig(100, 10, 503, false))

	log := a.Logger() != nil
	to := a.Timeout() != nil
	rate := a.RateLimiter() != nil
	retry := a.Retry() != nil
	fail := a.Failover() != nil
	fmt.Printf("test: newActuator() -> [logger:%v] [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", log, to, rate, retry, fail)

	_, d := a.Timeout().Duration()
	a1 := cloneActuator[*timeout](a, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500)))

	_, d1 := a1.Timeout().Duration()
	fmt.Printf("test: cloneActuator() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//Output:
	//test: newActuator() -> [logger:true] [timeout:true] [rateLimit:true] [retry:false] [failover:false]
	//test: cloneActuator() -> [prev-duration:1.5s] [curr-duration:500ms]

}
