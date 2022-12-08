package accesslog

import (
	"errors"
	"fmt"
	"strings"
)

var ingressEntries []Entry
var egressEntries []Entry

func CreateIngressEntries(config []Reference) error {
	return CreateEntries(&ingressEntries, config)
}

func CreateEgressEntries(config []Reference) error {
	return CreateEntries(&egressEntries, config)
}

func CreateEntries(items *[]Entry, config []Reference) error {
	if items == nil {
		return errors.New("invalid configuration : entries are nil")
	}
	if len(config) == 0 {
		return errors.New("invalid configuration : configuration is empty")
	}
	dup := make(map[string]string)
	for _, ref := range config {
		entry, err := createEntry(ref)
		if err != nil {
			return err
		}
		if IsEmpty(entry.Operator()) {
			return errors.New(fmt.Sprintf("invalid reference : operator is invalid %v", ref.Operator))
		}
		if IsEmpty(entry.Ref.Name) {
			return errors.New(fmt.Sprintf("invalid reference : name is empty %v", ref.Operator))
		}
		if _, ok := dup[entry.Name()]; ok {
			return errors.New(fmt.Sprintf("invalid reference : name is a duplicate [%v]", entry.Name()))
		}
		dup[entry.Name()] = entry.Name()
		*items = append(*items, entry)
	}
	return nil
}

func createEntry(ref Reference) (Entry, error) {
	if IsEmpty(ref.Operator) {
		return Entry{}, errors.New(fmt.Sprintf("invalid entry reference : operator is empty %v", ref.Operator))
	}
	if !strings.HasPrefix(ref.Operator, operatorPrefix) {
		if IsEmpty(ref.Name) {
			return Entry{}, errors.New(fmt.Sprintf("invalid entry reference : name is empty [operator=%v]", ref.Operator))
		}
		return NewEntry(directOperator, ref.Operator, ref.Name, true), nil
	}
	if entry, ok := directory[ref.Operator]; ok {
		item := NewEntry(entry.Operator(), entry.Name(), "", entry.StringValue)
		if !IsEmpty(ref.Name) {
			item.Ref.Name = ref.Name
		}
		return item, nil
	}
	if strings.HasPrefix(ref.Operator, requestReferencePrefix) {
		return createHeaderEntry(ref), nil
	}
	return Entry{}, errors.New(fmt.Sprintf("invalid entry reference : operator not found %v", ref.Operator))
}

func createHeaderEntry(ref Reference) Entry {
	if IsEmpty(ref.Operator) || !strings.HasPrefix(ref.Operator, requestReferencePrefix) || len(ref.Operator) <= len(requestReferencePrefix) {
		return Entry{}
	}
	s := ref.Operator[len(requestReferencePrefix):]
	tokens := strings.Split(s, ")")
	if len(tokens) == 1 || tokens[0] == "" {
		return Entry{}
	}
	op := fmt.Sprintf("%v:%v", headerPrefix, tokens[0])
	if IsEmpty(ref.Name) {
		return NewEntry(op, tokens[0], "", true)
	}
	return NewEntry(op, ref.Name, "", true)
}

