package template

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var unmarshalLoc = pkgPath + "/unmarshal"
var decodeLoc = pkgPath + "/decode"

/*

if decode {

			}
func (r *ResponseStatus) DecodeJson(i interface{}) error {
	if i == nil {
		return errors.New("invalid argument: interface{} is nil")
	}
	if !r.IsContent() {
		return nil
	}
	r.UnmarshalErr = json.NewDecoder(r.Response.Body).Decode(i)
	return nil
}

*/

//func Unmarshal[T any](resp *http.Response) (T, *Status) {
//	return unmarshal[T](resp, false)
//}

func Decode[T any](resp *http.Response) (T, *Status) {
	var t T
	if resp == nil || resp.Body == nil {
		return t, NewStatus(StatusInvalidContent, decodeLoc, errors.New("response or response body is nil"))
	}
	switch any(t).(type) {
	case []byte, string:
		return Unmarshal[T](resp)
	default:
		err := json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return t, NewStatus(StatusJsonDecodeError, decodeLoc, err)
		}
		return t, NewStatusOk()
	}
}

func Unmarshal[T any](resp *http.Response) (T, *Status) {
	var t T
	//var e E

	if resp == nil || resp.Body == nil {
		return t, NewStatus(StatusInvalidContent, unmarshalLoc, errors.New("response or response body is nil"))
	}
	buf, err := readAll(resp.Body)
	if err != nil {
		return t, NewStatus(StatusIOError, unmarshalLoc, err)
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
			return t, NewStatus(StatusInvalidContent, unmarshalLoc, err).SetContent(buf)
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
