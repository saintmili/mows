package mows

import "net/http"

// responseWriter wraps http.ResponseWriter and captures
// the status code and response size.
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// NewResponseWriter wraps http.ResponseWriter and tracks status code and size.
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         200,
	}
}

// WriteHeader captures the response status code.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the response body and tracks the response size.
func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
