package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type paginationRequest struct {
	Limit int
	Page  int
}

func bearerToken(request *http.Request) (string, error) {
	authorization := request.Header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("missing authorization")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization")
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == "" {
		return "", errors.New("invalid authorization")
	}

	return token, nil
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
