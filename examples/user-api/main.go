package main

import (
	"fmt"
	"errors"
	"net/http"

	"github.com/saintmili/mows"
)

// User represents a simple user model.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

var users = map[string]User{}

func main() {
	app := mows.New()

	// Global middleware
	app.Use(
		mows.Logger(),  // logs every request
		mows.Recover(), // recovers from panics
	)

	// Centralized error handler
	app.SetErrorHandler(func(c *mows.Context, err error) {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	})

	// Root route
	app.GET("/", func(c *mows.Context) error {
		c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to MOWS API",
		})
		return nil
	})

	// API group
	api := app.Group("/api")

	// Create a new user
	api.POST("/users", func(c *mows.Context) error {
		var req User
		if err := c.BindJSONAndValidate(&req); err != nil {
			return err
		}

		// generate a simple ID
		req.ID = fmt.Sprint(len(users) + 1 + '0')
		users[req.ID] = req

		return c.JSON(http.StatusCreated, req)
	})

	// Get all users
	api.GET("/users", func(c *mows.Context) error {
		result := make([]User, 0, len(users))
		for _, u := range users {
			result = append(result, u)
		}
		return c.JSON(http.StatusOK, result)
	})

	// Get single user
	api.GET("/users/:id", func(c *mows.Context) error {
		id := c.Param("id")
		user, ok := users[id]
		if !ok {
			return errors.New("user not found")
		}
		return c.JSON(http.StatusOK, user)
	})

	// Update user
	api.PUT("/users/:id", func(c *mows.Context) error {
		id := c.Param("id")
		user, ok := users[id]
		if !ok {
			return errors.New("user not found")
		}

		var req User
		if err := c.BindJSONAndValidate(&req); err != nil {
			return err
		}

		user.Name = req.Name
		user.Email = req.Email
		users[id] = user

		return c.JSON(http.StatusOK, user)
	})

	// Delete user
	api.DELETE("/users/:id", func(c *mows.Context) error {
		id := c.Param("id")
		_, ok := users[id]
		if !ok {
			return errors.New("user not found")
		}
		delete(users, id)
		return c.JSON(http.StatusNoContent, nil)
	})

	app.Run(":8080")
}
