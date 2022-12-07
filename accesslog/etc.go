package accesslog

import (
	"errors"
	"fmt"
	"strings"
)

var ingressAttrs []attribute
var egressAttrs []attribute

func AddIngressAttributes(attrs []Entry) error {
	return addAttributes(&ingressAttrs, attrs)
}

func AddEgressAttributes(attrs []Entry) error {
	return addAttributes(&egressAttrs, attrs)
}

func addAttributes(attrs *[]attribute, config []Entry) error {
	if len(config) == 0 {
		return errors.New("invalid : log entry configuration is empty")
	}
	dup := map[string]string{}
	for _, entry := range config {
		newAttr, err := createAttribute(entry)
		if err != nil {
			return err
		}
		if newAttr.operator == "" {
			return errors.New(fmt.Sprintf("invalid entry : operator is invalid %v", entry.Operator))
		}
		if newAttr.name == "" {
			return errors.New(fmt.Sprintf("invalid entry : name is empty %v", entry.Operator))
		}
		if _, ok := dup[newAttr.name]; ok {
			return errors.New(fmt.Sprintf("invalid entry : name is a duplicate %v", newAttr.name))
		}
		dup[newAttr.name] = ""
		*attrs = append(*attrs, newAttr)
	}
	return nil
}

func createAttribute(entry Entry) (attribute, error) {
	if entry.Operator == "" {
		return attribute{}, errors.New(fmt.Sprintf("invalid entry : operator is empty %v", entry.Operator))
	}
	if !strings.HasPrefix(entry.Operator, operatorPrefix) {
		return attribute{directOperator, entry.Operator, entry.Name, true}, nil
	}
	if config, ok := directory[entry.Operator]; ok {
		newAttr := attribute{config.operator, config.name, "", config.stringValue}
		if entry.Name != "" {
			newAttr.name = entry.Name
		}
		return newAttr, nil
	}
	if strings.HasPrefix(entry.Operator, requestReferencePrefix) {
		return parseHeaderAttribute(entry), nil
	}
	return attribute{}, errors.New(fmt.Sprintf("invalid operator : operator not found or not a valid reference %v", entry.Operator))
}

func parseHeaderAttribute(entry Entry) attribute {
	if entry.Operator == "" || !strings.HasPrefix(entry.Operator, requestReferencePrefix) || len(entry.Operator) <= len(requestReferencePrefix) {
		return attribute{}
	}
	s := entry.Operator[len(requestReferencePrefix):]
	tokens := strings.Split(s, ")")
	if len(tokens) == 1 || tokens[0] == "" {
		return attribute{}
	}
	op := fmt.Sprintf("%v:%v", headerPrefix, tokens[0])
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

func (a attribute) isHeader() bool {
	return strings.HasPrefix(a.operator, headerPrefix)
}

func (a attribute) isDirect() bool {
	return a.operator == directOperator
}

type attributes map[string]*attribute

const (
	headerPrefix            = "header"
	directOperator          = "direct"
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

var directory = attributes{
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
