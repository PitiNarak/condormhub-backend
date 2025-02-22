package utils

import (
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
