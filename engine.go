// Package mows provides a minimal HTTP web framework focused on
// simplicity, middleware-first design, and clean routing.
//
// MOWS is intended as a lightweight alternative for learning and
// building small services without heavy abstractions.
package mows

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
)

// Engine is the main application instance.
//
// Engine implements http.Handler and is responsible for:
//
//   - Route registration
//   - Global middleware
//   - Request dispatching
//
// Create a new engine using New().
type Engine struct {
	router       *Router
	middlewares  []Middleware
	server       *http.Server
	rootGroup    *RouterGroup
	validate     *validator.Validate
	errorHandler ErrorHandler
	templates    *TemplateEngine
	devMode      bool
}

// New creates and returns a new Engine instance.
//
// Example:
//
//	app := mows.New()
func New() *Engine {
	engine := &Engine{
		router:   NewRouter(),
		validate: validator.New(),
	}
	engine.rootGroup = &RouterGroup{
		engine: engine,
	}
	engine.errorHandler = defaultErrorHandler

	return engine
}

// Run starts the HTTP server and listens on the given address.
//
// This is a helper wrapper around http.ListenAndServe.
//
// Example:
//
//	app.Run(":8080")
func (e *Engine) Run(addr string) error {
	e.server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(e.handle),
	}

	// channel to listen to errors
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("🚀 Mows server running on %s", addr)
		serverErr <- e.server.ListenAndServe()
	}()

	// listen for ctrl+c / sigterm
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		log.Printf("🛑 Received signal: %s. Shutting down...", sig)
		return e.shutdown()
	}
}

// Run starts the HTTPS server and listens on the given address.
//
// This is a helper wrapper around http.ListenAndServeTLS.
//
// Example:
//
//	app.RunTLS(":443")
func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	e.server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(e.handle),
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("🔒 Mows HTTPS server running on %s", addr)
		serverErr <- e.server.ListenAndServeTLS(certFile, keyFile)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		log.Printf("🛑 Received signal: %s. Shutting down...", sig)
		return e.shutdown()
	}
}

// shutdown gracefully stops the HTTP server.
// It is used internally for graceful shutdown support.
func (e *Engine) shutdown() error {
	// wait 5 seconds for active requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.server.Shutdown(ctx); err != nil {
		log.Printf("❌ Graceful shutdown failed: %v", err)
		return err
	}

	log.Println("✅ Server stopped gracefully")
	return nil
}

// ServeHTTP implements the http.Handler interface.
// It should not be called directly by users.
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := newResponseWriter(w)
	e.handle(rw, r)
}

// handle is the internal request dispatcher for the Engine.
//
// It performs the following steps:
//
//  1. Wraps the ResponseWriter to track status and size.
//  2. Creates a new Context for the request.
//  3. Matches the request path and method against registered routes.
//  4. Returns 404 if no route is found.
//  5. Sets path parameters in the Context.
//  6. Wraps the route handler with route-specific middleware.
//  7. Wraps the result with global middleware via buildChain.
//  8. Executes the final handler and forwards any error to the configured error handler.
//
// Note: This function implements the core of the request lifecycle
// and should not be called directly by users.
func (e *Engine) handle(w http.ResponseWriter, r *http.Request) {
	rw := newResponseWriter(w)
	ctx := NewContext(rw, r, e)

	route, params := e.router.find(r.Method, r.URL.Path)
	if route == nil {
		http.NotFound(w, r)
		return
	}

	ctx.Params = params

	// route-specific chain
	h := route.handler
	for i := len(route.middlewares) - 1; i >= 0; i-- {
		h = route.middlewares[i](h)
	}

	// global middlewares
	final := e.buildChain(h)
	if err := final(ctx); err != nil {
		e.errorHandler(ctx, err)
	}
}

// Group creates a new RouterGroup with the provided path prefix.
//
// Example:
//
//	api := app.Group("/api")
func (e *Engine) Group(prefix string, m ...Middleware) *RouterGroup {
	return &RouterGroup{
		prefix:      prefix,
		middlewares: m,
		engine:      e,
	}
}

// addRoute registers a new route in the router.
// It is used internally by HTTP method helpers (GET, POST, etc).
func (e *Engine) addRoute(method string, path string, middlewares []Middleware, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("route must have at least one handler")
	}

	final := handlers[len(handlers)-1]
	routeHandlers := handlers[:len(handlers)-1]

	var routeMiddlewares []Middleware
	for _, h := range routeHandlers {
		routeMiddlewares = append(routeMiddlewares, wrapHandlerAsMiddleware(h))
	}

	allMiddlewares := append(middlewares, routeMiddlewares...)

	e.router.addWithMiddleware(method, path, final, allMiddlewares)
}

// ErrorHandler defines a centralized error handling function.
//
// It is invoked when a handler returns or triggers an error.
type ErrorHandler func(*Context, error)

// defaultErrorHandler is the fallback error handler used by the Engine.
func defaultErrorHandler(c *Context, err error) {
	c.JSON(http.StatusBadRequest, serverError{
		Error: err.Error(),
	})
}

// SetErrorHandler replaces the default error handler.
//
// This allows applications to customize how errors are returned.
func (e *Engine) SetErrorHandler(h ErrorHandler) {
	e.errorHandler = h
}

// Static serves files from a directory under a given URL prefix.
// Example:
//
//	app.Static("/static", "./public")
//
// Then /static/app.css -> ./public/app.css
func (e *Engine) Static(prefix string, root string) {
	fs := http.FileServer(http.Dir(root))

	// We use a wildcard route internally
	routePath := path.Join(prefix, "/*filepath")

	e.GET(routePath, func(c *Context) error {
		// remove prefix from URL path
		http.StripPrefix(prefix, fs).ServeHTTP(c.Writer, c.Request)
		return nil
	})
}

// DevMode enables or disables hot reload for templates.
//
// When enabled, Engine will reload filesystem templates before each render,
// allowing developers to see changes instantly without restarting the server.
//
// Note: Hot reload does NOT apply to templates loaded from embed.FS.
func (e *Engine) DevMode(enable bool) {
	e.devMode = enable
}
