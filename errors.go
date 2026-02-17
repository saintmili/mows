package mows

import "errors"

var ErrTemplatesNotLoaded = errors.New("templates not loaded")

// serverError represents a JSON structure for internal server errors.
//
// It is used to send a consistent error response when a handler fails
// or a panic is recovered.
//
// Example JSON response:
//
//	{ "error": "something went wrong" }
type serverError struct {
	Error string `json:"error"`
}
