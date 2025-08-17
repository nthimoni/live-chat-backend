package utils

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func FormatValidationError(err error) (msg map[string]string, httpStatus int) {
	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		fieldName := strings.ToLower(string(e.Field()[0])) + e.Field()[1:]
		switch e.Tag() {
		case "required":
			errors[fieldName] = "this field is required"
		case "email":
			errors[fieldName] = "invalid email format"
		case "min":
			errors[fieldName] = "too short"
		default:
			errors[fieldName] = "invalid value"
		}
	}

	return errors, fiber.StatusUnprocessableEntity
}

func FormatRegisterUserError(err error) (msg string, httpStatus int) {
	if strings.Contains(err.Error(), "password must") {
		// password policy error, we keep the message defined in auth service
		return err.Error(), fiber.StatusUnprocessableEntity
	}

	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		// postgre uniqueness error
		if strings.Contains(err.Error(), "users_email") {
			return "email is already in use", fiber.StatusConflict
		} else if strings.Contains(err.Error(), "users_username") {
			return "username is already taken", fiber.StatusConflict
		}
	}

	// unknown error
	log.Println("unexpected error during user registration:", err)
	return "an error occurred while registering the user, please try again later", fiber.StatusInternalServerError
}
