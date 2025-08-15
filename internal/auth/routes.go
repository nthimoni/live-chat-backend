package auth

import "github.com/gofiber/fiber/v2"

// SetupRoutes attaches all auth routes to the given fiber app or group
func SetupRoutes(app fiber.Router) {
	app.Post("/register", func(c *fiber.Ctx) error {
		return c.SendString("Register")
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})
}
