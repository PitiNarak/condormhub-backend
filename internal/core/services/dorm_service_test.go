package services

import (
	"testing"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDormRepo struct {
	saveFunc func(dorm *domain.Dorm) error
}

func (m *mockDormRepo) Create(dorm *domain.Dorm) error {
	return m.saveFunc(dorm)
}

func (m *mockDormRepo) GetAll(limit int, page int, search string, min_price int, max_price int, district string, subdistrict string, province string, zipcode string) ([]domain.Dorm, int, int, error) {
	panic("unimplemented")
}

func (m *mockDormRepo) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	panic("unimplemented")
}

func (m *mockDormRepo) Update(id uuid.UUID, dorm dto.DormUpdateRequestBody) error {
	panic("unimplemented")
}

func (m *mockDormRepo) Delete(dorm domain.Dorm) error {
	panic("unimplemented")
}

func (m *mockDormRepo) SaveDormImage(dormImage *domain.DormImage) error {
	panic("unimplemented")
}

func (m *mockDormRepo) GetByOwnerID(ownerID uuid.UUID, limit int, page int) ([]domain.Dorm, int, int, error) {
	panic("unimplemented")
}

func (m *mockDormRepo) DeleteImageByKey(imageKey string) error {
	panic("unimplemented")
}

func (m *mockDormRepo) GetImageByKey(imageKey string) (*domain.DormImage, error) {
	panic("unimplemented")
}

func TestCreateDorm(t *testing.T) {
	// Success case: User is lessor
	t.Run("lessor", func(t *testing.T) {
		repo := &mockDormRepo{
			saveFunc: func(dorm *domain.Dorm) error {
				return nil
			},
		}
		service := NewDormService(repo, nil)

		err := service.Create(domain.LessorRole, &domain.Dorm{Name: "SpaceDorm"})
		assert.NoError(t, err)
	})

	// Success case: User is admin
	t.Run("admin", func(t *testing.T) {
		repo := &mockDormRepo{
			saveFunc: func(dorm *domain.Dorm) error {
				return nil
			},
		}
		service := NewDormService(repo, nil)

		err := service.Create(domain.AdminRole, &domain.Dorm{Name: "SpaceDorm"})
		assert.NoError(t, err)
	})

	// Failure case: User is a lessee
	t.Run("lessee", func(t *testing.T) {
		repo := &mockDormRepo{
			saveFunc: func(dorm *domain.Dorm) error {
				return nil
			},
		}
		service := NewDormService(repo, nil)

		err := service.Create(domain.LesseeRole, &domain.Dorm{Name: "SpaceDorm"})
		assert.Error(t, err)
		assert.Equal(t, "You do not have permission to create a dorm - unauthorized action", err.Error())
	})
}
