package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
)

type LeasingHistoryHandler struct {
	service ports.LeasingHistoryService
}

func NewLeasingHistoryHandler(service ports.LeasingHistoryService) ports.LeasingHistoryHandler {
	return &LeasingHistoryHandler{service: service}
}
