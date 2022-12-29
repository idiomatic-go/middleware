package accesslog

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	EgressTraffic  = "egress"
	IngressTraffic = "ingress"
	PingTraffic    = "ping"
)

// Origin - attributes that uniquely identify a service instance
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

// Logd - struct for all logging information
type Logd struct {
	Traffic     string
	Start       time.Time
	Duration    time.Duration
	Origin      *Origin
	RouteName   string
	Timeout     []string
	RateLimiter []string
	Retry       []string
	Failover    []string

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

func NewLogd(traffic string, start time.Time, duration time.Duration, origin *Origin, routeName string, timeout []string, rateLimiter []string, failover []string, retry []string, req *http.Request, resp *http.Response, statusFlags string) *Logd {
	l := new(Logd)
	l.Traffic = traffic
	l.Start = start
	l.Duration = duration
	l.Origin = origin
	l.RouteName = routeName
	l.Timeout = timeout
	l.RateLimiter = rateLimiter
	l.Retry = retry
	l.Failover = failover
	l.AddRequest(req)
	l.AddResponse(resp)
	l.StatusFlags = statusFlags
	return l
}

func (l *Logd) IsIngress() bool {
	return l.Traffic == IngressTraffic
}

func (l *Logd) IsEgress() bool {
	return l.Traffic == EgressTraffic
}

func (l *Logd) IsPing() bool {
	return l.IsIngress() && isPingTraffic(l.RouteName)
}

func (l *Logd) AddResponse(resp *http.Response) {
	if resp == nil {
		return
	}
	l.StatusCode = resp.StatusCode
	l.BytesReceived = resp.ContentLength
}

func (l *Logd) AddRequest(req *http.Request) {
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

func (l *Logd) Value(entry Entry) string {
	if entry.IsClientHeader() {
		return l.HeaderValue(entry)
	}
	switch entry.Operator {
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
		return l.RouteName
		/*
			case TimeoutDurationOperator:
				if l.Timeout != nil {
					return l..
				}Enabled {
					return l.Act.Timeout.Value(0)
				}
			case RateLimitOperator:
				if l.Act.RateLimiter.Enabled {
					return l.Act.RateLimiter.Value(0)
				}
			case RateBurstOperator:
				if l.Act.RateLimiter.Enabled {
					return l.Act.RateLimiter.Value(1)
				}
			case FailoverOperator:
				if l.Act.Failover.Enabled {
					return l.Act.Failover.Value(0)
				}

		*/
	}

	return ""
}

func (l *Logd) HeaderValue(entry Entry) string {
	tokens := strings.Split(entry.Operator, ":")
	if len(tokens) == 1 {
		return ""
	}
	return l.Header.Get(tokens[1])
}
