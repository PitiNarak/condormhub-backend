package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type Receipt struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt      time.Time   `gorm:"autoCreateTime"`
	OwnerID       uuid.UUID   `gorm:"type:uuid;not null"`
	Owner         User        `gorm:"foreignKey:OwnerID;references:ID"`
	TransactionID uuid.UUID   `gorm:"type:uuid;not null"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID;references:ID"`
	FileKey       string
}

func (r *Receipt) ToDTO(url string) dto.ReceiptResponseBody {
	return dto.ReceiptResponseBody{
		ID:          r.ID,
		Owner:       r.Owner.ToDTO(),
		Transaction: r.Transaction.ToDTO(),
		Url:         url,
	}
}
