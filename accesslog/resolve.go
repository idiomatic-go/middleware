package accesslog

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"strconv"
	"time"
)

const (
	egressTraffic  = "egress"
	ingressTraffic = "ingress"
	pingTraffic    = "ping"
	startTimeName  = "start_time"
	regionName     = "region"
	zoneName       = "zone"
	subZoneName    = "sub_zone"
	serviceName    = "service"
	instanceIdName = "instance_id"
	trafficName    = "traffic"
	routeName      = "route_name"
	durationName   = "duration_ms"
	urlName        = "url"
	methodName     = "method"
	statusCodeName = "status_code"
	protocolName   = "proto"
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

func (l logd) IsIngress() bool {
	return l.traffic == ingressTraffic
}

func (l logd) IsEgress() bool {
	return l.traffic == egressTraffic
}

func (l logd) IsPing() bool {
	return l.IsIngress() && l.route.Ping
}

func resolve(name string, data *logd) (string, string, error) {
	if data == nil {
		return "", "", errors.New("data is nil")
	}
	switch name {
	case startTimeName:
		return FmtTimestamp(data.start), markupString, nil
	case regionName:
		return origin.Region, markupString, nil
	case zoneName:
		return origin.Zone, markupString, nil
	case subZoneName:
		return origin.SubZone, markupString, nil
	case serviceName:
		return origin.Service, markupString, nil
	case instanceIdName:
		return origin.InstanceId, markupString, nil
	case trafficName:
		if data.IsPing() {
			return pingTraffic, markupString, nil
		}
		return data.traffic, markupString, nil
	case routeName:
		return data.route.Name, markupString, nil
	case durationName:
		d := int(data.duration / time.Duration(1e6))
		return strconv.Itoa(d), markupValue, nil
	case urlName:
		return data.req.URL.String(), markupString, nil
	case methodName:
		return data.req.Method, markupString, nil
	case statusCodeName:
		if data.IsIngress() {
			return strconv.Itoa(data.code), markupValue, nil
		} else {
			return strconv.Itoa(data.resp.StatusCode), markupValue, nil
		}
	}
	return "", "", errors.New(fmt.Sprintf("INVALID REFERENCE:%v", name))
}
