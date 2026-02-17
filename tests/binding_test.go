package tests

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saintmili/mows"
)

func TestBindJSON(t *testing.T) {
	app := mows.New()

	app.POST("/users", func(c *mows.Context) error {
		var body struct {
			Name string `json:"name"`
		}

		if err := c.BindJSON(&body); err != nil {
			return err
		}

		c.Text(200, body.Name)
		return nil
	})

	req := httptest.NewRequest("POST", "/users",
		strings.NewReader(`{"name":"mows"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Body.String() != "mows" {
		t.Fatalf("expected mows got %s", w.Body.String())
	}
}
