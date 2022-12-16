package automation

type Action interface {
	Name() string
	IsEnabled() bool
	Reset()
	Disable()
	Configure(v ...any)
}

// different actions
// logging
