package extract

import (
	"errors"
	"github.com/idiomatic-go/middleware/accesslog"
	"net/http"
	urlpkg "net/url"
	"strings"
)

type messageHandler func(l *accesslog.Logd) bool

var (
	url     string
	c       chan *accesslog.Logd
	entries []accesslog.Entry
	client                 = http.DefaultClient
	handler messageHandler = do
	config                 = []accesslog.Reference{
		{Operator: accesslog.StartTimeOperator},
		{Operator: accesslog.DurationOperator},
		{Operator: accesslog.TrafficOperator},
		{Operator: accesslog.RouteNameOperator},

		{Operator: accesslog.OriginRegionOperator},
		{Operator: accesslog.OriginZoneOperator},
		{Operator: accesslog.OriginSubZoneOperator},
		{Operator: accesslog.OriginServiceOperator},
		{Operator: accesslog.OriginInstanceIdOperator},

		{Operator: accesslog.RequestMethodOperator},
		{Operator: accesslog.RequestHostOperator},
		{Operator: accesslog.RequestPathOperator},
		{Operator: accesslog.RequestProtocolOperator},
		{Operator: accesslog.RequestIdOperator},
		{Operator: accesslog.RequestForwardedForOperator},

		{Operator: accesslog.ResponseStatusCodeOperator},
		{Operator: accesslog.ResponseFlagsOperator},
		{Operator: accesslog.ResponseBytesReceivedOperator},
		{Operator: accesslog.ResponseBytesSentOperator},

		{Operator: accesslog.RouteTimeoutOperator},
		{Operator: accesslog.RouteLimitOperator},
		{Operator: accesslog.RouteBurstOperator},
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
	entries = []accesslog.Entry{}
	err = accesslog.CreateEntries(&entries, config)
	if err != nil {
		return err
	}
	c = make(chan *accesslog.Logd, 100)
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

func extract(l *accesslog.Logd) {
	if l != nil {
		c <- l
	}
}

func do(l *accesslog.Logd) bool {
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

	reader := strings.NewReader(accesslog.FormatJson(entries, l))
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
