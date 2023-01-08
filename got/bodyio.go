package got

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var unmarshalLoc = pkgPath + "/unmarshal"
var decodeLoc = pkgPath + "/decode"

func Decode[T any, E ErrorHandler](resp *http.Response) (T, *Status) {
	var t T
	var e E

	if resp == nil || resp.Body == nil {
		return t, e.Handle(decodeLoc, errors.New("response or response body is nil")).SetCode(StatusInvalidContent)
	}
	switch any(t).(type) {
	case []byte, string:
		return Unmarshal[T, E](resp)
	default:
		err := json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return t, e.Handle(decodeLoc, err).SetCode(StatusJsonDecodeError)
		}
		return t, NewStatusOk()
	}
}

func Unmarshal[T any, E ErrorHandler](resp *http.Response) (T, *Status) {
	var t T
	var e E

	if resp == nil || resp.Body == nil {
		return t, e.Handle(unmarshalLoc, errors.New("response or response body is nil")).SetCode(StatusInvalidContent)
	}
	buf, err := readAll(resp.Body)
	if err != nil {
		return t, e.Handle(unmarshalLoc, err).SetCode(StatusIOError)
	}
	switch any(t).(type) {
	case []byte:
		unmarshalBytes(&t, buf)
		return t, NewStatusOk()
	case string:
		unmarshalString(&t, buf)
		return t, NewStatusOk()
	default:
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, e.Handle(unmarshalLoc, err).SetContent(buf).SetCode(StatusInvalidContent)
		}
		return t, NewStatusOk()
	}
}

func unmarshalBytes(out any, in []byte) {
	if buf, ok := out.(*[]byte); ok {
		*buf = append(*buf, in...)
	}
}

func unmarshalString(out any, in []byte) {
	if ptr, ok := out.(*string); ok {
		*ptr = string(in)
	}
}

func readAll(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return io.ReadAll(body)
}
