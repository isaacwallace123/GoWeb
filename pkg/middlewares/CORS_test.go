package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isaacwallace123/GoWeb/app/types"
)

func setupRequest(method, origin string) *http.Request {
	req, _ := http.NewRequest(method, "/test", nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	return req
}

func runMiddleware(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	ctx := &types.MiddlewareContext{
		Request:        req,
		ResponseWriter: rr,
		ResponseEntity: nil,
		Index:          -1,
		Chain: []types.MiddlewareFunc{
			CORS.Func(), // only CORS middleware
		},
	}

	_ = ctx.Next()
	return rr
}

func TestCORS_AllowsOrigin(t *testing.T) {
	CORS.Config.AllowedOrigins = []string{"http://example.com"}
	CORS.Config.AllowedMethods = []string{"GET"}

	req := setupRequest("GET", "http://example.com")
	rr := runMiddleware(req)

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://example.com" {
		t.Errorf("expected Access-Control-Allow-Origin to be set, got %q", got)
	}
}

func TestCORS_RejectsDisallowedMethod(t *testing.T) {
	CORS.Config.AllowedOrigins = []string{"http://example.com"}
	CORS.Config.AllowedMethods = []string{"GET"}

	req := setupRequest("POST", "http://example.com")
	rr := runMiddleware(req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestCORS_PreflightOPTIONS(t *testing.T) {
	CORS.Config.AllowedOrigins = []string{"http://example.com"}
	CORS.Config.AllowedMethods = []string{"GET", "POST"}

	req := setupRequest("OPTIONS", "http://example.com")
	rr := runMiddleware(req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204 for OPTIONS, got %d", rr.Code)
	}
}

func TestCORS_MissingOrigin(t *testing.T) {
	CORS.Config.AllowedOrigins = []string{"http://example.com"}

	req := setupRequest("GET", "") // no Origin header
	rr := runMiddleware(req)

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no CORS header when Origin is missing, got %q", got)
	}
}
