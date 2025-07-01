# GoWeb — A Spring Boot-style Web Framework in Go

GoWeb is a lightweight Go web framework that mimics the structure, design, and fluent API of Java Spring Boot. It embraces convention over configuration using reflection-based routing, request/response DTOs, and a fluent `ResponseEntity` builder.

---

## ✨ Features

- ✅ Spring Boot–style `Controller` structure (`BasePath()` + `Routes()`)
- ✅ Automatic path variable and request body binding
- ✅ Strongly-typed `Request` and `Response` DTOs
- ✅ Fluent `ResponseEntity` response builder
- ✅ `HttpStatus` constants (e.g. `httpstatus.CREATED`)
- ✅ Reflection-powered routing using method signatures
- ✅ Minimal, clean design using Go standard library only

---

## 🚀 Usage

### 1. Define Request/Response DTOs

```go
// Models/UserRequest.go
type UserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Models/UserResponse.go
type UserResponse struct {
    Id    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```
### 2. Create a Controller
```go
type UsersController struct{}

func (c *UsersController) BasePath() string {
    return "/users"
}

func (c *UsersController) Routes() []core.RouteEntry {
    return []core.RouteEntry{
        {Method: "GET", Path: "/", Handler: "GetAll"},
        {Method: "GET", Path: "/{userid}", Handler: "Get"},
        {Method: "POST", Path: "/", Handler: "Post"},
        {Method: "PUT", Path: "/{userid}", Handler: "Put"},
        {Method: "DELETE", Path: "/{userid}", Handler: "Delete"},
    }
}

func (c *UsersController) Get(userid int) *response.ResponseEntity {
    // ...
}
```
### 3. Register Controllers in `main.go`
```go
func main() {
    router := core.RegisterControllers(
        &handlers.UsersController{},
    )

    core.Use(core.LoggingMiddleware)

    log.Println("🚀 Server listening on http://localhost:8080")
    err := router.Listen(":8080")
    if err != nil {
        log.Fatal("❌ Server failed:", err)
    }
}
```

---

### 🧱 Core Concepts

| Concept           | Description                                                     |
| ----------------- | --------------------------------------------------------------- |
| `Controller`      | A struct implementing `BasePath()` and `Routes()`               |
| `RouteEntry`      | Defines each HTTP method, subpath, and method handler name      |
| DTOs              | Define `UserRequest`, `UserResponse`, etc. in `Models/` package |
| `BindArguments()` | Automatically binds path vars and JSON body to method arguments |
| `ResponseEntity`  | Fluent builder for response body, status, and headers           |


### 🧪 Response Builder Example

```go
response.Status(httpstatus.OK).
    Body(map[string]string{"hello": "world"}).
    Header("X-Test", "Value").
    Send(w)
```

### 📌 Example JSON Response

```json
HTTP/1.1 201 Created
Content-Type: application/json
X-Custom: example

{
  "message": "User created",
  "id": 1
}
```

### ❤️ Inspired By

- Java Spring Boot
- Go’s minimalistic standard library
- Clean Architecture principles