package host

import (
	"github.com/idiomatic-go/middleware/host"
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

func GetRuntime() string {
	s := os.Getenv(runtimeEnvKey)
	if s == "" {
		return devEnvValue
	}
	return s
}

func matchEnvironment(env int) bool {
	s := GetRuntime()
	switch env {
	case host.DevEnv:
		return s == "" || strings.EqualFold(s, devEnvValue)
	case host.ReviewEnv:
		return strings.EqualFold(s, reviewEnvValue)
	case host.TestEnv:
		return strings.EqualFold(s, testEnvValue)
	case host.StageEnv:
		return strings.EqualFold(s, stageEnvValue)
	case host.ProdEnv:
		return strings.EqualFold(s, prodEnvValue)
	}
	return false
}

func isDevEnv() bool {
	return matchEnvironment(host.DevEnv)
}

func isReviewEnv() bool {
	return matchEnvironment(host.ReviewEnv)
}

func isTestEnv() bool {
	return matchEnvironment(host.TestEnv)
}

func isStageEnv() bool {
	return matchEnvironment(host.StageEnv)
}

func isProdEnv() bool {
	return matchEnvironment(host.ProdEnv)
}
