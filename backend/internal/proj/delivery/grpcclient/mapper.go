package grpcclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/internal/proj/delivery/models"
	pb "backend/pkg/grpc/delivery/v1"
)

// ErrNilShipment returned when nil shipment is passed to mapper
var ErrNilShipment = errors.New("nil shipment provided")

// MapShipmentFromProto конвертирует протобуф Shipment в модель БД
func MapShipmentFromProto(pbShipment *pb.Shipment) (*models.Shipment, error) {
	if pbShipment == nil {
		return nil, ErrNilShipment
	}

	// Парсим sender и recipient info
	senderInfo, err := json.Marshal(map[string]interface{}{
		"street":        pbShipment.FromAddress.Street,
		"city":          pbShipment.FromAddress.City,
		"state":         pbShipment.FromAddress.State,
		"postal_code":   pbShipment.FromAddress.PostalCode,
		"country":       pbShipment.FromAddress.Country,
		"contact_name":  pbShipment.FromAddress.ContactName,
		"contact_phone": pbShipment.FromAddress.ContactPhone,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sender info: %w", err)
	}

	recipientInfo, err := json.Marshal(map[string]interface{}{
		"street":        pbShipment.ToAddress.Street,
		"city":          pbShipment.ToAddress.City,
		"state":         pbShipment.ToAddress.State,
		"postal_code":   pbShipment.ToAddress.PostalCode,
		"country":       pbShipment.ToAddress.Country,
		"contact_name":  pbShipment.ToAddress.ContactName,
		"contact_phone": pbShipment.ToAddress.ContactPhone,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recipient info: %w", err)
	}

	// Парсим package info
	packageInfo, err := json.Marshal(map[string]interface{}{
		"weight":         pbShipment.Package.Weight,
		"length":         pbShipment.Package.Length,
		"width":          pbShipment.Package.Width,
		"height":         pbShipment.Package.Height,
		"description":    pbShipment.Package.Description,
		"declared_value": pbShipment.Package.DeclaredValue,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal package info: %w", err)
	}

	// Парсим стоимость
	cost, err := strconv.ParseFloat(pbShipment.Cost, 64)
	if err != nil {
		cost = 0
	}

	shipment := &models.Shipment{
		ExternalID:     &pbShipment.Id,
		TrackingNumber: &pbShipment.TrackingNumber,
		Status:         MapStatusFromProto(pbShipment.Status),
		SenderInfo:     senderInfo,
		RecipientInfo:  recipientInfo,
		PackageInfo:    packageInfo,
		DeliveryCost:   &cost,
		CreatedAt:      pbShipment.CreatedAt.AsTime(),
		UpdatedAt:      pbShipment.UpdatedAt.AsTime(),
	}

	if pbShipment.EstimatedDelivery != nil {
		t := pbShipment.EstimatedDelivery.AsTime()
		shipment.EstimatedDelivery = &t
	}

	if pbShipment.ActualDelivery != nil {
		t := pbShipment.ActualDelivery.AsTime()
		shipment.ActualDeliveryDate = &t
	}

	return shipment, nil
}

// MapAddressToProto конвертирует JSON адрес в протобуф
func MapAddressToProto(addressJSON json.RawMessage) (*pb.Address, error) {
	var addr map[string]interface{}
	if err := json.Unmarshal(addressJSON, &addr); err != nil {
		return nil, err
	}

	getString := func(key string) string {
		if val, ok := addr[key].(string); ok {
			return val
		}
		return ""
	}

	return &pb.Address{
		Street:       getString("street"),
		City:         getString("city"),
		State:        getString("state"),
		PostalCode:   getString("postal_code"),
		Country:      getString("country"),
		ContactName:  getString("contact_name"),
		ContactPhone: getString("contact_phone"),
	}, nil
}

// MapPackageToProto конвертирует JSON посылку в протобуф
func MapPackageToProto(packageJSON json.RawMessage) (*pb.Package, error) {
	var pkg map[string]interface{}
	if err := json.Unmarshal(packageJSON, &pkg); err != nil {
		return nil, err
	}

	getString := func(key string) string {
		if val, ok := pkg[key].(string); ok {
			return val
		}
		if val, ok := pkg[key].(float64); ok {
			return fmt.Sprintf("%.2f", val)
		}
		return ""
	}

	return &pb.Package{
		Weight:        getString("weight"),
		Length:        getString("length"),
		Width:         getString("width"),
		Height:        getString("height"),
		Description:   getString("description"),
		DeclaredValue: getString("declared_value"),
	}, nil
}

// MapProviderCodeToEnum конвертирует строковый код в enum
func MapProviderCodeToEnum(code string) pb.DeliveryProvider {
	switch code {
	case "post_express":
		return pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS
	case "bex_express":
		return pb.DeliveryProvider_DELIVERY_PROVIDER_BEX_EXPRESS
	case "aks_express":
		return pb.DeliveryProvider_DELIVERY_PROVIDER_AKS_EXPRESS
	case "d_express":
		return pb.DeliveryProvider_DELIVERY_PROVIDER_D_EXPRESS
	case "city_express":
		return pb.DeliveryProvider_DELIVERY_PROVIDER_CITY_EXPRESS
	default:
		return pb.DeliveryProvider_DELIVERY_PROVIDER_UNSPECIFIED
	}
}

// MapProviderEnumToCode конвертирует enum в строковый код
func MapProviderEnumToCode(provider pb.DeliveryProvider) string {
	switch provider {
	case pb.DeliveryProvider_DELIVERY_PROVIDER_UNSPECIFIED:
		return ""
	case pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS:
		return "post_express"
	case pb.DeliveryProvider_DELIVERY_PROVIDER_BEX_EXPRESS:
		return "bex_express"
	case pb.DeliveryProvider_DELIVERY_PROVIDER_AKS_EXPRESS:
		return "aks_express"
	case pb.DeliveryProvider_DELIVERY_PROVIDER_D_EXPRESS:
		return "d_express"
	case pb.DeliveryProvider_DELIVERY_PROVIDER_CITY_EXPRESS:
		return "city_express"
	default:
		return ""
	}
}

// MapStatusFromProto конвертирует протобуф статус в строку
func MapStatusFromProto(status pb.ShipmentStatus) string {
	switch status {
	case pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED:
		return models.ShipmentStatusPending
	case pb.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return models.ShipmentStatusPending
	case pb.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED:
		return models.ShipmentStatusProcessing
	case pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return models.ShipmentStatusInTransit
	case pb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY:
		return models.ShipmentStatusInTransit
	case pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return models.ShipmentStatusDelivered
	case pb.ShipmentStatus_SHIPMENT_STATUS_FAILED:
		return models.ShipmentStatusFailed
	case pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED: //nolint:misspell // CANCELLED is correct UK spelling used in protobuf
		return models.ShipmentStatusCancelled
	case pb.ShipmentStatus_SHIPMENT_STATUS_RETURNED:
		return models.ShipmentStatusFailed
	default:
		return models.ShipmentStatusPending
	}
}

// MapStatusToProto конвертирует строковый статус в протобуф
func MapStatusToProto(status string) pb.ShipmentStatus {
	switch status {
	case models.ShipmentStatusPending:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case models.ShipmentStatusProcessing:
		return pb.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED
	case models.ShipmentStatusShipped, models.ShipmentStatusInTransit:
		return pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case models.ShipmentStatusDelivered:
		return pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case models.ShipmentStatusFailed:
		return pb.ShipmentStatus_SHIPMENT_STATUS_FAILED
	case models.ShipmentStatusCancelled:
		return pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED //nolint:misspell // CANCELLED is correct UK spelling used in protobuf
	default:
		return pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}

// MapTrackingEventsFromProto конвертирует протобуф события в модели БД
func MapTrackingEventsFromProto(pbEvents []*pb.TrackingEvent) []models.TrackingEvent {
	events := make([]models.TrackingEvent, 0, len(pbEvents))
	for _, pbEvent := range pbEvents {
		event := models.TrackingEvent{
			EventTime:   pbEvent.Timestamp.AsTime(),
			Status:      MapStatusFromProto(pbEvent.Status),
			Location:    &pbEvent.Location,
			Description: &pbEvent.Description,
			CreatedAt:   time.Now(),
		}
		events = append(events, event)
	}
	return events
}

// TimeToProto конвертирует time.Time в protobuf Timestamp
func TimeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// ProtoToTime конвертирует protobuf Timestamp в time.Time
func ProtoToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
