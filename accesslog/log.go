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

func LogEgress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, err error) {
	if route == nil {
		egressWrite("error : route is nil")
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
	writeMarkup(&sb, startTimeName, data)
	writeMarkup(&sb, regionName, data)
	writeMarkup(&sb, zoneName, data)
	writeMarkup(&sb, subZoneName, data)
	writeMarkup(&sb, serviceName, data)
	writeMarkup(&sb, instanceIdName, data)
	writeMarkup(&sb, trafficName, data)
	writeMarkup(&sb, routeName, data)
	writeMarkup(&sb, durationName, data)
	if req != nil {
		writeMarkup(&sb, urlName, data)
		writeMarkup(&sb, methodName, data)
	}
	if resp != nil {
		writeMarkup(&sb, statusCodeName, data)
		writeMarkup(&sb, protocolName, data)
	}
	sb.WriteString("}")
	if route.WriteAccessLog {
		egressWrite(sb.String())
	}
}

func LogIngress(route *route.Route, start time.Time, duration time.Duration, req *http.Request, code int, written int, err error) {
	if route == nil {
		ingressWrite("error : route is nil")
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
	if route.WriteAccessLog {
		sb := strings.Builder{}
		writeMarkup(&sb, startTimeName, data)
		writeMarkup(&sb, regionName, data)
		writeMarkup(&sb, zoneName, data)
		writeMarkup(&sb, subZoneName, data)
		writeMarkup(&sb, serviceName, data)
		writeMarkup(&sb, instanceIdName, data)
		writeMarkup(&sb, trafficName, data)
		writeMarkup(&sb, routeName, data)
		writeMarkup(&sb, durationName, data)
		sb.WriteString("}")
		ingressWrite(sb.String())
	}
	if route.RedirectAccessLog {

	}
}
