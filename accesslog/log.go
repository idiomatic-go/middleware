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
		if len(ingressEntries) == 0 {
			ingressWrite(fmt.Sprintf(errorEmptyFmt, traffic))
			return
		}
		s := accessdata.WriteJson(ingressEntries, data)
		ingressWrite(s)
	} else {
		if !opt.writeEgress {
			return
		}
		if len(egressEntries) == 0 {
			egressWrite(fmt.Sprintf(errorEmptyFmt, traffic))
			return
		}
		s := accessdata.WriteJson(egressEntries, data)
		egressWrite(s)
	}
}
