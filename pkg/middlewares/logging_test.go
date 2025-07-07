package middlewares

import (
	"bytes"
	"github.com/isaacwallace123/GoWeb/app/types"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	// Save original stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function while writing to pipe
	fn()

	// Close and restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func dummyHandler(ctx *types.MiddlewareContext) error {
	ctx.ResponseEntity = &types.ResponseEntity{
		StatusCode: http.StatusOK,
		BodyData:   []byte(`OK`),
	}
	return ctx.Next()
}

func TestLoggingPre_Enabled(t *testing.T) {
	LoggingPre.Config.Enabled = true

	req := httptest.NewRequest("GET", "/test/pre", nil)
	res := httptest.NewRecorder()

	logs := captureStdout(func() {
		ctx := &types.MiddlewareContext{
			Request:        req,
			ResponseWriter: res,
			Chain: []types.MiddlewareFunc{
				LoggingPre.Func(),
				dummyHandler,
			},
			Index: -1,
		}
		_ = ctx.Next()
	})

	if !strings.Contains(logs, "GET") && !strings.Contains(logs, "/test/pre") {
		t.Errorf("expected log to contain 'GET /test/pre', got: %s", logs)
	}
}

func TestLoggingPost_Enabled(t *testing.T) {
	LoggingPost.Config.Enabled = true

	req := httptest.NewRequest("POST", "/test/post", nil)
	res := httptest.NewRecorder()

	logs := captureStdout(func() {
		ctx := &types.MiddlewareContext{
			Request:        req,
			ResponseWriter: res,
			Chain: []types.MiddlewareFunc{
				dummyHandler,
				LoggingPost.Func(),
			},
			Index: -1,
		}
		_ = ctx.Next()
	})

	if !strings.Contains(logs, "POST") && !strings.Contains(logs, "/test/post") {
		t.Errorf("expected log to contain 'POST /test/post 200', got: %s", logs)
	}
}

func TestLogging_Disabled(t *testing.T) {
	LoggingPre.Config.Enabled = false
	LoggingPost.Config.Enabled = false

	req := httptest.NewRequest("GET", "/no/logs", nil)
	res := httptest.NewRecorder()

	logs := captureStdout(func() {
		ctx := &types.MiddlewareContext{
			Request:        req,
			ResponseWriter: res,
			Chain: []types.MiddlewareFunc{
				LoggingPre.Func(),
				dummyHandler,
				LoggingPost.Func(),
			},
			Index: -1,
		}
		_ = ctx.Next()
	})

	if logs != "" {
		t.Errorf("expected no logs, got: %s", logs)
	}
}
