package delivery

import (
	"context"
	"fmt"

	"github.com/vondi-global/listings/internal/service"
)

// ServiceAdapter adapts the delivery Client to service.DeliveryClient interface
type ServiceAdapter struct {
	client *Client
}

// NewServiceAdapter creates a new adapter for the delivery client
func NewServiceAdapter(client *Client) *ServiceAdapter {
	return &ServiceAdapter{client: client}
}

// CreateShipment implements service.DeliveryClient
func (a *ServiceAdapter) CreateShipment(ctx context.Context, req *service.DeliveryCreateShipmentRequest) (*service.DeliveryShipment, error) {
	// Convert service types to client types
	clientReq := &CreateShipmentRequest{
		Provider:    mapServiceProviderToClient(req.Provider),
		FromAddress: mapServiceAddressToClient(req.FromAddress),
		ToAddress:   mapServiceAddressToClient(req.ToAddress),
		Package:     mapServicePackageToClient(req.Package),
		UserID:      req.UserID,
	}

	shipment, err := a.client.CreateShipment(ctx, clientReq)
	if err != nil {
		return nil, err
	}

	return mapClientShipmentToService(shipment), nil
}

// TrackShipment implements service.DeliveryClient
func (a *ServiceAdapter) TrackShipment(ctx context.Context, trackingNumber string) (*service.DeliveryTrackingInfo, error) {
	info, err := a.client.TrackShipment(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}

	return &service.DeliveryTrackingInfo{
		Shipment: mapClientShipmentToService(info.Shipment),
		Events:   mapClientEventsToService(info.Events),
	}, nil
}

// CalculateRate implements service.DeliveryClient
func (a *ServiceAdapter) CalculateRate(ctx context.Context, req *service.DeliveryCalculateRateRequest) (*service.DeliveryRateInfo, error) {
	clientReq := &CalculateRateRequest{
		Provider:    mapServiceProviderToClient(req.Provider),
		FromAddress: mapServiceAddressToClient(req.FromAddress),
		ToAddress:   mapServiceAddressToClient(req.ToAddress),
		Package:     mapServicePackageToClient(req.Package),
	}

	rate, err := a.client.CalculateRate(ctx, clientReq)
	if err != nil {
		return nil, err
	}

	return &service.DeliveryRateInfo{
		Cost:              rate.Cost,
		Currency:          rate.Currency,
		EstimatedDelivery: rate.EstimatedDelivery,
	}, nil
}

// mapServiceProviderToClient converts service provider to client provider
func mapServiceProviderToClient(p service.DeliveryProvider) DeliveryProvider {
	switch p {
	case service.DeliveryProviderPostExpress:
		return ProviderPostExpress
	case service.DeliveryProviderBexExpress:
		return ProviderBexExpress
	case service.DeliveryProviderAksExpress:
		return ProviderAksExpress
	case service.DeliveryProviderDExpress:
		return ProviderDExpress
	case service.DeliveryProviderCityExpress:
		return ProviderCityExpress
	default:
		return ProviderUnspecified
	}
}

// mapClientProviderToService converts client provider to service provider
func mapClientProviderToService(p DeliveryProvider) service.DeliveryProvider {
	switch p {
	case ProviderPostExpress:
		return service.DeliveryProviderPostExpress
	case ProviderBexExpress:
		return service.DeliveryProviderBexExpress
	case ProviderAksExpress:
		return service.DeliveryProviderAksExpress
	case ProviderDExpress:
		return service.DeliveryProviderDExpress
	case ProviderCityExpress:
		return service.DeliveryProviderCityExpress
	default:
		return service.DeliveryProviderUnspecified
	}
}

// mapServiceAddressToClient converts service address to client address
func mapServiceAddressToClient(addr *service.DeliveryAddress) *Address {
	if addr == nil {
		return nil
	}
	return &Address{
		Street:       addr.Street,
		City:         addr.City,
		State:        addr.State,
		PostalCode:   addr.PostalCode,
		Country:      addr.Country,
		ContactName:  addr.ContactName,
		ContactPhone: addr.ContactPhone,
	}
}

// mapClientAddressToService converts client address to service address
func mapClientAddressToService(addr *Address) *service.DeliveryAddress {
	if addr == nil {
		return nil
	}
	return &service.DeliveryAddress{
		Street:       addr.Street,
		City:         addr.City,
		State:        addr.State,
		PostalCode:   addr.PostalCode,
		Country:      addr.Country,
		ContactName:  addr.ContactName,
		ContactPhone: addr.ContactPhone,
	}
}

// mapServicePackageToClient converts service package to client package
func mapServicePackageToClient(pkg *service.DeliveryPackage) *Package {
	if pkg == nil {
		return nil
	}
	return &Package{
		Weight:        pkg.Weight,
		Length:        pkg.Length,
		Width:         pkg.Width,
		Height:        pkg.Height,
		Description:   pkg.Description,
		DeclaredValue: pkg.DeclaredValue,
	}
}

// mapClientPackageToService converts client package to service package
func mapClientPackageToService(pkg *Package) *service.DeliveryPackage {
	if pkg == nil {
		return nil
	}
	return &service.DeliveryPackage{
		Weight:        pkg.Weight,
		Length:        pkg.Length,
		Width:         pkg.Width,
		Height:        pkg.Height,
		Description:   pkg.Description,
		DeclaredValue: pkg.DeclaredValue,
	}
}

// mapClientShipmentToService converts client shipment to service shipment
func mapClientShipmentToService(s *Shipment) *service.DeliveryShipment {
	if s == nil {
		return nil
	}
	return &service.DeliveryShipment{
		ID:                s.ID,
		TrackingNumber:    s.TrackingNumber,
		Provider:          mapClientProviderToService(s.Provider),
		Status:            fmt.Sprintf("%d", s.Status),
		FromAddress:       mapClientAddressToService(s.FromAddress),
		ToAddress:         mapClientAddressToService(s.ToAddress),
		Package:           mapClientPackageToService(s.Package),
		Cost:              s.Cost,
		Currency:          s.Currency,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
		EstimatedDelivery: s.EstimatedDelivery,
		ActualDelivery:    s.ActualDelivery,
	}
}

// mapClientEventsToService converts client events to service events
func mapClientEventsToService(events []TrackingEvent) []service.DeliveryTrackingEvent {
	if events == nil {
		return nil
	}
	result := make([]service.DeliveryTrackingEvent, 0, len(events))
	for _, e := range events {
		result = append(result, service.DeliveryTrackingEvent{
			Status:      fmt.Sprintf("%d", e.Status),
			Location:    e.Location,
			Description: e.Description,
			Timestamp:   e.Timestamp,
		})
	}
	return result
}

// Ensure ServiceAdapter implements service.DeliveryClient
var _ service.DeliveryClient = (*ServiceAdapter)(nil)
