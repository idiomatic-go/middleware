package accesslog

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	markupNull   = "\"%v\":null"
	markupString = "\"%v\":\"%v\""
	markupValue  = "\"%v\":%v"
)

func writeMarkup(sb *strings.Builder, name, value string, quotes bool) {
	if sb.Len() == 0 {
		sb.WriteString("{")
	} else {
		sb.WriteString(",")
	}
	var format = markupString
	if !quotes {
		format = markupValue
	}
	if value == "" {
		format = markupNull
		sb.WriteString(fmt.Sprintf(format, name))
	} else {
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}

func writeLocation(sb *strings.Builder) {
	writeMarkup(sb, "region", origin.Region, true)
	writeMarkup(sb, "zone", origin.Zone, true)
	writeMarkup(sb, "sub_zone", origin.SubZone, true)
	writeMarkup(sb, "service", origin.Service, true)
	writeMarkup(sb, "instance_id", origin.InstanceId, true)
}

func writeStartTime(sb *strings.Builder, start time.Time) {
	writeMarkup(sb, "start_time", FmtTimestamp(start), true)
}

func writeDuration(sb *strings.Builder, duration time.Duration) {
	d := int(duration / time.Duration(1e6))
	writeMarkup(sb, "duration_ms", strconv.Itoa(d), false)
}
