package actuator

import (
	"fmt"
)

func Example_newLogger() {
	l := newLogger(nil)
	fmt.Printf("test: newLogger(nil) -> [enabled:%v] [accessInvoke:%v]\n", l.IsEnabled(), l.config.accessInvoke != nil)

	l = newLogger(NewLoggerConfig(nil))
	fmt.Printf("test: newLogger(nil) -> [enabled:%v] [accessInvoke:%v]\n", l.IsEnabled(), l.config.accessInvoke != nil)

	//l = newLogger(NewLoggerConfig(nil))
	//fmt.Printf("test: newLogger(nil) -> [enabled:%v] [accessFn:%v]\n", l.IsEnabled(), l.config.ingressInvoke != nil)

	l.Disable()
	fmt.Printf("test: Disable() -> [enabled:%v]\n", l.IsEnabled())

	l.Enable()
	fmt.Printf("test: Enabled() -> [enabled:%v]\n", l.IsEnabled())

	//Output:
	//test: newLogger(nil) -> [enabled:true] [accessInvoke:true]
	//test: newLogger(nil) -> [enabled:true] [accessInvoke:true]
	//test: Disable() -> [enabled:false]
	//test: Enabled() -> [enabled:true]

}

func Example_defaultLogger() {
	l := defaultLogger
	fmt.Printf("test: defaultLogger -> [enabled:%v] [accessInvoke:%v]\n", l.IsEnabled(), l.config.accessInvoke != nil)

	SetDefaultLogger(NewLoggerConfig(nil))
	l = defaultLogger
	fmt.Printf("test: defaultLogger -> [enabled:%v] [accessInvoke:%v]\n", l.IsEnabled(), l.config.accessInvoke != nil)

	//Output:
	//test: defaultLogger -> [enabled:true] [accessInvoke:true]
	//test: defaultLogger -> [enabled:true] [accessInvoke:true]

}

/*
func _Example_LogAccess() {
	start := time.Now()
	fn := func(act Actuator, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, respFlags string) {
		fmt.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"responseFlags\":\"%v\"}\n", traffic, start, duration, req, resp, respFlags)
	}

	defaultLogger.LogAccess(nil, "ingress", start, time.Since(start), nil, nil, "flags")
	time.Sleep(time.Second * 1)
	start = time.Now()
	l := newLogger(NewLoggerConfig(fn))
	l.LogAccess(nil, "egress", start, time.Since(start), nil, nil, "new-flags")

	//Output:
	//{"traffic":"egress","start_time":"2022-12-19 07:53:31.9875524 -0600 CST m=+0.006632001","duration_ms":0s,"request":"<nil>","response":"<nil>","responseFlags":"flags"}

}


*/
