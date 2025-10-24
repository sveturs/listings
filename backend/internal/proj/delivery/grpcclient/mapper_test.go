package grpcclient_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/internal/proj/delivery/grpcclient"
	"backend/internal/proj/delivery/models"
	pb "backend/pkg/grpc/delivery/v1"
)

func TestMapShipmentFromProto(t *testing.T) {
	tests := []struct {
		name        string
		pbShipment  *pb.Shipment
		wantErr     bool
		errContains string
		validate    func(t *testing.T, shipment *models.Shipment)
	}{
		{
			name:        "Nil shipment",
			pbShipment:  nil,
			wantErr:     true,
			errContains: "nil shipment",
		},
		{
			name: "Valid shipment with all fields",
			pbShipment: &pb.Shipment{
				Id:             "EXT123",
				TrackingNumber: "TRK123456",
				Status:         pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
				FromAddress: &pb.Address{
					Street:       "Kneza Milosa 10",
					City:         "Belgrade",
					State:        "Belgrade",
					PostalCode:   "11000",
					Country:      "RS",
					ContactName:  "John Sender",
					ContactPhone: "+381601234567",
				},
				ToAddress: &pb.Address{
					Street:       "Kralja Petra 5",
					City:         "Novi Sad",
					State:        "Vojvodina",
					PostalCode:   "21000",
					Country:      "RS",
					ContactName:  "Jane Recipient",
					ContactPhone: "+381607654321",
				},
				Package: &pb.Package{
					Weight:        "2.5",
					Length:        "30.0",
					Width:         "20.0",
					Height:        "15.0",
					Description:   "Books",
					DeclaredValue: "1000.00",
				},
				Cost:              "500.50",
				CreatedAt:         timestamppb.New(time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)),
				UpdatedAt:         timestamppb.New(time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)),
				EstimatedDelivery: timestamppb.New(time.Date(2025, 1, 20, 18, 0, 0, 0, time.UTC)),
				ActualDelivery:    timestamppb.New(time.Date(2025, 1, 19, 14, 30, 0, 0, time.UTC)),
			},
			wantErr: false,
			validate: func(t *testing.T, shipment *models.Shipment) {
				assert.Equal(t, "EXT123", *shipment.ExternalID)
				assert.Equal(t, "TRK123456", *shipment.TrackingNumber)
				assert.Equal(t, models.ShipmentStatusInTransit, shipment.Status)
				assert.Equal(t, 500.50, *shipment.DeliveryCost)

				// Проверяем sender info
				var senderInfo map[string]interface{}
				err := json.Unmarshal(shipment.SenderInfo, &senderInfo)
				require.NoError(t, err)
				assert.Equal(t, "Kneza Milosa 10", senderInfo["street"])
				assert.Equal(t, "Belgrade", senderInfo["city"])
				assert.Equal(t, "John Sender", senderInfo["contact_name"])

				// Проверяем recipient info
				var recipientInfo map[string]interface{}
				err = json.Unmarshal(shipment.RecipientInfo, &recipientInfo)
				require.NoError(t, err)
				assert.Equal(t, "Novi Sad", recipientInfo["city"])
				assert.Equal(t, "Jane Recipient", recipientInfo["contact_name"])

				// Проверяем package info
				var packageInfo map[string]interface{}
				err = json.Unmarshal(shipment.PackageInfo, &packageInfo)
				require.NoError(t, err)
				assert.Equal(t, "2.5", packageInfo["weight"])
				assert.Equal(t, "Books", packageInfo["description"])

				// Проверяем даты
				assert.NotNil(t, shipment.EstimatedDelivery)
				assert.NotNil(t, shipment.ActualDeliveryDate)
			},
		},
		{
			name: "Shipment with minimal fields",
			pbShipment: &pb.Shipment{
				Id:             "MIN123",
				TrackingNumber: "MIN-TRK",
				Status:         pb.ShipmentStatus_SHIPMENT_STATUS_PENDING,
				FromAddress: &pb.Address{
					Street:  "Street 1",
					City:    "City 1",
					Country: "RS",
				},
				ToAddress: &pb.Address{
					Street:  "Street 2",
					City:    "City 2",
					Country: "RS",
				},
				Package: &pb.Package{
					Weight: "1.0",
				},
				Cost:      "100.00",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: false,
			validate: func(t *testing.T, shipment *models.Shipment) {
				assert.Equal(t, "MIN123", *shipment.ExternalID)
				assert.Equal(t, models.ShipmentStatusPending, shipment.Status)
				assert.Nil(t, shipment.EstimatedDelivery)
				assert.Nil(t, shipment.ActualDeliveryDate)
			},
		},
		{
			name: "Invalid cost format",
			pbShipment: &pb.Shipment{
				Id:             "INV123",
				TrackingNumber: "INV-TRK",
				Status:         pb.ShipmentStatus_SHIPMENT_STATUS_PENDING,
				FromAddress:    &pb.Address{Street: "S1", City: "C1", Country: "RS"},
				ToAddress:      &pb.Address{Street: "S2", City: "C2", Country: "RS"},
				Package:        &pb.Package{Weight: "1.0"},
				Cost:           "invalid_cost",
				CreatedAt:      timestamppb.Now(),
				UpdatedAt:      timestamppb.Now(),
			},
			wantErr: false, // Invalid cost парсится как 0
			validate: func(t *testing.T, shipment *models.Shipment) {
				assert.Equal(t, 0.0, *shipment.DeliveryCost)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shipment, err := grpcclient.MapShipmentFromProto(tt.pbShipment)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, shipment)

			if tt.validate != nil {
				tt.validate(t, shipment)
			}
		})
	}
}

