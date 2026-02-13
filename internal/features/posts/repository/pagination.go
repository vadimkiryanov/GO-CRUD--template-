package repository

// PaginatedResult представляет результаты с пагинацией
type PaginatedResult struct {
	Data       []PostModel `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// PaginationParams содержит параметры для пагинации
type PaginationParams struct {
	Limit  int
	Offset int
}
