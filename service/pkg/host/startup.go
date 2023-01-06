package host

import (
	"net/http"
)

func Startup(r *http.ServeMux) (http.Handler, bool) {
	initLogging()
	initIngress()
	initEgress()
	initRoutes(r)

	return r, true
}

func Shutdown() {}

func initExtract() {
	//err := extract.Initialize(uri string, newClient *http.Client, fn ErrorHandler)

}
