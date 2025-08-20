package handler

import (
	"errors"
	"fmt"
	"live-chat-backend/internal/dto"
	"live-chat-backend/internal/service"
	errs "live-chat-backend/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		responseBody, status := errs.FormatBodyParserError(err)
		return c.Status(status).JSON(responseBody)
	}

	err = validate.Struct(payloadUser)
	if err != nil {
		responseBody, status := errs.FormatValidationError(err)
		return c.Status(status).JSON(responseBody)
	}

	registeredUser, err := h.authService.RegisterUser(payloadUser.Username, payloadUser.Email, payloadUser.Password)
	if err != nil {
		msg, status := h.formatRegisterUserError(err)
		return c.Status(status).JSON(fiber.Map{
			"error":   http.StatusText(status),
			"message": msg,
		})
	}

	token, err := h.authService.GenerateUserJWT(registeredUser.ID)
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
		msg, status := errs.FormatValidationError(err)
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

	token, err := h.authService.GenerateUserJWT(loggedUser.ID)
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

func (h *AuthHandler) formatRegisterUserError(err error) (msg string, httpStatus int) {
	if errors.Is(err, service.ErrPasswordPolicy) {
		// password policy error, we keep the message defined in auth service
		return err.Error(), fiber.StatusUnprocessableEntity
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		// db uniqueness error
		if strings.Contains(err.Error(), "users_email") {
			return "email is already in use", fiber.StatusConflict
		} else if strings.Contains(err.Error(), "users_username") {
			return "username is already taken", fiber.StatusConflict
		}
	}

	// unknown error
	log.Println(errs.UnknownErrorMessage, err)
	return errs.UnknownErrorMessage, fiber.StatusInternalServerError
}
