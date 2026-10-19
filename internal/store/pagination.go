package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		limit, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = limit
	}

	offset := qs.Get(("offset"))

	if offset != "" {
		offset, err := strconv.Atoi(offset)

		if err != nil {
			return fq, err
		}

		fq.Offset = offset
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}
