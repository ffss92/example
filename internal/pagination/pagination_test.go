package pagination

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		pagination     Pagination
		expectedOffset int
	}{
		{
			pagination: Pagination{
				page:  1,
				limit: 25,
			},
			expectedOffset: 0,
		},
		{
			pagination: Pagination{
				page:  2,
				limit: 25,
			},
			expectedOffset: 25,
		},
		{
			pagination: Pagination{
				page:  5,
				limit: 10,
			},
			expectedOffset: 40,
		},
	}

	for _, tc := range testCases {
		page := tc.pagination.Page()
		limit := tc.pagination.Limit()
		offset := tc.pagination.Offset()

		assert.Equal(tc.pagination.page, page)
		assert.Equal(tc.pagination.limit, limit)
		assert.Equal(tc.expectedOffset, offset)
	}
}

func TestFromRequest(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		name       string
		page       string
		limit      string
		pageValue  int
		limitValue int
	}{
		{
			name:       "should use base values for blank page / limit",
			page:       "",
			limit:      "",
			pageValue:  basePage,
			limitValue: baseLimit,
		},
		{
			name:       "should use base values for invalid page / limit",
			page:       "a",
			limit:      "b",
			pageValue:  basePage,
			limitValue: baseLimit,
		},
		{
			name:       "should use correct values for valid page / limit",
			page:       "10",
			limit:      "25",
			pageValue:  10,
			limitValue: 25,
		},
		{
			name:       "should use maxLimit for limits > maxLimit",
			limit:      "100",
			pageValue:  basePage,
			limitValue: maxLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse("/")
			q := u.Query()
			q.Set("page", tc.page)
			q.Set("limit", tc.limit)
			u.RawQuery = q.Encode()
			url := u.String()

			req := httptest.NewRequest(http.MethodGet, url, nil)
			p := FromRequest(req)

			assert.Equal(tc.pageValue, p.page)
			assert.Equal(tc.limitValue, p.limit)
		})
	}
}

func TestPaginatedCalulation(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		paginated   Paginated[string]
		totalPages  int
		hasPrevious bool
		hasNext     bool
	}{
		{
			paginated:   NewPaginated(1, 10, 100, []string{}),
			totalPages:  10,
			hasPrevious: false,
			hasNext:     true,
		},
		{
			paginated:   NewPaginated(3, 10, 100, []string{}),
			totalPages:  10,
			hasPrevious: true,
			hasNext:     true,
		},
		{
			paginated:   NewPaginated(10, 10, 100, []string{}),
			totalPages:  10,
			hasPrevious: true,
			hasNext:     false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(tc.hasPrevious, tc.paginated.HasPrevious)
		assert.Equal(tc.hasNext, tc.paginated.HasNext)
		assert.Equal(tc.totalPages, tc.paginated.TotalPages)
	}
}
