package accesslog

import (
	"strings"
)

func IsEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimLeft(s, " ") == ""
}
