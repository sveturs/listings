// TEMPORARY: Will be moved to microservice
// Этот модуль предоставляет минимальную функциональность marketplace
// для работы frontend до полной миграции на микросервис

package handler

import (
	"backend/internal/clients/listings"
	"backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Handler struct {
	storage              storage.MarketplaceStorage
	services             service.ServicesInterface
	jwtParserMW          fiber.Handler
	logger               zerolog.Logger
	listingsClient       *listings.Client
	useListingsMicroservice bool
}

func NewHandler(
	db *sqlx.DB,
	services service.ServicesInterface,
	jwtParserMW fiber.Handler,
	logger zerolog.Logger,
	listingsClient *listings.Client,
	useListingsMicroservice bool,
) *Handler {
	return &Handler{
		storage:              storage.NewPostgresMarketplaceStorage(db, logger),
		services:             services,
		jwtParserMW:          jwtParserMW,
		logger:               logger.With().Str("module", "marketplace_handler").Logger(),
		listingsClient:       listingsClient,
		useListingsMicroservice: useListingsMicroservice,
	}
}
