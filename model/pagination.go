package model

type Pagination struct {
	Total       int   `json:"total"`
	CurrentPage int   `json:"current_page"`
	TotalInPage int64 `json:"total_in_page"`
	PerPage     int   `json:"per_page"`
	LastPage    int   `json:"last_page"`
}
