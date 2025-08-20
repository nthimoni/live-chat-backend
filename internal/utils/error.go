package errs

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var ErrUnknown = errors.New("unknown error occurred")

var UnknownErrorMessage = "An Unexpected Error Occurred, please try again later"
var InvalidJSONMessage = "Invalid JSON body"

func FormatValidationError(err error) (responseBody fiber.Map, httpStatus int) {

	//type assertion to know if it's a validation error
	ve, isValidationError := err.(validator.ValidationErrors)
	if isValidationError {
		errorsMap := make(map[string]string)
		for _, e := range ve {
			fieldName := strings.ToLower(string(e.Field()[0])) + e.Field()[1:]
			switch e.Tag() {
			case "required":
				errorsMap[fieldName] = "this field is required"
			case "email":
				errorsMap[fieldName] = "invalid email format"
			case "min":
				errorsMap[fieldName] = "too short"
			default:
				errorsMap[fieldName] = "invalid value"
			}
		}

		return fiber.Map{
			"error":   http.StatusText(fiber.StatusUnprocessableEntity),
			"message": errorsMap,
		}, fiber.StatusUnprocessableEntity
	}

	// Fallback for unknown errors
	return fiber.Map{
		"error":   http.StatusText(fiber.StatusInternalServerError),
		"message": UnknownErrorMessage,
	}, fiber.StatusInternalServerError
}

// FormatBodyParserError formats errors returned by BodyParser
func FormatBodyParserError(err error) (responseBody fiber.Map, httpStatus int) {
	return fiber.Map{
		"error":   http.StatusText(fiber.StatusBadRequest),
		"message": InvalidJSONMessage,
	}, fiber.StatusBadRequest

}
