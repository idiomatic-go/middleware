package host

import (
	"fmt"
	"os"
	"strings"
)

const (
	runtimeEnvKey  = "RUNTIME_ENV"
	devEnvValue    = "dev"
	reviewEnvValue = "review"
	testEnvValue   = "test"
	stageEnvValue  = "stage"
	prodEnvValue   = "prod"
)

func ExampleAccess_Credentials() {

	msg := Message{To: "to-uri", From: "from-uri", Content: []any{
		"text content",
		500,
		Credentials(func() (username, password string, err error) { return "", "", nil }),
	}}

	fmt.Printf("test: AccessCredentials(nil) -> %v\n", AccessCredentials(nil) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", AccessCredentials(&Message{To: "to-uri"}) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", AccessCredentials(&msg) != nil)

	//Output:
	//test: AccessCredentials(nil) -> false
	//test: AccessCredentials(msg) -> false
	//test: AccessCredentials(msg) -> true
}

func ExampleAccess_IsEnvironment() {

	msg := Message{To: "to-uri", From: "from-uri", Content: []any{
		"text content",
		500,
		EnvironmentMatcher(func(env int) bool {
			s := os.Getenv(runtimeEnvKey)
			switch env {
			case DevEnv:
				return s == "" || strings.EqualFold(s, devEnvValue)
			case ReviewEnv:
				return strings.EqualFold(s, reviewEnvValue)
			case TestEnv:
				return strings.EqualFold(s, testEnvValue)
			case StageEnv:
				return strings.EqualFold(s, stageEnvValue)
			case ProdEnv:
				return strings.EqualFold(s, prodEnvValue)
			}
			return false
		}),
	}}

	fmt.Printf("test: AccessEnvironmentMatcher(nil) -> %v\n", AccessEnvironmentMatcher(nil) != nil)
	fmt.Printf("test: AccessEnvironmentMatcher(msg) -> %v\n", AccessEnvironmentMatcher(&Message{To: "to-uri"}) != nil)
	fmt.Printf("test: AccessEnvironmentMatcher(msg) -> %v\n", AccessEnvironmentMatcher(&msg) != nil)

	//Output:
	//test: AccessEnvironmentMatcher(nil) -> false
	//test: AccessEnvironmentMatcher(msg) -> false
	//test: AccessEnvironmentMatcher(msg) -> true
}
