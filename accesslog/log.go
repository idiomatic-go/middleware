package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"log"
	"net/http"
	"strings"
	"time"
)

var remapStatus = true

func DisableServiceUnavailableRemap() {
	remapStatus = false
}

type Write func(s string)

var ingressWrite Write
var egressWrite Write

func init() {
	SetIngressWrite(nil)
	SetEgressWrite(nil)
}

func SetIngressWrite(fn Write) {
	if fn != nil {
		ingressWrite = fn
	} else {
		ingressWrite = func(s string) {
			log.Println(s)
		}
	}
}

func SetEgressWrite(fn Write) {
	if fn != nil {
		egressWrite = fn
	} else {
		egressWrite = func(s string) {
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

func WriteEgress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, err error) {
	if route == nil {
		egressWrite("error : route is nil")
		return
	}
	if !route.IsLogging() {
		return
	}
	data := &logd{
		traffic:      egressTraffic,
		start:        start,
		duration:     duration,
		bytesWritten: 0,
		route:        route,
		req:          req,
		resp:         resp,
		err:          err,
		code:         0,
	}
	sb := strings.Builder{}
	for _, attr := range egressAttrs {
		if attr.IsDirect() {
			writeJsonMarkup(&sb, attr.name, attr.value, attr.stringValue)
			continue
		}
		writeJsonMarkup(&sb, attr.name, data.resolve(attr), attr.stringValue)
	}
	sb.WriteString("}")
	egressWrite(sb.String())
}

func WriteIngress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, code int, written int, err error) {
	if route == nil {
		ingressWrite("error : route is nil")
		return
	}
	if !route.IsLogging() {
		return
	}
	data := &logd{
		traffic:      ingressTraffic,
		start:        start,
		duration:     duration,
		bytesWritten: written,
		route:        route,
		req:          req,
		resp:         nil,
		err:          err,
		code:         code,
	}
	sb := strings.Builder{}
	for _, attr := range ingressAttrs {
		if attr.IsDirect() {
			writeJsonMarkup(&sb, attr.name, attr.value, attr.stringValue)
			continue
		}
		writeJsonMarkup(&sb, attr.name, data.resolve(attr), attr.stringValue)
	}
	sb.WriteString("}")
	ingressWrite(sb.String())
}
