package main

import (
	"live-chat-backend/internal/auth"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}

	app := fiber.New()

	// Global middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running !")
	})

	// routes
	authGroup := app.Group("/auth")
	auth.SetupRoutes(authGroup)

	app.Listen(":" + port)
}
