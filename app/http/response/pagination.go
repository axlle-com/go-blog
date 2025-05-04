package response

import "math"

type Pagination struct {
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
	Pages   int `json:"pages"`
}

func NewPagination(total, page, perPage int) *Pagination {
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	return &Pagination{
		Total:   total,
		PerPage: perPage,
		Page:    page,
		Pages:   pages,
	}
}
