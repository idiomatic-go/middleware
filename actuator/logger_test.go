package actuator

import (
	"fmt"
	"net/http"
	"time"
)

// TODO: test if can work with a nil range argument - works on a nil slice

func Example_newLogger() {
	l := newLogger(nil)
	fmt.Printf("test: newLogger(nil) -> [enabled:%v] [writeEgress:%v] [writeIngress:%v] [pingTraffic:%v] [accessFn:%v]\n", l.IsEnabled(), l.WriteEgress(), l.WriteIngress(), l.IsPingTraffic(""), l.config.accessInvoke != nil)

	l = newLogger(NewLoggerConfig(false, false, false, nil, nil))
	fmt.Printf("test: newLogger(nil) -> [enabled:%v] [writeEgress:%v] [writeIngress:%v] [pingTraffic:%v] [accessFn:%v]\n", l.IsEnabled(), l.WriteEgress(), l.WriteIngress(), l.IsPingTraffic(""), l.config.accessInvoke != nil)

	l = newLogger(NewLoggerConfig(false, false, false, nil, []string{"route-name"}))
	fmt.Printf("test: newLogger(nil) -> [enabled:%v] [writeEgress:%v] [writeIngress:%v] [pingTraffic:%v] [accessFn:%v]\n", l.IsEnabled(), l.WriteEgress(), l.WriteIngress(), l.IsPingTraffic("route-name"), l.config.accessInvoke != nil)

	l.Disable()
	fmt.Printf("test: Disable() -> [enabled:%v]\n", l.IsEnabled())

	l.Enable()
	fmt.Printf("test: Enabled() -> [enabled:%v]\n", l.IsEnabled())

	//Output:
	//test: newLogger(nil) -> [enabled:true] [writeEgress:true] [writeIngress:true] [pingTraffic:false] [accessFn:true]
	//test: newLogger(nil) -> [enabled:true] [writeEgress:false] [writeIngress:false] [pingTraffic:false] [accessFn:true]
	//test: newLogger(nil) -> [enabled:true] [writeEgress:false] [writeIngress:false] [pingTraffic:true] [accessFn:true]
	//test: Disable() -> [enabled:false]
	//test: Enabled() -> [enabled:true]

}

func Example_defaultLogger() {
	l := defaultLogger
	fmt.Printf("test: defaultLogger -> [enabled:%v] [writeEgress:%v] [writeIngress:%v] [pingTraffic:%v] [accessFn:%v]\n", l.IsEnabled(), l.WriteEgress(), l.WriteIngress(), l.IsPingTraffic(""), l.config.accessInvoke != nil)

	SetDefaultLogger(NewLoggerConfig(false, false, false, nil, []string{"route-name"}))
	l = defaultLogger
	fmt.Printf("test: defaultLogger -> [enabled:%v] [writeEgress:%v] [writeIngress:%v] [pingTraffic:%v] [accessFn:%v]\n", l.IsEnabled(), l.WriteEgress(), l.WriteIngress(), l.IsPingTraffic(""), l.config.accessInvoke != nil)

	//Output:
	//test: defaultLogger -> [enabled:true] [writeEgress:true] [writeIngress:true] [pingTraffic:false] [accessFn:true]
	//test: defaultLogger -> [enabled:true] [writeEgress:false] [writeIngress:false] [pingTraffic:false] [accessFn:true]

}

func _Example_LogAccess() {
	start := time.Now()
	fn := func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
		fmt.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"responseFlags\":\"%v\"}\n", traffic, start, duration, req, resp, respFlags)
	}

	defaultLogger.LogAccess(nil, "ingress", start, time.Since(start), nil, nil, "flags")
	time.Sleep(time.Second * 1)
	start = time.Now()
	l := newLogger(NewLoggerConfig(true, true, true, fn, nil))
	l.LogAccess(nil, "egress", start, time.Since(start), nil, nil, "new-flags")

	//Output:
	//{"traffic":"egress","start_time":"2022-12-19 07:53:31.9875524 -0600 CST m=+0.006632001","duration_ms":0s,"request":"<nil>","response":"<nil>","responseFlags":"flags"}

}
