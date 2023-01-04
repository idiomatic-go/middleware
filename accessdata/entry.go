package accessdata

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	EgressTraffic      = "egress"
	IngressTraffic     = "ingress"
	PingTraffic        = "ping"
	TimeoutName        = "timeout"
	FailoverName       = "failover"
	RetryName          = "retry"
	RetryRateLimitName = "retryRateLimit"
	RetryRateBurstName = "retryBurst"
	RateLimitName      = "rateLimit"
	RateBurstName      = "burst"
	ActName            = "name"
)

// Origin - attributes that uniquely identify a service instance
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

// Entry - struct for all logging information
type Entry struct {
	Traffic  string
	Start    time.Time
	Duration time.Duration
	Origin   *Origin
	ActState map[string]string

	// Request
	Url      string
	Path     string
	Host     string
	Protocol string
	Method   string
	Header   http.Header

	// Response
	StatusCode    int
	BytesSent     int   // ingress response
	BytesReceived int64 // egress response content length
	StatusFlags   string
}

func newEntry(traffic string, start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) *Entry {
	l := new(Entry)
	l.Traffic = traffic
	l.Start = start
	l.Duration = duration
	l.Origin = getOrigin()
	if actState == nil {
		actState = make(map[string]string, 1)
	}
	l.ActState = actState
	l.AddRequest(req)
	l.AddResponse(resp)
	l.StatusFlags = statusFlags
	return l
}

func NewIngressEntry(start time.Time, duration time.Duration, actState map[string]string, req *http.Request, statusCode int, written int, statusFlags string) *Entry {
	e := newEntry(IngressTraffic, start, duration, actState, req, nil, statusFlags)
	e.StatusCode = statusCode
	e.BytesSent = written
	return e
}

func NewEgressEntry(start time.Time, duration time.Duration, actState map[string]string, req *http.Request, resp *http.Response, statusFlags string) *Entry {
	return newEntry(EgressTraffic, start, duration, actState, req, resp, statusFlags)
}

func (l *Entry) IsIngress() bool {
	return l.Traffic == IngressTraffic
}

func (l *Entry) IsEgress() bool {
	return l.Traffic == EgressTraffic
}

func (l *Entry) IsPing() bool {
	return l.IsIngress() && IsPingTraffic(l.ActState[ActName])
}

func (l *Entry) AddResponse(resp *http.Response) {
	if resp == nil {
		return
	}
	l.StatusCode = resp.StatusCode
	l.BytesReceived = resp.ContentLength
}

func (l *Entry) AddRequest(req *http.Request) {
	if req == nil {
		return
	}
	l.Protocol = req.Proto
	l.Method = req.Method
	l.Header = req.Header.Clone()
	if req.URL != nil {
		l.Url = req.URL.String()
		l.Path = req.URL.Path
		if req.Host == "" {
			l.Host = req.URL.Host
		} else {
			l.Host = req.Host
		}
	}
}

func (l *Entry) SetActuatorState(m map[string]string) {
	if m != nil {
		l.ActState = m
	}

}
func (l *Entry) Value(value string) string {
	switch value {
	case TrafficOperator:
		if l.IsPing() {
			return PingTraffic
		}
		return l.Traffic
	case StartTimeOperator:
		return FmtTimestamp(l.Start)
	case DurationOperator:
		d := int(l.Duration / time.Duration(1e6))
		return strconv.Itoa(d)

		// Origin
	case OriginRegionOperator:
		if l.Origin != nil {
			return l.Origin.Region
		}
	case OriginZoneOperator:
		if l.Origin != nil {
			return l.Origin.Zone
		}
	case OriginSubZoneOperator:
		if l.Origin != nil {
			return l.Origin.SubZone
		}
	case OriginServiceOperator:
		if l.Origin != nil {
			return l.Origin.Service
		}
	case OriginInstanceIdOperator:
		if l.Origin != nil {
			return l.Origin.InstanceId
		}

		// Request
	case RequestMethodOperator:
		return l.Method
	case RequestProtocolOperator:
		return l.Protocol
	case RequestPathOperator:
		return l.Path
	case RequestUrlOperator:
		return l.Url
	case RequestHostOperator:
		return l.Host
	case RequestIdOperator:
		return l.Header.Get(RequestIdHeaderName)
	case RequestFromRouteOperator:
		return l.Header.Get(FromRouteHeaderName)
	case RequestUserAgentOperator:
		return l.Header.Get(UserAgentHeaderName)
	case RequestAuthorityOperator:
		return ""
	case RequestForwardedForOperator:
		return l.Header.Get(FordwardedForHeaderName)

		// Response
	case StatusFlagsOperator:
		return l.StatusFlags
	case ResponseBytesReceivedOperator:
		return strconv.Itoa(int(l.BytesReceived))
	case ResponseBytesSentOperator:
		return strconv.Itoa(l.BytesSent)
	case ResponseStatusCodeOperator:
		return strconv.Itoa(l.StatusCode)

	// Actuator State
	case RouteNameOperator:
		return l.ActState[ActName]
	case TimeoutDurationOperator:
		return l.ActState[TimeoutName]
	case RateLimitOperator:
		return l.ActState[RateLimitName]
	case RateBurstOperator:
		return l.ActState[RateBurstName]
	case FailoverOperator:
		return l.ActState[FailoverName]
	case RetryOperator:
		return l.ActState[RetryName]
	case RetryRateLimitOperator:
		return l.ActState[RetryRateLimitName]
	case RetryRateBurstOperator:
		return l.ActState[RetryRateBurstName]
	}
	if strings.HasPrefix(value, RequestReferencePrefix) {
		name := requestOperatorHeaderName(value)
		return l.Header.Get(name)
	}
	if !strings.HasPrefix(value, OperatorPrefix) {
		return value
	}

	return ""
}
