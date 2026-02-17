package tests

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saintmili/mows"
)

func TestValidationFails(t *testing.T) {
	app := mows.New()

	app.POST("/users", func(c *mows.Context) error {
		var body struct {
			Name string `json:"name" validate:"required"`
		}

		if err := c.BindJSONAndValidate(&body); err != nil {
			return err
		}

		c.Text(200, "ok")
		return nil
	})

	req := httptest.NewRequest("POST", "/users",
		strings.NewReader(`{"name":""}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatal("validation should fail")
	}
}
