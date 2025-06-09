// backend/internal/proj/geocode/handler/handler.go
package handler

import (
	globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Geocode *GeocodeHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		Geocode: NewGeocodeHandler(services.Geocode()),
	}
}
