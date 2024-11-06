package api

import (
	"github.com/ONSdigital/dis-search-upstream-stub/config"

	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router    *mux.Router
	Cfg       *config.Config
	DataStore DataStorer
}

// Setup function sets up the api and returns an api
func Setup(r *mux.Router, cfg *config.Config, dataStorer DataStorer) *API {
	api := &API{
		Router:    r,
		Cfg:       cfg,
		DataStore: dataStorer,
	}

	r.HandleFunc("/resource", GetResources(api)).Methods("GET")
	return api
}
