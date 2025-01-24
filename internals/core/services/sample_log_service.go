package services

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
)

type SampleLogService struct {
	sampleLogRepository ports.SampleLogRepository
}

func NewSampleLogService(sampleLogRepository ports.SampleLogRepository) *SampleLogService {
	return &SampleLogService{
		sampleLogRepository: sampleLogRepository,
	}
}

func (s *SampleLogService) Save(message string) error {
	return s.sampleLogRepository.Save(message)
}

func (s *SampleLogService) Delete(id int) error {
	return s.sampleLogRepository.Delete(id)
}

func (s *SampleLogService) GetAll() ([]domain.SampleLog, error) {
	return s.sampleLogRepository.GetAll()
}
