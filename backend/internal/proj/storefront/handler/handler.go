// Package handler
// backend/internal/proj/storefront/handler/handler.go
package handler

import (
	globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Storefront *StorefrontHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		Storefront: NewStorefrontHandler(services),
	}
}
