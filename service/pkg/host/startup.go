package host

import (
	"github.com/idiomatic-go/middleware/service/pkg/resource"
	"log"
	"net/http"
)

func Startup(r *http.ServeMux) (http.Handler, bool) {
	resource.ReadFile("")
	err := initLogging()
	if err != nil {
		log.Printf("startup logging error: %v", err)
		return nil, false
	}
	initIngress()
	initEgress()
	initRoutes(r)

	return r, true
}

func Shutdown() {}

func initExtract() {
	//err := extract.Initialize(uri string, newClient *http.Client, fn ErrorHandler)

}
