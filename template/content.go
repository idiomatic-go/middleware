package template

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type pkg struct{}

var pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()

func ProcessContent[T any, E ErrorHandler](content any) (T, *Status) {
	var t T
	var e E
	if IsNil(content) {
		return t, e.Handle(pkgPath+"/ProcessContent", errors.New("invalid argument: no content available")).SetCode(StatusInvalidArgument)
	}
	if status, ok := content.(*Status); ok {
		return t, status
	}
	// Code for err must be after Status as Status is an error
	if err, ok := content.(error); ok {
		return t, NewStatusError(pkgPath, err)
	}
	if t1, ok := content.(T); ok {
		return t1, NewStatusOk()
	}
	// TODO : update to reflect contained type.
	return t, e.Handle(pkgPath, errors.New(fmt.Sprintf("invalid argument: invalid content type : %v", reflect.TypeOf(content)))).SetCode(StatusInvalidArgument)
}

func ProcessContextContent[T any, E ErrorHandler](ctx context.Context) (T, *Status) {
	var t T
	var e E
	if ctx == nil {
		return t, e.Handle(pkgPath+"/ProcessContextContent", errors.New(fmt.Sprintf("invalid configuration: context is nil"))).SetCode(StatusInvalidArgument)
	}
	i := ctx.Value(contentKey)
	if IsNil(i) {
		return t, e.Handle(pkgPath, errors.New(fmt.Sprintf("invalid configuration: no content available"))).SetCode(StatusInvalidArgument)
	}
	return ProcessContent[T, E](i)
}
