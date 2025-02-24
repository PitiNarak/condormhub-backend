package dto

type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	LastPage    int   `json:"lastPage"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
}
