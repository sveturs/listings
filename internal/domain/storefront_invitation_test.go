package domain

import (
	"testing"
	"time"
)

func TestGenerateInviteCode(t *testing.T) {
	code1, err := GenerateInviteCode()
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	// Check prefix
	if len(code1) != len(InviteCodePrefix)+16 {
		t.Errorf("expected code length %d, got %d", len(InviteCodePrefix)+16, len(code1))
	}

	if code1[:len(InviteCodePrefix)] != InviteCodePrefix {
		t.Errorf("expected prefix %s, got %s", InviteCodePrefix, code1[:len(InviteCodePrefix)])
	}

	// Check uniqueness
	code2, err := GenerateInviteCode()
	if err != nil {
		t.Fatalf("failed to generate second code: %v", err)
	}

	if code1 == code2 {
		t.Error("expected unique codes, got duplicates")
	}
}

func TestStorefrontInvitation_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "nil expires_at (never expires)",
			expiresAt: nil,
			want:      false,
		},
		{
			name:      "expires in future",
			expiresAt: timePtr(time.Now().Add(1 * time.Hour)),
			want:      false,
		},
		{
			name:      "expired in past",
			expiresAt: timePtr(time.Now().Add(-1 * time.Hour)),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := &StorefrontInvitation{ExpiresAt: tt.expiresAt}
			if got := inv.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorefrontInvitation_CanAccept(t *testing.T) {
	tests := []struct {
		name       string
		status     StorefrontInvitationStatus
		expiresAt  *time.Time
		maxUses    *int32
		currentUses int32
		invType    StorefrontInvitationType
		want       bool
	}{
		{
			name:    "pending email invitation",
			status:  InvitationStatusPending,
			invType: InvitationTypeEmail,
			want:    true,
		},
		{
			name:    "already accepted",
			status:  InvitationStatusAccepted,
			invType: InvitationTypeEmail,
			want:    false,
		},
		{
			name:      "expired invitation",
			status:    InvitationStatusPending,
			expiresAt: timePtr(time.Now().Add(-1 * time.Hour)),
			invType:   InvitationTypeEmail,
			want:      false,
		},
		{
			name:        "link with remaining uses",
			status:      InvitationStatusPending,
			invType:     InvitationTypeLink,
			maxUses:     int32Ptr(5),
			currentUses: 3,
			want:        true,
		},
		{
			name:        "link at max uses",
			status:      InvitationStatusPending,
			invType:     InvitationTypeLink,
			maxUses:     int32Ptr(5),
			currentUses: 5,
			want:        false,
		},
		{
			name:        "link with unlimited uses",
			status:      InvitationStatusPending,
			invType:     InvitationTypeLink,
			maxUses:     nil,
			currentUses: 100,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := &StorefrontInvitation{
				Status:      tt.status,
				ExpiresAt:   tt.expiresAt,
				Type:        tt.invType,
				MaxUses:     tt.maxUses,
				CurrentUses: tt.currentUses,
			}
			if got := inv.CanAccept(); got != tt.want {
				t.Errorf("CanAccept() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorefrontInvitation_IncrementUses(t *testing.T) {
	tests := []struct {
		name        string
		invType     StorefrontInvitationType
		maxUses     *int32
		currentUses int32
		wantErr     bool
		wantUses    int32
	}{
		{
			name:    "email invitation (error)",
			invType: InvitationTypeEmail,
			wantErr: true,
		},
		{
			name:        "link with remaining uses",
			invType:     InvitationTypeLink,
			maxUses:     int32Ptr(5),
			currentUses: 3,
			wantErr:     false,
			wantUses:    4,
		},
		{
			name:        "link at max uses",
			invType:     InvitationTypeLink,
			maxUses:     int32Ptr(5),
			currentUses: 5,
			wantErr:     true,
			wantUses:    5,
		},
		{
			name:        "link with unlimited uses",
			invType:     InvitationTypeLink,
			maxUses:     nil,
			currentUses: 100,
			wantErr:     false,
			wantUses:    101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := &StorefrontInvitation{
				Type:        tt.invType,
				MaxUses:     tt.maxUses,
				CurrentUses: tt.currentUses,
			}
			err := inv.IncrementUses()
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrementUses() error = %v, wantErr %v", err, tt.wantErr)
			}
			if inv.CurrentUses != tt.wantUses {
				t.Errorf("CurrentUses = %d, want %d", inv.CurrentUses, tt.wantUses)
			}
		})
	}
}

func TestStorefrontInvitation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		inv     *StorefrontInvitation
		wantErr bool
	}{
		{
			name: "valid email invitation",
			inv: &StorefrontInvitation{
				StorefrontID:  1,
				InvitedByID:   1,
				Role:          "staff",
				Type:          InvitationTypeEmail,
				InvitedEmail:  strPtr("test@example.com"),
			},
			wantErr: false,
		},
		{
			name: "valid link invitation",
			inv: &StorefrontInvitation{
				StorefrontID: 1,
				InvitedByID:  1,
				Role:         "manager",
				Type:         InvitationTypeLink,
				InviteCode:   strPtr("sf-abc123"),
			},
			wantErr: false,
		},
		{
			name: "missing storefront_id",
			inv: &StorefrontInvitation{
				InvitedByID: 1,
				Role:        "staff",
				Type:        InvitationTypeEmail,
				InvitedEmail: strPtr("test@example.com"),
			},
			wantErr: true,
		},
		{
			name: "missing invited_by_id",
			inv: &StorefrontInvitation{
				StorefrontID: 1,
				Role:         "staff",
				Type:         InvitationTypeEmail,
				InvitedEmail: strPtr("test@example.com"),
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			inv: &StorefrontInvitation{
				StorefrontID: 1,
				InvitedByID:  1,
				Role:         "invalid",
				Type:         InvitationTypeEmail,
				InvitedEmail: strPtr("test@example.com"),
			},
			wantErr: true,
		},
		{
			name: "email invitation without email",
			inv: &StorefrontInvitation{
				StorefrontID: 1,
				InvitedByID:  1,
				Role:         "staff",
				Type:         InvitationTypeEmail,
			},
			wantErr: true,
		},
		{
			name: "link invitation without code",
			inv: &StorefrontInvitation{
				StorefrontID: 1,
				InvitedByID:  1,
				Role:         "staff",
				Type:         InvitationTypeLink,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inv.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorefrontInvitation_StatusMethods(t *testing.T) {
	inv := &StorefrontInvitation{Status: InvitationStatusPending}

	// Test MarkAsAccepted
	inv.MarkAsAccepted()
	if !inv.IsAccepted() {
		t.Error("expected IsAccepted() to be true after MarkAsAccepted()")
	}
	if inv.AcceptedAt == nil {
		t.Error("expected AcceptedAt to be set")
	}

	// Test MarkAsDeclined
	inv = &StorefrontInvitation{Status: InvitationStatusPending}
	inv.MarkAsDeclined()
	if !inv.IsDeclined() {
		t.Error("expected IsDeclined() to be true after MarkAsDeclined()")
	}
	if inv.DeclinedAt == nil {
		t.Error("expected DeclinedAt to be set")
	}

	// Test MarkAsExpired
	inv = &StorefrontInvitation{Status: InvitationStatusPending}
	inv.MarkAsExpired()
	if inv.Status != InvitationStatusExpired {
		t.Error("expected status to be expired")
	}

	// Test MarkAsRevoked
	inv = &StorefrontInvitation{Status: InvitationStatusPending}
	inv.MarkAsRevoked()
	if !inv.IsRevoked() {
		t.Error("expected IsRevoked() to be true after MarkAsRevoked()")
	}
}

// Helper functions
func timePtr(t time.Time) *time.Time {
	return &t
}

func int32Ptr(i int32) *int32 {
	return &i
}

func strPtr(s string) *string {
	return &s
}
