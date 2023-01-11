package actuator

import (
	"fmt"
	"net/http"
	"time"
)

var actuatorFn Actuate = func(act Actuator, events []Event) error {
	fmt.Printf("test: Actuate() -> [%v]\n", act.Name())
	return nil
}

func ExampleActuator_newActuator() {
	t := newTable(true)

	a, _ := newActuator("test", t, actuatorFn, NewTimeoutConfig(time.Millisecond*1500, 0), NewRateLimiterConfig(100, 10, 503))

	_, toOk := a.Timeout()
	_, rateOk := a.RateLimiter()
	_, retryOk := a.Retry()
	_, failOk := a.Failover()
	fmt.Printf("test: newActuator() -> [logger:%v] [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", a.Logger() != nil, toOk, rateOk, retryOk, failOk)

	d := a.timeout.Duration()
	a1 := cloneActuator[*timeout](a, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	d1 := a1.timeout.Duration()
	fmt.Printf("test: cloneActuator() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	a.Actuate(nil)

	//Output:
	//test: newActuator() -> [logger:true] [timeout:true] [rateLimit:true] [retry:false] [failover:false]
	//test: cloneActuator() -> [prev-duration:1.5s] [curr-duration:500ms]
	//test: Actuate() -> [test]

}

func ExampleActuator_newActuator_config() {
	t := newTable(true)

	a, _ := newActuator("test", t, actuatorFn, nil, NewTimeoutConfig(time.Millisecond*1500, 0), nil, NewRateLimiterConfig(100, 10, 503), nil)

	_, toOk := a.Timeout()
	_, rateOk := a.RateLimiter()
	_, retryOk := a.Retry()
	_, failOk := a.Failover()
	fmt.Printf("test: newActuator() -> [logger:%v] [timeout:%v] [rateLimit:%v] [retry:%v] [failover:%v]\n", a.Logger() != nil, toOk, rateOk, retryOk, failOk)

	//d := a.timeout.Duration()
	//a1 := cloneActuator[*timeout](a, newTimeout("new-timeout", t, NewTimeoutConfig(time.Millisecond*500, http.StatusGatewayTimeout)))

	//d1 := a1.timeout.Duration()
	//fmt.Printf("test: cloneActuator() -> [prev-duration:%v] [curr-duration:%v]\n", d, d1)

	//a.Actuate(nil)

	//Output:
	//test: newActuator() -> [logger:true] [timeout:true] [rateLimit:true] [retry:false] [failover:false]

}

func ExampleActuator_newActuator_Error() {
	t := newTable(false)

	_, errs := newActuator("test", t, actuatorFn, NewTimeoutConfig(time.Millisecond*1500, 0), NewRateLimiterConfig(100, 10, 503))
	fmt.Printf("test: newActuator() -> [errs:%v]\n", errs)

	_, errs = newActuator("test", t, actuatorFn, NewTimeoutConfig(time.Millisecond*1500, 0), NewRetryConfig(nil, 100, 10, 0))
	fmt.Printf("test: newActuator() -> [errs:%v]\n", errs)

	_, errs = newActuator("test", t, actuatorFn, NewTimeoutConfig(0, 0))
	fmt.Printf("test: newActuator() -> [errs:%v]\n", errs)

	_, errs = newActuator("test", t, actuatorFn, NewTimeoutConfig(10, 0), NewFailoverConfig(nil))
	fmt.Printf("test: newActuator() -> [errs:%v]\n", errs)

	_, errs = newActuator("test", t, actuatorFn, NewRateLimiterConfig(-1, 10, 504))
	fmt.Printf("test: newActuator() -> [errs:%v]\n", errs)

	//Output:
	//test: newActuator() -> [errs:[]]
	//test: newActuator() -> [errs:[invalid configuration: RetryController status codes are empty]]
	//test: newActuator() -> [errs:[invalid configuration: TimeoutController duration is <= 0]]
	//test: newActuator() -> [errs:[invalid configuration: FailoverController FailureInvoke function is nil]]
	//test: newActuator() -> [errs:[invalid configuration: RateLimiterController limit is < 0]]

}
