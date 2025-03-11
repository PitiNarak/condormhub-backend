package databases

import (
	"errors"
	"math"

	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"gorm.io/gorm"
)

func Paginate(value interface{}, db *gorm.DB, limit, page int, order string) (func(db *gorm.DB) *gorm.DB, int, int, error) {
	var totalRows int64
	err := db.Model(value).Count(&totalRows).Error
	if err != nil {
		return func(db *gorm.DB) *gorm.DB { return nil }, 0, 0, apperror.InternalServerError(err, "cannot paginate the given value")
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	offset := (page - 1) * limit
	if page > totalPages {
		return func(db *gorm.DB) *gorm.DB { return nil }, totalPages, int(totalRows), apperror.BadRequestError(errors.New("page exceeded"), "page exceeded")
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit).Order(order)
	}, totalPages, int(totalRows), nil

}
