package delivery

import (
	"time"
)

// DeliveryProvider represents delivery provider type
type DeliveryProvider int32

const (
	ProviderUnspecified DeliveryProvider = iota
	ProviderPostExpress
	ProviderBexExpress
	ProviderAksExpress
	ProviderDExpress
	ProviderCityExpress
)

// ShipmentStatus represents shipment status
type ShipmentStatus int32

const (
	StatusUnspecified ShipmentStatus = iota
	StatusPending
	StatusConfirmed
	StatusInTransit
	StatusOutForDelivery
	StatusDelivered
	StatusFailed
	StatusCancelled
	StatusReturned
)

// Address represents a delivery address
type Address struct {
	Street       string
	City         string
	State        string
	PostalCode   string
	Country      string
	ContactName  string
	ContactPhone string
}

// Package represents package dimensions and details
type Package struct {
	Weight        string // Weight in kg
	Length        string // Length in cm
	Width         string // Width in cm
	Height        string // Height in cm
	Description   string
	DeclaredValue string // Value in currency
}

// Shipment represents a delivery shipment
type Shipment struct {
	ID                string
	TrackingNumber    string
	Provider          DeliveryProvider
	Status            ShipmentStatus
	FromAddress       *Address
	ToAddress         *Address
	Package           *Package
	Cost              string
	Currency          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EstimatedDelivery time.Time
	ActualDelivery    time.Time
}

// TrackingEvent represents a single tracking event
type TrackingEvent struct {
	Status      ShipmentStatus
	Location    string
	Description string
	Timestamp   time.Time
}

// TrackingInfo contains shipment and its tracking events
type TrackingInfo struct {
	Shipment *Shipment
	Events   []TrackingEvent
}

// CreateShipmentRequest contains data for creating a shipment
type CreateShipmentRequest struct {
	Provider    DeliveryProvider
	FromAddress *Address
	ToAddress   *Address
	Package     *Package
	UserID      string
}

// CalculateRateRequest contains data for rate calculation
type CalculateRateRequest struct {
	Provider    DeliveryProvider
	FromAddress *Address
	ToAddress   *Address
	Package     *Package
}

// RateInfo contains delivery rate information
type RateInfo struct {
	Cost              string
	Currency          string
	EstimatedDelivery time.Time
}

// Settlement represents a city/settlement
type Settlement struct {
	ID      int32
	Name    string
	ZipCode string
	Country string
}

// Street represents a street
type Street struct {
	ID             int32
	Name           string
	SettlementName string
}

// ParcelLocker represents a parcel locker location
type ParcelLocker struct {
	ID        int32
	Code      string
	Name      string
	Address   string
	City      string
	ZipCode   string
	Latitude  float64
	Longitude float64
	Available bool
}
