package models

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGetScrubberParams(t *testing.T) {
// 	// Test case 1: query has one "q" parameter
// 	query1 := url.Values{}
// 	query1.Set("q", "test")
// 	expected1 := &ScrubberParams{
// 		Query: "test",
// 	}
// 	assert.Equal(t, expected1, GetScrubberParams(query1))

// 	// Test case 2: query has no "q" parameter
// 	query2 := url.Values{}
// 	expected2 := &ScrubberParams{
// 		Query: "",
// 	}
// 	assert.Equal(t, expected2, GetScrubberParams(query2))

// 	// Test case 3: query has multiple "q" parameters
// 	query3 := url.Values{}
// 	query3.Set("q", "test1")
// 	query3.Add("q", "test2")
// 	expected3 := &ScrubberParams{
// 		Query: "test1",
// 	}
// 	assert.Equal(t, expected3, GetScrubberParams(query3))
// }

func TestGetScrubberParams(t *testing.T) {
	tests := []struct {
		name     string
		query    url.Values
		expected *ScrubberParams
	}{
		{
			name:  "empty query",
			query: url.Values{},
			expected: &ScrubberParams{
				Query: "",
				SIC:   []string{},
				OAC:   []string{},
			},
		},
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
				"q": []string{"#1 dental-care!"},
			},
			expected: &ScrubberParams{
				Query: " 1 dental care ",
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
				Query: "12345 dentists",
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
				Query: "X12345678 dentists",
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
				Query: "12345 X12345678 dentists 12345 X12345678",
				SIC:   []string{"12345"},
				OAC:   []string{"X12345678"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := GetScrubberParams(tt.query)
			assert.Equal(t, tt.expected, params)
		})
	}
}
