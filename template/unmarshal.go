package template

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Unmarshaler[T any] interface {
	Unmarshal(resp *http.Response) (T, error)
}

func UnmarshalInterface[T any, E ErrorHandler](resp *http.Response, u Unmarshaler[T]) (T, *Status) {
	var t T
	var e E
	if resp == nil {
		return t, e.Handle(pkgPath+"/unmarshal", errors.New("response is nil")).SetCode(StatusInvalidArgument)
	}
	if resp.Body == nil {
		return t, e.Handle(pkgPath+"/unmarshal", errors.New("response body is nil")).SetCode(StatusInvalidContent)
	}
	t, err := u.Unmarshal(resp)
	return t, e.Handle(pkgPath+"/unmarshal", err)
}

type StringUnmarshal struct{}

func (StringUnmarshal) Unmarshal(resp *http.Response) (string, error) {
	if resp == nil || resp.Body == nil {
		return "", errors.New("response or response body is nil")
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Unmarshal[T any](resp *http.Response) (T, *Status) {
	var t T
	//var e E

	if resp == nil || resp.Body == nil {
		return t, NewStatus(StatusInvalidContent, pkgPath+"/unmarshal", errors.New("response or response body is nil"))
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return t, NewStatus(StatusInvalidContent, pkgPath+"/unmarshal", errors.New("no content"))
	}
	switch any(t).(type) {
	case []byte:
		unmarshalBytes(&t, bytes)
		return t, NewStatusOk()
	case string:
		unmarshalString(&t, bytes)
		return t, NewStatusOk()
	default:
		err = json.Unmarshal(bytes, &t)
	}
	return t, NewStatus(StatusInvalidContent, "loc", errors.New("unmapped content type")).SetContent(bytes)
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
