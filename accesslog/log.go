package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"net/http"
	"time"
)

const (
	errorNilRouteFmt = "{\"error\": \"%v route name is empty\"}"
	errorEmptyFmt    = "{\"error\": \"%v log entries are empty\"}"
)

func Log(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) {
	if actState == nil || actState[accessdata.ActName] == "" {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, traffic))
		return
	}
	data := accessdata.NewEntry(traffic, start, duration, actState, req, resp, statusFlags)
	callExtract(data)
	if traffic == accessdata.IngressTraffic {
		if !opt.writeIngress {
			return
		}
		if len(ingressOperators) == 0 {
			ingressWrite(fmt.Sprintf(errorEmptyFmt, traffic))
			return
		}
		s := accessdata.WriteJson(ingressOperators, data)
		ingressWrite(s)
	} else {
		if !opt.writeEgress {
			return
		}
		if len(egressOperators) == 0 {
			egressWrite(fmt.Sprintf(errorEmptyFmt, traffic))
			return
		}
		s := accessdata.WriteJson(egressOperators, data)
		egressWrite(s)
	}
}
