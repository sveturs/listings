//backend/internal/handlers/car.go
package handlers

import (
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
    "backend/internal/services"
    "strconv"
    "log"

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
    var carData struct {
        Make        string   `json:"make"`
        Model       string   `json:"model"`
        Year        int      `json:"year"`
        PricePerDay float64  `json:"price_per_day"`
        Location    string   `json:"location"`
        Latitude    float64  `json:"latitude"`
        Longitude   float64  `json:"longitude"`
        Description string   `json:"description"`
        Availability bool    `json:"availability"`
        Transmission string  `json:"transmission"`
        FuelType    string   `json:"fuel_type"`
        Seats       int      `json:"seats"`
        Features    []string `json:"features"`
    }

    if err := c.BodyParser(&carData); err != nil {
        log.Printf("Error parsing request body: %v", err)
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
    }

    log.Printf("Received car data: %+v", carData)

    // Валидация обязательных полей
    if carData.Make == "" || carData.Model == "" || carData.PricePerDay == 0 {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing required fields")
    }

    car := &models.Car{
        Make:         carData.Make,
        Model:        carData.Model,
        Year:         carData.Year,
        PricePerDay:  carData.PricePerDay,
        Location:     carData.Location,
        Latitude:     carData.Latitude,
        Longitude:    carData.Longitude,
        Description:  carData.Description,
        Availability: true,
        Transmission: carData.Transmission,
        FuelType:     carData.FuelType,
        Seats:        carData.Seats,
        Features:     carData.Features,
    }

    carID, err := h.services.Car().AddCar(c.Context(), car)
    if err != nil {
        log.Printf("Error adding car to database: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding car")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "id": carID,
        "message": "Car added successfully",
    })
}
func (h *CarHandler) CreateBooking(c *fiber.Ctx) error {
    var booking models.CarBooking
    if err := c.BodyParser(&booking); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
    }

    // Получаем ID пользователя из сессии
    userID := c.Locals("user_id").(int)
    booking.UserID = userID

    err := h.services.Car().CreateBooking(c.Context(), &booking)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating booking")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "message": "Booking created successfully",
        "booking": booking,
    })
}
func (h *CarHandler) UploadImages(c *fiber.Ctx) error {
    log.Printf("Starting image upload for car")
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
    log.Printf("Getting available cars")  // добавим для отладки
    cars, err := h.services.Car().GetAvailableCars(c.Context())
    if err != nil {
        log.Printf("Error getting available cars: %v", err)  // добавим для отладки
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching cars")
    }

    log.Printf("Found %d available cars", len(cars))  // добавим для отладки
    return utils.SuccessResponse(c, cars)
}