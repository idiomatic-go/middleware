package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
)

const (
	errorNilEntryFmt    = "{\"error\": \"access data entry is nil\"}"
	errorNilActuatorFmt = "{\"error\": \"actuator is nil or %v route name is empty\"}"
	errorEmptyFmt       = "{\"error\": \"%v log entries are empty\"}"
)

func Log(entry *accessdata.Entry) {
	if entry == nil {
		egressWrite(errorNilEntryFmt)
		return
	}
	//if entry.ActState == nil || entry.ActState[accessdata.ActName] == "" {
	//	egressWrite(fmt.Sprintf(errorNilActuatorFmt, entry.Traffic))
	//	return
	//}
	//data := accessdata.NewEntry(traffic, start, duration, actState, req, resp, statusFlags)
	//callExtract(data)
	if entry.IsIngress() {
		if !opt.writeIngress {
			return
		}
		if len(ingressOperators) == 0 {
			ingressWrite(fmt.Sprintf(errorEmptyFmt, entry.Traffic))
			return
		}
		s := accessdata.WriteJson(ingressOperators, entry)
		ingressWrite(s)
	} else {
		if !opt.writeEgress {
			return
		}
		if len(egressOperators) == 0 {
			egressWrite(fmt.Sprintf(errorEmptyFmt, entry.Traffic))
			return
		}
		s := accessdata.WriteJson(egressOperators, entry)
		egressWrite(s)
	}
}
