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
	url     string
	c       chan *accessdata.Entry
	entries []accessdata.Operator
	client                 = http.DefaultClient
	handler messageHandler = do
	config                 = []accessdata.Operator{
		{Value: accessdata.StartTimeOperator},
		{Value: accessdata.DurationOperator},
		{Value: accessdata.TrafficOperator},
		{Value: accessdata.RouteNameOperator},

		{Value: accessdata.OriginRegionOperator},
		{Value: accessdata.OriginZoneOperator},
		{Value: accessdata.OriginSubZoneOperator},
		{Value: accessdata.OriginServiceOperator},
		{Value: accessdata.OriginInstanceIdOperator},

		{Value: accessdata.RequestMethodOperator},
		{Value: accessdata.RequestHostOperator},
		{Value: accessdata.RequestPathOperator},
		{Value: accessdata.RequestProtocolOperator},
		{Value: accessdata.RequestIdOperator},
		{Value: accessdata.RequestForwardedForOperator},

		{Value: accessdata.ResponseStatusCodeOperator},
		{Value: accessdata.StatusFlagsOperator},
		{Value: accessdata.ResponseBytesReceivedOperator},
		{Value: accessdata.ResponseBytesSentOperator},

		{Value: accessdata.TimeoutDurationOperator},
		{Value: accessdata.RateLimitOperator},
		{Value: accessdata.RateBurstOperator},
		{Value: accessdata.RetryOperator},
		{Value: accessdata.RetryRateLimitOperator},
		{Value: accessdata.RetryRateBurstOperator},
		{Value: accessdata.FailoverOperator},
	}
)

func Initialize(uri string, newClient *http.Client, fn ErrorHandler) error {
	var err error

	if accesslog.IsEmpty(uri) {
		return errors.New("invalid argument : uri is empty")
	}
	u, err1 := urlpkg.Parse(uri)
	if err1 != nil {
		return err1
	}
	url = u.String()
	entries = []accessdata.Operator{}
	err = accesslog.CreateEntries(&entries, config)
	if err != nil {
		return err
	}
	c = make(chan *accessdata.Entry, 100)
	go receive()
	if newClient != nil {
		client = newClient
	}
	SetErrorHandler(fn)
	accesslog.SetExtract(extract)
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

	reader := strings.NewReader(accessdata.WriteJson(entries, l))
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
