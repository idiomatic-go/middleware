package accesslog

import (
	"strconv"
	"strings"
	"time"
	"unicode"
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
			// TODO: normalize Http header names. There is code in the net/http package
			var name = ""
			tokens := strings.Split(entry.Operator(), ":")
			if !unicode.IsUpper(rune(tokens[1][0])) {
				var s = string(unicode.ToUpper(rune(tokens[1][0])))
				name = s + tokens[1][1:]
			} else {
				name = tokens[1]
			}
			values := l.Req.Header[name]
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
