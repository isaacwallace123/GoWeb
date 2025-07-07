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

## 📚 Table of Contents

- [✨ Features](#-features)
- [🚀 Usage](#-usage)
  1. [Define Request/Response DTOs](#1-define-requestresponse-dtos)
  2. [Create a Controller](#2-create-a-controller)
  3. [Register Controllers in main.go](#3-register-controllers-in-maingo)
- [🧱 Core Concepts](#-core-concepts)
- [💡 Why These Matter](#-why-these-matter)
- [🧪 Response Builder](#-response-builder)
- [📌 Example JSON Response](#-example-json-response)
- [🌐 CORS Middleware](#-cors-middleware)
  - [✅ Registering the CORS Middleware](#-registering-the-cors-middleware)
  - [⚙️ Behavior](#-behavior)
  - [🛡 Example: Block all but GET](#-example-block-all-but-get)
- [❤️ Inspired By](#-inspired-by)

---

## 🚀 Usage

### 1. Define Request/Response DTOs

```go
// models/UserRequest.go
type UserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// models/UserResponse.go
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

func (c *UsersController) Get(userid int) *types.ResponseEntity {
     return ResponseEntity.Status(...).Body(...)
}
```
### 3. Register Controllers in `main.go`
```go
func main() {
    // Create router
    router := app.NewRouter()
    
    // Register controllers
    router.RegisterControllers(
        &controllers.UsersController{}, // You can create a controller through OOP 
        controllers.New(),              // You can create a controller using the builder
    )
	
    // Register pre-middleware 
    app.Use(middlewares.LoggingPre)
    
    // Register post-middleware
    app.UseAfter(middlewares.LoggingPost)
    
    // Launch the server
    logger.Info("Server listening on http://localhost:8080")
    err := router.Listen(":8080")
    if err != nil {
        logger.Fatal("Server failed to launch:", err)
    }
}
```

---

## 🧱 Core Concepts
GoWeb is built on a clean and extendable foundation inspired by Spring Boot, but optimized for Go. Below are the key architectural components of the framework:

| Concept                     | Description                                                                                                                                                                  |
| --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`Controller`**            | A struct implementing `BasePath()` and `Routes()` to define a route group. It may embed `ControllerBase` to enable optional controller-specific middleware.                  |
| **`Route`**                 | Defines a single HTTP endpoint via `Method`, `Path`, and a string-based `Handler` name that maps to a function in the controller.                                            |
| **DTOs**                    | Plain structs representing request or response data (e.g., `UserRequest`, `UserResponse`). GoWeb automatically binds path/query/body values to arguments.                    |
| **`ResponseEntity`**        | A fluent builder for setting status, body, and headers. Example: `ResponseEntity.Status(200).Body(data).Send(w)` or return it directly from controllers.                     |
| **`HttpStatus`**            | Enum-style constants for all HTTP status codes, e.g., `HttpStatus.OK`, `HttpStatus.CREATED`, etc., making your response code more readable.                                  |
| **`HttpMethod`**            | Enum-like constants for HTTP methods (`GET`, `POST`, etc.) and helpers like `IsValid(method)` to validate custom usage.                                                      |
| **`exception`**             | Standardized error response utilities like `BadRequestException(...)` or `InternalServerException(...)` that send JSON error responses with status codes.                    |
| **Middleware**              | Middleware objects implement the `Middleware` interface. They're registered globally or per controller using `app.Use(...)` or `controller.Use(...)`.                        |
| **Middleware Builder**      | Use `NewMiddlewareBuilder(...)` to create strongly-typed, reusable middleware with config (`.Config`), init logic (`.WithInit()`), and error hooks (`.OnError()`).           |
| **Request Context Helpers** | Access path params, query strings, and headers using `types.PathVar(ctx, "id")`, `QueryParam(ctx, "q")`, and `Header(ctx, "X-Token")`. Injected automatically by the router. |

### 💡 Why These Matter

- ✅ **Minimal and clean:** Core concepts like `Controller` and `Route` are simple, composable structs.

- ✅ **Extensible:** The middleware system uses generics and fluent chaining to support per-middleware config and lifecycle hooks.

- ✅ **Type-safe binding:** Reflect-based argument resolution injects only what your handler expects — nothing more.

- ✅ **Production-ready responses:** Use `ResponseEntity` and `exception` to consistently shape output without boilerplate.

- ✅ **Testable architecture:** Middleware and controllers can be unit tested with standard Go tools (`httptest`).

---

## 🧪 Response Builder
Example of just sending a message:
```go
ResponseEntity.Status(HttpStatus.OK).
    Body(map[string]string{"hello": "world"}).
    Header("X-Test", "Value").
    Send(w)
```
An example of sending an entity
```go
newUser := Models.UserResponseModel{
    Id:    newId,
    Name:  req.Name,
    Email: req.Email,
}

return ResponseEntity.Status(201).Body(newUser)
```

### 📌 Example JSON Response

```json
HTTP/1.1 201 Created
Content-Type: application/json
X-Custom: example

{
  "Id": 1
  "Name": "Test"
  "Email": "Test@example.com"
}
```
---
## 🌐 CORS Middleware

GoWeb includes a built-in CORS middleware that allows you to control which origins, methods, and headers are allowed to access your server across different domains. This is especially useful when building frontend-backend systems or public APIs.

### ✅ Registering the CORS Middleware

To enable it globally:
```go
app.Use(middlewares.CORS)
```

Then configure it as needed:
```go
middlewares.CORS.Config.AllowedOrigins = []string{"https://example.com"}
middlewares.CORS.Config.AllowedMethods = []string{"GET", "POST"}
middlewares.CORS.Config.AllowedHeaders = []string{"Content-Type", "Authorization"}
middlewares.CORS.Config.AllowCredentials = true
```
### ⚙️ Behavior
| Feature                   | Description                                                             |
| ------------------------- | ----------------------------------------------------------------------- |
| `AllowedOrigins`          | List of allowed domains (use `"*"` for all)                             |
| `AllowedMethods`          | List of allowed HTTP methods (`GET`, `POST`, etc.)                      |
| `AllowedHeaders`          | List of allowed request headers                                         |
| `AllowCredentials`        | Enables `Access-Control-Allow-Credentials: true`                        |
| Auto-Handles `OPTIONS`    | Returns `204 No Content` and skips route logic                          |
| Blocks Disallowed Methods | Returns `405 Method Not Allowed` if the request method is not permitted |

### 🛡 Example: Block all but GET

```go
app.Use(middlewares.CORS)
middlewares.CORS.Config.AllowedMethods = []string{"GET"}
```

If a client sends a `POST` request, the server will respond with:

```http request
HTTP/1.1 405 Method Not Allowed
Access-Control-Allow-Methods: GET
Access-Control-Allow-Origin: https://example.com
```

---

### ❤️ Inspired By

- Java Spring Boot
- Go’s minimalistic standard library
- Clean Architecture principles