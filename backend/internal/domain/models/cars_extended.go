package models

import "time"

// CarMakeWithStats марка автомобиля со статистикой
type CarMakeWithStats struct {
	CarMake
	ModelCount   int `json:"model_count"`
	ListingCount int `json:"listing_count"`
}

// VehicleInfo полная информация об автомобиле
type VehicleInfo struct {
	ID           int       `json:"id"`
	ListingID    int       `json:"listing_id"`
	MakeID       int       `json:"make_id"`
	MakeName     string    `json:"make_name"`
	ModelID      int       `json:"model_id"`
	ModelName    string    `json:"model_name"`
	GenerationID *int      `json:"generation_id,omitempty"`
	Generation   string    `json:"generation,omitempty"`
	Year         int       `json:"year"`
	Price        float64   `json:"price"`
	Mileage      *int      `json:"mileage,omitempty"`
	FuelType     *string   `json:"fuel_type,omitempty"`
	Transmission *string   `json:"transmission,omitempty"`
	BodyType     *string   `json:"body_type,omitempty"`
	EngineSize   *float64  `json:"engine_size,omitempty"`
	Color        *string   `json:"color,omitempty"`
	Location     string    `json:"location"`
	ImageURL     *string   `json:"image_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// VINDecodeResult результат декодирования VIN
type VINDecodeResult struct {
	VIN          string                 `json:"vin"`
	Valid        bool                   `json:"valid"`
	MakeID       *int                   `json:"make_id,omitempty"`
	MakeName     string                 `json:"make_name"`
	ModelID      *int                   `json:"model_id,omitempty"`
	ModelName    string                 `json:"model_name"`
	Year         int                    `json:"year"`
	Trim         string                 `json:"trim,omitempty"`
	BodyType     string                 `json:"body_type,omitempty"`
	Engine       VINEngineInfo          `json:"engine"`
	Transmission string                 `json:"transmission,omitempty"`
	DriveType    string                 `json:"drive_type,omitempty"`
	FuelType     string                 `json:"fuel_type,omitempty"`
	Manufacturer VINManufacturerInfo    `json:"manufacturer"`
	Source       string                 `json:"source"` // nhtsa, vehicle_databases, etc
	RawData      map[string]interface{} `json:"raw_data,omitempty"`
	DecodedAt    time.Time              `json:"decoded_at"`
}

// VINEngineInfo информация о двигателе из VIN
type VINEngineInfo struct {
	Type          string   `json:"type,omitempty"`
	Displacement  *float64 `json:"displacement,omitempty"` // в литрах
	Cylinders     *int     `json:"cylinders,omitempty"`
	Configuration string   `json:"configuration,omitempty"` // inline, V, etc
	PowerHP       *int     `json:"power_hp,omitempty"`
	PowerKW       *int     `json:"power_kw,omitempty"`
}

// VINManufacturerInfo информация о производителе из VIN
type VINManufacturerInfo struct {
	Name      string `json:"name"`
	Country   string `json:"country,omitempty"`
	PlantCity string `json:"plant_city,omitempty"`
	PlantCode string `json:"plant_code,omitempty"`
}
