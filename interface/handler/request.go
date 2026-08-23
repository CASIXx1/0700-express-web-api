package handler

import (
	"net/url"
	"strconv"
)

type paginationRequest struct {
	Limit int
	Page  int
}

func parsePaginationParams(values url.Values) (*paginationRequest, error) {
	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil {
		return nil, err
	}

	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		return nil, err
	}

	return &paginationRequest{
		Limit: limit,
		Page:  page,
	}, nil
}
