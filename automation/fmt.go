package automation

import "strings"

func IsEmpty(s string) bool {
	if s != "" {
		return strings.TrimLeft(s, " ") == ""
	}
	return true
}
