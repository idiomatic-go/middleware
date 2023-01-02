package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"log"
)

// Extract - optionally allows extraction of log data
type Extract func(l *accessdata.Entry)

// Write - override log output disposition, default is log.Println
type Write func(s string)

type options struct {
	writeIngress bool
	writeEgress  bool
	extractFn    Extract
	ingressWrite Write
	egressWrite  Write
}

var opt options

func init() {
	opt.writeIngress = true
	opt.writeEgress = true
	SetIngressWrite(nil)
	SetEgressWrite(nil)
}

func IsExtract() bool {
	return opt.extractFn != nil
}

func SetExtract(fn Extract) {
	opt.extractFn = fn
}

func callExtract(l *accessdata.Entry) {
	if IsExtract() {
		opt.extractFn(l)
	}
}

func SetIngressWriteStatus(enabled bool) {
	opt.writeIngress = enabled
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

func SetEgressWriteStatus(enabled bool) {
	opt.writeEgress = enabled
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
