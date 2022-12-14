package accesslog

import "log"

// Extract - optionally allows extraction of log data
type Extract func(l *Logd)

// Write - override log output disposition, default is log.Println
type Write func(s string)

type options struct {
	extractFn    Extract
	ingressWrite Write
	egressWrite  Write
}

var opt options

func init() {
	SetIngressWrite(nil)
	SetEgressWrite(nil)
}

func IsExtract() bool {
	return opt.extractFn != nil
}

func SetExtract(fn Extract) {
	opt.extractFn = fn
}

func callExtract(l *Logd) {
	if IsExtract() {
		opt.extractFn(l)
	}
}

func SetIngressWrite(fn Write) {
	if fn != nil {
		opt.ingressWrite = fn
	} else {
		opt.ingressWrite = func(s string) {
			log.Println(s)
		}
	}
}

func SetEgressWrite(fn Write) {
	if fn != nil {
		opt.egressWrite = fn
	} else {
		opt.egressWrite = func(s string) {
			log.Println(s)
		}
	}
}

func ingressWrite(s string) {
	if opt.ingressWrite != nil {
		opt.ingressWrite(s)
	}
}

func egressWrite(s string) {
	if opt.egressWrite != nil {
		opt.egressWrite(s)
	}
}
