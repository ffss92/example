package pagination

import (
	"math"
	"net/http"
	"strconv"
)

const (
	basePage  = 1
	baseLimit = 50
	maxLimit  = 50
)

type Pagination struct {
	page  int
	limit int
}

func (p Pagination) Offset() int {
	return (p.page - 1) * p.limit
}

func (p Pagination) Page() int {
	return p.page
}

func (p Pagination) Limit() int {
	return p.limit
}

func FromRequest(r *http.Request) Pagination {
	var pagination Pagination
	p := r.URL.Query().Get("page")
	page, err := strconv.Atoi(p)
	switch {
	case err != nil, page < 1:
		pagination.page = basePage
	default:
		pagination.page = page
	}

	l := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(l)
	switch {
	case err != nil, limit < 1:
		pagination.limit = baseLimit
	case limit > maxLimit:
		pagination.limit = maxLimit
	default:
		pagination.limit = limit
	}

	return pagination
}

type Paginated[T any] struct {
	Limit       int  `json:"limit"`
	CurrentPage int  `json:"current_page"`
	TotalCount  int  `json:"total_count"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
	Data        []T  `json:"data"`
}

func NewPaginated[T any](page, limit, totalCount int, data []T) Paginated[T] {
	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	return Paginated[T]{
		Limit:       limit,
		CurrentPage: page,
		TotalCount:  totalCount,
		TotalPages:  int(totalPages),
		HasPrevious: page > 1,
		HasNext:     page < int(totalPages),
		Data:        data,
	}
}
