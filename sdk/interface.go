package sdk

import (
	"context"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	apiError "github.com/ONSdigital/dis-search-upstream-stub/sdk/errors"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter
type Clienter interface {
	Checker(ctx context.Context, check *health.CheckState) error
	GetResources(ctx context.Context, options Options) (*models.Resources, apiError.Error)
}
