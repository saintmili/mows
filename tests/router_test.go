package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/saintmili/mows"
)

func TestGETRoute(t *testing.T) {
	app := mows.New()

	app.GET("/hello", func(c *mows.Context) error {
		c.Text(200, "world")
		return nil
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	// call internal handler
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	if w.Body.String() != "world" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestGetRouteWithParam(t *testing.T) {
	app := mows.New()

	app.GET("/param/:param", func(c *mows.Context) error {
		c.Text(200, c.Param("param"))
		return nil
	})

	req := httptest.NewRequest("GET", "/param/abcd", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	if w.Body.String() != "abcd" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestGetRouteWithParams(t *testing.T) {
	app := mows.New()

	app.GET("/param/:param1/p/:param2", func(c *mows.Context) error {
		c.Text(200, c.Param("param1")+"/"+c.Param("param2"))
		return nil
	})

	req := httptest.NewRequest("GET", "/param/abcd/p/efgh", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	if w.Body.String() != "abcd/efgh" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestFullFlow(t *testing.T) {
	app := mows.New()
	app.Use(mows.Logger())

	app.GET("/ping", func(c *mows.Context) error {
		c.JSON(200, map[string]string{"pong": "yes"})
		return nil
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal("expected 200")
	}
}
