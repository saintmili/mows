package mows

import (
	"fmt"
	"time"
)

// Logger returns a middleware that logs HTTP requests.
//
// Logged information includes:
//
//   - Status code
//   - Request latency
//   - HTTP method
//   - Request path
//   - Client IP
func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				return err
			}

			latency := time.Since(start)
			method := c.Request.Method
			path := c.Request.URL.Path
			status := c.Writer.status
			ip := c.Request.RemoteAddr

			fmt.Printf(
				"[%s] %d | %v | %s %s | %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				status,
				latency,
				method,
				path,
				ip,
			)
			return nil
		}
	}
}
