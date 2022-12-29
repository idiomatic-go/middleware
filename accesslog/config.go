package accesslog

import (
	"errors"
	"fmt"
	"strings"
)

const (
	operatorPrefix          = "%"
	requestReferencePrefix  = "%REQ("
	RequestIdHeaderName     = "X-REQUEST-ID"
	UserAgentHeaderName     = "USER-AGENT"
	FordwardedForHeaderName = "X-FORWARDED-FOR"

	TrafficOperator   = "%TRAFFIC%"    //  ingress, egress, ping
	StartTimeOperator = "%START_TIME%" // start time
	DurationOperator  = "%DURATION%"   // Total duration in milliseconds of the request from the start time to the last byte out.

	OriginRegionOperator     = "%REGION%"      // origin region
	OriginZoneOperator       = "%ZONE%"        // origin zone
	OriginSubZoneOperator    = "%SUB_ZONE%"    // origin sub zone
	OriginServiceOperator    = "%SERVICE%"     // origin service
	OriginInstanceIdOperator = "%INSTANCE_ID%" // origin instance id

	RouteNameOperator       = "%ROUTE_NAME%"
	TimeoutDurationOperator = "%TIMEOUT_DURATION%"
	RateLimitOperator       = "%RATE_LIMIT%"
	RateBurstOperator       = "%RATE_BURST%"
	RetryOperator           = "%RETRY"
	RetryRateLimitOperator  = "%RETRY_RATE_LIMIT%"
	RetryRateBurstOperator  = "%RETRY_RATE_BURST%"
	FailoverOperator        = "%FAILOVER%"

	ResponseStatusCodeOperator    = "%STATUS_CODE%"    // HTTP status code
	ResponseBytesReceivedOperator = "%BYTES_RECEIVED%" // bytes received
	ResponseBytesSentOperator     = "%BYTES_SENT%"     // bytes sent
	StatusFlagsOperator           = "%STATUS_FLAGS%"   // status flags
	//UpstreamHostOperator  = "%UPSTREAM_HOST%"  // Upstream host URL (e.g., tcp://ip:port for TCP connections).

	RequestProtocolOperator = "%PROTOCOL%" // HTTP Protocol
	RequestMethodOperator   = "%METHOD%"   // HTTP method
	RequestUrlOperator      = "%URL%"
	RequestPathOperator     = "%PATH%"
	RequestHostOperator     = "%HOST%"

	RequestIdOperator           = "%X-REQUEST-ID%"    // X-REQUEST-ID request header value
	RequestUserAgentOperator    = "%USER-AGENT%"      // user agent request header value
	RequestAuthorityOperator    = "%AUTHORITY%"       // authority request header value
	RequestForwardedForOperator = "%X-FORWARDED-FOR%" // client IP address (X-FORWARDED-FOR request header value)

	GRPCStatusOperator       = "%GRPC_STATUS(X)%"     // gRPC status code formatted according to the optional parameter X, which can be CAMEL_STRING, SNAKE_STRING and NUMBER. X-REQUEST-ID request header value
	GRPCStatusNumberOperator = "%GRPC_STATUS_NUMBER%" // gRPC status code.

)

var ingressEntries []Entry
var egressEntries []Entry

func CreateIngressEntries(config []Reference) error {
	ingressEntries = []Entry{}
	return CreateEntries(&ingressEntries, config)
}

func CreateEgressEntries(config []Reference) error {
	egressEntries = []Entry{}
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
		if IsEmpty(entry.Operator) {
			return errors.New(fmt.Sprintf("invalid reference : operator is invalid %v", ref.Operator))
		}
		if IsEmpty(entry.Name) {
			return errors.New(fmt.Sprintf("invalid reference : name is empty %v", ref.Operator))
		}
		if _, ok := dup[entry.Name]; ok {
			return errors.New(fmt.Sprintf("invalid reference : name is a duplicate [%v]", entry.Name))
		}
		dup[entry.Name] = entry.Name
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
		item := NewEntry(entry.Operator, entry.Name, "", entry.StringValue)
		if !IsEmpty(ref.Name) {
			item.Name = ref.Name
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
