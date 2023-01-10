package resource

import (
	"github.com/idiomatic-go/middleware/actuator"
)

type Route struct {
	Name        string
	Timeout     actuator.TimeoutConfig
	RateLimiter actuator.RateLimiterConfig
	Retry       actuator.RetryConfig
}
