package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	GoogleID   string    `json:"google_id"`
	PictureURL string    `json:"picture_url"`
	CreatedAt  time.Time `json:"created_at"`
}

type Room struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Capacity           int       `json:"capacity"`
	PricePerNight      float64   `json:"price_per_night"`
	AddressStreet      string    `json:"address_street"`
	AddressCity        string    `json:"address_city"`
	AddressState       string    `json:"address_state"`
	AddressCountry     string    `json:"address_country"`
	AddressPostalCode  string    `json:"address_postal_code"`
	AccommodationType  string    `json:"accommodation_type"`
	IsShared           bool      `json:"is_shared"`
	TotalBeds          int       `json:"total_beds"`
	AvailableBeds      int       `json:"available_beds"`
	HasPrivateBathroom bool      `json:"has_private_bathroom"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	FormattedAddress   string    `json:"formatted_address"`
	CreatedAt          time.Time `json:"created_at"`
}

type RoomImage struct {
    ID          int       `json:"id"`
    RoomID      int       `json:"room_id"`
    BedID       int       `json:"bed_id,omitempty"`  // Добавляем это поле
    FilePath    string    `json:"file_path"`
    FileName    string    `json:"file_name"`
    FileSize    int       `json:"file_size"`
    ContentType string    `json:"content_type"`
    IsMain      bool      `json:"is_main"`
    CreatedAt   time.Time `json:"created_at"`
}

type Booking struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	RoomID            int       `json:"room_id"`
	BedID             *int      `json:"bed_id,omitempty"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Status            string    `json:"status"`
	RoomName          string    `json:"room_name"`
	AccommodationType string    `json:"type"`
	UserName          string    `json:"user_name"`
	UserEmail         string    `json:"user_email"`
	BookingType       string    `json:"booking_type"`
}

type Bed struct {
    ID            int     `json:"id"`
    RoomID        int     `json:"room_id"`
    BedNumber     string  `json:"bed_number"`
    IsAvailable   bool    `json:"is_available"`
    PricePerNight float64 `json:"price_per_night"`
    HasOutlet     bool    `json:"has_outlet"`
    HasLight      bool    `json:"has_light"`
    HasShelf      bool    `json:"has_shelf"` 
    BedType       string  `json:"bed_type"`
}