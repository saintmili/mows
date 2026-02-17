# Contributing to MOWS

Thanks for your interest in contributing ❤️

MOWS is a learning-driven framework.

The main goals are **clarity**, **simplicity**, and **small abstractions**.

Before contributing, please read this guide.

# Project Philosophy

The goals are:
- Educational
- Minimal
- Clean codebase
- Easy to read in one sitting
- Zero unnecessary magic

If a feature adds complexity without clear value -> it will likely be rejected.

# Project Structure

```go
mows/
│
├── engine.go
├── router.go
├── context.go
├── group.go
├── middleware.go
├── logger.go
├── recover.go
├── response_writer.go
└── internal tests (*.go)
```

Guidelines:
- Keep files small and focused
- Prefer multiple small files over one large file
- Avoid circular dependencies


# Code Style

We follow standard Go conventions:

```bash
go fmt ./...
go vet ./...
```

Naming rules:

| Type     | Convention      |
| -------- | --------------- |
| exported | PascalCase      |
| private  | camelCase       |
| acronyms | JSON, HTTP, URL |

# Commit Style

Use clear commit message.

Good:
```
add route groups middleware support
fix panic in JSON binding
refactor router param matching
```

Bad:
```
fix stuff
update code
changes
```

# Branching

```
main → stable
feature/* → new features
fix/* → bug fixes
```

# How to Propose a Feature

Open an issue describing:
- The problem
- Why MOWS should support it
- Proposed design (optional)

We prefer discussion before large changes.

# Running Tests

```bash
go test ./...
```

All tests must pass before PR is accepted.

# Writing Tests for MOWS

Testing the framework is **very important** before adding new features.

We rely on Go's standard testing tools:
- `testing`
- `net/http/httptest`

# Testing Strategy

We test the framework like a user would use it.

Each feature must have tests for:

| Scenario         | Example               |
| ---------------- | --------------------- |
| Happy path       | valid request works   |
| Error path       | bad input handled     |
| Edge cases       | missing params, panic |
| Middleware chain | order is correct      |

# Example Test File

create: `engine_test.go`

```go
package mows

import (
	"net/http"
	"net/http/httptest"
	"testing"
)
```

## Test 1 - Basic Route

```go
func TestBasicRoute(t *testing.T) {
	app := New()

	app.GET("/ping", func(c *Context) {
		c.String(200, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if w.Body.String() != "pong" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
```

## Test 2 - Path Params

Create `router_test.go`

```go
func TestPathParams(t *testing.T) {
	app := New()

	app.GET("/users/:id", func(c *Context) {
		c.String(200, c.Param("id"))
	})

	req := httptest.NewRequest("GET", "/users/42", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Body.String() != "42" {
		t.Fatalf("expected param 42, got %s", w.Body.String())
	}
}
```

## Test 3 — Middleware Execution Order

```go
func TestMiddlewareOrder(t *testing.T) {
	app := New()
	order := ""

	m1 := func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			order += "A"
			next(c)
		}
	}

	m2 := func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			order += "B"
			next(c)
		}
	}

	app.Use(m1, m2)

	app.GET("/test", func(c *Context) {
		order += "H"
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if order != "ABH" {
		t.Fatalf("wrong order: %s", order)
	}
}
```

## Test 4 — Recovery Middleware

```go
func TestRecoveryMiddleware(t *testing.T) {
	app := New()
	app.Use(Recover())

	app.GET("/panic", func(c *Context) {
		panic("boom")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
```

## Test 5 — JSON Binding

```go
func TestJSONBinding(t *testing.T) {
	app := New()

	type User struct {
		Name string `json:"name"`
	}

	app.POST("/users", func(c *Context) {
		var u User
		err := c.BindJSON(&u)
		if err != nil {
			t.Fatal(err)
		}
		c.String(200, u.Name)
	})

	body := `{"name":"mows"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Body.String() != "mows" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
```

## Test Coverage

Run coverage:

``` bash
go test -cover ./...
```

Goal for framework core: >80%

# PR Checklist

Before submitting:

- [ ] Code formatted
- [ ] Tests added/updated
- [ ] No breaking changes (or documented)
- [ ] README updated if needed

# Thank You ❤️

MOWS is a small project, and every contribution matters.
