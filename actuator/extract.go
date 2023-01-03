package actuator

import "github.com/idiomatic-go/middleware/accessdata"

var defaultExtract Extract

/*
= func(entry *accessdata.Entry) {
	if entry != nil {
		log.Printf("{\"traffic\":\"%v\",\"start_time\":\"%v\",\"duration_ms\":%v,\"request\":\"%v\",\"response\":\"%v\",\"statusFlags\":\"%v\"}\n", entry.Traffic, entry.Start, entry.Duration, nil, nil, entry.StatusFlags)
	}
}
*/

type ExtractController interface {
	Extract(entry *accessdata.Entry)
}

type extract struct {
}

func newExtract() *extract {
	return new(extract)
}

func (e *extract) Extract(entry *accessdata.Entry) {
	if defaultExtract != nil {
		defaultExtract(entry)
	}
}
