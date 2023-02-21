package api

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-nlp-search-scrubber/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{
		AreaDataFile:     "data/2011 OAC Clusters and Names csv v2.csv",
		IndustryDataFile: "data/SIC07_CH_condensed_list_en.csv",
	}

	// Create a mock router
	r := mux.NewRouter()

	// Setup the API
	api := Setup(context.Background(), r, cfg)

	// Assert that the Router was set correctly
	assert.Equal(t, r, api.Router)

	// Assert that the "/scrubber/search" route was added
	route := r.Get("PrefixSearchHandler")
	assert.NotNil(t, route, "Expected PrefixSearchHandler to be added")
}
