package storage

import (
	"context"
	"time"

	"backend/internal/proj/postexpress/models"
)

// Repository представляет интерфейс репозитория для Post Express
type Repository interface {
	// Управление настройками
	GetSettings(ctx context.Context) (*models.PostExpressSettings, error)
	UpdateSettings(ctx context.Context, settings *models.PostExpressSettings) error
	CreateSettings(ctx context.Context, settings *models.PostExpressSettings) error

	// Локации
	GetLocationByID(ctx context.Context, id int) (*models.PostExpressLocation, error)
	GetLocationByPostExpressID(ctx context.Context, postExpressID int) (*models.PostExpressLocation, error)
	SearchLocations(ctx context.Context, query string, limit int) ([]*models.PostExpressLocation, error)
	CreateLocation(ctx context.Context, location *models.PostExpressLocation) error
	UpdateLocation(ctx context.Context, location *models.PostExpressLocation) error
	BulkUpsertLocations(ctx context.Context, locations []*models.PostExpressLocation) error

	// Отделения
	GetOfficeByID(ctx context.Context, id int) (*models.PostExpressOffice, error)
	GetOfficeByCode(ctx context.Context, code string) (*models.PostExpressOffice, error)
	GetOfficesByLocationID(ctx context.Context, locationID int) ([]*models.PostExpressOffice, error)
	CreateOffice(ctx context.Context, office *models.PostExpressOffice) error
	UpdateOffice(ctx context.Context, office *models.PostExpressOffice) error
	BulkUpsertOffices(ctx context.Context, offices []*models.PostExpressOffice) error

	// Тарифы
	GetRates(ctx context.Context) ([]*models.PostExpressRate, error)
	GetRateForWeight(ctx context.Context, weight float64) (*models.PostExpressRate, error)
	CreateRate(ctx context.Context, rate *models.PostExpressRate) error
	UpdateRate(ctx context.Context, rate *models.PostExpressRate) error

	// Отправления
	CreateShipment(ctx context.Context, shipment *models.PostExpressShipment) (*models.PostExpressShipment, error)
	GetShipmentByID(ctx context.Context, id int) (*models.PostExpressShipment, error)
	GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*models.PostExpressShipment, error)
	GetShipmentsByOrderID(ctx context.Context, orderID int, isStorefront bool) ([]*models.PostExpressShipment, error)
	UpdateShipment(ctx context.Context, shipment *models.PostExpressShipment) error
	UpdateShipmentStatus(ctx context.Context, id int, status models.ShipmentStatus) error
	ListShipments(ctx context.Context, filters ShipmentFilters) ([]*models.PostExpressShipment, int, error)

	// Отслеживание
	CreateTrackingEvent(ctx context.Context, event *models.TrackingEvent) error
	GetTrackingEventsByShipmentID(ctx context.Context, shipmentID int) ([]*models.TrackingEvent, error)
	BulkCreateTrackingEvents(ctx context.Context, events []*models.TrackingEvent) error

	// API логи
	CreateAPILog(ctx context.Context, log *APILog) error
	GetAPILogsByShipmentID(ctx context.Context, shipmentID int) ([]*APILog, error)
	CleanupOldAPILogs(ctx context.Context, daysToKeep int) error

	// Склад
	GetWarehouses(ctx context.Context) ([]*models.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id int) (*models.Warehouse, error)
	GetWarehouseByCode(ctx context.Context, code string) (*models.Warehouse, error)

	// Заказы на самовывоз
	CreatePickupOrder(ctx context.Context, order *models.WarehousePickupOrder) (*models.WarehousePickupOrder, error)
	GetPickupOrderByID(ctx context.Context, id int) (*models.WarehousePickupOrder, error)
	GetPickupOrderByCode(ctx context.Context, code string) (*models.WarehousePickupOrder, error)
	UpdatePickupOrder(ctx context.Context, order *models.WarehousePickupOrder) error
	ListPickupOrders(ctx context.Context, filters PickupOrderFilters) ([]*models.WarehousePickupOrder, int, error)

	// Статистика
	GetShipmentStatistics(ctx context.Context, filters StatisticsFilters) (*ShipmentStatistics, error)
	GetWarehouseStatistics(ctx context.Context, warehouseID int) (*WarehouseStatistics, error)
}

