package handler

import (
    "github.com/gofiber/fiber/v2"
    globalService "backend/internal/proj/global/service"
    "backend/internal/types"
    "backend/internal/domain/models"
    "backend/pkg/utils"
)


type BookingHandler struct {
    services globalService.ServicesInterface
}


func NewBookingHandler(services globalService.ServicesInterface) *BookingHandler {
    return &BookingHandler{
        services: services,
    }
}

func (h *BookingHandler) Create(c *fiber.Ctx) error {
	sessionData := c.Locals("user").(*types.SessionData)
	if sessionData == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Необходима авторизация")
	}

	var booking models.BookingRequest
	if err := c.BodyParser(&booking); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	err := h.services.Booking().CreateBooking(c.Context(), sessionData.UserID, &booking)
	if err != nil {
		switch err.Error() {
		case "bed ID is required for bed booking":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Для койко-места необходимо указать ID кровати")
		case "bed is not available":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Койко-место недоступно")
		case "bed is already booked for these dates":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Койко-место уже забронировано на эти даты")
		case "room is already booked for these dates":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Помещение занято на указанные даты")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка создания бронирования")
		}
	}

	return utils.SuccessResponse(c, "Бронирование создано успешно")
}

func (h *BookingHandler) List(c *fiber.Ctx) error {
	bookings, err := h.services.Booking().GetAllBookings(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения списка бронирований")
	}
	return utils.SuccessResponse(c, bookings)
}

func (h *BookingHandler) Delete(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	bookingType := c.Query("type", "room") // по умолчанию "room"

	err := h.services.Booking().DeleteBooking(c.Context(), bookingID, bookingType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка удаления бронирования")
	}
	return utils.SuccessResponse(c, "Бронирование удалено")
}
