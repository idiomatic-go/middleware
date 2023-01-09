package actuator

const (
	WatchAlertStatus   = "watch"
	WarningAlertStatus = "warning"
	CancelStatus       = "cancel"
)

type Event struct {
	SLOName        string
	ActuatorName   string
	Traffic        string // ingress, handler
	SLOCategory    string // latency, status codes, traffic, saturation
	SLOAlertStatus string // watch, warning, canceled
}

func (e Event) IsWatch() bool {
	return e.SLOAlertStatus == WatchAlertStatus
}

func (e Event) IsWarning() bool {
	return e.SLOAlertStatus == WarningAlertStatus
}

func (e Event) IsCancel() bool {
	return e.SLOAlertStatus == CancelStatus
}

func AdjustRateLimiter(act Actuator, percentage int) bool {
	if rlc, ok := act.RateLimiter(); ok {
		rlc.AdjustRateLimiter(percentage)
		return true
	}
	return false
}

func AdjustRetryRateLimiter(act Actuator, percentage int) bool {
	if rc, ok := act.Retry(); ok {
		rc.AdjustRateLimiter(percentage)
		return true
	}
	return false
}
