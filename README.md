# MOWS - My Own Web Server

A **minimal HTTP web framework for Go**, designed for **learning**, **simplicity**, and **zero boilerplate**.

It provides:
- Lightweight routing
- Middleware support (global & route-specific)
- Path parameters and route groups
- JSON binding & validation
- Centralized error handling
- Logging and panic recovery

This project is intentionally minimal and will grow over time.

## Installation

```bash
go get github.com/saintmili/mows
```

## Quick Start

```go
package main

import "mows"

func main() {
	app := mows.New()

	app.GET("/hello", func(c *mows.Context) error {
		c.JSON(200, map[string]string{
			"message": "hello world",
		})
        return nil
	})

	app.Run(":8080")
}
```

Run:

```bash
go run main.go
```

Visit `http://localhost:8080/hello`:

```json
{"message":"hello world"}
```

## Full Example: User API

See full example in `examples/user-api/main.go`

Features demonstrated:
- Global middleware
- Error handling with SetErrorHandler
- Route groups (/api)
- CRUD operations with path params (/users/:id)
- JSON binding and validation
- HTTP status codes (200, 201, 204, 400, 404)

## Features

| Feature          | Description                             |
| ---------------- | --------------------------------------- |
| Routing          | GET, POST, PUT, DELETE with path params |
| Route Groups     | Nested prefixes and middleware          |
| Middleware       | Global or route-specific middleware     |
| Logging          | Built-in request logger middleware      |
| Recovery         | Panic recovery middleware               |
| JSON Binding     | `BindJSON` & `BindAndValidate`          |
| Validation       | Struct validation using tags            |
| Error Handling   | Centralized `ErrorHandler`              |
| Response Helpers | `JSON()`, `String()`, `Status()`        |


## Engine

The Engine is the main application container.

```go
app := mows.New()
app.Run(":8080")
```

The engine implements `http.Handler` internaly.

## Basic routes

```go
app.GET("/users", handler)
app.POST("/users", handler)
app.PUT("/users/:id", handler)
app.DELETE("/users/:id", handler)
```

Handlers use this signature:

```go
func(c *mows.Context) error
```

## Route Groups

Groups allow shared prefixes and middleware.

```go
api := app.Group("/api")

api.GET("/users", listUsers)
api.GET("/users/:id", getUser)
```

Nested groups are supported:

```go
admin := api.Group("/admin")
admin.GET("/stats", statsHandler)
```

## Path Parameters

Define params using `:name`.

```go
app.GET("/users/:id", func(c *mows.Context) error {
    id := c.Param("id")

    c.JSON(200, map[string]string{
        "user_id": id,
    })

    return nil
})
```

## Sending Responses

### JSON response

```go
c.JSON(200, map[string]string{
    "status": "ok",
})
```

### Plain text

```go
c.String(200, "hello")
```

Access params:

```go
id := c.Param("id")
```

## Middleware

Middleware is the **heart of MOWS**

Middleware signature:

```go
type Middleware func(HandlerFunc) HandlerFunc
```

### Global middleware

```go
app.Use(mows.Logger(), mows.Recover())
```

Runs for **every request**.

### Route-specific middleware

```go
auth := func(next mows.HandlerFunc) mows.HandlerFunc {
    return func(c *mows.Context) error {
        if c.Request.Header.Get("Authorization") == "" {
            c.JSON(401, map[string]string{"error":"unauthorized"})
            return nil
        }
        return next(c)
    }
}

app.GET("/private", handler, auth)
```

Middleware order:

```go
Global → Group → Route → Handler
```

## Built-in Middleware

### Logger

Logs every request.

```go
app.Use(mows.Logger())
```

Example output:

```bash
[2026-02-17 19:55:01] 200 | 245µs | GET /hello | 127.0.0.1:54321
```

### Recovery

Prevents server crash on panic.

```go
app.Use(mows.Recover())
```

Returns:

```json
{ "error": "panic message" }
```

## JSON Binding

Bind request JSON to struct.

```go
type CreateUser struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

app.POST("/users", func(c *mows.Context) error {
    var req CreateUser

    if err := c.BindJSON(&req); err != nil {
        return err
    }

    c.JSON(200, req)
    return nil
})
```

## Request Lifecycle

For each request:
1. Request enters Engine
2. Global middleware runs
3. Group middleware runs
4. Route middleware runs
5. Handler executes
6. Logger prints result

## Project Structure Suggestion

For apps using MOWS:

```go
myapp/
│
├── main.go
├── handlers/
│   ├── users.go
│   └── auth.go
├── middleware/
│   └── auth.go
└── models/
```

## Testing Example

Use Go's built-in test server.

```go
func TestHelloRoute(t *testing.T) {
	app := mows.New()

	app.GET("/hello", func(c *mows.Context) error {
		c.String(200, "ok")
        return nil
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal("expected 200")
	}
}
```

## Roadmap (Upcoming)

Planned features:
- Static file serving
- HTML templates
- Request ID middleware
- Error handling / abort system
- Validation helpers
- OpenAPI generation

