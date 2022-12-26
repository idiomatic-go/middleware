package accesslog

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	errorNilRouteFmt = "{\"error\": \"%v route name is empty\"}"
	errorEmptyFmt    = "{\"error\": \"%v log entries are empty\"}"
)

func WriteEgress(start time.Time, duration time.Duration, act ActuatorState, req *http.Request, resp *http.Response, responseFlags string) {
	if act.Name == "" {
		egressWrite(fmt.Sprintf(errorNilRouteFmt, EgressTraffic))
		return
	}
	data := NewLogd(EgressTraffic, start, duration, getOrigin(), act, req, resp, responseFlags)
	callExtract(data)
	if !opt.writeEgress {
		return
	}
	if len(egressEntries) == 0 {
		egressWrite(fmt.Sprintf(errorEmptyFmt, EgressTraffic))
		return
	}
	s := FormatJson(egressEntries, data)
	egressWrite(s)
}

func WriteIngress(start time.Time, duration time.Duration, act ActuatorState, req *http.Request, resp *http.Response, responseFlags string) {
	if act.Name == "" {
		ingressWrite(fmt.Sprintf(errorNilRouteFmt, IngressTraffic))
		return
	}
	data := NewLogd(IngressTraffic, start, duration, getOrigin(), act, req, resp, responseFlags)
	//data.StatusCode = code
	//data.BytesSent = bytesSent
	callExtract(data)
	if !opt.writeIngress {
		return
	}
	if len(ingressEntries) == 0 {
		ingressWrite(fmt.Sprintf(errorEmptyFmt, IngressTraffic))
		return
	}
	s := FormatJson(ingressEntries, data)
	ingressWrite(s)
}

func FormatJson(items []Entry, data *Logd) string {
	if len(items) == 0 || data == nil {
		return "{}"
	}
	sb := strings.Builder{}
	for _, entry := range items {
		if entry.IsDirect() {
			writeJson(&sb, entry.Name, entry.Value, entry.StringValue)
			continue
		}
		writeJson(&sb, entry.Name, data.Value(entry), entry.StringValue)
	}
	sb.WriteString("}")
	return sb.String()
}
