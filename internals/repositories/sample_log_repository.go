package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SampleLogRepository struct {
	db *gorm.DB
}

func NewSampleLogRepository(db *gorm.DB) *SampleLogRepository {
	return &SampleLogRepository{
		db: db,
	}
}

func (r *SampleLogRepository) Save(message string) error {
	if err := r.db.Create(&domain.SampleLog{
		Message: message,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *SampleLogRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.SampleLog{}, "id = ?", id).Error
}

func (r *SampleLogRepository) GetAll() ([]domain.SampleLog, error) {
	var logs []domain.SampleLog
	err := r.db.Find(&logs).Error
	return logs, err
}

func (r *SampleLogRepository) EditMessage(id uuid.UUID, message string) error {
	return r.db.Model(&domain.SampleLog{}).Where("id = ?", id).Update("message", message).Error
}
