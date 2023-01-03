package actuator

type Event struct {
	SLOName        string
	ActuatorName   string
	Traffic        string // ingress, egress
	SLOCategory    string // latency, status codes, traffic, saturation
	SLOAlertStatus string // watch, warning, canceled

}
