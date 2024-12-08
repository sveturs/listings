package handlers

import (
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
    "backend/internal/services"

)
type CarHandler struct {
    services services.ServicesInterface
}

func NewCarHandler(services services.ServicesInterface) *CarHandler {
    return &CarHandler{
        services: services,
    }
}
func (h *CarHandler) AddCar(c *fiber.Ctx) error {
    var car models.Car
    if err := c.BodyParser(&car); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
    }

    carID, err := h.services.Car().AddCar(c.Context(), &car)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding car")
    }

    return utils.SuccessResponse(c, fiber.Map{"id": carID})
}

func (h *CarHandler) GetAvailableCars(c *fiber.Ctx) error {
    cars, err := h.services.Car().GetAvailableCars(c.Context())
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching cars")
    }

    return utils.SuccessResponse(c, cars)
}
