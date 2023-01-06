package template

import (
	"context"
	"time"
)

type contentid struct{}

var contentKey contentid

// ContextWithContent - creates a new Context with content
func ContextWithContent(ctx context.Context, content any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if content == nil {
		return ctx
	}
	return &valueCtx{ctx, contentKey, content}
}

func ContextContent(ctx context.Context) any {
	if ctx == nil {
		return nil
	}
	i := ctx.Value(contentKey)
	if IsNil(i) {
		return nil
	}
	return i
}

func IsContextContent(c context.Context) bool {
	if c == nil {
		return false
	}
	for {
		switch c.(type) {
		case *valueCtx:
			/*
				if ctx != nil {
					if ctx.val != nil {
						if reflect.TypeOf(ctx.val) == reflect.TypeOf(t) {
							return true
						} else {
							return false
						}
					}
				}

			*/
			return true
		default:
			return false
		}
	}
}

type valueCtx struct {
	ctx      context.Context
	key, val any
}

func (*valueCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*valueCtx) Done() <-chan struct{} {
	return nil
}

func (*valueCtx) Err() error {
	return nil
}

func (v *valueCtx) Value(key any) any {
	return v.val
}
