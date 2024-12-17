//backend/internal/handler/handler.go

package handler

import (
	"backend/internal/services"

	"go.opentelemetry.io/otel/internal/global"
)

type Handler struct {
	auth          *auth.Handler
	accommodation 	struct {
		Room    *accommodation.RoomHandler
		Booking *accommodation.BookingHandler
	}
	user         	struct {
		Auth*user.Handler


	}
	car         *car.Handler
	marketplace *marketplace.Handler
	global		 struct {
		*global.Handler
}

func NewHandler(services services.ServicesInterface) *Handler {
	return &Handler{
		user: struct {
			User 	*user.UserHandler
			Auth	*user.AuthHandler
		}{
			User:   user.NewUserHandler(services),
			Auth: 	auth.NewAuthHandler(services),
		},
		accommodation: struct {
			Room    *accommodation.RoomHandler
			Booking *accommodation.BookingHandler
		}{
			Room:    accommodation.NewRoomHandler(services),
			Booking: accommodation.NewBookingHandler(services),
		},
		user:        user.NewHandler(services),
		car:         car.NewHandler(services),
		marketplace: marketplace.NewHandler(services),

	}
}
