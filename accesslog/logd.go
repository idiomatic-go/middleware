package accesslog

import (
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (l *Logd) IsIngress() bool {
	return l.Traffic == IngressTraffic
}

func (l *Logd) IsEgress() bool {
	return l.Traffic == EgressTraffic
}

func (l *Logd) IsPing() bool {
	return l.IsIngress() && l.Route.IsPingTraffic()
}

func (l *Logd) AddResponse(resp *http.Response) {
	if resp == nil {
		return
	}
	l.RespCode = resp.StatusCode
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
	case RouteNameOperator:
		return l.Route.Name()
	case StartTimeOperator:
		return FmtTimestamp(l.Start)
	case DurationOperator:
		d := int(l.Duration / time.Duration(1e6))
		return strconv.Itoa(d)

		// Origin
	case OriginRegionOperator:
		return l.Origin.Region
	case OriginZoneOperator:
		return l.Origin.Zone
	case OriginSubZoneOperator:
		return l.Origin.SubZone
	case OriginServiceOperator:
		return l.Origin.Service
	case OriginInstanceIdOperator:
		return l.Origin.InstanceId

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
	case ResponseFlagsOperator:
		return l.ResponseFlags
	case ResponseBytesReceivedOperator:
		return strconv.Itoa(int(l.BytesReceived))
	case ResponseBytesSentOperator:
		return strconv.Itoa(l.BytesSent)
	case ResponseCodeOperator:
		return strconv.Itoa(l.RespCode)

		// Timeout
	case TimeoutOperator:
		return strconv.Itoa(l.Route.Timeout())

		// Rate Limiting
	case RateLimitOperator:
		if l.Route.Limit() == rate.Inf {
			return "INF"
		}
		return strconv.Itoa(int(l.Route.Limit()))
	case RateBurstOperator:
		return strconv.Itoa(l.Route.Burst())

	}
	return ""
}

func (l *Logd) HeaderValue(entry Entry) string {
	tokens := strings.Split(entry.Operator, ":")
	if len(tokens) == 1 {
		return ""
	}
	//name := NormalizeHttpHeaderName(tokens[1])
	//values := l.Header.Get(tokens[1])
	//if values != nil {
	//	return values[0]
	//}
	return l.Header.Get(tokens[1])
}
