package mows

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Context carries request and response data across handlers and middleware.
//
// It provides helper methods for:
//
//   - Sending responses
//   - Reading request data
//   - Accessing path parameters
type Context struct {
	Writer  *responseWriter
	Request *http.Request
	Params  map[string]string
	Status  int
	engine  *Engine
}

// NewContext creates a new Context for the incoming HTTP request.
// It is used internally by the Engine during request dispatch.
func NewContext(w *responseWriter, r *http.Request, engine *Engine) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Params:  make(map[string]string),
		Status:  200,
		engine:  engine,
	}
}

// JSON sends a JSON response with the provided status code.
func (c *Context) JSON(code int, v any) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer.ResponseWriter).Encode(v)
}

// String sends a plain text response.
func (c *Context) Text(code int, s string) error {
	c.Writer.WriteHeader(code)
	_, err := c.Writer.Write([]byte(s))
	return err
}

// Param returns the value of a path parameter.
//
// Example:
//
//	app.GET("/users/:id", func(c *Context) error {
//	    id := c.Param("id")
//		return nil
//	})
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// BindJSON parses the request body as JSON into the provided struct.
//
// Returns an error if:
//
//   - Content-Type is not application/json
//   - JSON is malformed
//   - Decoding fails
func (c *Context) BindJSON(v any) error {
	if c.Request.Body == nil {
		return errors.New("request body is empty")
	}

	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return errors.New("content-type must be application/json")
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return errors.New("empty json body")
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}

	return nil
}

// Query returns the value of a query parameter.
//
// Example:
//
//	/users?page=2 → c.Query("page") == "2"
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// DefaultQuery returns the query value or a default value if the key is missing.
func (c *Context) DefaultQuery(key string, defaultValue string) string {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// QueryInt returns a query parameter parsed as an int.
// Returns an error if the value cannot be converted.
func (c *Context) QueryInt(key string) (int, error) {
	val := c.Query(key)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// DefaultQueryInt returns a query parameter as int or a default value if missing or invalid.
func (c *Context) DefaultQueryInt(key string, defaultValue int) int {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return i
}

// QueryBool returns a query parameter parsed as a boolean.
// Accepted values: true, false, 1, 0.
func (c *Context) QueryBool(key string) (bool, error) {
	val := c.Query(key)
	if val == "" {
		return false, nil
	}
	return strconv.ParseBool(val)
}

// DefaultQueryBool returns a query parameter as bool or a default value if missing or invalid.
func (c *Context) DefaultQueryBool(key string, defaultValue bool) bool {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return b
}

// Validate validates a struct using the configured validator.
// Returns an error if validation fails.
func (c *Context) Validate(v any) error {
	return c.engine.validate.Struct(v)
}

// BindJSONAndValidate binds JSON request body into the struct and validates it.
// This is a convenience helper combining BindJSON and Validate.
func (c *Context) BindJSONAndValidate(v any) error {
	if err := c.BindJSON(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	return nil
}

// HTML renders an HTML template.
func (c *Context) HTML(code int, name string, data any) error {
	if c.engine.templates == nil {
		return ErrTemplatesNotLoaded
	}

	if c.engine.devMode {
		if err := c.engine.templates.load(); err != nil {
			return err
		}
	}

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(code)

	return c.engine.templates.tmpl.ExecuteTemplate(c.Writer, name, data)
}
