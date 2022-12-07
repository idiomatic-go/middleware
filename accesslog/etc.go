package accesslog

import (
	"errors"
	"fmt"
	"strings"
)

type Entry struct {
	Operator string
	Name     string
}

var ingressAttrs []attribute
var egressAttrs []attribute

func AddIngressAttributes(attrs []Entry) error {
	return addAttributes(&ingressAttrs, attrs)
}

func AddEgressAttributes(attrs []Entry) error {
	return addAttributes(&egressAttrs, attrs)
}

func addAttributes(src *[]attribute, attrs []Entry) error {
	if attrs == nil || len(attrs) == 0 {
		return errors.New("log entry slice is empty")
	}
	for _, a := range attrs {
		if a.Operator == "" {
			return errors.New(fmt.Sprintf("invalid operator : operator is empty %v", a.Operator))
		}
		// User defined attribute constant
		if !strings.HasPrefix(a.Operator, operatorPrefix) {
			*src = append(*src, attribute{"direct", a.Operator, a.Name, true})
			continue
		}
		config, ok := defaultConfig[a.Operator]
		if ok {
			newAttr := attribute{config.operator, config.name, "", config.stringValue}
			if a.Name != "" {
				newAttr.name = a.Name
			}
			*src = append(*src, newAttr)
			continue
		}
		if strings.HasPrefix(a.Operator, requestReferencePrefix) {
			*src = append(*src, parseHeaderAttribute(a.Operator[len(requestReferencePrefix):], a))
			continue
		}
		return errors.New(fmt.Sprintf("invalid operator : operator not found or not a valid reference %v", a.Operator))
	}
	return nil
}

func parseHeaderAttribute(s string, entry Entry) attribute {
	tokens := strings.Split(s, ")")
	op := fmt.Sprintf("header:%v", tokens[0])
	if entry.Name == "" {
		return attribute{operator: op, name: tokens[0], value: "", stringValue: true}
	}
	return attribute{operator: op, name: entry.Name, value: "", stringValue: true}
}

type attribute struct {
	operator    string
	name        string
	value       string
	stringValue bool
}

func (a attribute) IsHeader() bool {
	return strings.HasPrefix(a.operator, "header")
}

func (a attribute) IsDirect() bool {
	return a.operator == "direct"
}

type attributes map[string]*attribute

const (
	operatorPrefix          = "%"
	requestReferencePrefix  = "%REQ("
	responseReferencePrefix = "%RESP("

	// Application
	trafficOperator     = "%TRAFFIC%"     //  ingress, egress, ping
	regionOperator      = "%REGION%"      //, origin region
	zoneOperator        = "%ZONE%"        //, origin zone
	subZoneOperator     = "%SUB_ZONE%"    // origin sub zone
	serviceNameOperator = "%SERVICE%"     // origin service
	instanceIdOperator  = "%INSTANCE_ID%" // origin instance id

	// Envoy
	routeNameOperator = "%ROUTE_NAME%" // route name
	startTimeOperator = "%START_TIME%" // start time
	durationOperator  = "%DURATION%"   // Total duration in milliseconds of the request from the start time to the last byte out.

	// Response
	responseCodeOperator  = "%RESPONSE_CODE%"  // HTTP status code
	bytesReceivedOperator = "%BYTES_RECEIVED%" // bytes received
	bytesSentOperator     = "%BYTES_SENT%"     // bytes sent
	responseFlagsOperator = "%RESPONSE_FLAGS%" // response flags
	upstreamHostOperator  = "%UPSTREAM_HOST%"  // Upstream host URL (e.g., tcp://ip:port for TCP connections).

	// Request
	protocolOperator     = "%PROTOCOL%"          // HTTP Protocol
	requestIdOperator    = "%REQ(X-REQUEST-ID)%" // X-REQUEST-ID request header value
	userAgentOperator    = "%REQ(USER-AGENT)%"   // user agent request header value
	authorityOperator    = "%REQ(:AUTHORITY)%"   // authority request header value
	httpMethodOperator   = "%REQ(:METHOD)%"      // HTTP method
	pathOperator         = "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"
	forwardedForOperator = "%REQ(X-FORWARDED-FOR)%" // client IP address (X-FORWARDED-FOR request header value)

	// gRPC
	grpcStatusOperator       = "%GRPC_STATUS(X)%"     // gRPC status code formatted according to the optional parameter X, which can be CAMEL_STRING, SNAKE_STRING and NUMBER. X-REQUEST-ID request header value
	grpcStatusNumberOperator = "%GRPC_STATUS_NUMBER%" // gRPC status code.

	//%REQUESTED_SERVER_NAME%: SNI host
	//%DYNAMIC_METADATA(envoy.lua)%: Apigee dynamic metadata
	//%DOWNSTREAM_TLS_VERSION%: TLS protocol
	//%DOWNSTREAM_DIRECT_REMOTE_ADDRESS%: remote address
	//%RESPONSE_DURATION%: response duration
	//%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%: X-ENVOY-UPSTREAM-SERVICE-TIME response header value
	//%UPSTREAM_CLUSTER%: upstream cluster
	//%RESPONSE_CODE_DETAILS%: HTTP status details
)

var defaultConfig = attributes{
	trafficOperator:     &attribute{trafficOperator, "traffic", "", true},
	regionOperator:      &attribute{regionOperator, "region", "", true},
	zoneOperator:        &attribute{zoneOperator, "zone", "", true},
	subZoneOperator:     &attribute{subZoneOperator, "sub_zone", "", true},
	serviceNameOperator: &attribute{serviceNameOperator, "service", "", true},
	instanceIdOperator:  &attribute{instanceIdOperator, "instance_id", "", true},

	routeNameOperator: &attribute{routeNameOperator, "route_name", "", true},
	startTimeOperator: &attribute{startTimeOperator, "start_time", "", true},
	durationOperator:  &attribute{durationOperator, "duration_ms", "", false},

	// Response
	responseCodeOperator:  &attribute{responseCodeOperator, "status_code", "", true},
	bytesReceivedOperator: &attribute{bytesReceivedOperator, "bytes_received", "", true},
	bytesSentOperator:     &attribute{bytesSentOperator, "bytes_sent", "", true},
	responseFlagsOperator: &attribute{responseFlagsOperator, "response_flags", "", true},
	upstreamHostOperator:  &attribute{upstreamHostOperator, "upstream_host", "", true},
	pathOperator:          &attribute{pathOperator, "path", "", true},
	forwardedForOperator:  &attribute{forwardedForOperator, "forwarded", "", true},

	// Request
	protocolOperator:   &attribute{protocolOperator, "protocol", "", true},
	requestIdOperator:  &attribute{requestIdOperator, "request_id", "", true},
	userAgentOperator:  &attribute{userAgentOperator, "user_agent", "", true},
	authorityOperator:  &attribute{authorityOperator, "authority", "", true},
	httpMethodOperator: &attribute{httpMethodOperator, "method", "", true},

	// gRPC
	grpcStatusOperator:       &attribute{grpcStatusOperator, "grpc_status", "", true},
	grpcStatusNumberOperator: &attribute{grpcStatusNumberOperator, "grpc_number", "", true},
}
