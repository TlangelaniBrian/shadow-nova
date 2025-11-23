package validator

import (
	"encoding/json"
	"net/http"
	"strings"

	"shadow-nova/backend/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Register custom validators
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidateRequest decodes JSON and validates the request
func ValidateRequest(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return ValidateStruct(v)
}

// WriteValidationError writes formatted validation errors to response
func WriteValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]string)
		for _, e := range validationErrs {
			details[strings.ToLower(e.Field())] = getErrorMessage(e)
		}
		
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Validation failed",
			Details: details,
		})
		return
	}
	
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: err.Error(),
	})
}

// Custom validator for strong passwords
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}
	
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Get user-friendly error messages
func getErrorMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	
	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "url":
		return field + " must be a valid URL"
	case "gte":
		return field + " must be greater than or equal to " + e.Param()
	case "strong_password":
		return field + " must contain uppercase, lowercase, number, and special character"
	default:
		return field + " is invalid"
	}
}
