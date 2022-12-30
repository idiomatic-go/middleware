package actuator

// Need to create an SLO based on status flags.
//An agonist is a chemical that activates a receptor to produce a biological response.

// Muscles in the torso, arms, and legs are arranged in opposing pairs. The main muscle that's moving is called the primer, or
// agonist. For example, if you pick up a coffee mug from the table, the agonist muscle is your bicep. The opposing muscle is
// the triceps, which is referred to as the antagonist.

type Event struct {
	// Identity
	SLOName      string
	ActuatorName string
	Traffic      string

	// Classification
	SLOCategory string //latency, status codes, traffic, saturation-latency, saturation-utilization, counter, profile

	LifeCycleStates string // Status of the SLO alerting  : watch,warning,canceled

	// Context specific attributes
	RPS        int // Used to set an appropriate RateLimiter
	StatusCode int // Specifically which status code caused the notification
	// Could the set of actions be : enable, disable only for retry, increase decrease
}

// Saturation - Latency SLO on the Host route
// Ingress Only - on the host route
// With  a saturation-latency - host RateLimiter should be adjusted, based on the host current RPS
// There needs to be a counter SLO on status code 429, too many requests
// There is only action when a watch or warning is issued, a canceled will not affect an action

// Latency
// Ingress only on host route, not applicable on other ingress routes as their latency is related to upstream traffic
// Egress - 1. Determine if the latency is due to upstream saturation, or transient latency.
//             - To do this enable Retry
//               actions are to enable, disable, amplify and attenuate
//               Need an opposing SLO on the existence of too many retries
//
// if the latency is saturation, then we will increase, decrease te rate limiter, with a corresponding SLO on status codes

// Status Codes
// Egress Only - if the code is 429, then rate limiting, 503, 504, rate limiting
//
