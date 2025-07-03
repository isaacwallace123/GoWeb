package tests

import (
	"bytes"
	"github.com/isaacwallace123/GoWeb/ResponseEntity"
	"github.com/isaacwallace123/GoWeb/app/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isaacwallace123/GoWeb/app"
)

type TestUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MockController struct {
	db map[int]TestUser
}

func (c *MockController) BasePath() string {
	return "/mock"
}

func (c *MockController) Routes() []app.RouteEntry {
	return []app.RouteEntry{
		{Method: "GET", Path: "/{id}", Handler: "Get"},
		{Method: "POST", Path: "/", Handler: "Post"},
		{Method: "PUT", Path: "/{id}", Handler: "Put"},
		{Method: "DELETE", Path: "/{id}", Handler: "Delete"},
	}
}

func (c *MockController) Get(id int) *ResponseEntity.ResponseEntity {
	user, ok := c.db[id]
	if !ok {
		return ResponseEntity.Status(http.StatusNotFound).Body(map[string]string{"error": "Not found"})
	}
	return ResponseEntity.Status(http.StatusOK).Body(user)
}

func (c *MockController) Post(req TestUser) *ResponseEntity.ResponseEntity {
	newID := len(c.db) + 1
	c.db[newID] = req
	return ResponseEntity.Status(http.StatusCreated).Body(map[string]any{"id": newID, "user": req})
}

func (c *MockController) Put(id int, req TestUser) *ResponseEntity.ResponseEntity {
	if _, ok := c.db[id]; !ok {
		return ResponseEntity.Status(http.StatusNotFound).Body(map[string]string{"error": "Not found"})
	}
	c.db[id] = req
	return ResponseEntity.Status(http.StatusOK).Body(req)
}

func (c *MockController) Delete(id int) *ResponseEntity.ResponseEntity {
	if _, ok := c.db[id]; !ok {
		return ResponseEntity.Status(http.StatusNotFound).Body(map[string]string{"error": "Not found"})
	}
	delete(c.db, id)
	return ResponseEntity.Status(http.StatusNoContent)
}

func setupTestRouter() *types.Router {
	ctrl := &MockController{db: make(map[int]TestUser)}
	return types.RegisterControllers(ctrl)
}

func TestPostUser(t *testing.T) {
	router := setupTestRouter()
	body := `{"name":"Isaac","email":"isaac@example.com"}`
	req := httptest.NewRequest("POST", "/mock/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", w.Code)
	}
}

func TestGetUser(t *testing.T) {
	router := setupTestRouter()

	// Create a user first
	body := `{"name":"Isaac","email":"isaac@example.com"}`
	postReq := httptest.NewRequest("POST", "/mock/", bytes.NewBufferString(body))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	router.ServeHTTP(postW, postReq)

	// Get the user back
	getReq := httptest.NewRequest("GET", "/mock/1", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", getW.Code)
	}
}

func TestPutUser(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"Isaac","email":"isaac@example.com"}`
	postReq := httptest.NewRequest("POST", "/mock/", bytes.NewBufferString(body))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	router.ServeHTTP(postW, postReq)

	update := `{"name":"Updated","email":"updated@example.com"}`
	putReq := httptest.NewRequest("PUT", "/mock/1", bytes.NewBufferString(update))
	putReq.Header.Set("Content-Type", "application/json")
	putW := httptest.NewRecorder()
	router.ServeHTTP(putW, putReq)

	if putW.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", putW.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"Isaac","email":"isaac@example.com"}`
	postReq := httptest.NewRequest("POST", "/mock/", bytes.NewBufferString(body))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	router.ServeHTTP(postW, postReq)

	delReq := httptest.NewRequest("DELETE", "/mock/1", nil)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)

	if delW.Code != http.StatusNoContent {
		t.Fatalf("Expected 204 No Content, got %d", delW.Code)
	}
}