func TestMapAddressToProto(t *testing.T) {
	tests := []struct {
		name        string
		addressJSON json.RawMessage
		want        *pb.Address
		wantErr     bool
	}{
		{
			name: "Complete address",
			addressJSON: json.RawMessage(`{
				"street": "Kneza Milosa 10",
				"city": "Belgrade",
				"state": "Belgrade",
				"postal_code": "11000",
				"country": "RS",
				"contact_name": "John Doe",
				"contact_phone": "+381601234567"
			}`),
			want: &pb.Address{
				Street:       "Kneza Milosa 10",
				City:         "Belgrade",
				State:        "Belgrade",
				PostalCode:   "11000",
				Country:      "RS",
				ContactName:  "John Doe",
				ContactPhone: "+381601234567",
			},
			wantErr: false,
		},
		{
			name: "Minimal address",
			addressJSON: json.RawMessage(`{
				"street": "Street",
				"city": "City",
				"country": "RS"
			}`),
			want: &pb.Address{
				Street:  "Street",
				City:    "City",
				Country: "RS",
			},
			wantErr: false,
		},
		{
			name:        "Invalid JSON",
			addressJSON: json.RawMessage(`{invalid json`),
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := grpcclient.MapAddressToProto(tt.addressJSON)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapPackageToProto(t *testing.T) {
	tests := []struct {
		name        string
		packageJSON json.RawMessage
		want        *pb.Package
		wantErr     bool
	}{
		{
			name: "Complete package",
			packageJSON: json.RawMessage(`{
				"weight": "2.5",
				"length": "30",
				"width": "20",
				"height": "15",
				"description": "Electronics",
				"declared_value": "5000.00"
			}`),
			want: &pb.Package{
				Weight:        "2.5",
				Length:        "30",
				Width:         "20",
				Height:        "15",
				Description:   "Electronics",
				DeclaredValue: "5000.00",
			},
			wantErr: false,
		},
		{
			name: "Numeric values converted to strings",
			packageJSON: json.RawMessage(`{
				"weight": 2.5,
				"length": 30,
				"width": 20,
				"height": 15
			}`),
			want: &pb.Package{
				Weight: "2.50",
				Length: "30.00",
				Width:  "20.00",
				Height: "15.00",
			},
			wantErr: false,
		},
		{
			name:        "Invalid JSON",
			packageJSON: json.RawMessage(`{invalid`),
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := grpcclient.MapPackageToProto(tt.packageJSON)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapProviderCodeToEnum(t *testing.T) {
	tests := []struct {
		code string
		want pb.DeliveryProvider
	}{
		{"post_express", pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS},
		{"bex_express", pb.DeliveryProvider_DELIVERY_PROVIDER_BEX_EXPRESS},
		{"aks_express", pb.DeliveryProvider_DELIVERY_PROVIDER_AKS_EXPRESS},
		{"d_express", pb.DeliveryProvider_DELIVERY_PROVIDER_D_EXPRESS},
		{"city_express", pb.DeliveryProvider_DELIVERY_PROVIDER_CITY_EXPRESS},
		{"unknown", pb.DeliveryProvider_DELIVERY_PROVIDER_UNSPECIFIED},
		{"", pb.DeliveryProvider_DELIVERY_PROVIDER_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got := grpcclient.MapProviderCodeToEnum(tt.code)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapProviderEnumToCode(t *testing.T) {
	tests := []struct {
		provider pb.DeliveryProvider
		want     string
	}{
		{pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS, "post_express"},
		{pb.DeliveryProvider_DELIVERY_PROVIDER_BEX_EXPRESS, "bex_express"},
		{pb.DeliveryProvider_DELIVERY_PROVIDER_AKS_EXPRESS, "aks_express"},
		{pb.DeliveryProvider_DELIVERY_PROVIDER_D_EXPRESS, "d_express"},
		{pb.DeliveryProvider_DELIVERY_PROVIDER_CITY_EXPRESS, "city_express"},
		{pb.DeliveryProvider_DELIVERY_PROVIDER_UNSPECIFIED, ""},
	}

	for _, tt := range tests {
		t.Run(tt.provider.String(), func(t *testing.T) {
			got := grpcclient.MapProviderEnumToCode(tt.provider)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapStatusFromProto(t *testing.T) {
	tests := []struct {
		status pb.ShipmentStatus
		want   string
	}{
		{pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED, models.ShipmentStatusPending},
		{pb.ShipmentStatus_SHIPMENT_STATUS_PENDING, models.ShipmentStatusPending},
		{pb.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED, models.ShipmentStatusProcessing},
		{pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT, models.ShipmentStatusInTransit},
		{pb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY, models.ShipmentStatusInTransit},
		{pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED, models.ShipmentStatusDelivered},
		{pb.ShipmentStatus_SHIPMENT_STATUS_FAILED, models.ShipmentStatusFailed},
		{pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED, models.ShipmentStatusCancelled},
		{pb.ShipmentStatus_SHIPMENT_STATUS_RETURNED, models.ShipmentStatusFailed},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			got := grpcclient.MapStatusFromProto(tt.status)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapStatusToProto(t *testing.T) {
	tests := []struct {
		status string
		want   pb.ShipmentStatus
	}{
		{models.ShipmentStatusPending, pb.ShipmentStatus_SHIPMENT_STATUS_PENDING},
		{models.ShipmentStatusProcessing, pb.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED},
		{models.ShipmentStatusShipped, pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT},
		{models.ShipmentStatusInTransit, pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT},
		{models.ShipmentStatusDelivered, pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED},
		{models.ShipmentStatusFailed, pb.ShipmentStatus_SHIPMENT_STATUS_FAILED},
		{models.ShipmentStatusCancelled, pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED},
		{"unknown", pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := grpcclient.MapStatusToProto(tt.status)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapTrackingEventsFromProto(t *testing.T) {
	now := time.Now()
	pbEvents := []*pb.TrackingEvent{
		{
			Timestamp:   timestamppb.New(now),
			Status:      pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
			Location:    "Belgrade Sorting Center",
			Description: "Package sorted",
		},
		{
			Timestamp:   timestamppb.New(now.Add(time.Hour)),
			Status:      pb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY,
			Location:    "Novi Sad Distribution Center",
			Description: "Out for delivery",
		},
	}

	events := grpcclient.MapTrackingEventsFromProto(pbEvents)

	assert.Len(t, events, 2)
	assert.Equal(t, models.ShipmentStatusInTransit, events[0].Status)
	assert.Equal(t, "Belgrade Sorting Center", *events[0].Location)
	assert.Equal(t, "Package sorted", *events[0].Description)
	assert.Equal(t, models.ShipmentStatusInTransit, events[1].Status)
}

func TestTimeToProto(t *testing.T) {
	now := time.Now()

	// Test with valid time
	pbTime := grpcclient.TimeToProto(&now)
	require.NotNil(t, pbTime)
	assert.Equal(t, now.Unix(), pbTime.AsTime().Unix())

	// Test with nil
	pbTime = grpcclient.TimeToProto(nil)
	assert.Nil(t, pbTime)
}

func TestProtoToTime(t *testing.T) {
	now := time.Now()
	pbTime := timestamppb.New(now)

	// Test with valid timestamp
	goTime := grpcclient.ProtoToTime(pbTime)
	require.NotNil(t, goTime)
	assert.Equal(t, now.Unix(), goTime.Unix())

	// Test with nil
	goTime = grpcclient.ProtoToTime(nil)
	assert.Nil(t, goTime)
}
