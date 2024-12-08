package handlers

import (
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
    "backend/internal/services"
    "strconv"

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
    var car struct {
        Make        string   `json:"make"`
        Model       string   `json:"model"`
        Year        int      `json:"year"`
        PricePerDay float64  `json:"price_per_day"`
        Location    string   `json:"location"`
        Availability bool    `json:"availability"`
        Transmission string  `json:"transmission"`
        FuelType    string   `json:"fuel_type"`
        Seats       int      `json:"seats"`
        Features    []string `json:"features"`
    }

    if err := c.BodyParser(&car); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
    }

    // Валидация данных
    if car.Make == "" || car.Model == "" || car.Year == 0 || car.PricePerDay == 0 || car.Location == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "All required fields must be filled")
    }

    // Создаем модель Car для базы данных
    carModel := &models.Car{
        Make:         car.Make,
        Model:        car.Model,
        Year:         car.Year,
        PricePerDay:  car.PricePerDay,
        Location:     car.Location,
        Availability: car.Availability,
        Transmission: car.Transmission,
        FuelType:     car.FuelType,
        Seats:        car.Seats,
        Features:     car.Features,
    }

    carID, err := h.services.Car().AddCar(c.Context(), carModel)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding car: "+err.Error())
    }

    return utils.SuccessResponse(c, fiber.Map{
        "id": carID,
        "message": "Car added successfully",
    })
}
func (h *CarHandler) UploadImages(c *fiber.Ctx) error {
    carID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid car ID")
    }

    form, err := c.MultipartForm()
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Error getting files")
    }

    files := form.File["images"]
    isMain := len(files) > 0

    var uploadedImages []models.CarImage
    for _, file := range files {
        fileName, err := h.services.Car().ProcessImage(file)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error processing image")
        }

        image := models.CarImage{
            CarID:       carID,
            FilePath:    fileName,
            FileName:    file.Filename,
            FileSize:    int(file.Size),
            ContentType: file.Header.Get("Content-Type"),
            IsMain:      isMain,
        }

        imageID, err := h.services.Car().AddCarImage(c.Context(), &image)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error saving image information")
        }

        image.ID = imageID
        uploadedImages = append(uploadedImages, image)
        isMain = false
    }

    return utils.SuccessResponse(c, uploadedImages)
}

func (h *CarHandler) GetImages(c *fiber.Ctx) error {
    images, err := h.services.Car().GetCarImages(c.Context(), c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error getting images")
    }
    return utils.SuccessResponse(c, images)
}
func (h *CarHandler) GetAvailableCars(c *fiber.Ctx) error {
    cars, err := h.services.Car().GetAvailableCars(c.Context())
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching cars")
    }

    return utils.SuccessResponse(c, cars)
}
