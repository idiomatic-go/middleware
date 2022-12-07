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

func WriteEgress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, err error) {
	if route == nil {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, EgressTraffic))
		return
	}
	data := &Logd{
		traffic:      EgressTraffic,
		start:        start,
		duration:     duration,
		bytesWritten: 0,
		route:        route,
		req:          req,
		resp:         resp,
		err:          err,
		code:         0,
		remapStatus:  remapStatus,
	}
	if extractFn != nil {
		extractFn(data)
	}
	if !route.IsLogging() {
		return
	}
	if len(egressAttrs) == 0 {
		egressWrite(fmt.Sprintf(errorEmptyFmt, EgressTraffic))
		return
	}
	s := formatJson(egressAttrs, data)
	egressWrite(s)
}

func WriteIngress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, code int, written int, err error) {
	if route == nil {
		ingressWrite(fmt.Sprintf(errorNilRouteFmt, IngressTraffic))
		return
	}
	data := &Logd{
		traffic:      IngressTraffic,
		start:        start,
		duration:     duration,
		bytesWritten: written,
		route:        route,
		req:          req,
		resp:         nil,
		err:          err,
		code:         code,
		remapStatus:  remapStatus,
	}
	if extractFn != nil {
		extractFn(data)
	}
	if !route.IsLogging() {
		return
	}
	if len(ingressAttrs) == 0 {
		ingressWrite(fmt.Sprintf(errorEmptyFmt, IngressTraffic))
		return
	}
	s := formatJson(ingressAttrs, data)
	ingressWrite(s)
}

func formatJson(attrs []attribute, data *Logd) string {
	sb := strings.Builder{}
	for _, attr := range attrs {
		if attr.isDirect() {
			writeJson(&sb, attr.name, attr.value, attr.stringValue)
			continue
		}
		writeJson(&sb, attr.name, data.value(attr), attr.stringValue)
	}
	sb.WriteString("}")
	return sb.String()
}
