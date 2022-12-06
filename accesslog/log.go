package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var remapStatus = true

func DisableServiceUnavailableRemap() {
	remapStatus = false
}

type LogFn func(s string)

var ingressLogFn LogFn
var egressLogFn LogFn

func init() {
	SetIngressLogFn(nil)
	SetEgressLogFn(nil)
}

func SetIngressLogFn(fn LogFn) {
	if fn != nil {
		ingressLogFn = fn
	} else {
		ingressLogFn = func(s string) {
			log.Println(s)
		}
	}
}

func SetEgressLogFn(fn LogFn) {
	if fn != nil {
		egressLogFn = fn
	} else {
		egressLogFn = func(s string) {
			log.Println(s)
		}
	}
}

type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

var origin Origin

func SetOrigin(o Origin) {
	origin = o
}

func LogEgress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, err error) {
	if route == nil {
		if egressLogFn != nil {
			egressLogFn("error : route is nil")
		}
		return
	}
	sb := strings.Builder{}
	sb.WriteString("{")
	writeStartTime(&sb, start)
	writeLocation(&sb)
	writeMarkup(&sb, "traffic", "egress", "")
	writeMarkup(&sb, "route_name", route.Name, "")
	writeDuration(&sb, duration)
	if req != nil {
		writeMarkup(&sb, "url", req.URL.String(), "")
		writeMarkup(&sb, "method", req.Method, "")
	}
	if resp != nil {
		writeMarkup(&sb, "status", strconv.Itoa(resp.StatusCode), markupValueComma)
		writeMarkup(&sb, "protocol", resp.Proto, "")
	}
	sb.WriteString("}")
	if route.WriteAccessLog && egressLogFn != nil {
		egressLogFn(sb.String())
	}
}

func LogIngress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, code int, written int) {
	// TODO : determine what to log if route is ping traffic
}
