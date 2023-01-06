package google

import (
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"io"
	"net/http"
	"reflect"
)

const (
	uri    = "https://www.google.com/search?q=test"
	search = "/Search"
)

type pkg struct{}

var pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()

func Search[E template.ErrorHandler](req *http.Request) ([]byte, *template.Status) {
	var e E

	if req == nil {
		return nil, e.Handle(pkgPath+search, errors.New("request is nil")).SetCode(template.StatusInvalidArgument)
	}
	if template.IsContextContent(req.Context()) {
		return template.ProcessContextContent[[]byte, template.NoOpHandler](req.Context())
	}
	newReq, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, e.Handle(pkgPath+search, err)
	}
	resp, err2 := http.DefaultClient.Do(newReq)
	if err2 != nil {
		return nil, e.Handle(pkgPath+search, err)
	}
	defer resp.Body.Close()
	bytes, err3 := io.ReadAll(resp.Body)
	return bytes, e.Handle(pkgPath+search, err3)
}
