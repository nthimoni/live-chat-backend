package main

import (
	"fmt"
	"live-chat-backend/internal/handler"
	"live-chat-backend/internal/models"
	"live-chat-backend/internal/repository"
	"live-chat-backend/internal/service"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET not set")
	}

	app := fiber.New()

	db, err := connectToDatabase()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Instanciate the repositories
	userRepo := repository.NewUserRepository(db)

	// Instanciate the services
	authService := service.NewAuthService(jwtSecret, userRepo)

	// Instanciate the handlers
	authHandler := handler.NewAuthHandler(authService)

	// Global middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New())

	// Register routes
	api := app.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running !")
	})

	app.Listen(":" + port)
}

func connectToDatabase() (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
