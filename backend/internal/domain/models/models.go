// backend/internal/domain/models/models.go
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
type CarCategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CarImage struct {
	ID          int       `json:"id"`
	CarID       int       `json:"car_id"`
	FilePath    string    `json:"file_path"`
	FileName    string    `json:"file_name"`
	FileSize    int       `json:"file_size"`
	ContentType string    `json:"content_type"`
	IsMain      bool      `json:"is_main"`
	CreatedAt   time.Time `json:"created_at"`
}

type Car struct {
	ID                int        `json:"id"`
	Make              string     `json:"make"`
	Model             string     `json:"model"`
	Year              int        `json:"year"`
	PricePerDay       float64    `json:"price_per_day"`
	Location          string     `json:"location"`
	Latitude          float64    `json:"latitude"`
	Longitude         float64    `json:"longitude"`
	Description       string     `json:"description,omitempty"`
	Availability      bool       `json:"availability"`
	Transmission      string     `json:"transmission"`
	FuelType          string     `json:"fuel_type"`
	Seats             int        `json:"seats"`
	CategoryID        int        `json:"category_id"`
	Category          string     `json:"category"`
	Features          []string   `json:"features"`
	DailyMileageLimit *int       `json:"daily_mileage_limit"` // Изменить на указатель
	InsuranceIncluded bool       `json:"insurance_included"`
	Images            []CarImage `json:"images,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type CarBooking struct {
	ID              int       `json:"id"`
	CarID           int       `json:"car_id"`
	UserID          int       `json:"user_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	PickupLocation  string    `json:"pickup_location"`
	DropoffLocation string    `json:"dropoff_location"`
	Status          string    `json:"status"`
	TotalPrice      float64   `json:"total_price"`
	CreatedAt       time.Time `json:"created_at"`
	Car             *Car      `json:"car,omitempty"`
	UserName        string    `json:"user_name,omitempty"`
	UserEmail       string    `json:"user_email,omitempty"`
}
type CarFeature struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
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
    TotalBeds          *int      `json:"total_beds,omitempty"`      // Сделали указателем
    AvailableBeds      *int      `json:"available_beds,omitempty"`  // Сделали указателем
    HasPrivateBathroom bool      `json:"has_private_bathroom"`
    Latitude           *float64  `json:"latitude,omitempty"`        // Сделали указателем
    Longitude          *float64  `json:"longitude,omitempty"`       // Сделали указателем
    FormattedAddress   string    `json:"formatted_address"`
    CreatedAt          time.Time `json:"created_at"`
    ActualPrice        float64   `json:"actual_price"`
    AvailableCount     int       `json:"available_count"`
    Rating             float64   `json:"rating"`
    Images             []RoomImage `json:"images,omitempty"`
}

type RoomImage struct {
	ID          int       `json:"id"`
	RoomID      int       `json:"room_id"`
	BedID       int       `json:"bed_id,omitempty"` // Добавляем это поле
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
type MarketplaceListing struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    CategoryID  int       `json:"category_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Condition   string    `json:"condition"`
    Status      string    `json:"status"`
    Location    string    `json:"location"`
    Latitude    *float64  `json:"latitude,omitempty"`
    Longitude   *float64  `json:"longitude,omitempty"`
    City        string    `json:"city"`
    Country     string    `json:"country"`
    ViewsCount  int       `json:"views_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Images      []MarketplaceImage `json:"images,omitempty"`
    User        *User     `json:"user,omitempty"`
    Category    *MarketplaceCategory `json:"category,omitempty"`
}

type MarketplaceCategory struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Slug      string    `json:"slug"`
    ParentID  *int      `json:"parent_id,omitempty"`
    Icon      string    `json:"icon,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}
type MarketplaceImage struct {
    ID          int       `json:"id"`
    ListingID   int       `json:"listing_id"`
    FilePath    string    `json:"file_path"`
    FileName    string    `json:"file_name"`
    FileSize    int       `json:"file_size"`
    ContentType string    `json:"content_type"`
    IsMain      bool      `json:"is_main"`
    CreatedAt   time.Time `json:"created_at"`
}