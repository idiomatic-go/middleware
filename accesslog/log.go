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

var ingressWrite Write
var egressWrite Write

func init() {
	SetIngressWrite(nil)
	SetEgressWrite(nil)
}

func WriteEgress(start time.Time, duration time.Duration, route route.Route, req *http.Request, resp *http.Response, responseFlags string) {
	if route == nil {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, EgressTraffic))
		return
	}
	data := &Logd{
		Traffic:       EgressTraffic,
		Start:         start,
		Duration:      duration,
		Origin:        &origin,
		Route:         route,
		ResponseFlags: responseFlags,
	}
	data.AddResponse(resp)
	data.AddRequest(req)
	if extractFn != nil {
		extractFn(data)
	}
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
	data := &Logd{
		Traffic:       IngressTraffic,
		Start:         start,
		Duration:      duration,
		Origin:        &origin,
		Route:         route,
		RespCode:      code,
		BytesSent:     bytesSent,
		ResponseFlags: responseFlags,
	}
	data.AddRequest(req)
	if extractFn != nil {
		extractFn(data)
	}
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
