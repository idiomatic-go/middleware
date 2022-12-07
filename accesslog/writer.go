package accesslog

import (
	"fmt"
	"strings"
)

const (
	markupNull   = "\"%v\":null"
	markupString = "\"%v\":\"%v\""
	markupValue  = "\"%v\":%v"
)

func writeJson(sb *strings.Builder, name, value string, stringValue bool) {
	if sb.Len() == 0 {
		sb.WriteString("{")
	} else {
		sb.WriteString(",")
	}
	if value == "" {
		sb.WriteString(fmt.Sprintf(markupNull, name))
	} else {
		format := markupString
		if !stringValue {
			format = markupValue
		}
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}
