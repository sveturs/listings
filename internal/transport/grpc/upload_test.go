package grpc

import (
	"testing"

	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
)

// TestValidateImageMetadata tests metadata validation
func TestValidateImageMetadata(t *testing.T) {
	s := &Server{}

	tests := []struct {
		name      string
		metadata  *listingspb.UploadImageMetadata
		wantError bool
	}{
		{
			name: "valid metadata",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:    1,
				UserId:       1,
				Filename:     "test.jpg",
				ContentType:  "image/jpeg",
				FileSize:     1024,
				DisplayOrder: 0,
				IsPrimary:    true,
			},
			wantError: false,
		},
		{
			name: "invalid listing_id",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   0,
				UserId:      1,
				Filename:    "test.jpg",
				ContentType: "image/jpeg",
				FileSize:    1024,
			},
			wantError: true,
		},
		{
			name: "invalid user_id",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   1,
				UserId:      0,
				Filename:    "test.jpg",
				ContentType: "image/jpeg",
				FileSize:    1024,
			},
			wantError: true,
		},
		{
			name: "empty filename",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   1,
				UserId:      1,
				Filename:    "",
				ContentType: "image/jpeg",
				FileSize:    1024,
			},
			wantError: true,
		},
		{
			name: "no extension",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   1,
				UserId:      1,
				Filename:    "test",
				ContentType: "image/jpeg",
				FileSize:    1024,
			},
			wantError: true,
		},
		{
			name: "invalid extension",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   1,
				UserId:      1,
				Filename:    "test.txt",
				ContentType: "text/plain",
				FileSize:    1024,
			},
			wantError: true,
		},
		{
			name: "file too large",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   1,
				UserId:      1,
				Filename:    "test.jpg",
				ContentType: "image/jpeg",
				FileSize:    11 * 1024 * 1024, // 11MB > 10MB limit
			},
			wantError: true,
		},
		{
			name: "invalid content type",
			metadata: &listingspb.UploadImageMetadata{
				ListingId:   1,
				UserId:      1,
				Filename:    "test.jpg",
				ContentType: "application/octet-stream",
				FileSize:    1024,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.validateImageMetadata(tt.metadata)
			if (err != nil) != tt.wantError {
				t.Errorf("validateImageMetadata() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestAllowedExtensions tests extension validation
func TestAllowedExtensions(t *testing.T) {
	tests := []struct {
		ext     string
		allowed bool
	}{
		{".jpg", true},
		{".jpeg", true},
		{".png", true},
		{".gif", true},
		{".webp", true},
		{".txt", false},
		{".pdf", false},
		{".mp4", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			if allowedExtensions[tt.ext] != tt.allowed {
				t.Errorf("extension %s: got %v, want %v", tt.ext, allowedExtensions[tt.ext], tt.allowed)
			}
		})
	}
}

// TestConstants verifies upload limits are reasonable
func TestConstants(t *testing.T) {
	if maxImageSize != 10*1024*1024 {
		t.Errorf("maxImageSize = %d, want 10MB", maxImageSize)
	}
	if maxTotalSize != 50*1024*1024 {
		t.Errorf("maxTotalSize = %d, want 50MB", maxTotalSize)
	}
	if maxFiles != 10 {
		t.Errorf("maxFiles = %d, want 10", maxFiles)
	}
	if thumbnailSize != 200 {
		t.Errorf("thumbnailSize = %d, want 200", thumbnailSize)
	}
	if thumbnailQuality != 85 {
		t.Errorf("thumbnailQuality = %d, want 85", thumbnailQuality)
	}
}
