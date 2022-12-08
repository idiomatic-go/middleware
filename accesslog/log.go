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

func WriteEgress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, responseFlags string, err error) {
	if route == nil {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, EgressTraffic))
		return
	}
	data := &Logd{
		Origin:        origin,
		Traffic:       EgressTraffic,
		Start:         start,
		Duration:      duration,
		BytesWritten:  0,
		Route:         route,
		Req:           req,
		Resp:          resp,
		Err:           err,
		Code:          0,
		RemapStatus:   remapStatus,
		ResponseFlags: responseFlags,
	}
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

func WriteIngress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, code int, written int, responseFlags string, err error) {
	if route == nil {
		ingressWrite(fmt.Sprintf(errorNilRouteFmt, IngressTraffic))
		return
	}
	data := &Logd{
		Origin:        origin,
		Traffic:       IngressTraffic,
		Start:         start,
		Duration:      duration,
		BytesWritten:  written,
		Route:         route,
		Req:           req,
		Resp:          nil,
		Err:           err,
		Code:          code,
		RemapStatus:   remapStatus,
		ResponseFlags: responseFlags,
	}
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
			writeJson(&sb, entry.Name(), entry.Value, entry.StringValue)
			continue
		}
		writeJson(&sb, entry.Name(), data.Value(entry), entry.StringValue)
	}
	sb.WriteString("}")
	return sb.String()
}
