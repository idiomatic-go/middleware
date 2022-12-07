package accesslog

import (
	"strconv"
	"strings"
	"time"
)

func (l Logd) isIngress() bool {
	return l.traffic == IngressTraffic
}

func (l Logd) isEgress() bool {
	return l.traffic == EgressTraffic
}

func (l Logd) isPing() bool {
	return l.isIngress() && l.route.Ping
}

func (l Logd) value(attr attribute) string {
	if attr.isHeader() {
		if l.req != nil {
			tokens := strings.Split(attr.operator, ":")
			values := l.req.Header[tokens[0]]
			if values != nil {
				return values[0]
			}
			return ""
		}
		return ""
	}
	switch attr.operator {
	case trafficOperator:
		if l.isPing() {
			return PingTraffic
		}
		return l.traffic
	case regionOperator:
		return origin.Region
	case zoneOperator:
		return origin.Zone
	case subZoneOperator:
		return origin.SubZone
	case serviceNameOperator:
		return origin.Service
	case instanceIdOperator:
		return origin.InstanceId

	case startTimeOperator:
		return FmtTimestamp(l.start)
	case routeNameOperator:
		return l.route.Name
	case durationOperator:
		d := int(l.duration / time.Duration(1e6))
		return strconv.Itoa(d)

	case httpMethodOperator:
		return l.req.Method
	case responseCodeOperator:
		if l.isIngress() {
			return strconv.Itoa(l.code)
		} else {
			return strconv.Itoa(l.resp.StatusCode)
		}
	}
	return ""
}
