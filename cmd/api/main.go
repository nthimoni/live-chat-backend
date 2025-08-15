package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	// middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running !")
	})

	app.Listen(":" + port)
}
