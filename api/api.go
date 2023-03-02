package api

import (
	"context"

	"github.com/ONSdigital/dp-nlp-search-scrubber/config"
	"github.com/ONSdigital/dp-nlp-search-scrubber/db"
	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config) *API {
	api := &API{
		Router: r,
	}

	db := db.LoadCsvData(ctx, cfg)
	r.HandleFunc("/scrubber/search", PrefixSearchHandler(db)).Methods("GET").Name("PrefixSearchHandler")
	return api
}
