package handler

import (
	"context"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

type contextKey int

const (
	userIDKey contextKey = iota
)

type paginationRequest struct {
	Limit int
	Page  int
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
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
