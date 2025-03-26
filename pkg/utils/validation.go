package utils

import (
	"errors"
	"regexp"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/go-playground/validator/v10"
)

// Custom validator for Lifestyles
func ValidateLifestyles(fl validator.FieldLevel) bool {
	lifestyles, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, l := range lifestyles {
		if !domain.Lifestyle(l).IsValid() {
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

func ValidateRole(fl validator.FieldLevel) bool {
	if fl.Field().IsValid() {
		role := fl.Field().String()
		return role == "ADMIN" || role == "LESSOR" || role == "LESSEE"
	}
	return false

}

func ValidateUserForReview(user *domain.User, history *domain.LeasingHistory, create bool) error {
	if create && history.ReviewFlag {
		return apperror.BadRequestError(errors.New("review already exist"), "review already exist")
	} else if !create && !history.ReviewFlag {
		return apperror.BadRequestError(errors.New("review not exist"), "review not exist")
	}
	if user.Role == domain.LessorRole {
		return apperror.UnauthorizedError(errors.New("user is unauthorized"), "user is unauthorized")
	}
	if history.LesseeID != user.ID {
		return apperror.UnauthorizedError(errors.New("user is unauthorized"), "user is unauthorized")
	}
	return nil
}
