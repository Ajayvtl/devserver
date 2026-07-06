package common

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type QueryParams struct {
	Page    int
	PerPage int
	Sort    string
	Dir     string
	Query   string
}

func ParseQueryParams(r *http.Request, allowedSortFields []string) (QueryParams, error) {
	q := r.URL.Query()

	page := 1
	if pStr := q.Get("page"); pStr != "" {
		p, err := strconv.Atoi(pStr)
		if err != nil || p < 1 {
			return QueryParams{}, fmt.Errorf("invalid page: must be a positive integer")
		}
		page = p
	}

	perPage := 50
	if ppStr := q.Get("per_page"); ppStr != "" {
		pp, err := strconv.Atoi(ppStr)
		if err != nil || pp < 1 {
			return QueryParams{}, fmt.Errorf("invalid per_page: must be a positive integer")
		}
		if pp > 100 {
			return QueryParams{}, fmt.Errorf("invalid per_page: exceeds maximum limit of 100")
		}
		perPage = pp
	}

	sort := q.Get("sort")
	if sort == "" {
		sort = "created_at"
	} else if len(allowedSortFields) > 0 {
		validSort := false
		for _, f := range allowedSortFields {
			if sort == f {
				validSort = true
				break
			}
		}
		if !validSort {
			return QueryParams{}, fmt.Errorf("invalid sort field: %s", sort)
		}
	}

	dir := strings.ToUpper(q.Get("dir"))
	if dir != "" && dir != "ASC" && dir != "DESC" {
		return QueryParams{}, fmt.Errorf("invalid sort direction: must be ASC or DESC")
	}
	if dir == "" {
		dir = "DESC"
	}

	query := q.Get("q")

	return QueryParams{
		Page:    page,
		PerPage: perPage,
		Sort:    sort,
		Dir:     dir,
		Query:   query,
	}, nil
}

func (qp QueryParams) Offset() int {
	return (qp.Page - 1) * qp.PerPage
}

func (qp QueryParams) Limit() int {
	return qp.PerPage
}
