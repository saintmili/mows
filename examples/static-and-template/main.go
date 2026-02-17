package main

import (
	"log"

	"github.com/saintmili/mows"
)

func main() {
	app := mows.New()

	app.DevMode(true)

	app.Static("/static", "./public")

	err := app.LoadTemplates("views/*.html")
	if err != nil {
		log.Fatal(err)
	}

	app.GET("/", func(c *mows.Context) error {
		return c.HTML(200, "home.html", mows.H{
			"title": "Mows",
			"name":  "User",
		})
	})

	app.Run(":8080")
}

