package tests

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/saintmili/mows"
)

func TestQueryHelpers(t *testing.T) {
	app := mows.New()

	app.GET("/test", func(c *mows.Context) error {
		page := c.DefaultQueryInt("page", 1)
		c.Text(200, strconv.Itoa(page))
		return nil
	})

	req := httptest.NewRequest("GET", "/test?page=5", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Body.String() != "5" {
		t.Fatalf("expected 5 got %s", w.Body.String())
	}
}
