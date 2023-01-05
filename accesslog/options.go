package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"log"
)

// CreateIngressOperators - allows configuration of access log attributes for ingress traffic
func CreateIngressOperators(config []accessdata.Operator) error {
	ingressOperators = []accessdata.Operator{}
	return CreateOperators(&ingressOperators, config)
}

// CreateEgressOperators - allows configuration of access log attributes for egress traffic
func CreateEgressOperators(config []accessdata.Operator) error {
	egressOperators = []accessdata.Operator{}
	return CreateOperators(&egressOperators, config)
}

// LogFn - override log output disposition, default is log.Println
type LogFn func(s string)

func SetIngressLogStatus(enabled bool) {
	opt.ingress = enabled
}

func SetIngressLogFn(fn LogFn) {
	if fn != nil {
		opt.ingressFn = fn
	} else {
		opt.ingressFn = func(s string) {
			log.Println(s)
		}
	}
}

func SetEgressLogStatus(enabled bool) {
	opt.egress = enabled
}

func SetEgressLogFn(fn LogFn) {
	if fn != nil {
		opt.egressFn = fn
	} else {
		opt.egressFn = func(s string) {
			log.Println(s)
		}
	}
}

func SetTestIngressLogFn() {
	SetIngressLogFn(func(s string) {
		fmt.Printf("test: WriteIngress() -> [%v]\n", s)
	})
}

func SetTestEgressLogFn() {
	SetEgressLogFn(func(s string) {
		fmt.Printf("test: WriteEgress() -> [%v]\n", s)
	})
}

type options struct {
	ingress   bool
	egress    bool
	ingressFn LogFn
	egressFn  LogFn
}

var opt options

func init() {
	opt.ingress = true
	opt.egress = true
	SetIngressLogFn(nil)
	SetEgressLogFn(nil)
}
