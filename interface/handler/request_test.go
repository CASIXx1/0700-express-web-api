package handler

import (
	"net/url"
	"testing"
)

func TestParsePaginationParams(t *testing.T) {
	t.Run("it can parse limit and page", func(t *testing.T) {
		values := url.Values{}
		values.Set("limit", "20")
		values.Set("page", "1")

		result, err := parsePaginationParams(values)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.Limit != 20 {
			t.Fatalf("expected limit 20, got %d", result.Limit)
		}

		if result.Page != 1 {
			t.Fatalf("expected page 1, got %d", result.Page)
		}
	})

	t.Run("it returns an error if limit is not a number", func(t *testing.T) {
		values := url.Values{}
		values.Set("limit", "abc")
		values.Set("page", "1")

		result, err := parsePaginationParams(values)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if result != nil {
			t.Fatalf("expected result nil, got %+v", result)
		}
	})

	t.Run("it returns an error if page is not a number", func(t *testing.T) {
		values := url.Values{}
		values.Set("limit", "20")
		values.Set("page", "abc")

		result, err := parsePaginationParams(values)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if result != nil {
			t.Fatalf("expected result nil, got %+v", result)
		}
	})
}
