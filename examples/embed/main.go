package main

import (
	"embed"
	"github.com/saintmili/mows"
)

//go:embed public/*
var publicFS embed.FS

//go:embed views/*
var viewsFS embed.FS

func main() {
	app := mows.New()

	// embedded static files
	app.StaticFS("/static", publicFS, "public")

	// embedded templates
	app.LoadTemplatesFS(viewsFS, "views/*.html")

	app.GET("/", func(c *mows.Context) error {
		return c.HTML(200, "home.html", mows.H{
			"title": "MOWS with embed",
			"name":  "User",
		})
	})

	app.Run(":8080")
}

