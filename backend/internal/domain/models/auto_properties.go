// backend/internal/domain/models/auto_properties.go
package models

import (
	"time"
)

// AutoProperties содержит дополнительные свойства для объявлений с автомобилями
type AutoProperties struct {
	ListingID         int       `json:"listing_id"`
	Brand             string    `json:"brand"`
	Model             string    `json:"model"`
	Year              int       `json:"year"`
	Mileage           int       `json:"mileage"`
	FuelType          string    `json:"fuel_type"`
	Transmission      string    `json:"transmission"`
	EngineCapacity    float64   `json:"engine_capacity"`
	Power             int       `json:"power"`
	Color             string    `json:"color"`
	BodyType          string    `json:"body_type"`
	DriveType         string    `json:"drive_type"`
	NumberOfDoors     int       `json:"number_of_doors"`
	NumberOfSeats     int       `json:"number_of_seats"`
	AdditionalFeatures string    `json:"additional_features"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AutoConstants предоставляет константы для свойств автомобилей
type AutoConstants struct {
	FuelTypes     []string `json:"fuel_types"`
	Transmissions []string `json:"transmissions"`
	BodyTypes     []string `json:"body_types"`
	DriveTypes    []string `json:"drive_types"`
	Brands        []string `json:"brands"`
}

// AutoFilter содержит параметры фильтрации для автомобилей
type AutoFilter struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	YearFrom     int    `json:"year_from"`
	YearTo       int    `json:"year_to"`
	MileageFrom  int    `json:"mileage_from"`
	MileageTo    int    `json:"mileage_to"`
	FuelType     string `json:"fuel_type"`
	Transmission string `json:"transmission"`
	BodyType     string `json:"body_type"`
	DriveType    string `json:"drive_type"`
}

// GetAutoConstants возвращает константы для автомобильных свойств
func GetAutoConstants() AutoConstants {
	return AutoConstants{
		FuelTypes: []string{
			"Бензин",
			"Дизель",
			"Электро",
			"Гибрид",
			"Газ",
			"Газ/Бензин",
		},
		Transmissions: []string{
			"Механическая",
			"Автоматическая",
			"Робот",
			"Вариатор",
		},
		BodyTypes: []string{
			"Седан",
			"Хэтчбек",
			"Универсал",
			"Внедорожник",
			"Купе",
			"Кабриолет",
			"Пикап",
			"Минивэн",
			"Фургон",
		},
		DriveTypes: []string{
			"Передний",
			"Задний",
			"Полный",
		},
		Brands: []string{
			"Audi",
			"BMW",
			"Chevrolet",
			"Citroen",
			"Dacia",
			"Fiat",
			"Ford",
			"Honda",
			"Hyundai",
			"Kia",
			"Lada",
			"Mazda",
			"Mercedes-Benz",
			"Mitsubishi",
			"Nissan",
			"Opel",
			"Peugeot",
			"Renault",
			"Skoda",
			"Suzuki",
			"Toyota",
			"Volkswagen",
			"Volvo",
			"Zastava",
		},
	}
}

// Расширим модель MarketplaceListing для включения автомобильных свойств
type AutoListing struct {
	MarketplaceListing
	AutoProperties *AutoProperties `json:"auto_properties,omitempty"`
}