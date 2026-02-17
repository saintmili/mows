package mows

// RouterGroup represents a group of routes sharing a common prefix
// and optional middleware.
type RouterGroup struct {
	prefix      string
	middlewares []Middleware
	engine         *Engine
}

// Group creates a nested RouterGroup with an additional path prefix.
func (rg *RouterGroup) Group(prefix string, m ...Middleware) *RouterGroup {
	return &RouterGroup{
		prefix:      rg.prefix + prefix,
		middlewares: append(rg.middlewares, m...),
		engine:         rg.engine,
	}
}

// Use attaches middleware to the RouterGroup.
//
// Group middleware runs after global middleware but before
// route-specific middleware.
func (rg *RouterGroup) Use(m ...Middleware) {
	rg.middlewares = append(rg.middlewares, m...)
}

// GET registers a GET route inside the RouterGroup.
func (rg *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.engine.addRoute("GET", fullPath, rg.middlewares, handlers...)
}

// POST registers a POST route inside the RouterGroup.
func (rg *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.engine.addRoute("POST", fullPath, rg.middlewares, handlers...)
}

// PUT registers a PUT route inside the RouterGroup.
func (rg *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.engine.addRoute("PUT", fullPath, rg.middlewares, handlers...)
}

// DELETE registers a DELETE route inside the RouterGroup.
func (rg *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	fullPath := rg.prefix + path
	rg.engine.addRoute("DELETE", fullPath, rg.middlewares, handlers...)
}
