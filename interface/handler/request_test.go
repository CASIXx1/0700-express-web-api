package handler

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			expectedError:  strconv.ErrSyntax,
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
			expectedError:  strconv.ErrSyntax,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input()

			result, err := parsePaginationParams(input)
			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}
