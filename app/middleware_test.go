package app

import (
	"github.com/isaacwallace123/GoWeb/app/types"
	"testing"
)

func clearMiddleware() {
	types.PreMiddlewares = nil
	types.PostMiddlewares = nil
}

// Dummy middleware implementation
type dummyMiddleware struct {
	name string
}

func (d *dummyMiddleware) Func() types.MiddlewareFunc {
	return func(ctx *types.MiddlewareContext) error {
		return nil
	}
}

func TestUseRegistersPreMiddleware(t *testing.T) {
	clearMiddleware()
	d := &dummyMiddleware{name: "pre"}
	Use(d)

	if len(types.PreMiddlewares) != 1 {
		t.Errorf("expected 1 pre-middleware, got %d", len(types.PreMiddlewares))
	}

	fn := Pre()[0]
	if fn == nil {
		t.Error("Pre() did not return a valid middleware func")
	}
}

func TestUseAfterRegistersPostMiddleware(t *testing.T) {
	clearMiddleware()
	d := &dummyMiddleware{name: "post"}
	UseAfter(d)

	if len(types.PostMiddlewares) != 1 {
		t.Errorf("expected 1 post-middleware, got %d", len(types.PostMiddlewares))
	}

	fn := Post()[0]
	if fn == nil {
		t.Error("Post() did not return a valid middleware func")
	}
}

func TestMultipleMiddlewares(t *testing.T) {
	clearMiddleware()
	d1 := &dummyMiddleware{name: "a"}
	d2 := &dummyMiddleware{name: "b"}

	Use(d1, d2)
	UseAfter(d1, d2)

	if len(types.PreMiddlewares) != 2 {
		t.Errorf("expected 2 pre-middlewares, got %d", len(types.PreMiddlewares))
	}
	if len(types.PostMiddlewares) != 2 {
		t.Errorf("expected 2 post-middlewares, got %d", len(types.PostMiddlewares))
	}
}
