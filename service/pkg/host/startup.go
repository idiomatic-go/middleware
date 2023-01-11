package host

import (
	mhost "github.com/idiomatic-go/middleware/host"
	"github.com/idiomatic-go/middleware/service/pkg/facebook"
	"github.com/idiomatic-go/middleware/service/pkg/google"
	"github.com/idiomatic-go/middleware/service/pkg/resource"
	"github.com/idiomatic-go/middleware/service/pkg/twitter"
	"github.com/idiomatic-go/middleware/template"
	"net/http"
	"time"
)

func Startup[E template.ErrorHandler](r *http.ServeMux) (http.Handler, *template.Status) {
	var e E
	resource.ReadFile("")
	err := initLogging()
	if err != nil {
		return nil, e.Handle("/host/startup/logging", err)
	}
	err = initEgress()
	if err != nil {
		return nil, e.Handle("/host/startup/egress", err)
	}
	initIngress()
	initRoutes(r)
	status := mhost.Startup[E](time.Second*5, map[string][]any{
		google.Uri: {matchEnvironment}, twitter.Uri: {matchEnvironment}, facebook.Uri: {matchEnvironment}},
	)
	return r, status
}

func Shutdown() {}

func initExtract() {
	//err := extract.Initialize(uri string, newClient *http.Client, fn ErrorHandler)

}
