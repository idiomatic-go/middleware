package host

import (
	"github.com/idiomatic-go/middleware/accessdata"
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

func Startup(r *http.ServeMux) bool {
	// accessdata options
	//   SetOrigin() - part of the access log data, and will show on each log entry
	//   SetPingRoutes() - determine which routes/actuator are health liveness check routes
	accessdata.SetOrigin(accessdata.Origin{
		Region:     "region",
		Zone:       "zone",
		SubZone:    "sub-zone",
		Service:    "example-middleware",
		InstanceId: "1234-567-8901",
	})
	accessdata.SetPingRoutes(nil)

	addRoutes(r)

	return true
}

func Shutdown() {}

func addRoutes(r *http.ServeMux) {
	r.Handle(IndexPattern, http.HandlerFunc(pprof.Index))
	r.Handle(CmdLinePattern, http.HandlerFunc(pprof.Cmdline))
	r.Handle(ProfilePattern, http.HandlerFunc(pprof.Profile))
	r.Handle(SymbolPattern, http.HandlerFunc(pprof.Symbol))
	r.Handle(TracePattern, http.HandlerFunc(pprof.Trace))

	//r.HandleFunc(indexPattern, pprof.Index)
	//r.HandleFunc(cmdLinePattern, pprof.Cmdline)
	//r.HandleFunc(profilePattern, pprof.Profile)
	//r.HandleFunc(symbolPattern, pprof.Symbol)
	//r.HandleFunc(tracePattern, pprof.Trace)
}
