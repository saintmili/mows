package mows

// HandlerFunc defines a request handler used by MOWS.
type HandlerFunc func(*Context) error

// Middleware defines a function that wraps a HandlerFunc.
//
// Middleware can run logic before and/or after the next handler.
type Middleware func(HandlerFunc) HandlerFunc

// Use registers global middleware that runs for every request.
//
// Middleware execution order:
//
//	Global → Group → Route → Handler
func (e *Engine) Use(m ...Middleware) {
	e.middlewares = append(e.middlewares, m...)
}

// buildChain combines middleware and handler into a single HandlerFunc.
// Middleware are applied in reverse order to preserve execution order.
func (e *Engine) buildChain(final HandlerFunc) HandlerFunc {
	h := final

	// reverse order so first middleware runs first
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		h = e.middlewares[i](h)
	}

	return h
}
