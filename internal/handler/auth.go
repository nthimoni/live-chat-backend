package handler

import (
	"fmt"
	"live-chat-backend/internal/dto"
	"live-chat-backend/internal/service"
	"live-chat-backend/internal/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

var validate = validator.New()

func (h *AuthHandler) Register(c *fiber.Ctx) error {

	var user dto.RegisterUserRequest

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON body",
		})
	}

	err = validate.Struct(user)
	if err != nil {
		msg, status := utils.FormatValidationError(err)
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": msg,
		})
	}

	userModel, err := h.authService.RegisterUser(user.Username, user.Email, user.Password)
	if err != nil {
		msg, status := utils.FormatRegisterUserError(err)
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": msg,
		})
	}

	token, err := h.authService.GenerateUserJWT(userModel.ID, userModel.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   http.StatusText(fiber.StatusInternalServerError),
			"message": "Failed to generate JWT",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.RegisterUserResponse{
		User: dto.UserDTO{
			ID:       fmt.Sprint(userModel.ID),
			Email:    userModel.Email,
			Username: userModel.Username,
		},
		Token: token,
	})
}

func (authHandler *AuthHandler) Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}
