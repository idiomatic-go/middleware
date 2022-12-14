package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"strings"
	"time"
)

const (
	errorNilRouteFmt = "{\"error\": \"%v route is nil\"}"
	errorEmptyFmt    = "{\"error\": \"%v log entries are empty\"}"
)

func WriteEgress(start time.Time, duration time.Duration, route route.Route, req *http.Request, resp *http.Response, responseFlags string) {
	if route == nil {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, EgressTraffic))
		return
	}
	data := NewLogd(EgressTraffic, start, duration, getOrigin(), route, req, resp, responseFlags)
	callExtract(data)
	if !route.IsLogging() {
		return
	}
	if len(egressEntries) == 0 {
		egressWrite(fmt.Sprintf(errorEmptyFmt, EgressTraffic))
		return
	}
	s := FormatJson(egressEntries, data)
	egressWrite(s)
}

func WriteIngress(start time.Time, duration time.Duration, route route.Route, req *http.Request, code int, bytesSent int, responseFlags string) {
	if route == nil {
		ingressWrite(fmt.Sprintf(errorNilRouteFmt, IngressTraffic))
		return
	}
	data := NewLogd(IngressTraffic, start, duration, getOrigin(), route, req, nil, responseFlags)
	data.StatusCode = code
	data.BytesSent = bytesSent
	callExtract(data)
	if !route.IsLogging() {
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
