package actuator

type Config struct {
	Name        string
	Host        string
	Timeout     *TimeoutConfig
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfig
	Failover    *FailoverConfig
}
