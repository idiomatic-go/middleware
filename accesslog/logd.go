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
	return l.IsIngress() && l.PingTraffic
}

func (l Logd) Value(entry Entry) string {
	if entry.IsHeader() {
		return l.HeaderValue(entry)
	}
	switch entry.Operator {
	case TrafficOperator:
		if l.IsPing() {
			return PingTraffic
		}
		return l.Traffic
	case RouteNameOperator:
		return l.RouteName
	case StartTimeOperator:
		return FmtTimestamp(l.Start)
	case DurationOperator:
		d := int(l.Duration / time.Duration(1e6))
		return strconv.Itoa(d)
		// Route - check for nil

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

	// Http Request - check for nil
	case RequestMethodOperator:
		if l.Req == nil {
			return ""
		}
		return l.Req.Method

	case ResponseFlagsOperator:
		return l.ResponseFlags
	case ResponseBytesSentOperator:
		return strconv.Itoa(l.BytesSent)
	case ResponseCodeOperator:
		return strconv.Itoa(l.RespCode)
	}
	return ""
}

func (l Logd) HeaderValue(entry Entry) string {
	if l.Req == nil {
		return ""
	}
	tokens := strings.Split(entry.Operator, ":")
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
