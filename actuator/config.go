package actuator

type Config struct {
	Name        string
	Timeout     *TimeoutConfig
	RateLimiter *RateLimiterConfig
	Retry       *RetryConfig
	Failover    *FailoverConfig
}

type ConfigList struct {
	Package string
	Config  Config
}
