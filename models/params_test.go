package models

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetScrubberParamsHappyCases(t *testing.T) {
	tests := []struct {
		name     string
		query    url.Values
		expected *ScrubberParams
	}{
		{
			name: "simple query",
			query: url.Values{
				"q": []string{"dentists"},
			},
			expected: &ScrubberParams{
				Query: "dentists",
				SIC:   []string{},
				OAC:   []string{},
			},
		},
		{
			name: "query with special characters",
			query: url.Values{
				"q": []string{"1 dental-care!"},
			},
			expected: &ScrubberParams{
				Query: "dental care",
				SIC:   []string{},
				OAC:   []string{},
			},
		},
		{
			name: "query with SIC code",
			query: url.Values{
				"q": []string{"12345 dentists"},
			},
			expected: &ScrubberParams{
				Query: "dentists",
				SIC:   []string{"12345"},
				OAC:   []string{},
			},
		},
		{
			name: "query with OAC code",
			query: url.Values{
				"q": []string{"X12345678 dentists"},
			},
			expected: &ScrubberParams{
				Query: "dentists",
				SIC:   []string{},
				OAC:   []string{"X12345678"},
			},
		},
		{
			name: "query with repeated codes",
			query: url.Values{
				"q": []string{"12345 X12345678 dentists 12345 X12345678"},
			},
			expected: &ScrubberParams{
				Query: "dentists",
				SIC:   []string{"12345"},
				OAC:   []string{"X12345678"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := GetScrubberParams(tt.query)
			assert.Empty(t, err)
			assert.Equal(t, tt.expected, params)
		})
	}
}

func TestGetScrubberParamsReturnsError(t *testing.T) {
	tests := []struct {
		name     string
		query    url.Values
		expected error
	}{
		{
			name: "wrong query name",
			query: url.Values{
				"query": []string{},
			},
			expected: fmt.Errorf("no query provided or wrong query name"),
		},
		{
			name: "multiple queries q",
			query: url.Values{
				"q": []string{"dentists", "IN", "london"},
			},
			expected: fmt.Errorf("one query expected, found multiple queries with the same name "),
		},
		{
			name: "query with SIC code",
			query: url.Values{
				"q":     []string{"12345 dentists"},
				"quer":  []string{"12345 dentists"},
				"query": []string{"12345 dentists"},
			},
			expected: fmt.Errorf("one query expected, found multiple queries "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := GetScrubberParams(tt.query)
			assert.Empty(t, params)
			assert.Equal(t, tt.expected, err)
		})
	}
}
