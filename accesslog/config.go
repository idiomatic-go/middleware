package accesslog

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
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

var ingressEntries []accessdata.Operator
var egressEntries []accessdata.Operator

func CreateIngressEntries(config []accessdata.Operator) error {
	ingressEntries = []accessdata.Operator{}
	return CreateEntries(&ingressEntries, config)
}

func CreateEgressEntries(config []accessdata.Operator) error {
	egressEntries = []accessdata.Operator{}
	return CreateEntries(&egressEntries, config)
}

func CreateEntries(items *[]accessdata.Operator, config []accessdata.Operator) error {
	if items == nil {
		return errors.New("invalid configuration : entries are nil")
	}
	if len(config) == 0 {
		return errors.New("invalid configuration : configuration is empty")
	}
	dup := make(map[string]string)
	for _, op := range config {
		op2, err := createOperator(op)
		if err != nil {
			return err
		}
		if IsEmpty(op2.Value) {
			return errors.New(fmt.Sprintf("invalid operator : operator is invalid %v", op2.Value))
		}
		if IsEmpty(op2.Name) {
			return errors.New(fmt.Sprintf("invalid reference : name is empty %v", op2.Name))
		}
		if _, ok := dup[op2.Value]; ok {
			return errors.New(fmt.Sprintf("invalid reference : name is a duplicate [%v]", op2.Value))
		}
		dup[op2.Value] = op2.Value
		*items = append(*items, op2)
	}
	return nil
}

func createOperator(op accessdata.Operator) (accessdata.Operator, error) {
	if IsEmpty(op.Value) {
		return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator: value is empty %v", op.Name))
	}
	if !strings.HasPrefix(op.Value, operatorPrefix) {
		if IsEmpty(op.Name) {
			return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : name is empty [%v]", op.Value))
		}
		return accessdata.Operator{Name: accessdata.CreateDirect(op.Name), Value: op.Value}, nil
	}
	if op2, ok := accessdata.Operators[op.Value]; ok {
		newOp := accessdata.Operator{Name: op2.Name, Value: op.Value}
		if !IsEmpty(op.Name) {
			newOp.Name = op.Name
		}
		return newOp, nil
	}
	if strings.HasPrefix(op.Value, requestReferencePrefix) {
		return createHeaderOperator(op), nil
	}
	return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : value not found %v", op.Value))
}

func createHeaderOperator(op accessdata.Operator) accessdata.Operator {
	if IsEmpty(op.Value) || !strings.HasPrefix(op.Value, requestReferencePrefix) || len(op.Value) <= len(requestReferencePrefix) {
		return accessdata.Operator{}
	}
	s := op.Value[len(requestReferencePrefix):]
	tokens := strings.Split(s, ")")
	if len(tokens) == 1 || tokens[0] == "" {
		return accessdata.Operator{}
	}
	op1 := fmt.Sprintf("%v:%v", headerPrefix, tokens[0])
	if IsEmpty(op.Name) {
		return accessdata.Operator{Name: tokens[0], Value: op1}
	}
	return accessdata.Operator{Name: op.Name, Value: op1}
}
