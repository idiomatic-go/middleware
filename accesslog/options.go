package accesslog

import (
	"fmt"
	"log"
)

// Extract - optionally allows extraction of log data
type Extract func(l *Logd)

// Write - override log output disposition, default is log.Println
type Write func(s string)

type options struct {
	extractFn    Extract
	ingressWrite Write
	egressWrite  Write
	origin       Origin
	pingRoutes   []string
}

var opt options

func init() {
	SetIngressWrite(nil)
	SetEgressWrite(nil)
}

func IsExtract() bool {
	return opt.extractFn != nil
}

func SetExtract(fn Extract) {
	opt.extractFn = fn
}

func callExtract(l *Logd) {
	if IsExtract() {
		opt.extractFn(l)
	}
}

func SetIngressWrite(fn Write) {
	if fn != nil {
		opt.ingressWrite = fn
	} else {
		opt.ingressWrite = func(s string) {
			log.Println(s)
		}
	}
}

func SetEgressWrite(fn Write) {
	if fn != nil {
		opt.egressWrite = fn
	} else {
		opt.egressWrite = func(s string) {
			log.Println(s)
		}
	}
}

func SetTestIngressWrite() {
	SetIngressWrite(func(s string) {
		fmt.Printf("test: WriteIngress() -> [%v]\n", s)
	})
}

func SetTestEgressWrite() {
	SetEgressWrite(func(s string) {
		fmt.Printf("test: WriteEgress() -> [%v]\n", s)
	})
}

func ingressWrite(s string) {
	if opt.ingressWrite != nil {
		opt.ingressWrite(s)
	}
}

func egressWrite(s string) {
	if opt.egressWrite != nil {
		opt.egressWrite(s)
	}
}

// SetOrigin - required to track service identification
func SetOrigin(o Origin) {
	opt.origin = o
}

func getOrigin() *Origin {
	return &opt.origin
}

// SetPingRoutes - initialize the ping routes
func SetPingRoutes(routes []string) {
	opt.pingRoutes = routes
}

func isPingTraffic(name string) bool {
	for _, n := range opt.pingRoutes {
		if n == name {
			return true
		}
	}
	return false
}
