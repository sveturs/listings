package delivery

import (
	deliveryv1 "github.com/sveturs/delivery/gen/go/delivery/v1"
)

// mapProviderToProto converts domain provider to proto
func mapProviderToProto(p DeliveryProvider) deliveryv1.DeliveryProvider {
	switch p {
	case ProviderPostExpress:
		return deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS
	case ProviderBexExpress:
		return deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_BEX_EXPRESS
	case ProviderAksExpress:
		return deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_AKS_EXPRESS
	case ProviderDExpress:
		return deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_D_EXPRESS
	case ProviderCityExpress:
		return deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_CITY_EXPRESS
	default:
		return deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_UNSPECIFIED
	}
}

// mapProviderFromProto converts proto provider to domain
func mapProviderFromProto(p deliveryv1.DeliveryProvider) DeliveryProvider {
	switch p {
	case deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS:
		return ProviderPostExpress
	case deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_BEX_EXPRESS:
		return ProviderBexExpress
	case deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_AKS_EXPRESS:
		return ProviderAksExpress
	case deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_D_EXPRESS:
		return ProviderDExpress
	case deliveryv1.DeliveryProvider_DELIVERY_PROVIDER_CITY_EXPRESS:
		return ProviderCityExpress
	default:
		return ProviderUnspecified
	}
}

// mapStatusFromProto converts proto status to domain
func mapStatusFromProto(s deliveryv1.ShipmentStatus) ShipmentStatus {
	switch s {
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return StatusPending
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED:
		return StatusConfirmed
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return StatusInTransit
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY:
		return StatusOutForDelivery
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return StatusDelivered
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_FAILED:
		return StatusFailed
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:
		return StatusCancelled
	case deliveryv1.ShipmentStatus_SHIPMENT_STATUS_RETURNED:
		return StatusReturned
	default:
		return StatusUnspecified
	}
}

// mapAddressToProto converts domain address to proto
func mapAddressToProto(a *Address) *deliveryv1.Address {
	if a == nil {
		return nil
	}
	return &deliveryv1.Address{
		Street:       a.Street,
		City:         a.City,
		State:        a.State,
		PostalCode:   a.PostalCode,
		Country:      a.Country,
		ContactName:  a.ContactName,
		ContactPhone: a.ContactPhone,
	}
}

// mapAddressFromProto converts proto address to domain
func mapAddressFromProto(a *deliveryv1.Address) *Address {
	if a == nil {
		return nil
	}
	return &Address{
		Street:       a.Street,
		City:         a.City,
		State:        a.State,
		PostalCode:   a.PostalCode,
		Country:      a.Country,
		ContactName:  a.ContactName,
		ContactPhone: a.ContactPhone,
	}
}

// mapPackageToProto converts domain package to proto
func mapPackageToProto(p *Package) *deliveryv1.Package {
	if p == nil {
		return nil
	}
	return &deliveryv1.Package{
		Weight:        p.Weight,
		Length:        p.Length,
		Width:         p.Width,
		Height:        p.Height,
		Description:   p.Description,
		DeclaredValue: p.DeclaredValue,
	}
}

// mapPackageFromProto converts proto package to domain
func mapPackageFromProto(p *deliveryv1.Package) *Package {
	if p == nil {
		return nil
	}
	return &Package{
		Weight:        p.Weight,
		Length:        p.Length,
		Width:         p.Width,
		Height:        p.Height,
		Description:   p.Description,
		DeclaredValue: p.DeclaredValue,
	}
}

// mapShipmentFromProto converts proto shipment to domain
func mapShipmentFromProto(s *deliveryv1.Shipment) *Shipment {
	if s == nil {
		return nil
	}

	shipment := &Shipment{
		ID:             s.Id,
		TrackingNumber: s.TrackingNumber,
		Provider:       mapProviderFromProto(s.Provider),
		Status:         mapStatusFromProto(s.Status),
		FromAddress:    mapAddressFromProto(s.FromAddress),
		ToAddress:      mapAddressFromProto(s.ToAddress),
		Package:        mapPackageFromProto(s.Package),
		Cost:           s.Cost,
		Currency:       s.Currency,
	}

	if s.CreatedAt != nil {
		shipment.CreatedAt = s.CreatedAt.AsTime()
	}
	if s.UpdatedAt != nil {
		shipment.UpdatedAt = s.UpdatedAt.AsTime()
	}
	if s.EstimatedDelivery != nil {
		shipment.EstimatedDelivery = s.EstimatedDelivery.AsTime()
	}
	if s.ActualDelivery != nil {
		shipment.ActualDelivery = s.ActualDelivery.AsTime()
	}

	return shipment
}

// mapTrackingEventsFromProto converts proto tracking events to domain
func mapTrackingEventsFromProto(events []*deliveryv1.TrackingEvent) []TrackingEvent {
	if events == nil {
		return nil
	}

	result := make([]TrackingEvent, 0, len(events))
	for _, e := range events {
		event := TrackingEvent{
			Status:      mapStatusFromProto(e.Status),
			Location:    e.Location,
			Description: e.Description,
		}
		if e.Timestamp != nil {
			event.Timestamp = e.Timestamp.AsTime()
		}
		result = append(result, event)
	}
	return result
}

// mapSettlementsFromProto converts proto settlements to domain
func mapSettlementsFromProto(settlements []*deliveryv1.Settlement) []Settlement {
	if settlements == nil {
		return nil
	}

	result := make([]Settlement, 0, len(settlements))
	for _, s := range settlements {
		result = append(result, Settlement{
			ID:      s.Id,
			Name:    s.Name,
			ZipCode: s.ZipCode,
			Country: s.Country,
		})
	}
	return result
}

// mapStreetsFromProto converts proto streets to domain
func mapStreetsFromProto(streets []*deliveryv1.Street) []Street {
	if streets == nil {
		return nil
	}

	result := make([]Street, 0, len(streets))
	for _, s := range streets {
		result = append(result, Street{
			ID:             s.Id,
			Name:           s.Name,
			SettlementName: s.SettlementName,
		})
	}
	return result
}

// mapParcelLockersFromProto converts proto parcel lockers to domain
func mapParcelLockersFromProto(lockers []*deliveryv1.ParcelLocker) []ParcelLocker {
	if lockers == nil {
		return nil
	}

	result := make([]ParcelLocker, 0, len(lockers))
	for _, l := range lockers {
		result = append(result, ParcelLocker{
			ID:        l.Id,
			Code:      l.Code,
			Name:      l.Name,
			Address:   l.Address,
			City:      l.City,
			ZipCode:   l.ZipCode,
			Latitude:  l.Latitude,
			Longitude: l.Longitude,
			Available: l.Available,
		})
	}
	return result
}
