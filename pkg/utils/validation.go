package utils

import (
	"regexp"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/go-playground/validator/v10"
)

// Custom validator for Lifestyles
func ValidateLifestyles(fl validator.FieldLevel) bool {
	lifestyles, ok := fl.Field().Interface().(domain.LifestyleArray)
	if !ok {
		return false
	}

	for _, l := range lifestyles {
		if !l.IsValid() {
			return false
		}
	}
	return true
}

// Custom validator for Phone
func ValidatePhone(fl validator.FieldLevel) bool {
	var phoneRegex = regexp.MustCompile(`^0\d{2}-\d{3}-\d{4}$`)
	return phoneRegex.MatchString(fl.Field().String())
}
