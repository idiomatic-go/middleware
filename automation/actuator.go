package automation

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type Actuator interface {
	Action(name string) Action
	Actuate(events string) error
}
