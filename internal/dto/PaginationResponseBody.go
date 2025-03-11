package dto

type PaginationResponseBody struct {
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
}
