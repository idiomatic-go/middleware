package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"net/http"
	"strings"
	"time"
)

const (
	errorNilRouteFmt = "{\"error\": \"%v route is nil\"}"
	errorEmptyFmt    = "{\"error\": \"%v log entries are empty\"}"
)

func WriteEgress(start time.Time, duration time.Duration, act actuator.Actuator, req *http.Request, resp *http.Response, responseFlags string) {
	if act == nil {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, EgressTraffic))
		return
	}
	data := NewLogd(EgressTraffic, start, duration, getOrigin(), act, req, resp, responseFlags)
	callExtract(data)
	if !act.Logger().WriteEgress() {
		return
	}
	if len(egressEntries) == 0 {
		egressWrite(fmt.Sprintf(errorEmptyFmt, EgressTraffic))
		return
	}
	s := FormatJson(egressEntries, data)
	egressWrite(s)
}

func WriteIngress(start time.Time, duration time.Duration, act actuator.Actuator, req *http.Request, code int, bytesSent int, responseFlags string) {
	if act == nil {
		ingressWrite(fmt.Sprintf(errorNilRouteFmt, IngressTraffic))
		return
	}
	data := NewLogd(IngressTraffic, start, duration, getOrigin(), act, req, nil, responseFlags)
	data.StatusCode = code
	data.BytesSent = bytesSent
	callExtract(data)
	if !act.Logger().WriteIngress() {
		return
	}
	if len(ingressEntries) == 0 {
		ingressWrite(fmt.Sprintf(errorEmptyFmt, IngressTraffic))
		return
	}
	s := FormatJson(ingressEntries, data)
	ingressWrite(s)
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
