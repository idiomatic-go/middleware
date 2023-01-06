package template

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

func ProcessContent[T any](content any) (T, *Status) {
	var t T
	if IsNil(content) {
		status := NewStatusInvalidArgument(errors.New("invalid argument: no content available"))
		return t, status
	}
	if status, ok := content.(*Status); ok {
		return t, status
	}
	// Code for err must be after Status as Status is an error
	if err, ok := content.(error); ok {
		return t, NewStatusError(err)
	}
	if t1, ok := content.(T); ok {
		return t1, NewStatusOk()
	}
	// TODO : update to reflect contained type.
	status := NewStatusInvalidArgument(errors.New(fmt.Sprintf("vhost.ProcessContent internal error : invalid content type : %v", reflect.TypeOf(content))))
	return t, status
}

func ProcessContextContent[T any](ctx context.Context) (T, *Status) {
	var t T
	if ctx == nil {
		return t, NewStatusInvalidArgument(errors.New(fmt.Sprintf("invalid configuration: context is nil")))
	}
	i := ctx.Value(contentKey)
	if IsNil(i) {
		return t, NewStatusInvalidArgument(errors.New(fmt.Sprintf("vhost.ProcessContent internal error : no content available")))
	}
	return ProcessContent[T](i)
}
