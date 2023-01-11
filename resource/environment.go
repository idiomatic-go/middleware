package resource

import (
	"os"
	"strings"
)

const (
	DevEnv = iota
	ReviewEnv
	TestEnv
	StageEnv
	ProdEnv
	runtimeEnvKey  = "RUNTIME_ENV"
	devEnvValue    = "dev"
	reviewEnvValue = "review"
	testEnvValue   = "test"
	stageEnvValue  = "stage"
	prodEnvValue   = "prod"
)

type EnvironmentMatcher func(env int) bool

func SetRuntimeKey(s string) {
	if s != "" {
		runtimeKey = s
	}
}

func SetEnvironmentMatcher(fn EnvironmentMatcher) {
	if fn != nil {
		envMatcher = fn
	} else {
		envMatcher = matchEnvironment
	}
}

var runtimeKey = runtimeEnvKey
var envMatcher EnvironmentMatcher

func init() {
	SetEnvironmentMatcher(nil)
}

func matchEnvironment(env int) bool {
	s := os.Getenv(runtimeKey)
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
}

func IsDevEnv() bool {
	return envMatcher(DevEnv)
}

func IsReviewEnv() bool {
	return envMatcher(ReviewEnv)
}

func IsTestEnv() bool {
	return envMatcher(TestEnv)
}

func IsStageEnv() bool {
	return envMatcher(StageEnv)
}

func IsProdEnv() bool {
	return envMatcher(ProdEnv)
}

/*
// FuncBool - type for niladic functions, functions with no parameters
type FuncBool func() bool

func OverrideIsDevEnv(fn FuncBool) {
	if fn != nil {
		isDevEnv = fn
		dev = IsDevEnv()
	}
}

var runtimeKey = runtimeEnvKey
var isDevEnv FuncBool
var dev = true

func IsDevEnvironment() bool {
	return dev
}

func init() {
	isDevEnv = func() bool {
		target := GetEnv()
		if len(target) == 0 || strings.EqualFold(target, devEnv) {
			return true
		}
		return false
	}
	dev = isDevEnv()
}



//func IsDevEnv() bool {
//	return isDevEnv()
//}
*/

// GetEnv - function to get the resource runtime environment
func GetEnv() string {
	s := os.Getenv(runtimeKey)
	if s == "" {
		return "dev"
	}
	return s
}

// SetEnv - function to set the resource runtime environment
func SetEnv(s string) {
	os.Setenv(runtimeKey, s)
}

//var IsDevEnv EnvValid = func() bool {
//	return IsEnvMatch(RUNTIME_ENV, DEV_ENV)
//}

/*
var IsReviewEnv EnvValid = func() bool {
	return IsEnvMatch(REVIEW_ENV)
}

var IsTestEnv EnvValid = func() bool {
	return IsEnvMatch(TEST_ENV)
}

var IsStageEnv EnvValid = func() bool {
	return IsEnvMatch(STAGE_ENV)
}

var IsProdEnv EnvValid = func() bool {
	return IsEnvMatch(PROD_ENV)
}


func IsEnvMatch(val string) bool {
	target := os.Getenv(runtimeKey)
	if len(target) == 0 || !strings.EqualFold(target, val) {
		return false
	}
	return true
}

*/
