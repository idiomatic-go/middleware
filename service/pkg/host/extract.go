package host

import "github.com/idiomatic-go/middleware/extract"

func initExtract() {
	extract.SetErrorHandler(nil)
	//err := extract.Initialize(uri string, newClient *http.Client, fn ErrorHandler)
}
