package middleware

import (
	"live-chat-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func NewJWTAuth(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		tokenStr := authHeader[len("Bearer "):]

		user, err := authService.ValidateTokenAndGetUser(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		// TODO: handle the different error types from ValidateTokenAndGetUser (500 error if  db connection error)

		c.Locals("user", user)

		return c.Next()
	}
}
