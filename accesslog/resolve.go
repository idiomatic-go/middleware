package accesslog

import (
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	egressTraffic  = "egress"
	ingressTraffic = "ingress"
	pingTraffic    = "ping"
)

type logd struct {
	traffic      string
	start        time.Time
	duration     time.Duration
	bytesWritten int
	route        *route.Route
	req          *http.Request
	resp         *http.Response
	err          error
	code         int
}

func (l logd) isIngress() bool {
	return l.traffic == ingressTraffic
}

func (l logd) isEgress() bool {
	return l.traffic == egressTraffic
}

func (l logd) isPing() bool {
	return l.isIngress() && l.route.Ping
}

func (l logd) resolve(attr attribute) string {
	if attr.IsHeader() {
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
			return pingTraffic
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
