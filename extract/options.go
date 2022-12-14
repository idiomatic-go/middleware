package extract

import (
	"fmt"
	"log"
)

// ErrorHandler - allows handling of extract errors, default is to log.Println
type ErrorHandler func(err error)

type options struct {
	handler ErrorHandler
}

var opt options

func init() {
	SetErrorHandler(nil)
}

func SetErrorHandler(fn ErrorHandler) {
	if fn != nil {
		opt.handler = fn
	} else {
		opt.handler = func(err error) {
			log.Println(err)
		}
	}
}

func SetTestErrorHandler() {
	opt.handler = func(err error) {
		fmt.Printf("test: extract(logd) -> [err:%v]\n", err)
	}
}

func OnError(err error) {
	if opt.handler != nil {
		opt.handler(err)
	}
}
