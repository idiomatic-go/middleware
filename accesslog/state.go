package accesslog

import "strings"

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
