// backend/internal/domain/models/geo_location.go
package models

// GeoLocation представляет информацию о местоположении
type GeoLocation struct {
	City    string  `json:"city"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}