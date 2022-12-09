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
	EgressTraffic      = "egress"
	IngressTraffic     = "ingress"
	PingTraffic        = "ping"
	RateLimitFlag      = "FL"
	UpstreamTmeoutFlag = "UT"
)

// Logd - struct for all logging information
type Logd struct {
	Traffic  string
	Start    time.Time
	Duration time.Duration
	Origin   *Origin
	Route    route.Route

	// Request
	Url      string
	Path     string
	Protocol string
	Method   string
	Header   http.Header

	// Response
	RespCode      int
	BytesSent     int   // ingress response
	BytesReceived int64 // egress response content length
	ResponseFlags string
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
