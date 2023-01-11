package host

import (
	"github.com/idiomatic-go/middleware/service/pkg/resource"
	"net/http"
)

func Startup(r *http.ServeMux) (http.Handler, bool) {
	resource.ReadFile("")
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
