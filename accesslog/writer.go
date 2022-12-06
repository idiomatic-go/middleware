package accesslog

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	markupNull      = "\"%v\":null"
	markupNullComma = "\"%v\":null,"

	markupString      = "\"%v\":\"%v\""
	markupStringComma = "\"%v\":\"%v\","

	markupValue      = "\"%v\":%v"
	markupValueComma = "\"%v\":%v,"
)

func writeMarkup(sb *strings.Builder, name, value string, format string) {
	if format == "" {
		format = markupStringComma
	}
	if value == "" {
		if format == markupStringComma {
			format = markupNullComma
		} else {
			format = markupNull
		}
		sb.WriteString(fmt.Sprintf(format, name))
	} else {
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}

func writeLocation(sb *strings.Builder) {
	writeMarkup(sb, "region", origin.Region, "")
	writeMarkup(sb, "zone", origin.Zone, "")
	writeMarkup(sb, "sub_zone", origin.SubZone, "")
	writeMarkup(sb, "service", origin.Service, "")
	writeMarkup(sb, "instance_id", origin.InstanceId, "")
}

func writeStartTime(sb *strings.Builder, start time.Time) {
	writeMarkup(sb, "start_time", FmtTimestamp(start), "")
}

func writeDuration(sb *strings.Builder, duration time.Duration) {
	d := int(duration / time.Duration(1e6))
	writeMarkup(sb, "duration_ms", strconv.Itoa(d), markupValueComma)
}
