package extract

import (
	"log"
)

// ErrorHandler - allows handling of extract errors, default is to log.Println
type ErrorHandler func(err error)

func SetErrorHandler(fn ErrorHandler) {
	if fn != nil {
		opt.handler = fn
	} else {
		opt.handler = func(err error) {
			log.Println(err)
		}
	}
}

func OnError(err error) {
	if opt.handler != nil {
		opt.handler(err)
	}
}

type options struct {
	handler ErrorHandler
}

var opt options

func init() {
	SetErrorHandler(nil)
}
