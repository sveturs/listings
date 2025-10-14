package service

import (
	"context"

	"backend/internal/proj/postexpress"
	"backend/internal/proj/postexpress/models"
	"backend/internal/proj/postexpress/storage"
)

// Service представляет интерфейс сервиса Post Express
type Service interface {
	// Управление настройками
	GetSettings(ctx context.Context) (*models.PostExpressSettings, error)
	UpdateSettings(ctx context.Context, settings *models.PostExpressSettings) error

	// Работа с локациями
	SearchLocations(ctx context.Context, query string) ([]*models.PostExpressLocation, error)
	GetLocationByID(ctx context.Context, id int) (*models.PostExpressLocation, error)
	SyncLocations(ctx context.Context) error

	// Работа с офисами
	GetOfficesByLocation(ctx context.Context, locationID int) ([]*models.PostExpressOffice, error)
	GetOfficeByCode(ctx context.Context, code string) (*models.PostExpressOffice, error)
	SyncOffices(ctx context.Context) error

	// Расчет стоимости
	CalculateRate(ctx context.Context, req *models.CalculateRateRequest) (*models.CalculateRateResponse, error)
	GetRates(ctx context.Context) ([]*models.PostExpressRate, error)

	// Управление отправлениями
	CreateShipment(ctx context.Context, req *models.CreateShipmentRequest) (*models.PostExpressShipment, error)
	GetShipment(ctx context.Context, id int) (*models.PostExpressShipment, error)
	GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*models.PostExpressShipment, error)
	ListShipments(ctx context.Context, filters storage.ShipmentFilters) ([]*models.PostExpressShipment, int, error)
	UpdateShipmentStatus(ctx context.Context, id int, status models.ShipmentStatus) error
	CancelShipment(ctx context.Context, id int, reason string) error

	// Документы
	GetShipmentLabel(ctx context.Context, id int) ([]byte, error)
	GetShipmentInvoice(ctx context.Context, id int) ([]byte, error)
	PrintShipmentDocuments(ctx context.Context, id int) error

	// Отслеживание
	TrackShipment(ctx context.Context, trackingNumber string) ([]*models.TrackingEvent, error)
	UpdateTrackingStatus(ctx context.Context, trackingNumber string) error
	SyncAllActiveShipments(ctx context.Context) error

	// Склад и самовывоз
	GetWarehouses(ctx context.Context) ([]*models.Warehouse, error)
	GetWarehouseByCode(ctx context.Context, code string) (*models.Warehouse, error)
	CreatePickupOrder(ctx context.Context, req *models.CreatePickupOrderRequest) (*models.WarehousePickupOrder, error)
	GetPickupOrder(ctx context.Context, id int) (*models.WarehousePickupOrder, error)
	GetPickupOrderByCode(ctx context.Context, code string) (*models.WarehousePickupOrder, error)
	ConfirmPickup(ctx context.Context, id int, confirmedBy string, documentType string, documentNumber string) error
	CancelPickupOrder(ctx context.Context, id int, reason string) error

	// Статистика
	GetShipmentStatistics(ctx context.Context, filters storage.StatisticsFilters) (*storage.ShipmentStatistics, error)
	GetWarehouseStatistics(ctx context.Context, warehouseID int) (*storage.WarehouseStatistics, error)
}

// WSPClient представляет интерфейс клиента WSP API
type WSPClient interface {
	// Базовый метод для всех транзакций
	Transaction(ctx context.Context, req *models.TransactionRequest) (*models.TransactionResponse, error)

	// Специализированные методы
	GetLocations(ctx context.Context, search string) ([]WSPLocation, error)
	GetOffices(ctx context.Context, locationID int) ([]WSPOffice, error)
	CreateShipment(ctx context.Context, shipment *WSPShipmentRequest) (*WSPShipmentResponse, error)
	CreateShipmentViaManifest(ctx context.Context, shipment *WSPShipmentRequest) (*postexpress.ManifestResponse, error)
	CreateManifest(ctx context.Context, manifest *postexpress.ManifestRequest) (*postexpress.ManifestResponse, error)
	GetShipmentStatus(ctx context.Context, trackingNumber string) (*WSPTrackingResponse, error)
	PrintLabel(ctx context.Context, shipmentID string) ([]byte, error)
	CancelShipment(ctx context.Context, shipmentID string) error
}

// WSP API структуры данных

// WSPLocation представляет населенный пункт в WSP API
type WSPLocation struct {
	ID           int    `json:"Id"`
	Name         string `json:"Naziv"`
	PostalCode   string `json:"PostanskiBroj"`
	Municipality string `json:"Opstina"`
}

// WSPOffice представляет почтовое отделение в WSP API
type WSPOffice struct {
	Code         string `json:"Sifra"`
	Name         string `json:"Naziv"`
	Address      string `json:"Adresa"`
	Phone        string `json:"Telefon"`
	WorkingHours string `json:"RadnoVreme"`
}

// WSPShipmentRequest представляет запрос на создание отправления в WSP API
type WSPShipmentRequest struct {
	SenderName          string  `json:"NazivPosiljaoca"`
	SenderAddress       string  `json:"AdresaPosiljaoca"`
	SenderCity          string  `json:"MestoPosiljaoca"`
	SenderPostalCode    string  `json:"PostanskiBrojPosiljaoca"`
	SenderPhone         string  `json:"TelefonPosiljaoca"`
	RecipientName       string  `json:"NazivPrimaoca"`
	RecipientAddress    string  `json:"AdresaPrimaoca"`
	RecipientCity       string  `json:"MestoPrimaoca"`
	RecipientPostalCode string  `json:"PostanskiBrojPrimaoca"`
	RecipientPhone      string  `json:"TelefonPrimaoca"`
	Weight              float64 `json:"Tezina"`
	CODAmount           float64 `json:"Otkupnina"`
	InsuranceAmount     float64 `json:"VrednostPosiljke"`
	ServiceType         string  `json:"VrstaUsluge"`
	Content             string  `json:"Sadrzaj"`
	Note                string  `json:"Napomena"`
}

// WSPShipmentResponse представляет ответ на создание отправления
type WSPShipmentResponse struct {
	Success         bool   `json:"Uspesno"`
	TrackingNumber  string `json:"BrojPosiljke"`
	Barcode         string `json:"Barkod"`
	LabelURL        string `json:"URLNalepnice"`
	ErrorMessage    string `json:"Greska"`
	ReferenceNumber string `json:"ReferencniBroj"` // Наш внутренний номер посылки
}

// WSPTrackingResponse представляет ответ отслеживания
type WSPTrackingResponse struct {
	TrackingNumber string             `json:"BrojPosiljke"`
	Status         string             `json:"Status"`
	Events         []WSPTrackingEvent `json:"Dogadjaji"`
}

// WSPTrackingEvent представляет событие отслеживания
type WSPTrackingEvent struct {
	Date        string `json:"Datum"`
	Time        string `json:"Vreme"`
	Location    string `json:"Mesto"`
	Description string `json:"Opis"`
	Code        string `json:"Kod"`
}
