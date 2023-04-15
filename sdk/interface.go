package sdk

import (
	"context"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-nlp-search-scrubber/models"
	"github.com/ONSdigital/dp-nlp-search-scrubber/sdk/errors"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter

type Clienter interface {
	Checker(ctx context.Context, check *health.CheckState) error
	GetSearch(ctx context.Context, options Options) (*models.ScrubberResp, errors.Error)
	Health() *healthcheck.Client
	URL() string
}
