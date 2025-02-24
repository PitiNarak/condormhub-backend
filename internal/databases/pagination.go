package databases

import (
	"errors"
	"math"

	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"gorm.io/gorm"
)

func Paginate(value interface{}, db *gorm.DB, limit, page int, order string) (func(db *gorm.DB) *gorm.DB, int, int, error) {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	offset := (page - 1) * limit
	if page > totalPages {
		offset = (totalPages - 1) * limit
		return func(db *gorm.DB) *gorm.DB { return db.Offset(offset).Limit(limit).Order(order) }, totalPages, int(totalRows), errorHandler.BadRequestError(errors.New("page exceeded"), "page exceeded")
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit).Order(order)
	}, totalPages, int(totalRows), nil

}
