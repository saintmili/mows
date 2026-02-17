package mows

import "strings"

// route represents a static route.
type route struct {
	handler     HandlerFunc
	middlewares []Middleware
}

// paramRoute represents a route containing path parameters.
type paramRoute struct {
	pathParts []string
	route     route
}

// Router stores registered routes and performs route matching.
//
// It supports static routes and parameterized paths such as:
//
//	/users/:id
type Router struct {
	staticRoutes map[string]map[string]route
	paramRoutes  map[string][]paramRoute
}

// NewRouter creates and initializes a new Router instance.
func NewRouter() *Router {
	return &Router{
		staticRoutes: make(map[string]map[string]route),
		paramRoutes:  make(map[string][]paramRoute),
	}
}

// GET registers a route that responds to HTTP GET requests.
func (e *Engine) GET(path string, handlers ...HandlerFunc) {
	e.rootGroup.GET(path, handlers...)
}

// POST registers a route that responds to HTTP POST requests.
func (e *Engine) POST(path string, handlers ...HandlerFunc) {
	e.rootGroup.POST(path, handlers...)
}

// PUT registers a route that responds to HTTP PUT requests.
func (e *Engine) PUT(path string, handlers ...HandlerFunc) {
	e.rootGroup.PUT(path, handlers...)
}

// DELETE registers a route that responds to HTTP DELETE requests.
func (e *Engine) DELETE(path string, handlers ...HandlerFunc) {
	e.rootGroup.DELETE(path, handlers...)
}

// find matches an incoming request path and returns the handler and params.
// Returns nil if no route matches.
func (r *Router) find(method, path string) (*route, map[string]string) {
	// check for static route
	if m, ok := r.staticRoutes[method]; ok {
		if route, ok := m[path]; ok {
			return &route, nil
		}
	}

	// check for param route
	requestParts := strings.Split(strings.Trim(path, "/"), "/")
	for _, paramRoute := range r.paramRoutes[method] {
		if len(paramRoute.pathParts) != len(requestParts) {
			continue
		}

		params := make(map[string]string)
		match := true

		for i, part := range paramRoute.pathParts {
			if strings.HasPrefix(part, ":") {
				key := part[1:]
				params[key] = requestParts[i]
				continue
			}

			if part != requestParts[i] {
				match = false
				break
			}
		}

		if match {
			return &paramRoute.route, params
		}
	}

	return nil, nil
}

// wrapHandlerAsMiddleware converts a HandlerFunc into a Middleware.
func wrapHandlerAsMiddleware(h HandlerFunc) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			err := h(c)
			if err != nil {
				return err
			}
			err = next(c)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

// hasParams checks whether a route path contains parameters (":id").
func hasParams(path string) bool {
	return strings.Contains(path, ":")
}

// addWithMiddleware registers a route and applies middleware chain.
func (r *Router) addWithMiddleware(method string, path string, handler HandlerFunc, middlewares []Middleware) {
	rt := route{
		handler:     handler,
		middlewares: middlewares,
	}

	// static route
	if !hasParams(path) {
		if r.staticRoutes[method] == nil {
			r.staticRoutes[method] = make(map[string]route)
		}
		r.staticRoutes[method][path] = rt
		return
	}

	// param route
	parts := strings.Split(strings.Trim(path, "/"), "/")
	r.paramRoutes[method] = append(r.paramRoutes[method], paramRoute{
		pathParts: parts,
		route:     rt,
	})
}
