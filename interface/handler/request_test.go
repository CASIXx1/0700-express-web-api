package handler

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestParsePaginationParams(t *testing.T) {
	tests := []struct {
		name           string
		input          func() url.Values
		expectedResult *paginationRequest
		expectedError  error
	}{
		{
			name: "normal case: parse limit and page",
			input: func() url.Values {
				values := url.Values{}
				values.Set("limit", "20")
				values.Set("page", "1")
				return values
			},
			expectedResult: &paginationRequest{
				Limit: 20,
				Page:  1,
			},
			expectedError: nil,
		},
		{
			name: "error case: limit is not a number",
			input: func() url.Values {
				values := url.Values{}
				values.Set("limit", "abc")
				values.Set("page", "1")
				return values
			},
			expectedResult: nil,
			expectedError: &strconv.NumError{
				Func: "Atoi",
				Num:  "abc",
				Err:  strconv.ErrSyntax,
			},
		},
		{
			name: "error case: page is not a number",
			input: func() url.Values {
				values := url.Values{}
				values.Set("limit", "20")
				values.Set("page", "abc")
				return values
			},
			expectedResult: nil,
			expectedError: &strconv.NumError{
				Func: "Atoi",
				Num:  "abc",
				Err:  strconv.ErrSyntax,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input()

			result, err := parsePaginationParams(input)
			if err != nil {
				if test.expectedError == nil {
					t.Fatalf("got unexpected error %v", err)
				}

				if !errors.Is(err, test.expectedError) && !reflect.DeepEqual(test.expectedError, err) {
					t.Fatalf("expected error %v but got %v", test.expectedError, err)
				}

				return
			}

			if test.expectedError != nil {
				t.Fatal("expected error, got nil")
			}

			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Fatalf("expected %+v but got %+v", test.expectedResult, result)
			}
		})
	}
}
