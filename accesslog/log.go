package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
)

const (
	errorNilEntryFmt = "{\"error\": \"access data entry is nil\"}"
	//errorNilActuatorFmt = "{\"error\": \"actuator is nil or %v route name is empty\"}"
	errorEmptyFmt = "{\"error\": \"%v log entries are empty\"}"
)

func ingressLog(s string) {
	if opt.ingressFn != nil {
		opt.ingressFn(s)
	}
}

func egressLog(s string) {
	if opt.egressFn != nil {
		opt.egressFn(s)
	}
}

func Log(entry *accessdata.Entry) {
	if entry == nil {
		egressLog(errorNilEntryFmt)
		return
	}
	if entry.IsIngress() {
		if !opt.ingress {
			return
		}
		if len(ingressOperators) == 0 {
			ingressLog(fmt.Sprintf(errorEmptyFmt, entry.Traffic))
			return
		}
		s := accessdata.WriteJson(ingressOperators, entry)
		ingressLog(s)
	} else {
		if !opt.egress {
			return
		}
		if len(egressOperators) == 0 {
			egressLog(fmt.Sprintf(errorEmptyFmt, entry.Traffic))
			return
		}
		s := accessdata.WriteJson(egressOperators, entry)
		egressLog(s)
	}
}
