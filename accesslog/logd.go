package accesslog

import (
	"strconv"
	"strings"
	"time"
)

func (l Logd) IsIngress() bool {
	return l.Traffic == IngressTraffic
}

func (l Logd) IsEgress() bool {
	return l.Traffic == EgressTraffic
}

func (l Logd) IsPing() bool {
	return l.IsIngress() && l.Route != nil && l.Route.Ping
}

func (l Logd) Value(entry Entry) string {
	if entry.IsHeader() {
		return l.HeaderValue(entry)
	}
	switch entry.Operator() {
	case TrafficOperator:
		if l.IsPing() {
			return PingTraffic
		}
		return l.Traffic
	case RegionOperator:
		return l.Origin.Region
	case ZoneOperator:
		return l.Origin.Zone
	case SubZoneOperator:
		return l.Origin.SubZone
	case ServiceNameOperator:
		return l.Origin.Service
	case InstanceIdOperator:
		return l.Origin.InstanceId
	case StartTimeOperator:
		return FmtTimestamp(l.Start)
	case DurationOperator:
		d := int(l.Duration / time.Duration(1e6))
		return strconv.Itoa(d)

	// Route - check for nil
	case RouteNameOperator:
		if l.Route != nil {
			return l.Route.Name
		}

	// Http Request - check for nil
	case HttpMethodOperator:
		if l.Req != nil {
			return ""
		}

	// Http Response - check for nil
	case ResponseFlagsOperator:
		return l.ResponseFlags

	case ResponseCodeOperator:
		if l.IsIngress() {
			return strconv.Itoa(l.Code)
		} else {
			if l.Resp != nil {
				return strconv.Itoa(l.Resp.StatusCode)
			}
		}
	}
	return ""
}

func (l Logd) HeaderValue(entry Entry) string {
	if l.Req == nil {
		return ""
	}
	tokens := strings.Split(entry.Operator(), ":")
	if len(tokens) == 1 {
		return ""
	}
	name := NormalizeHttpHeaderName(tokens[1])
	values := l.Req.Header[name]
	if values != nil {
		return values[0]
	}
	return ""
}
