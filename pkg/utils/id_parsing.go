package utils

import (
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ParseIdParam(c *fiber.Ctx) (uuid.UUID, error) {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return uuid.Nil, apperror.BadRequestError(err, "Incorrect UUID format")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return uuid.Nil, err
		}
		return uuid.Nil, apperror.InternalServerError(err, "Can not parse UUID")
	}
	return parsedID, nil
}
