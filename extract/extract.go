package extract

import (
	"errors"
	"github.com/idiomatic-go/middleware/accesslog"
	"log"
	"net/http"
	neturl "net/url"
	"strings"
)

type LogError func(err error)

var client = http.DefaultClient
var url string
var logError = func(err error) {
	log.Println(err)
}
var c chan *accesslog.Logd
var entries []accesslog.Entry
var config = []accesslog.Reference{
	{Operator: accesslog.StartTimeOperator},
	{Operator: accesslog.DurationOperator},
	{Operator: accesslog.RouteNameOperator},
	{Operator: accesslog.TrafficOperator},

	{Operator: accesslog.RegionOperator},
	{Operator: accesslog.ZoneOperator},
	{Operator: accesslog.SubZoneOperator},
	{Operator: accesslog.ServiceNameOperator},
	{Operator: accesslog.InstanceIdOperator},

	{Operator: accesslog.HttpMethodOperator},
	{Operator: accesslog.AuthorityOperator},
	{Operator: accesslog.PathOperator},
	{Operator: accesslog.ProtocolOperator},
	{Operator: accesslog.RequestIdOperator},
	{Operator: accesslog.ForwardedForOperator},

	{Operator: accesslog.ResponseCodeOperator},
	{Operator: accesslog.ResponseFlagsOperator},
	{Operator: accesslog.BytesReceivedOperator},
}

func Initialize(uri string, newClient *http.Client, fn LogError) error {
	var err error

	if accesslog.IsEmpty(uri) {
		return errors.New("invalid argument : uri is empty")
	}
	u, err1 := neturl.Parse(uri)
	if err1 != nil {
		return err1
	}
	url = u.String()

	err = accesslog.CreateEntries(&entries, config)
	if err != nil {
		return err
	}

	c = make(chan *accesslog.Logd, 100)
	go receive()
	if newClient != nil {
		client = newClient
	}
	if fn != nil {
		logError = fn
	}
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

func sameUrl(l *accesslog.Logd) bool {
	if l == nil || l.Req == nil || l.Req.URL == nil {
		return false
	}
	s := l.Req.URL.String()
	return url == s
}

func do(l *accesslog.Logd) {
	if l == nil {
		logError(errors.New("invalid argument : acessLog data is nil"))
		return
	}
	// Guard against infinite loop
	if sameUrl(l) {
		return
	}
	var req *http.Request
	var err error

	reader := strings.NewReader(accesslog.FormatJson(entries, l))
	req, err = http.NewRequest(http.MethodPost, url, reader)
	if err == nil {
		_, err = client.Do(req)
	}
	if err != nil {
		logError(err)
	}
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			// Exit on a closed channel
			if !open {
				return
			}
			do(msg)
		}
	}
}
