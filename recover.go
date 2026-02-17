package mows

import "net/http"

// Recover returns a middleware that recovers from panics.
//
// If a panic occurs, the middleware:
//
//   - Prevents the server from crashing
//   - Returns HTTP 500
//   - Sends the panic message as JSON
func Recover() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			defer func() {
				if err := recover(); err != nil {
					c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "internal server error",
					})
				}
			}()
			return next(c)
		}
	}
}
