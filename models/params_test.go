package models

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetScrubberParams(t *testing.T) {
	// Test case 1: query has one "q" parameter
	query1 := url.Values{}
	query1.Set("q", "test")
	expected1 := &ScrubberParams{
		Query: "test",
	}
	assert.Equal(t, expected1, GetScrubberParams(query1))

	// Test case 2: query has no "q" parameter
	query2 := url.Values{}
	expected2 := &ScrubberParams{
		Query: "",
	}
	assert.Equal(t, expected2, GetScrubberParams(query2))

	// Test case 3: query has multiple "q" parameters
	query3 := url.Values{}
	query3.Set("q", "test1")
	query3.Add("q", "test2")
	expected3 := &ScrubberParams{
		Query: "test1",
	}
	assert.Equal(t, expected3, GetScrubberParams(query3))
}
