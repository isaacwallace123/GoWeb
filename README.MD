# GoWeb — A Spring Boot-style Web Framework in Go

GoWeb is a lightweight Go web framework that mimics the structure, design, and fluent API of Java Spring Boot. It includes decorators like `@Controller`, `@GetMapping`, `@PostMapping`, and a fluent `ResponseEntity` builder.

---

## ✨ Features

- ✅ `@Controller`-style handler registration
- ✅ `@GetMapping`, `@PostMapping` decorators
- ✅ Automatic `@RequestBody` binding to structs
- ✅ Fluent `ResponseEntity` response builder (like Spring Boot)
- ✅ `HttpStatus` constants (e.g. `httpstatus.CREATED`)
- ✅ Minimal, clean design using the Go standard library

---

## 🚀 Usage

### 1. Define a Controller

```go
package controllers

type CreateRequest struct {
    Name string `json:"name"`
}

type ContainerController struct {
    // Injected service, etc.
}

func (c *ContainerController) Create(req CreateRequest) *response.ResponseEntity {
    return response.Status(httpstatus.CREATED).
        Body(map[string]string{"message": "Created " + req.Name})
}
```

### 2. Register Routes in main.go

```go
func main() {
    router := core.NewRouter()

    decorators.Controller("/containers", &controllers.ContainerController{})
    decorators.PostMapping("", (*controllers.ContainerController).Create)

    router.RegisterControllers()
    router.Listen(":8080")
}

```

## 🧱 Components
- ✅ `@Controller(basePath string, instance any)` Registers a controller instance with a base route.
- ✅ `@GetMapping(path string, handler any)` Registers a GET route method.
- ✅ `@PostMapping(path string, handler any)` Registers a POST route method.
- ✅ `ResponseEntity` Fluent builder for HTTP responses, like in Spring Boot:

```go
response.Status(httpstatus.OK).
    Body(map[string]string{"hello": "world"}).
    Header("X-Test", "Value").
    Send(w)
```
---
## 📌 Example Response

```json
HTTP/1.1 201 Created
Content-Type: application/json
X-Custom: example

{
  "message": "Created handlers"
}
```
---
## ❤️ Inspired By
- Java Spring Boot
- Go's minimalistic standard library