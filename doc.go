/*
Package mows is a minimal HTTP web framework designed for learning
and building small web services.

Example:

	app := mows.New()
	app.Use(mows.Logger(), mows.Recover())

	app.GET("/hello", func(c *mows.Context) error {
	    c.JSON(200, map[string]string{"message":"hello"})
		return error
	})

	app.Run(":8080")
*/
package mows
