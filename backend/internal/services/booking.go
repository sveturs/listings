package services

import (
    "context"
    "fmt"
    "backend/internal/domain/models"
    "backend/internal/storage"
)

type BookingService struct {
    storage storage.Storage
}

func NewBookingService(storage storage.Storage) *BookingService {
    return &BookingService{
        storage: storage,
    }
}

func (s *BookingService) CreateBooking(ctx context.Context, userID int, booking *models.BookingRequest) error {
	// Проверка пользователя
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	booking.UserID = user.ID
	return s.storage.CreateBooking(ctx, booking)
}

func (s *BookingService) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	return s.storage.GetAllBookings(ctx)
}

func (s *BookingService) DeleteBooking(ctx context.Context, bookingID string, bookingType string) error {
	return s.storage.DeleteBooking(ctx, bookingID, bookingType)
}

func (s *BookingService) GetAvailableBeds(ctx context.Context, roomID, startDate, endDate string) ([]models.Bed, error) {
	return s.storage.GetAvailableBeds(ctx, roomID, startDate, endDate)
}

func (s *BookingService) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    return s.storage.AddBed(ctx, roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType)
}
