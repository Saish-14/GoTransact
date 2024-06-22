package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidatePassword checks if the password meets the complexity requirements
func ValidatePassword(fl validator.FieldLevel) bool {
	Password := fl.Field().String()

	var (
		hasMinLen    = len(Password) >= 8
		hasUpperCase = regexp.MustCompile(`[A-Z]`).MatchString(Password)
		hasLowerCase = regexp.MustCompile(`[a-z]`).MatchString(Password)
		hasNumber    = regexp.MustCompile(`[0-9]`).MatchString(Password)
		hasSpecial   = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(Password)
	)

	return hasMinLen && hasUpperCase && hasLowerCase && hasNumber && hasSpecial
}

var validate *validator.Validate

// Init initializes the custom validator
func Init() {
	validate = validator.New()
	validate.RegisterValidation("password_complexity", ValidatePassword)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
