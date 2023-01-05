package host

import (
	"net/http"
)

func Startup(r *http.ServeMux) (http.Handler, bool) {
	initLogging()
	initIngress(r)
	initEgress()
	initExtract()

	return r, true
}

func Shutdown() {}
