package dto

import "github.com/PitiNarak/condormhub-backend/internal/core/domain"

type PaginationResponseBody struct {
	Currentpage    int                     `json:"currentPage"`
	Lastpage       int                     `json:"lastPage"`
	Limit          int                     `json:"limit"`
	Total          int                     `json:"total"`
	LeasingHistory []domain.LeasingHistory `json:"leasingHistory"`
}