// APILog представляет лог взаимодействия с API
type APILog struct {
	ID              int                     `json:"id" db:"id"`
	TransactionID   string                  `json:"transaction_id" db:"transaction_id"`
	TransactionType int                     `json:"transaction_type" db:"transaction_type"`
	RequestData     *map[string]interface{} `json:"request_data,omitempty" db:"request_data"`
	ResponseData    *map[string]interface{} `json:"response_data,omitempty" db:"response_data"`
	Status          string                  `json:"status" db:"status"`
	ErrorMessage    *string                 `json:"error_message,omitempty" db:"error_message"`
	ShipmentID      *int                    `json:"shipment_id,omitempty" db:"shipment_id"`
	ExecutionTimeMs *int                    `json:"execution_time_ms,omitempty" db:"execution_time_ms"`
	CreatedAt       time.Time               `json:"created_at" db:"created_at"`
}

// PickupOrderFilters представляет фильтры для заказов на самовывоз
type PickupOrderFilters struct {
	WarehouseID *int                      `json:"warehouse_id,omitempty"`
	Status      *models.PickupOrderStatus `json:"status,omitempty"`
	CustomerID  *int                      `json:"customer_id,omitempty"`
	DateFrom    *string                   `json:"date_from,omitempty"`
	DateTo      *string                   `json:"date_to,omitempty"`
	Page        int                       `json:"page"`
	PageSize    int                       `json:"page_size"`
}

// ShipmentFilters представляет фильтры для отправлений
type ShipmentFilters struct {
	Status     *models.ShipmentStatus `json:"status,omitempty"`
	DateFrom   *string                `json:"date_from,omitempty"`
	DateTo     *string                `json:"date_to,omitempty"`
	Search     *string                `json:"search,omitempty"`
	OrderID    *int                   `json:"order_id,omitempty"`
	CustomerID *int                   `json:"customer_id,omitempty"`
	City       *string                `json:"city,omitempty"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
}

// StatisticsFilters представляет фильтры для статистики
type StatisticsFilters struct {
	Status   *models.ShipmentStatus `json:"status,omitempty"`
	DateFrom *string                `json:"date_from,omitempty"`
	DateTo   *string                `json:"date_to,omitempty"`
	GroupBy  string                 `json:"group_by,omitempty"`
}

// ShipmentStatistics представляет статистику отправлений
type ShipmentStatistics struct {
	Total      int     `json:"total"`
	Pending    int     `json:"pending"`
	InTransit  int     `json:"in_transit"`
	Delivered  int     `json:"delivered"`
	Cancelled  int     `json:"cancelled"`
	TotalValue float64 `json:"total_value"`
	TotalCOD   float64 `json:"total_cod"`
	// Дополнительные поля для совместимости с repository
	TotalShipments      int            `json:"total_shipments"`
	DeliveredShipments  int            `json:"delivered_shipments"`
	InTransitShipments  int            `json:"in_transit_shipments"`
	FailedShipments     int            `json:"failed_shipments"`
	TotalRevenue        float64        `json:"total_revenue"`
	AverageDeliveryTime float64        `json:"average_delivery_time"`
	DeliverySuccessRate float64        `json:"delivery_success_rate"`
	ByStatus            map[string]int `json:"by_status,omitempty"`
	ByCity              map[string]int `json:"by_city,omitempty"`
}

// WarehouseStatistics представляет статистику склада
type WarehouseStatistics struct {
	WarehouseID     int     `json:"warehouse_id"`
	WarehouseName   string  `json:"warehouse_name"`
	TotalOrders     int     `json:"total_orders"`
	PendingOrders   int     `json:"pending_orders"`
	ReadyOrders     int     `json:"ready_orders"`
	CompletedOrders int     `json:"completed_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
	TotalValue      float64 `json:"total_value"`
	// Дополнительные поля для совместимости с repository
	TotalPickupOrders     int     `json:"total_pickup_orders"`
	PendingPickupOrders   int     `json:"pending_pickup_orders"`
	CompletedPickupOrders int     `json:"completed_pickup_orders"`
	AveragePickupTime     float64 `json:"average_pickup_time"`
	TotalVolumeM3         float64 `json:"total_volume_m3"`
	OccupancyPercent      float64 `json:"occupancy_percent"`
	TotalInventoryItems   int     `json:"total_inventory_items"`
}
