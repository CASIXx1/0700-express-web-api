package handler

import (
	"context"
	"net/url"
	"strconv"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
)

type paginationRequest struct {
	Limit int
	Page  int
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
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
