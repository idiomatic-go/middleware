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

func writeMarkup(sb *strings.Builder, name string, data *logd) {
	if sb.Len() == 0 {
		sb.WriteString("{")
	} else {
		sb.WriteString(",")
	}
	value, format, err := resolve(name, data)
	if err != nil {
		value = err.Error()
	}
	if value == "" {
		sb.WriteString(fmt.Sprintf(markupNull, name))
	} else {
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}
