package api

import (
	_ "embed"
	"github.com/ONSdigital/dis-search-upstream-stub/apierrors"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/pagination"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

//go:embed example_resource.json
var jsonResponse []byte

var (
	serverErrorMessage = apierrors.ErrInternalServer.Error()
)

// GetResources returns all resources that are wanted to be indexed in search
func GetResources(api *API) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		offsetParam := req.URL.Query().Get("offset")
		limitParam := req.URL.Query().Get("limit")
		logData := log.Data{}

		// initialise pagination
		offset, limit, err := pagination.InitialisePagination(api.Cfg, offsetParam, limitParam)
		if err != nil {
			logData["offset_parameter"] = offsetParam
			logData["limit_parameter"] = limitParam

			log.Error(ctx, "pagination validation failed", err, logData)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		options := data.Options{
			Offset: offset,
			Limit:  limit,
		}

		// get resources from datastore
		resources, err := api.DataStore.GetResources(ctx, options)
		if err != nil {
			logData["options"] = options
			log.Error(ctx, "getting list of resources failed", err, logData)
			http.Error(w, serverErrorMessage, http.StatusInternalServerError)
			return
		}

		logData["resources_count"] = resources.Count
		logData["resources_limit"] = resources.Limit
		logData["resources_offset"] = resources.Offset
		logData["resources_total_count"] = resources.TotalCount

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Error(ctx, "writing response failed", err)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
