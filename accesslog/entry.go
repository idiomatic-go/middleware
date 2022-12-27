package accesslog

import "strings"

const (
	headerPrefix   = "header"
	directOperator = "direct"
)

// Reference - configuration of logging entries
type Reference struct {
	Operator string
	Name     string
}

// Entry - information for log field generation
type Entry struct {
	Operator    string
	Name        string
	Value       string
	StringValue bool
}

func (e Entry) IsClientHeader() bool {
	return strings.HasPrefix(e.Operator, headerPrefix)
}

func (e Entry) IsDirect() bool {
	return e.Operator == directOperator
}

func NewEntry(operator, name, value string, stringValue bool) Entry {
	return Entry{Operator: operator, Name: name, Value: value, StringValue: stringValue}
}

type directoryT map[string]*Entry

var directory = directoryT{
	TrafficOperator:   &Entry{TrafficOperator, "traffic", "", true},
	StartTimeOperator: &Entry{StartTimeOperator, "start_time", "", true},
	DurationOperator:  &Entry{DurationOperator, "duration_ms", "", false},

	OriginRegionOperator:     &Entry{OriginRegionOperator, "region", "", true},
	OriginZoneOperator:       &Entry{OriginZoneOperator, "zone", "", true},
	OriginSubZoneOperator:    &Entry{OriginSubZoneOperator, "sub_zone", "", true},
	OriginServiceOperator:    &Entry{OriginServiceOperator, "service", "", true},
	OriginInstanceIdOperator: &Entry{OriginInstanceIdOperator, "instance_id", "", true},

	// Route
	RouteNameOperator:       &Entry{RouteNameOperator, "route_name", "", true},
	TimeoutDurationOperator: &Entry{TimeoutDurationOperator, "timeout_ms", "", false},
	RateLimitOperator:       &Entry{RateLimitOperator, "rate_limit_s", "", false},
	RateBurstOperator:       &Entry{RateBurstOperator, "rate_burst", "", false},
	FailoverOperator:        &Entry{FailoverOperator, "failover", "", false},

	// Response
	ResponseStatusCodeOperator:    &Entry{ResponseStatusCodeOperator, "status_code", "", true},
	ResponseBytesReceivedOperator: &Entry{ResponseBytesReceivedOperator, "bytes_received", "", true},
	ResponseBytesSentOperator:     &Entry{ResponseBytesSentOperator, "bytes_sent", "", true},
	StatusFlagsOperator:           &Entry{StatusFlagsOperator, "status_flags", "", true},
	//UpstreamHostOperator:  &Entry{UpstreamHostOperator, "upstream_host", "", true},

	// Request
	RequestProtocolOperator: &Entry{RequestProtocolOperator, "protocol", "", true},
	RequestUrlOperator:      &Entry{RequestUrlOperator, "url", "", true},
	RequestMethodOperator:   &Entry{RequestMethodOperator, "method", "", true},
	RequestPathOperator:     &Entry{RequestPathOperator, "path", "", true},
	RequestHostOperator:     &Entry{RequestHostOperator, "host", "", true},

	RequestIdOperator:           &Entry{RequestIdOperator, "request_id", "", true},
	RequestUserAgentOperator:    &Entry{RequestUserAgentOperator, "user_agent", "", true},
	RequestAuthorityOperator:    &Entry{RequestAuthorityOperator, "authority", "", true},
	RequestForwardedForOperator: &Entry{RequestForwardedForOperator, "forwarded", "", true},

	// gRPC
	GRPCStatusOperator:       &Entry{GRPCStatusOperator, "grpc_status", "", true},
	GRPCStatusNumberOperator: &Entry{GRPCStatusNumberOperator, "grpc_number", "", true},
}
