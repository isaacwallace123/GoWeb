package app

import (
	"github.com/isaacwallace123/GoWeb/pkg/HttpStatus"
	"github.com/isaacwallace123/GoWeb/pkg/ResponseEntity"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/isaacwallace123/GoUtils/jsonutil"
	"github.com/isaacwallace123/GoWeb/app/types"
)

type TestResponse struct {
	Method string `json:"method"`
	ID     string `json:"id,omitempty"`
}

type DummyController struct{}

func (c *DummyController) BasePath() string { return "/api/v1/test" }
func (c *DummyController) Routes() []types.Route {
	return []types.Route{
		{Method: "GET", Path: "/", Handler: "GetAll"},
		{Method: "POST", Path: "/", Handler: "Post"},
		{Method: "PUT", Path: "/{id}", Handler: "Put"},
		{Method: "DELETE", Path: "/{id}", Handler: "Delete"},
	}
}
func (c *DummyController) GetAll() *ResponseEntity.ResponseEntity {
	return ResponseEntity.Status(HttpStatus.OK).Body(TestResponse{Method: "GET"})
}
func (c *DummyController) Post() *ResponseEntity.ResponseEntity {
	return ResponseEntity.Status(HttpStatus.CREATED).Body(TestResponse{Method: "POST"})
}
func (c *DummyController) Put(id string) *ResponseEntity.ResponseEntity {
	return ResponseEntity.Status(HttpStatus.OK).Body(TestResponse{Method: "PUT", ID: id})
}
func (c *DummyController) Delete(id string) *ResponseEntity.ResponseEntity {
	return ResponseEntity.Status(HttpStatus.NO_CONTENT).Body(nil)
}

func clearAllGlobalState() {
	types.PreMiddlewares = nil
	types.PostMiddlewares = nil
}

// For test isolation
func setupRouter() *Router {
	clearAllGlobalState()

	router := NewRouter()
	router.RegisterControllers(&DummyController{})

	return router
}

// GET /api/v1/test/
func TestRouter_GET(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest("GET", "/api/v1/test/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != HttpStatus.OK {
		t.Fatalf("GET: want status %d, got %d", HttpStatus.OK, resp.StatusCode)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(bodyBytes))
	if bodyStr == "" {
		t.Fatalf("GET: expected non-empty body")
	}
	var tr TestResponse
	if err := jsonutil.FromString(bodyStr, &tr); err != nil {
		t.Fatalf("GET: failed to decode JSON: %v\nRaw body: %s", err, bodyStr)
	}
	if tr.Method != "GET" {
		t.Errorf("GET: want method 'GET', got %q", tr.Method)
	}
}

// POST /api/v1/test/
func TestRouter_POST(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest("POST", "/api/v1/test/", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != HttpStatus.CREATED {
		t.Fatalf("POST: want status %d, got %d", HttpStatus.CREATED, resp.StatusCode)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(bodyBytes))
	if bodyStr == "" {
		t.Fatalf("POST: expected non-empty body")
	}
	var tr TestResponse
	if err := jsonutil.FromString(bodyStr, &tr); err != nil {
		t.Fatalf("POST: failed to decode JSON: %v\nRaw body: %s", err, bodyStr)
	}
	if tr.Method != "POST" {
		t.Errorf("POST: want method 'POST', got %q", tr.Method)
	}
}

// PUT /api/v1/test/{id}
func TestRouter_PUT(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest("PUT", "/api/v1/test/abc123", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != HttpStatus.OK {
		t.Fatalf("PUT: want status %d, got %d", HttpStatus.OK, resp.StatusCode)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(bodyBytes))
	if bodyStr == "" {
		t.Fatalf("PUT: expected non-empty body")
	}
	var tr TestResponse
	if err := jsonutil.FromString(bodyStr, &tr); err != nil {
		t.Fatalf("PUT: failed to decode JSON: %v\nRaw body: %s", err, bodyStr)
	}
	if tr.Method != "PUT" {
		t.Errorf("PUT: want method 'PUT', got %q", tr.Method)
	}
	if tr.ID != "abc123" {
		t.Errorf("PUT: want ID 'abc123', got %q", tr.ID)
	}
}

// DELETE /api/v1/test/{id}
func TestRouter_DELETE(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest("DELETE", "/api/v1/test/xyz789", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != HttpStatus.NO_CONTENT {
		t.Fatalf("DELETE: want status %d, got %d", HttpStatus.NO_CONTENT, resp.StatusCode)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(bodyBytes))
	if bodyStr != "" {
		t.Fatalf("DELETE: expected empty body for 204, got %q", bodyStr)
	}
}
