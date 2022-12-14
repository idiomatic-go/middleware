package host

import (
	"net/http"
	"net/http/pprof"
)

const (
	IndexPattern   = "/debug/pprof/"
	CmdLinePattern = "/debug/pprof/cmdline"
	ProfilePattern = "/debug/pprof/profile" // ?seconds=30
	SymbolPattern  = "/debug/pprof/symbol"
	TracePattern   = "/debug/pprof/trace"

	IndexRouteName   = "index"
	CmdLineRouteName = "cmdline"
	ProfileRouteName = "profile"
	SymbolRouteName  = "symbol"
	TraceRouteName   = "trace"
)

func initRoutes(r *http.ServeMux) {
	addRoutes(r)
}

func addRoutes(r *http.ServeMux) {
	r.Handle(IndexPattern, http.HandlerFunc(pprof.Index))
	r.Handle(CmdLinePattern, http.HandlerFunc(pprof.Cmdline))
	r.Handle(ProfilePattern, http.HandlerFunc(pprof.Profile))
	r.Handle(SymbolPattern, http.HandlerFunc(pprof.Symbol))
	r.Handle(TracePattern, http.HandlerFunc(pprof.Trace))

}
