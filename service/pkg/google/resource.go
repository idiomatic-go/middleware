package google

import (
	"github.com/idiomatic-go/middleware/template"
	"net/http"
)

const (
	uri = "https://www.google.com/search?q=test"
)

var searchLocation = pkgPath + "/search"

func Search[E template.ErrorHandler](req *http.Request) ([]byte, *template.Status) {
	var e E

	if req != nil && template.IsContextContent(req.Context()) {
		return template.ProcessContextContent[[]byte, E](req.Context())
	}
	newReq, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, e.Handle(searchLocation, err)
	}
	resp, status := template.Do(newReq)
	if status.IsErrors() {
		return nil, e.HandleStatus(status)
	}
	if !status.Ok() {
		return nil, status
	}
	return template.Unmarshal[[]byte, E](resp)
}
