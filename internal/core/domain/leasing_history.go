package domain

import "github.com/google/uuid"

type LeasingHistory struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DormID   uuid.UUID `json:"dormId" `
	LesseeID uuid.UUID `json:"lesseeId" `
	Orders   []Order   `json:"orders" gorm:"foreignKey:LeasingHistoryID"`
}
