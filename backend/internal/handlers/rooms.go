package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/services"
	"backend/pkg/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	services services.ServicesInterface
}

func NewRoomHandler(services services.ServicesInterface) *RoomHandler {
	return &RoomHandler{
		services: services,
	}
}

// Create создает новую комнату
// Create создает новую комнату
func (h *RoomHandler) Create(c *fiber.Ctx) error {
    var room models.Room
    if err := c.BodyParser(&room); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
    }

    roomID, err := h.services.Room().CreateRoom(c.Context(), &room)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка создания комнаты")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "id": roomID,
        "message": "Room created successfully",
    })
}
// UploadBedImages загружает изображения для койко-места
func (h *RoomHandler) UploadBedImages(c *fiber.Ctx) error {
    roomID, err := strconv.Atoi(c.Params("roomId"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID комнаты")
    }
    
    bedID, err := strconv.Atoi(c.Params("bedId"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID койко-места")
    }

    form, err := c.MultipartForm()
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка получения файлов")
    }

    files := form.File["images"]
    isMain := len(files) > 0

    var uploadedImages []models.RoomImage
    for _, file := range files {
        fileName, err := h.services.Room().ProcessImage(file)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обработки изображения")
        }

        image := models.RoomImage{
            RoomID:      roomID,
            BedID:       bedID,
            FilePath:    fileName,
            FileName:    file.Filename,
            FileSize:    int(file.Size),
            ContentType: file.Header.Get("Content-Type"),
            IsMain:      isMain,
        }

        imageID, err := h.services.Room().AddBedImage(c.Context(), &image)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка сохранения информации об изображении")
        }

        image.ID = imageID
        uploadedImages = append(uploadedImages, image)
        isMain = false
    }

    return utils.SuccessResponse(c, uploadedImages)
}
// List получает список комнат
func (h *RoomHandler) List(c *fiber.Ctx) error {
    filters := map[string]string{
        "capacity":           c.Query("capacity"),
        "start_date":         c.Query("start_date"),
        "end_date":           c.Query("end_date"),
        "min_price":          c.Query("min_price"),
        "max_price":          c.Query("max_price"),
        "city":               c.Query("city"),
        "country":            c.Query("country"),
        "accommodation_type": c.Query("accommodation_type"),
        "has_private_rooms":  c.Query("has_private_rooms"),
    }

    log.Printf("Getting rooms with filters: %+v", filters)
    
    rooms, err := h.services.Room().GetRooms(c.Context(), filters)
    if err != nil {
        log.Printf("Error getting rooms: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Ошибка получения списка комнат: %v", err))
    }

    log.Printf("Found %d rooms", len(rooms))
    for i, room := range rooms {
        log.Printf("Room %d: %+v", i, room)
    }

    return utils.SuccessResponse(c, rooms)
}
// В RoomHandler struct (handlers/rooms.go)
func (h *RoomHandler) ListBedImages(c *fiber.Ctx) error {
    bedID := c.Params("id")
    if bedID == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID койко-места не указан")
    }

    images, err := h.services.Room().GetBedImages(c.Context(), bedID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения изображений")
    }

    return utils.SuccessResponse(c, images)
}
// Get получает информацию о конкретной комнате
func (h *RoomHandler) Get(c *fiber.Ctx) error {
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID комнаты")
	}

	room, err := h.services.Room().GetRoomByID(c.Context(), roomID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Комната не найдена")
	}

	return utils.SuccessResponse(c, room)
}

// UploadImages загружает изображения для комнаты
func (h *RoomHandler) UploadImages(c *fiber.Ctx) error {
	log.Println("UploadImages handler called")
    roomID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
		log.Printf("Invalid room ID: %v", err)
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID комнаты")
    }

    form, err := c.MultipartForm()
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка получения файлов")
    }

	files := form.File["images"]
	isMain := len(files) > 0

	var uploadedImages []models.RoomImage
	for _, file := range files {
		fileName, err := h.services.Room().ProcessImage(file)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обработки изображения")
		}

		image := models.RoomImage{
			RoomID:      roomID,
			FilePath:    fileName,
			FileName:    file.Filename,
			FileSize:    int(file.Size),
			ContentType: file.Header.Get("Content-Type"),
			IsMain:      isMain,
		}
		log.Printf("Calling AddRoomImage service method for image: %+v", image)
		imageID, err := h.services.Room().AddRoomImage(c.Context(), &image)
		if err != nil {
			log.Printf("Error saving image information: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка сохранения информации об изображении")
		}
		log.Printf("Image saved successfully with ID: %d", imageID)
		image.ID = imageID
		uploadedImages = append(uploadedImages, image)
		isMain = false
	}
    log.Println("UploadImages handler completed successfully")
    return utils.SuccessResponse(c, uploadedImages)
}

// ListImages получает список изображений комнаты
func (h *RoomHandler) ListImages(c *fiber.Ctx) error {
	images, err := h.services.Room().GetRoomImages(c.Context(), c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения изображений")
	}

	return utils.SuccessResponse(c, images)
}

// DeleteImage удаляет изображение комнаты
func (h *RoomHandler) DeleteImage(c *fiber.Ctx) error {
	filePath, err := h.services.Room().DeleteRoomImage(c.Context(), c.Params("imageId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Изображение не найдено")
	}

	os.Remove(filepath.Join("uploads", filePath))
	return utils.SuccessResponse(c, "Изображение удалено")
}

// AddBed добавляет кровать в комнату
func (h *RoomHandler) AddBed(c *fiber.Ctx) error {
    roomID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID комнаты")
    }

    var bedReq models.BedRequest
    if err := c.BodyParser(&bedReq); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
    }

    bedID, err := h.services.Room().AddBed(c.Context(), roomID, bedReq.BedNumber, bedReq.PricePerNight, bedReq.HasOutlet, bedReq.HasLight, bedReq.HasShelf, bedReq.BedType)
    if err != nil {
        if err.Error() == "room not found" {
            return utils.ErrorResponse(c, fiber.StatusNotFound, "Комната не найдена")
        }
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка добавления кровати")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "id":              bedID,
        "room_id":         roomID,
        "bed_number":      bedReq.BedNumber,
        "price_per_night": bedReq.PricePerNight,
        "has_outlet":      bedReq.HasOutlet,
        "has_light":       bedReq.HasLight,
        "has_shelf":       bedReq.HasShelf,
        "bed_type":        bedReq.BedType,
        "is_available":    true,
    })
}

// GetAvailableBeds получает список доступных кроватей
func (h *RoomHandler) GetAvailableBeds(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Необходимо указать даты")
	}

	beds, err := h.services.Room().GetAvailableBeds(c.Context(), c.Params("id"), startDate, endDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения списка кроватей")
	}

	return utils.SuccessResponse(c, beds)
}
