package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/saintmili/mows"
)

func TestRouteGroupPrefix(t *testing.T) {
	app := mows.New()

	api := app.Group("/api")

	api.GET("/ping", func(c *mows.Context) error {
		c.Text(200, "pong")
		return nil
	})

	req := httptest.NewRequest("GET", "/api/ping", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Body.String() != "pong" {
		t.Fatalf("expected pong got %s", w.Body.String())
	}
}
