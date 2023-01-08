package google

import (
	"github.com/idiomatic-go/middleware/got"
	"net/http"
	"reflect"
)

const (
	uri    = "https://www.google.com/search?q=test"
	search = "/Search"
)

type pkg struct{}

var pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()
var searchLocation = pkgPath + search

func Search[E got.ErrorHandler](req *http.Request) ([]byte, *got.Status) {
	var e E

	if req != nil && got.IsContextContent(req.Context()) {
		return got.ProcessContextContent[[]byte, E](req.Context())
	}
	newReq, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, e.Handle(searchLocation, err)
	}
	resp, status := got.Do(newReq)
	if status.IsErrors() {
		return nil, e.HandleStatus(status)
	}
	if !status.Ok() {
		return nil, status
	}
	return got.Unmarshal[[]byte, E](resp)
}
