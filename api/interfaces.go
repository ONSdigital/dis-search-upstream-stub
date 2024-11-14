package api

import (
	"context"

	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
)

//go:generate moq -out ./mock/data_storer.go -pkg mock . DataStorer

// DataStorer is an interface for a type that can store and retrieve resources
type DataStorer interface {
	GetResources(ctx context.Context, options data.Options) (resource *models.Resources, err error)
}

// Paginator defines the required methods from the paginator package
type Paginator interface {
	ValidateParameters(offsetParam string, limitParam string, totalCount int) (offset int, limit int, err error)
}
