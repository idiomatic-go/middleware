package host

import (
	"github.com/idiomatic-go/middleware/service/pkg/resource"
	"github.com/idiomatic-go/middleware/template"
	"net/http"
)

func Startup[E template.ErrorHandler](r *http.ServeMux) (http.Handler, *template.Status) {
	var e E
	resource.ReadFile("")
	err := initLogging()
	if err != nil {
		return nil, e.Handle("startup logging:", err)
	}
	err = initEgress()
	if err != nil {
		return nil, e.Handle("startup egress:", err)
	}
	initIngress()
	initRoutes(r)

	return r, template.NewStatusOk()
}

func Shutdown() {}

func initExtract() {
	//err := extract.Initialize(uri string, newClient *http.Client, fn ErrorHandler)

}
