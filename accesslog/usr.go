package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"log"
	"net/http"
	"time"
)

// Origin - attributes that uniquely identify a service instance
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

var origin Origin

// SetOrigin - required to track service identification
func SetOrigin(o Origin) {
	origin = o
}

const (
	EgressTraffic  = "egress"
	IngressTraffic = "ingress"
	PingTraffic    = "ping"
)

// Logd - struct for all logging information
type Logd struct {
	Origin       Origin
	Traffic      string
	Start        time.Time
	Duration     time.Duration
	BytesWritten int
	Route        *route.Route
	Req          *http.Request
	Resp         *http.Response
	Err          error
	Code         int
	RemapStatus  bool
}

// Extract - optionally allows extraction of log data
type Extract func(l *Logd)

var extractFn Extract

func SetExtract(fn Extract) {
	extractFn = fn
}

// DisableServiceUnavailableRemap - optionally disables HTTP status code remapping
func DisableServiceUnavailableRemap() {
	remapStatus = false
}

var remapStatus = true

// Write - required configuration of log output
type Write func(s string)

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