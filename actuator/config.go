package actuator

type Route struct {
	Name        string
	Host        string
	Timeout     *TimeoutConfig
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfig
	Failover    *FailoverConfig
}
