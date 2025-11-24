package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// InventoryReservation Validation Tests
// =============================================================================

func TestInventoryReservation_Validate_Success(t *testing.T) {
	reservation := &InventoryReservation{
		ListingID: 100,
		OrderID:   1,
		Quantity:  5,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err := reservation.Validate()
	assert.NoError(t, err)
}

func TestInventoryReservation_Validate_Failures(t *testing.T) {
	tests := []struct {
		name        string
		reservation *InventoryReservation
		wantErr     string
	}{
		{
			name:        "nil reservation",
			reservation: nil,
			wantErr:     "reservation cannot be nil",
		},
		{
			name: "invalid listing_id",
			reservation: &InventoryReservation{
				ListingID: 0,
				OrderID:   1,
				Quantity:  5,
				ExpiresAt: time.Now().Add(30 * time.Minute),
			},
			wantErr: "listing_id must be greater than 0",
		},
		{
			name: "invalid order_id",
			reservation: &InventoryReservation{
				ListingID: 100,
				OrderID:   0,
				Quantity:  5,
				ExpiresAt: time.Now().Add(30 * time.Minute),
			},
			wantErr: "order_id must be greater than 0",
		},
		{
			name: "invalid quantity",
			reservation: &InventoryReservation{
				ListingID: 100,
				OrderID:   1,
				Quantity:  0,
				ExpiresAt: time.Now().Add(30 * time.Minute),
			},
			wantErr: "quantity must be greater than 0",
		},
		{
			name: "zero expires_at",
			reservation: &InventoryReservation{
				ListingID: 100,
				OrderID:   1,
				Quantity:  5,
				ExpiresAt: time.Time{},
			},
			wantErr: "expires_at is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.reservation.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// =============================================================================
// Expiration Tests
// =============================================================================

func TestInventoryReservation_IsExpired_NotExpired(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	assert.False(t, reservation.IsExpired())
}

func TestInventoryReservation_IsExpired_Expired(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(-1 * time.Minute), // 1 minute ago
	}

	assert.True(t, reservation.IsExpired())
}

func TestInventoryReservation_IsExpired_AlreadyCommitted(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusCommitted,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired time but status is committed
	}

	assert.False(t, reservation.IsExpired()) // Not expired because not active
}

func TestInventoryReservation_IsExpired_NilReservation(t *testing.T) {
	var reservation *InventoryReservation
	assert.True(t, reservation.IsExpired())
}

// =============================================================================
// State Transition Tests
// =============================================================================

func TestInventoryReservation_CanCommit_Active(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	assert.True(t, reservation.CanCommit())
}

func TestInventoryReservation_CanCommit_Expired(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	assert.False(t, reservation.CanCommit()) // Can't commit expired
}

func TestInventoryReservation_CanCommit_AlreadyCommitted(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusCommitted,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	assert.False(t, reservation.CanCommit())
}

func TestInventoryReservation_CanCommit_NilReservation(t *testing.T) {
	var reservation *InventoryReservation
	assert.False(t, reservation.CanCommit())
}

func TestInventoryReservation_CanRelease_Active(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	assert.True(t, reservation.CanRelease())
}

func TestInventoryReservation_CanRelease_ActiveButExpired(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	assert.True(t, reservation.CanRelease()) // Can release even if expired
}

func TestInventoryReservation_CanRelease_Committed(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusCommitted,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	assert.False(t, reservation.CanRelease())
}

func TestInventoryReservation_CanRelease_NilReservation(t *testing.T) {
	var reservation *InventoryReservation
	assert.False(t, reservation.CanRelease())
}

func TestInventoryReservation_Commit_Success(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	oldUpdatedAt := reservation.UpdatedAt

	err := reservation.Commit()
	require.NoError(t, err)
	assert.Equal(t, ReservationStatusCommitted, reservation.Status)
	assert.NotNil(t, reservation.CommittedAt)
	assert.True(t, reservation.UpdatedAt.After(oldUpdatedAt))
}

func TestInventoryReservation_Commit_AlreadyCommitted(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusCommitted,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err := reservation.Commit()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reservation cannot be committed")
}

func TestInventoryReservation_Commit_Expired(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	err := reservation.Commit()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reservation cannot be committed")
}

func TestInventoryReservation_Release_Success(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	oldUpdatedAt := reservation.UpdatedAt

	err := reservation.Release()
	require.NoError(t, err)
	assert.Equal(t, ReservationStatusReleased, reservation.Status)
	assert.NotNil(t, reservation.ReleasedAt)
	assert.True(t, reservation.UpdatedAt.After(oldUpdatedAt))
}

func TestInventoryReservation_Release_NotActive(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusCommitted,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err := reservation.Release()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reservation cannot be released")
}

func TestInventoryReservation_Expire_Success(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusActive,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	oldUpdatedAt := reservation.UpdatedAt

	err := reservation.Expire()
	require.NoError(t, err)
	assert.Equal(t, ReservationStatusExpired, reservation.Status)
	assert.True(t, reservation.UpdatedAt.After(oldUpdatedAt))
}

func TestInventoryReservation_Expire_NotActive(t *testing.T) {
	reservation := &InventoryReservation{
		Status:    ReservationStatusCommitted,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	err := reservation.Expire()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only active reservations can be expired")
}

func TestInventoryReservation_Expire_NilReservation(t *testing.T) {
	var reservation *InventoryReservation
	err := reservation.Expire()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reservation cannot be nil")
}

// =============================================================================
// TTL Calculation Tests
// =============================================================================

func TestInventoryReservation_CalculateTTL_NotExpired(t *testing.T) {
	reservation := &InventoryReservation{
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	ttl := reservation.CalculateTTL()
	assert.GreaterOrEqual(t, ttl, int64(29)) // At least 29 minutes
	assert.LessOrEqual(t, ttl, int64(30))    // At most 30 minutes
}

func TestInventoryReservation_CalculateTTL_Expired(t *testing.T) {
	reservation := &InventoryReservation{
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	ttl := reservation.CalculateTTL()
	assert.Equal(t, int64(0), ttl)
}

func TestInventoryReservation_CalculateTTL_NilReservation(t *testing.T) {
	var reservation *InventoryReservation
	ttl := reservation.CalculateTTL()
	assert.Equal(t, int64(0), ttl)
}

func TestInventoryReservation_CalculateTTL_ZeroExpiresAt(t *testing.T) {
	reservation := &InventoryReservation{
		ExpiresAt: time.Time{},
	}

	ttl := reservation.CalculateTTL()
	assert.Equal(t, int64(0), ttl)
}

// =============================================================================
// Factory Method Tests
// =============================================================================

func TestNewInventoryReservation_Success(t *testing.T) {
	variantID := int64(10)
	reservation := NewInventoryReservation(100, &variantID, 1, 5)

	require.NotNil(t, reservation)
	assert.Equal(t, int64(100), reservation.ListingID)
	assert.Equal(t, int64(10), *reservation.VariantID)
	assert.Equal(t, int64(1), reservation.OrderID)
	assert.Equal(t, int32(5), reservation.Quantity)
	assert.Equal(t, ReservationStatusActive, reservation.Status)
	assert.False(t, reservation.ExpiresAt.IsZero())
	assert.False(t, reservation.CreatedAt.IsZero())
	assert.False(t, reservation.UpdatedAt.IsZero())

	// Should expire in approximately 30 minutes
	ttl := reservation.CalculateTTL()
	assert.GreaterOrEqual(t, ttl, int64(29))
	assert.LessOrEqual(t, ttl, int64(30))
}

func TestNewInventoryReservation_NilVariant(t *testing.T) {
	reservation := NewInventoryReservation(100, nil, 1, 5)

	require.NotNil(t, reservation)
	assert.Nil(t, reservation.VariantID)
}

func TestNewInventoryReservationWithTTL_CustomTTL(t *testing.T) {
	variantID := int64(10)
	reservation := NewInventoryReservationWithTTL(100, &variantID, 1, 5, 60)

	require.NotNil(t, reservation)
	assert.Equal(t, int64(100), reservation.ListingID)
	assert.Equal(t, ReservationStatusActive, reservation.Status)

	// Should expire in approximately 60 minutes
	ttl := reservation.CalculateTTL()
	assert.GreaterOrEqual(t, ttl, int64(59))
	assert.LessOrEqual(t, ttl, int64(60))
}

func TestNewInventoryReservationWithTTL_ZeroTTL(t *testing.T) {
	reservation := NewInventoryReservationWithTTL(100, nil, 1, 5, 0)

	require.NotNil(t, reservation)
	assert.True(t, reservation.IsExpired()) // Should be immediately expired
}

// =============================================================================
// ReservationStatus Proto Conversion Tests
// =============================================================================

func TestReservationStatusFromProto(t *testing.T) {
	tests := []struct {
		pbStatus pb.ReservationStatus
		expected ReservationStatus
	}{
		{pb.ReservationStatus_RESERVATION_STATUS_ACTIVE, ReservationStatusActive},
		{pb.ReservationStatus_RESERVATION_STATUS_COMMITTED, ReservationStatusCommitted},
		{pb.ReservationStatus_RESERVATION_STATUS_RELEASED, ReservationStatusReleased},
		{pb.ReservationStatus_RESERVATION_STATUS_EXPIRED, ReservationStatusExpired},
		{pb.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED, ReservationStatusUnspecified},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			result := ReservationStatusFromProto(tt.pbStatus)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReservationStatus_ToProtoReservationStatus(t *testing.T) {
	tests := []struct {
		status   ReservationStatus
		expected pb.ReservationStatus
	}{
		{ReservationStatusActive, pb.ReservationStatus_RESERVATION_STATUS_ACTIVE},
		{ReservationStatusCommitted, pb.ReservationStatus_RESERVATION_STATUS_COMMITTED},
		{ReservationStatusReleased, pb.ReservationStatus_RESERVATION_STATUS_RELEASED},
		{ReservationStatusExpired, pb.ReservationStatus_RESERVATION_STATUS_EXPIRED},
		{ReservationStatusUnspecified, pb.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := tt.status.ToProtoReservationStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// =============================================================================
// InventoryReservation Proto Conversion Tests
// =============================================================================

func TestInventoryReservationFromProto_Success(t *testing.T) {
	now := time.Now()
	variantID := int64(10)
	committedAt := now.Add(-1 * time.Hour)

	pbReservation := &pb.InventoryReservation{
		Id:          1,
		ListingId:   100,
		VariantId:   &variantID,
		OrderId:     1,
		Quantity:    5,
		Status:      pb.ReservationStatus_RESERVATION_STATUS_COMMITTED,
		ExpiresAt:   timestamppb.New(now.Add(30 * time.Minute)),
		CreatedAt:   timestamppb.New(now.Add(-2 * time.Hour)),
		UpdatedAt:   timestamppb.New(now),
		CommittedAt: timestamppb.New(committedAt),
	}

	reservation := InventoryReservationFromProto(pbReservation)
	require.NotNil(t, reservation)
	assert.Equal(t, int64(1), reservation.ID)
	assert.Equal(t, int64(100), reservation.ListingID)
	assert.Equal(t, int64(10), *reservation.VariantID)
	assert.Equal(t, int64(1), reservation.OrderID)
	assert.Equal(t, int32(5), reservation.Quantity)
	assert.Equal(t, ReservationStatusCommitted, reservation.Status)
	assert.NotNil(t, reservation.CommittedAt)
}

func TestInventoryReservationFromProto_Nil(t *testing.T) {
	reservation := InventoryReservationFromProto(nil)
	assert.Nil(t, reservation)
}

func TestInventoryReservation_ToProto_Success(t *testing.T) {
	now := time.Now()
	variantID := int64(10)
	committedAt := now.Add(-1 * time.Hour)

	reservation := &InventoryReservation{
		ID:          1,
		ListingID:   100,
		VariantID:   &variantID,
		OrderID:     1,
		Quantity:    5,
		Status:      ReservationStatusCommitted,
		ExpiresAt:   now.Add(30 * time.Minute),
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now,
		CommittedAt: &committedAt,
	}

	pbReservation := reservation.ToProto()
	require.NotNil(t, pbReservation)
	assert.Equal(t, int64(1), pbReservation.Id)
	assert.Equal(t, int64(100), pbReservation.ListingId)
	assert.Equal(t, int64(10), *pbReservation.VariantId)
	assert.Equal(t, int64(1), pbReservation.OrderId)
	assert.Equal(t, int32(5), pbReservation.Quantity)
	assert.Equal(t, pb.ReservationStatus_RESERVATION_STATUS_COMMITTED, pbReservation.Status)
	assert.NotNil(t, pbReservation.CommittedAt)
}

func TestInventoryReservation_ToProto_Nil(t *testing.T) {
	var reservation *InventoryReservation
	pbReservation := reservation.ToProto()
	assert.Nil(t, pbReservation)
}

func TestInventoryReservation_ToProto_NoOptionalFields(t *testing.T) {
	now := time.Now()

	reservation := &InventoryReservation{
		ID:        1,
		ListingID: 100,
		OrderID:   1,
		Quantity:  5,
		Status:    ReservationStatusActive,
		ExpiresAt: now.Add(30 * time.Minute),
		CreatedAt: now,
		UpdatedAt: now,
		// No VariantID, CommittedAt, ReleasedAt
	}

	pbReservation := reservation.ToProto()
	require.NotNil(t, pbReservation)
	assert.Nil(t, pbReservation.VariantId)
	assert.Nil(t, pbReservation.CommittedAt)
	assert.Nil(t, pbReservation.ReleasedAt)
}
