package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/isaacwallace123/GoWeb/core"
	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

type MockController struct{}

func (c *MockController) BasePath() string {
	return "/test"
}

func (c *MockController) Routes() []core.RouteEntry {
	return []core.RouteEntry{
		{Method: "GET", Path: "/{name}", Handler: "Greet"},
		{Method: "POST", Path: "/", Handler: "Echo"},
	}
}

func (c *MockController) Greet(name string) *response.ResponseEntity {
	return response.Status(httpstatus.OK).Body(map[string]string{"message": "Hello, " + name})
}

func (c *MockController) Echo(req map[string]string) *response.ResponseEntity {
	return response.Status(httpstatus.CREATED).Body(req)
}

func TestPathParamBinding(t *testing.T) {
	router := core.RegisterControllers(&MockController{})

	req := httptest.NewRequest("GET", "/test/GoWeb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}

	expected := `{"message":"Hello, GoWeb"}`
	if strings.TrimSpace(w.Body.String()) != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestJsonBodyBinding(t *testing.T) {
	router := core.RegisterControllers(&MockController{})

	jsonBody := `{"key":"value"}`
	req := httptest.NewRequest("POST", "/test/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", w.Code)
	}

	expected := `{"key":"value"}`
	if strings.TrimSpace(w.Body.String()) != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestRouteNotFound(t *testing.T) {
	router := core.RegisterControllers(&MockController{})

	req := httptest.NewRequest("GET", "/nonexistent/path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", w.Code)
	}
}
