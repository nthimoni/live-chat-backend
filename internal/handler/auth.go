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

	var payloadUser dto.RegisterUserRequest

	err := c.BodyParser(&payloadUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   http.StatusText(fiber.StatusBadRequest),
			"message": "Invalid JSON body",
		})
	}

	err = validate.Struct(payloadUser)
	if err != nil {
		msg, status := utils.FormatValidationError(err)
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": msg,
		})
	}

	registeredUser, err := h.authService.RegisterUser(payloadUser.Username, payloadUser.Email, payloadUser.Password)
	if err != nil {
		msg, status := utils.FormatRegisterUserError(err)
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": msg,
		})
	}

	token, err := h.authService.GenerateUserJWT(registeredUser.ID, registeredUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   http.StatusText(fiber.StatusInternalServerError),
			"message": "Failed to generate JWT",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.RegisterUserResponse{
		User: dto.UserDTO{
			ID:       fmt.Sprint(registeredUser.ID),
			Email:    registeredUser.Email,
			Username: registeredUser.Username,
		},
		Token: token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	var payloadUser dto.LoginRequest

	err := c.BodyParser(&payloadUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   http.StatusText(fiber.StatusBadRequest),
			"message": "Invalid JSON body",
		})
	}

	err = validate.Struct(payloadUser)
	if err != nil {
		msg, status := utils.FormatValidationError(err)
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": msg,
		})
	}

	loggedUser, err := h.authService.LoginUser(payloadUser.Email, payloadUser.Password)
	if err != nil {
		var status int
		if err.Error() == "invalid credentials" {
			status = fiber.StatusUnauthorized
		} else {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": err.Error(),
		})
	}

	token, err := h.authService.GenerateUserJWT(loggedUser.ID, loggedUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   http.StatusText(fiber.StatusInternalServerError),
			"message": "Failed to generate JWT, retry later",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.RegisterUserResponse{
		User: dto.UserDTO{
			ID:       fmt.Sprint(loggedUser.ID),
			Email:    loggedUser.Email,
			Username: loggedUser.Username,
		},
		Token: token,
	})
}
