package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}
	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = o
	}
	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}
	tags := qs.Get("tags")
	if tags != "" {
		tags_list := strings.Split(tags, ",")
		fq.Tags = tags_list
	} else {
		fq.Tags = []string{}

	}
	return fq, nil
}
