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
		if l.Req != nil {
			tokens := strings.Split(entry.Operator(), ":")
			values := l.Req.Header[tokens[0]]
			if values != nil {
				return values[0]
			}
			return ""
		}
		return ""
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
	case RouteNameOperator:
		return l.Route.Name
	case DurationOperator:
		d := int(l.Duration / time.Duration(1e6))
		return strconv.Itoa(d)

	case HttpMethodOperator:
		return l.Req.Method
	case ResponseCodeOperator:
		if l.IsIngress() {
			return strconv.Itoa(l.Code)
		} else {
			return strconv.Itoa(l.Resp.StatusCode)
		}
	}
	return ""
}
