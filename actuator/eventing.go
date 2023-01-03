package actuator

// Need to create an SLO based on status flags.
//An agonist is a chemical that activates a receptor to produce a biological response.

// Muscles in the torso, arms, and legs are arranged in opposing pairs. The main muscle that's moving is called the primer, or
// agonist. For example, if you pick up a coffee mug from the table, the agonist muscle is your bicep. The opposing muscle is
// the triceps, which is referred to as the antagonist.

type Event struct {
	SLOName        string
	ActuatorName   string
	Traffic        string // ingress, egress
	SLOCategory    string // latency, status codes, traffic, saturation
	SLOAlertStatus string // watch, warning, canceled

}
