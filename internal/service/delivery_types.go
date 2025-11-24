// Package service provides business logic layer for the listings microservice.
package service

import "time"

// DeliveryProvider represents delivery provider type
type DeliveryProvider int32

const (
	DeliveryProviderUnspecified DeliveryProvider = iota
	DeliveryProviderPostExpress
	DeliveryProviderBexExpress
	DeliveryProviderAksExpress
	DeliveryProviderDExpress
	DeliveryProviderCityExpress
)

// DeliveryAddress represents a delivery address
type DeliveryAddress struct {
	Street       string
	City         string
	State        string
	PostalCode   string
	Country      string
	ContactName  string
	ContactPhone string
}

// DeliveryPackage represents package dimensions
type DeliveryPackage struct {
	Weight        string // Weight in kg
	Length        string // Length in cm
	Width         string // Width in cm
	Height        string // Height in cm
	Description   string
	DeclaredValue string
}

// DeliveryCreateShipmentRequest contains parameters for creating a shipment via delivery microservice
type DeliveryCreateShipmentRequest struct {
	Provider    DeliveryProvider
	FromAddress *DeliveryAddress
	ToAddress   *DeliveryAddress
	Package     *DeliveryPackage
	UserID      string
}

// DeliveryShipment represents a shipment from delivery microservice
type DeliveryShipment struct {
	ID                string
	TrackingNumber    string
	Provider          DeliveryProvider
	Status            string
	FromAddress       *DeliveryAddress
	ToAddress         *DeliveryAddress
	Package           *DeliveryPackage
	Cost              string
	Currency          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EstimatedDelivery time.Time
	ActualDelivery    time.Time
}

// DeliveryTrackingEvent represents a tracking event
type DeliveryTrackingEvent struct {
	Status      string
	Location    string
	Description string
	Timestamp   time.Time
}

// DeliveryTrackingInfo contains tracking information
type DeliveryTrackingInfo struct {
	Shipment *DeliveryShipment
	Events   []DeliveryTrackingEvent
}

// DeliveryCalculateRateRequest contains parameters for rate calculation
type DeliveryCalculateRateRequest struct {
	Provider    DeliveryProvider
	FromAddress *DeliveryAddress
	ToAddress   *DeliveryAddress
	Package     *DeliveryPackage
}

// DeliveryRateInfo contains delivery rate information
type DeliveryRateInfo struct {
	Cost              string
	Currency          string
	EstimatedDelivery time.Time
}
