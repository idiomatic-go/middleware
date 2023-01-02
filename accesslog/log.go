package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"net/http"
	"strings"
	"time"
)

const (
	errorNilRouteFmt = "{\"error\": \"%v route name is empty\"}"
	errorEmptyFmt    = "{\"error\": \"%v log entries are empty\"}"
)

func Log(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) {
	if actState == nil || actState[ActName] == "" {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, traffic))
		return
	}
	data := NewLogd(traffic, start, duration, getOrigin(), actState, req, resp, statusFlags)
	callExtract(data)
	if traffic == IngressTraffic {
		if !opt.writeIngress {
			return
		}
		if len(ingressEntries) == 0 {
			ingressWrite(fmt.Sprintf(errorEmptyFmt, traffic))
			return
		}
		s := accessdata.WriteJson(ingressEntries, nil)
		ingressWrite(s)
	} else {
		if !opt.writeEgress {
			return
		}
		if len(egressEntries) == 0 {
			egressWrite(fmt.Sprintf(errorEmptyFmt, traffic))
			return
		}
		s := accessdata.WriteJson(egressEntries, nil)
		egressWrite(s)
	}
}

func FormatJson(items []Entry, data *Logd) string {
	if len(items) == 0 || data == nil {
		return "{}"
	}
	sb := strings.Builder{}
	for _, entry := range items {
		if entry.IsDirect() {
			writeJson(&sb, entry.Name, entry.Value, entry.StringValue)
			continue
		}
		writeJson(&sb, entry.Name, data.Value(entry), entry.StringValue)
	}
	sb.WriteString("}")
	return sb.String()
}
