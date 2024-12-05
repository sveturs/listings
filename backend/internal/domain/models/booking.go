package models

type BookingRequest struct {
	UserID    int    `json:"user_id"`
	RoomID    int    `json:"room_id"`
	BedID     *int   `json:"bed_id,omitempty"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type BedRequest struct {
    BedNumber     string  `json:"bed_number"`
    PricePerNight float64 `json:"price_per_night"`
    HasOutlet     bool    `json:"has_outlet"`
    HasLight      bool    `json:"has_light"` 
    HasShelf      bool    `json:"has_shelf"`
    BedType       string  `json:"bed_type"` // 'top', 'bottom', 'single'
}