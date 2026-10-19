package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int        `json:"limit" validate:"gte=1,lte=20"`
	Offset int        `json:"offset" validate:"gte=0"`
	Sort   string     `json:"sort"`
	Tags   []string   `json:"tags" validate:"max=5"`
	Search string     `json:"search" validate:"max=100"`
	Since  *time.Time `json:"since"`
	Until  *time.Time `json:"until"`
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

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		t, err := parseTime(since)
		if err != nil {
			return fq, err
		}
		fq.Since = t
	}

	until := qs.Get("until")
	if until != "" {
		t, err := parseTime(until)
		if err != nil {
			return fq, err
		}
		fq.Until = t
	}

	return fq, nil
}

func parseTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
