package vhost

import (
	"os"
	"strings"
)

const (
	devEnv        = "dev"
	runtimeEnvKey = "RUNTIME_ENV"
)

// FuncBool - type for niladic functions, functions with no parameters
type FuncBool func() bool

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

func IsDevEnv() bool {
	return isDevEnv()
}

// GetEnv - function to get the vhost runtime environment
func GetEnv() string {
	s := os.Getenv(runtimeKey)
	if s == "" {
		return devEnv
	}
	return s
}

// SetEnv - function to set the vhost runtime environment
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
