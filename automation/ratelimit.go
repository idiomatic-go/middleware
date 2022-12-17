package automation

import "golang.org/x/time/rate"

const (
	RateLimitName = "ratelimit"
)

type RateLimitController interface {
	//Controller
	Allow() bool
}

type RateLimitConfig struct {
	limit       rate.Limit
	burst       int
	canaryLimit rate.Limit
	canaryBurst int
}

func NewRateLimitConfig(max rate.Limit, burst int, canaryMax rate.Limit, canaryBurst int) *RateLimitConfig {
	c := new(RateLimitConfig)
	return c
}

type rateLimit struct {
	Default int
	current int
	canary  int
}

func (a *rateLimit) Name() string {
	return TimeoutName
}

func (a *rateLimit) IsEnabled() bool {
	return a.current != NilValue
}

func (a *rateLimit) Reset() {

}

func (a *rateLimit) Disable() {
}

func (a *rateLimit) Configure(v ...any) {
}

func (a *rateLimit) Allow() bool {
	return false
}
