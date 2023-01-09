package actuator

import (
	"fmt"
)

func Example_newLogger() {
	l := newLogger(nil)
	fmt.Printf("test: newLogger(nil) -> [accessInvoke:%v]\n", l.config.accessInvoke != nil)

	l = newLogger(NewLoggerConfig(nil))
	fmt.Printf("test: newLogger(nil) -> [accessInvoke:%v]\n", l.config.accessInvoke != nil)

	//l = newLogger(NewLoggerConfig(nil))
	//fmt.Printf("test: newLogger(nil) -> [enabled:%v] [accessFn:%v]\n", l.IsEnabled(), l.config.ingressInvoke != nil)

	//Output:
	//test: newLogger(nil) -> [accessInvoke:true]
	//test: newLogger(nil) -> [accessInvoke:true]

}

func Example_defaultLogger() {
	l := defaultLogger
	fmt.Printf("test: defaultLogger -> [accessInvoke:%v]\n", l.config.accessInvoke != nil)

	SetLoggerAccess(nil)
	l = defaultLogger
	fmt.Printf("test: defaultLogger -> [accessInvoke:%v]\n", l.config.accessInvoke != nil)

	//Output:
	//test: defaultLogger -> [accessInvoke:true]
	//test: defaultLogger -> [accessInvoke:true]

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
	l.LogAccess(nil, "handler", start, time.Since(start), nil, nil, "new-flags")

	//Output:
	//{"traffic":"handler","start_time":"2022-12-19 07:53:31.9875524 -0600 CST m=+0.006632001","duration_ms":0s,"request":"<nil>","response":"<nil>","responseFlags":"flags"}

}


*/
