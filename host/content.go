package host

type ContentMap map[string][]any

type Credentials func() (username string, password string, err error)

// CreateCredentialsMessage - functions
//func CreateCredentialsMessage(to, from, event string, fn Credentials) Message {
//	return Message{To:to, From:from,Event:event,Status:StatusNotProvided,Content: fn}
//}

func AccessCredentials(msg *Message) Credentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(Credentials); ok {
			return fn
		}
	}
	return nil
}

const (
	DevEnv = iota + 1
	ReviewEnv
	TestEnv
	StageEnv
	ProdEnv
)

type EnvironmentMatcher func(env int) bool

func AccessEnvironmentMatcher(msg *Message) EnvironmentMatcher {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(EnvironmentMatcher); ok {
			return fn
		}
	}
	return nil
}
