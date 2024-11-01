package api

import (
	"context"
	_ "embed"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

//go:embed example_resource.json
var jsonResponse []byte

// GetResources returns all resources that are wanted to be indexed in search
func GetResources(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err := w.Write(jsonResponse)
		if err != nil {
			log.Error(ctx, "writing response failed", err)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
