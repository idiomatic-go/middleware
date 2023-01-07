package template

import (
	"errors"
	"github.com/idiomatic-go/middleware/template/internal"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	EchoScheme = "echo"
)

var doLocation = pkgPath + "/Do"

func Do[E ErrorHandler](req *http.Request) (resp *http.Response, status *Status) {
	return DoClient[E](req, http.DefaultClient)
}

func DoClient[E ErrorHandler](req *http.Request, client *http.Client) (resp *http.Response, status *Status) {
	var e E
	if req == nil {
		return nil, e.Handle(doLocation, errors.New("invalid argument: request is nil"))
	}
	if client == nil {
		return nil, e.Handle(doLocation, errors.New("invalid argument: client is nil"))
	}
	var err error
	switch req.URL.Scheme {
	case EchoScheme:
		resp, err = createEchoResponse(req)
	default:
		resp, err = client.Do(req)
	}
	return resp, e.Handle(doLocation, err)
}

func createEchoResponse(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("invalid argument: Request is nil")
	}
	var resp = http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Request: req}
	for key, element := range req.URL.Query() {
		switch key {
		case "httpError":
			return nil, http.ErrHijacked
		case "status":
			sc, err := strconv.Atoi(element[0])
			if err == nil {
				resp.StatusCode = sc
			} else {
				resp.StatusCode = http.StatusInternalServerError
			}
		case "body":
			if len(element[0]) > 0 && resp.Body == nil {
				// Handle escaped path? See notes on the url.URL struct
				resp.Body = &internal.ReaderCloser{Reader: strings.NewReader(element[0]), Err: nil}
			}
		case "ioError":
			// resp.StatusCode = http.StatusInternalServerError
			resp.Body = &internal.ReaderCloser{Reader: nil, Err: io.ErrUnexpectedEOF}
		default:
			if len(element[0]) > 0 {
				resp.Header.Add(key, element[0])
			}
		}
	}
	return &resp, nil
}
