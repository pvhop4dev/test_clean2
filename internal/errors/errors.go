package errors

import (
	"clean-arch-go/internal/pkg/i18n"
	"fmt"
	"golang.org/x/text/language"
)

type AppError struct {
	Code    string
	Message string
	Detail  interface{}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new AppError
func NewAppError(code string, message string, detail interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// NewValidationError creates a validation error
func NewValidationError(field string, message string) *AppError {
	return NewAppError("VALIDATION_ERROR", message, map[string]interface{}{"field": field})
}

// NewNotFoundError creates a not found error
func NewNotFoundError(entity string) *AppError {
	return NewAppError("NOT_FOUND", fmt.Sprintf("%s not found", entity), nil)
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppError("BAD_REQUEST", message, nil)
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(message string) *AppError {
	return NewAppError("INTERNAL_ERROR", message, nil)
}

// Translate translates the error message using i18n
func (e *AppError) Translate(lang string) string {
	languageTag := language.MustParse(lang)
	localizer := i18n.GetLocalizer()
	
	templateData := map[string]interface{}{
		"message": e.Message,
	}
	
	if e.Detail != nil {
		templateData["detail"] = e.Detail
	}
	
	translated, err := localizer.Translate(languageTag, e.Code, templateData)
	if err != nil {
		// Fallback to the original message if translation fails
		return e.Message
	}
	return translated
}
