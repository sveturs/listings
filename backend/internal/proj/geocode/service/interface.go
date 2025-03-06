package service

import (
    "backend/internal/domain/models"
    "context"
)

type GeocodeServiceInterface interface {
    ReverseGeocode(ctx context.Context, lat, lon float64) (*models.GeoLocation, error)
    GetCitySuggestions(ctx context.Context, query string, limit int) ([]models.GeoLocation, error)
}