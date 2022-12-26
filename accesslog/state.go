package accesslog

import (
	"fmt"
	"golang.org/x/time/rate"
	"strings"
)

type ControllerState struct {
	Enabled bool
	Tags    []string
}

func (c ControllerState) Value(index int) string {
	if len(c.Tags) == 0 || index < 0 || index >= len(c.Tags) {
		return ""
	}
	tokens := strings.Split(c.Tags[index], ":")
	if len(tokens) <= 1 {
		return tokens[0]
	}
	return tokens[1]
}

type ActuatorState struct {
	Name         string
	WriteIngress bool
	WriteEgress  bool
	Timeout      ControllerState
	RateLimiter  ControllerState
	Failover     ControllerState
}

func NewActuatorStateWithTimeout(name string, timeout int) ActuatorState {
	s := fmt.Sprintf("%v:%v", "duration", timeout)
	return ActuatorState{Name: name, WriteEgress: true, WriteIngress: true, Timeout: ControllerState{Enabled: true, Tags: []string{s}}}
}

func NewActuatorStateWithRateLimiter(name string, limit rate.Limit, burst int) ActuatorState {
	var s string
	if limit == rate.Inf {
		s = fmt.Sprintf("%v:%v", "rateLimit", -1)
	} else {
		s = fmt.Sprintf("%v:%v", "rateLimit", limit)
	}
	s1 := fmt.Sprintf("%v:%v", "rateBurst", burst)
	return ActuatorState{Name: name, WriteEgress: true, WriteIngress: true, RateLimiter: ControllerState{Enabled: true, Tags: []string{s, s1}}}
}

func NewActuatorStateWithFailover(name string, failover bool) ActuatorState {
	s := fmt.Sprintf("%v:%v", "failover", failover)
	return ActuatorState{Name: name, WriteEgress: true, WriteIngress: true, Failover: ControllerState{Enabled: true, Tags: []string{s}}}

}
