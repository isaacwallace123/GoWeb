package app

import (
	"github.com/isaacwallace123/GoWeb/app/types"
	"testing"
)

// Clears global middleware slices before each test for isolation.
func clearMiddleware() {
	types.PreMiddlewares = nil
	types.PostMiddlewares = nil
}

func dummyPre(ctx *types.MiddlewareContext) error  { return nil }
func dummyPost(ctx *types.MiddlewareContext) error { return nil }

func TestUseRegistersPreMiddleware(t *testing.T) {
	clearMiddleware()
	Use(dummyPre)
	if len(types.PreMiddlewares) != 1 {
		t.Errorf("expected 1 pre-middleware, got %d", len(types.PreMiddlewares))
	}
	if Pre()[0] == nil {
		t.Error("Pre() did not return a valid middleware func")
	}
}

func TestUseAfterRegistersPostMiddleware(t *testing.T) {
	clearMiddleware()
	UseAfter(dummyPost)
	if len(types.PostMiddlewares) != 1 {
		t.Errorf("expected 1 post-middleware, got %d", len(types.PostMiddlewares))
	}
	if Post()[0] == nil {
		t.Error("Post() did not return a valid middleware func")
	}
}

func TestMultipleMiddlewares(t *testing.T) {
	clearMiddleware()
	Use(dummyPre, dummyPre)
	UseAfter(dummyPost, dummyPost)
	if len(types.PreMiddlewares) != 2 {
		t.Errorf("expected 2 pre-middlewares, got %d", len(types.PreMiddlewares))
	}
	if len(types.PostMiddlewares) != 2 {
		t.Errorf("expected 2 post-middlewares, got %d", len(types.PostMiddlewares))
	}
}
