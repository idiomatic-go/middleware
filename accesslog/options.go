package accesslog

import (
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

// Write - override log output disposition, default is log.Println
type Write func(s string)

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

type options struct {
	writeIngress bool
	writeEgress  bool
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