const (
	headerPrefix            = "header"
	directOperator          = "direct"
	operatorPrefix          = "%"
	requestReferencePrefix  = "%REQ("
	responseReferencePrefix = "%RESP("

	// Application
	TrafficOperator     = "%TRAFFIC%"     //  ingress, egress, ping
	RegionOperator      = "%REGION%"      //, origin region
	ZoneOperator        = "%ZONE%"        //, origin zone
	SubZoneOperator     = "%SUB_ZONE%"    // origin sub zone
	ServiceNameOperator = "%SERVICE%"     // origin service
	InstanceIdOperator  = "%INSTANCE_ID%" // origin instance id

	// Envoy
	RouteNameOperator = "%ROUTE_NAME%" // route name
	StartTimeOperator = "%START_TIME%" // start time
	DurationOperator  = "%DURATION%"   // Total duration in milliseconds of the request from the start time to the last byte out.

	// Response
	ResponseCodeOperator  = "%RESPONSE_CODE%"  // HTTP status code
	BytesReceivedOperator = "%BYTES_RECEIVED%" // bytes received
	BytesSentOperator     = "%BYTES_SENT%"     // bytes sent
	ResponseFlagsOperator = "%RESPONSE_FLAGS%" // response flags
	UpstreamHostOperator  = "%UPSTREAM_HOST%"  // Upstream host URL (e.g., tcp://ip:port for TCP connections).

	// Request
	ProtocolOperator     = "%PROTOCOL%"          // HTTP Protocol
	RequestIdOperator    = "%REQ(X-REQUEST-ID)%" // X-REQUEST-ID request header value
	UserAgentOperator    = "%REQ(USER-AGENT)%"   // user agent request header value
	AuthorityOperator    = "%REQ(:AUTHORITY)%"   // authority request header value
	HttpMethodOperator   = "%REQ(:METHOD)%"      // HTTP method
	PathOperator         = "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"
	ForwardedForOperator = "%REQ(X-FORWARDED-FOR)%" // client IP address (X-FORWARDED-FOR request header value)

	// gRPC
	GRPCStatusOperator       = "%GRPC_STATUS(X)%"     // gRPC status code formatted according to the optional parameter X, which can be CAMEL_STRING, SNAKE_STRING and NUMBER. X-REQUEST-ID request header value
	GRPCStatusNumberOperator = "%GRPC_STATUS_NUMBER%" // gRPC status code.

	//%REQUESTED_SERVER_NAME%: SNI host
	//%DYNAMIC_METADATA(envoy.lua)%: Apigee dynamic metadata
	//%DOWNSTREAM_TLS_VERSION%: TLS protocol
	//%DOWNSTREAM_DIRECT_REMOTE_ADDRESS%: remote address
	//%RESPONSE_DURATION%: response duration
	//%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%: X-ENVOY-UPSTREAM-SERVICE-TIME response header value
	//%UPSTREAM_CLUSTER%: upstream cluster
	//%RESPONSE_CODE_DETAILS%: HTTP status details
)

var directory = Directory{
	TrafficOperator: &Entry{Reference{TrafficOperator, "traffic"}, "", true},

	RegionOperator:      &Entry{Reference{RegionOperator, "region"}, "", true},
	ZoneOperator:        &Entry{Reference{ZoneOperator, "zone"}, "", true},
	SubZoneOperator:     &Entry{Reference{SubZoneOperator, "sub_zone"}, "", true},
	ServiceNameOperator: &Entry{Reference{ServiceNameOperator, "service"}, "", true},
	InstanceIdOperator:  &Entry{Reference{InstanceIdOperator, "instance_id"}, "", true},

	RouteNameOperator: &Entry{Reference{RouteNameOperator, "route_name"}, "", true},
	StartTimeOperator: &Entry{Reference{StartTimeOperator, "start_time"}, "", true},
	DurationOperator:  &Entry{Reference{DurationOperator, "duration_ms"}, "", false},

	// Response
	ResponseCodeOperator:  &Entry{Reference{ResponseCodeOperator, "status_code"}, "", true},
	BytesReceivedOperator: &Entry{Reference{BytesReceivedOperator, "bytes_received"}, "", true},
	BytesSentOperator:     &Entry{Reference{BytesSentOperator, "bytes_sent"}, "", true},
	ResponseFlagsOperator: &Entry{Reference{ResponseFlagsOperator, "response_flags"}, "", true},
	UpstreamHostOperator:  &Entry{Reference{UpstreamHostOperator, "upstream_host"}, "", true},
	PathOperator:          &Entry{Reference{PathOperator, "path"}, "", true},
	ForwardedForOperator:  &Entry{Reference{ForwardedForOperator, "forwarded"}, "", true},

	// Request
	ProtocolOperator:   &Entry{Reference{ProtocolOperator, "protocol"}, "", true},
	RequestIdOperator:  &Entry{Reference{RequestIdOperator, "request_id"}, "", true},
	UserAgentOperator:  &Entry{Reference{UserAgentOperator, "user_agent"}, "", true},
	AuthorityOperator:  &Entry{Reference{AuthorityOperator, "authority"}, "", true},
	HttpMethodOperator: &Entry{Reference{HttpMethodOperator, "method"}, "", true},

	// gRPC
	GRPCStatusOperator:       &Entry{Reference{GRPCStatusOperator, "grpc_status"}, "", true},
	GRPCStatusNumberOperator: &Entry{Reference{GRPCStatusNumberOperator, "grpc_number"}, "", true},
}
