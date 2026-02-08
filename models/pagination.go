package models

type PaginationParams struct {
	Search string
	Page   int
	Limit  int
}

type PaginationMeta struct {
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
