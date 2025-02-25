package dto

type PaginationResponseBody struct {
	Currentpage int `json:"currentPage"`
	Lastpage    int `json:"lastPage"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
}
