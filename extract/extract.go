package extract

import (
	"errors"
	"github.com/idiomatic-go/middleware/accessdata"
	"github.com/idiomatic-go/middleware/accesslog"
	"net/http"
	urlpkg "net/url"
	"strings"
)

type messageHandler func(l *accessdata.Entry) bool

var (
	url string
	c   chan *accessdata.Entry
	//entries []accessdata.Operator
	client                   = http.DefaultClient
	handler   messageHandler = do
	operators                = []accessdata.Operator{
		{Name: "start_time", Value: accessdata.StartTimeOperator},
		{Name: "duration_ms", Value: accessdata.DurationOperator},
		{Name: "traffic", Value: accessdata.TrafficOperator},
		{Name: "route_name", Value: accessdata.RouteNameOperator},

		{Name: "region", Value: accessdata.OriginRegionOperator},
		{Name: "zone", Value: accessdata.OriginZoneOperator},
		{Name: "sub_zone", Value: accessdata.OriginSubZoneOperator},
		{Name: "service", Value: accessdata.OriginServiceOperator},
		{Name: "instance_id", Value: accessdata.OriginInstanceIdOperator},

		{Name: "method", Value: accessdata.RequestMethodOperator},
		{Name: "host", Value: accessdata.RequestHostOperator},
		{Name: "path", Value: accessdata.RequestPathOperator},
		{Name: "protocol", Value: accessdata.RequestProtocolOperator},
		{Name: "request_id", Value: accessdata.RequestIdOperator},
		{Name: "forwarded", Value: accessdata.RequestForwardedForOperator},

		{Name: "status_code", Value: accessdata.ResponseStatusCodeOperator},
		{Name: "status_flags", Value: accessdata.StatusFlagsOperator},
		//{Name: "start_time", Value: accessdata.ResponseBytesReceivedOperator},
		//{}Name: "start_time", Value: accessdata.ResponseBytesSentOperator},

		{Name: "timeout_ms", Value: accessdata.TimeoutDurationOperator},
		{Name: "rate_limit", Value: accessdata.RateLimitOperator},
		{Name: "rate_burst", Value: accessdata.RateBurstOperator},
		{Name: "retry", Value: accessdata.RetryOperator},
		{Name: "retry_rate_limit", Value: accessdata.RetryRateLimitOperator},
		{Name: "retry_rate_burst", Value: accessdata.RetryRateBurstOperator},
		{Name: "failover", Value: accessdata.FailoverOperator},
	}
)

func Initialize(uri string, newClient *http.Client, fn ErrorHandler) error {
	//var err error

	if accesslog.IsEmpty(uri) {
		return errors.New("invalid argument : uri is empty")
	}
	u, err1 := urlpkg.Parse(uri)
	if err1 != nil {
		return err1
	}
	url = u.String()
	//entries = []accessdata.Operator{}
	//err = accesslog.CreateOperators(&entries, config)
	//if err != nil {
	//	return err
	//}
	c = make(chan *accessdata.Entry, 100)
	go receive()
	if newClient != nil {
		client = newClient
	}
	SetErrorHandler(fn)
	//accesslog.SetExtract(extract)
	return nil
}

func Shutdown() {
	if c != nil {
		close(c)
	}
}

func extract(l *accessdata.Entry) {
	if l != nil {
		c <- l
	}
}

func do(l *accessdata.Entry) bool {
	if l == nil {
		OnError(errors.New("invalid argument : access log data is nil"))
		return false
	}
	// let's not extract the extract, the extract, the extract ...
	if l.Url == url {
		return false
	}
	var req *http.Request
	var err error

	reader := strings.NewReader(accessdata.WriteJson(operators, l))
	req, err = http.NewRequest(http.MethodPost, url, reader)
	if err == nil {
		_, err = client.Do(req)
	}
	if err != nil {
		OnError(err)
		return false
	}
	return true
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			// Exit on a closed channel
			if !open {
				return
			}
			handler(msg)
		}
	}
}
