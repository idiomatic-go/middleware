package automation

import "golang.org/x/time/rate"

const (
	RateLimitName = "ratelimit"
)

type RateLimitAction interface {
	Allow() bool
}

type RateLimitConfig struct {
	limit rate.Limit
	burst int
}

type rateLimitAction struct {
	Default int
	current int
	canary  int
}

func (a *rateLimitAction) Name() string {
	return TimeoutName
}

func (a *rateLimitAction) IsEnabled() bool {
	return a.current != NilValue
}

func (a *rateLimitAction) Reset() {

}

func (a *rateLimitAction) Disable() {
}

func (a *rateLimitAction) Configure(v ...any) {
}

func (a *rateLimitAction) Allow() bool {
	return false
}
