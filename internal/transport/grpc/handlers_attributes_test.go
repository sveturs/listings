package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/attributes/v1"
	"github.com/vondi-global/listings/internal/domain"
)

// MockAttributeService is a mock for service.AttributeService
type MockAttributeService struct {
	mock.Mock
}

func (m *MockAttributeService) CreateAttribute(ctx context.Context, input *domain.CreateAttributeInput) (*domain.Attribute, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeService) UpdateAttribute(ctx context.Context, id int32, input *domain.UpdateAttributeInput) (*domain.Attribute, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeService) DeleteAttribute(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAttributeService) GetAttribute(ctx context.Context, identifier string) (*domain.Attribute, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeService) GetAttributeByID(ctx context.Context, id int32) (*domain.Attribute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeService) GetAttributeByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeService) ListAttributes(ctx context.Context, filter *domain.ListAttributesFilter) ([]*domain.Attribute, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Attribute), args.Get(1).(int64), args.Error(2)
}

func (m *MockAttributeService) LinkAttributeToCategory(ctx context.Context, categoryID string, attributeID int32, settings *domain.CategoryAttributeSettings) error {
	args := m.Called(ctx, categoryID, attributeID, settings)
	return args.Error(0)
}

func (m *MockAttributeService) UpdateCategoryAttribute(ctx context.Context, catAttrID int32, settings *domain.CategoryAttributeSettings) error {
	args := m.Called(ctx, catAttrID, settings)
	return args.Error(0)
}

func (m *MockAttributeService) UnlinkAttributeFromCategory(ctx context.Context, categoryID string, attributeID int32) error {
	args := m.Called(ctx, categoryID, attributeID)
	return args.Error(0)
}

func (m *MockAttributeService) GetCategoryAttributes(ctx context.Context, categoryID string, filter *domain.GetCategoryAttributesFilter) ([]*domain.CategoryAttribute, error) {
	args := m.Called(ctx, categoryID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CategoryAttribute), args.Error(1)
}

func (m *MockAttributeService) GetListingAttributes(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ListingAttributeValue), args.Error(1)
}

func (m *MockAttributeService) SetListingAttributes(ctx context.Context, listingID int32, values []domain.SetListingAttributeValue) error {
	args := m.Called(ctx, listingID, values)
	return args.Error(0)
}

func (m *MockAttributeService) ValidateAttributeValues(ctx context.Context, categoryID string, values []domain.SetListingAttributeValue) error {
	args := m.Called(ctx, categoryID, values)
	return args.Error(0)
}

func (m *MockAttributeService) GetCategoryVariantAttributes(ctx context.Context, categoryID string) ([]*domain.VariantAttribute, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.VariantAttribute), args.Error(1)
}

func (m *MockAttributeService) InvalidateAttributeCache(ctx context.Context, attributeID int32) error {
	args := m.Called(ctx, attributeID)
	return args.Error(0)
}

func (m *MockAttributeService) InvalidateCategoryCache(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

func (m *MockAttributeService) InvalidateListingCache(ctx context.Context, listingID int32) error {
	args := m.Called(ctx, listingID)
	return args.Error(0)
}

// setupTestAttributeServer creates a test server with mocked attribute service
func setupTestAttributeServer() (*Server, *MockAttributeService) {
	mockAttrService := new(MockAttributeService)
	logger := zerolog.Nop()

	server := &Server{
		attrService: mockAttrService,
		logger:      logger,
	}

	return server, mockAttrService
}

// ============================================================================
// Test: CreateAttribute
// ============================================================================

func TestCreateAttribute_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.CreateAttributeRequest{
		Code:          "test_attr",
		AttributeType: pb.AttributeType_ATTRIBUTE_TYPE_TEXT,
		Purpose:       pb.AttributePurpose_ATTRIBUTE_PURPOSE_REGULAR,
	}

	expectedAttr := &domain.Attribute{
		ID:            1,
		Code:          "test_attr",
		Name:          map[string]string{"en": "Test"},
		DisplayName:   map[string]string{"en": "Test Attribute"},
		AttributeType: domain.AttributeTypeText,
		Purpose:       domain.AttributePurposeRegular,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockService.On("CreateAttribute", ctx, mock.AnythingOfType("*domain.CreateAttributeInput")).
		Return(expectedAttr, nil)

	resp, err := server.CreateAttribute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Attribute.Id)
	assert.Equal(t, "test_attr", resp.Attribute.Code)
	mockService.AssertExpectations(t)
}

func TestCreateAttribute_MissingCode(t *testing.T) {
	ctx := context.Background()
	server, _ := setupTestAttributeServer()

	req := &pb.CreateAttributeRequest{
		Code:          "",
		AttributeType: pb.AttributeType_ATTRIBUTE_TYPE_TEXT,
	}

	resp, err := server.CreateAttribute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateAttribute_ServiceError(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.CreateAttributeRequest{
		Code:          "test_attr",
		AttributeType: pb.AttributeType_ATTRIBUTE_TYPE_TEXT,
	}

	mockService.On("CreateAttribute", ctx, mock.AnythingOfType("*domain.CreateAttributeInput")).
		Return(nil, errors.New("database error"))

	resp, err := server.CreateAttribute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

// ============================================================================
// Test: GetAttribute
// ============================================================================

func TestGetAttribute_ByID_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.GetAttributeRequest{
		Identifier: &pb.GetAttributeRequest_Id{Id: 1},
	}

	expectedAttr := &domain.Attribute{
		ID:            1,
		Code:          "test_attr",
		Name:          map[string]string{"en": "Test"},
		DisplayName:   map[string]string{"en": "Test Attribute"},
		AttributeType: domain.AttributeTypeText,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockService.On("GetAttributeByID", ctx, int32(1)).
		Return(expectedAttr, nil)

	resp, err := server.GetAttribute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Attribute.Id)
	mockService.AssertExpectations(t)
}

func TestGetAttribute_ByCode_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.GetAttributeRequest{
		Identifier: &pb.GetAttributeRequest_Code{Code: "test_attr"},
	}

	expectedAttr := &domain.Attribute{
		ID:            1,
		Code:          "test_attr",
		Name:          map[string]string{"en": "Test"},
		AttributeType: domain.AttributeTypeText,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockService.On("GetAttributeByCode", ctx, "test_attr").
		Return(expectedAttr, nil)

	resp, err := server.GetAttribute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test_attr", resp.Attribute.Code)
	mockService.AssertExpectations(t)
}

func TestGetAttribute_NotFound(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.GetAttributeRequest{
		Identifier: &pb.GetAttributeRequest_Id{Id: 999},
	}

	mockService.On("GetAttributeByID", ctx, int32(999)).
		Return(nil, errors.New("attribute not found"))

	resp, err := server.GetAttribute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
	mockService.AssertExpectations(t)
}

// ============================================================================
// Test: ListAttributes
// ============================================================================

func TestListAttributes_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.ListAttributesRequest{
		Page:     1,
		PageSize: 10,
	}

	attrs := []*domain.Attribute{
		{
			ID:            1,
			Code:          "attr1",
			Name:          map[string]string{"en": "Attribute 1"},
			AttributeType: domain.AttributeTypeText,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            2,
			Code:          "attr2",
			Name:          map[string]string{"en": "Attribute 2"},
			AttributeType: domain.AttributeTypeNumber,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockService.On("ListAttributes", ctx, mock.AnythingOfType("*domain.ListAttributesFilter")).
		Return(attrs, int64(2), nil)

	resp, err := server.ListAttributes(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Attributes))
	assert.Equal(t, int32(2), resp.TotalCount)
	mockService.AssertExpectations(t)
}

// ============================================================================
// Test: DeleteAttribute
// ============================================================================

func TestDeleteAttribute_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.DeleteAttributeRequest{
		Id: 1,
	}

	mockService.On("DeleteAttribute", ctx, int32(1)).
		Return(nil)

	resp, err := server.DeleteAttribute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockService.AssertExpectations(t)
}

func TestDeleteAttribute_InvalidID(t *testing.T) {
	ctx := context.Background()
	server, _ := setupTestAttributeServer()

	req := &pb.DeleteAttributeRequest{
		Id: 0,
	}

	resp, err := server.DeleteAttribute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

// ============================================================================
// Test: GetCategoryAttributes
// ============================================================================

func TestGetCategoryAttributes_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.GetCategoryAttributesRequest{
		CategoryId:      "3b4246cc-9970-403c-af01-c142a4178dc6",
		IncludeInactive: false,
	}

	catAttrs := []*domain.CategoryAttribute{
		{
			ID:          1,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 10,
			IsEnabled:   true,
			SortOrder:   1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockService.On("GetCategoryAttributes", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", mock.AnythingOfType("*domain.GetCategoryAttributesFilter")).
		Return(catAttrs, nil)

	resp, err := server.GetCategoryAttributes(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.CategoryAttributes))
	mockService.AssertExpectations(t)
}

// ============================================================================
// Test: GetListingAttributes
// ============================================================================

func TestGetListingAttributes_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.GetListingAttributesRequest{
		ListingId: 1,
	}

	textValue := "test value"
	attrValues := []*domain.ListingAttributeValue{
		{
			ID:          1,
			ListingID:   1,
			AttributeID: 10,
			ValueText:   &textValue,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockService.On("GetListingAttributes", ctx, int32(1)).
		Return(attrValues, nil)

	resp, err := server.GetListingAttributes(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.AttributeValues))
	mockService.AssertExpectations(t)
}

// ============================================================================
// Test: ValidateAttributeValues
// ============================================================================

func TestValidateAttributeValues_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.ValidateAttributeValuesRequest{
		CategoryId: "3b4246cc-9970-403c-af01-c142a4178dc6",
		Values: []*pb.AttributeValueInput{
			{
				AttributeId: 10,
				Value:       &pb.AttributeValueInput_ValueText{ValueText: "test"},
			},
		},
	}

	mockService.On("ValidateAttributeValues", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", mock.AnythingOfType("[]domain.SetListingAttributeValue")).
		Return(nil)

	resp, err := server.ValidateAttributeValues(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsValid)
	assert.Equal(t, 0, len(resp.Errors))
	mockService.AssertExpectations(t)
}

func TestValidateAttributeValues_ValidationFailed(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupTestAttributeServer()

	req := &pb.ValidateAttributeValuesRequest{
		CategoryId: "3b4246cc-9970-403c-af01-c142a4178dc6",
		Values: []*pb.AttributeValueInput{
			{
				AttributeId: 10,
				Value:       &pb.AttributeValueInput_ValueText{ValueText: "invalid"},
			},
		},
	}

	mockService.On("ValidateAttributeValues", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", mock.AnythingOfType("[]domain.SetListingAttributeValue")).
		Return(errors.New("validation failed: invalid value"))

	resp, err := server.ValidateAttributeValues(ctx, req)

	assert.NoError(t, err) // gRPC call succeeds
	assert.NotNil(t, resp)
	assert.False(t, resp.IsValid)
	assert.Greater(t, len(resp.Errors), 0)
	mockService.AssertExpectations(t)
}

// ============================================================================
// Test: Error Conversion
// ============================================================================

func TestConvertServiceError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode codes.Code
	}{
		{
			name:         "NotFound error",
			err:          errors.New("attribute not found"),
			expectedCode: codes.NotFound,
		},
		{
			name:         "Validation error",
			err:          errors.New("validation failed"),
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "Already exists error",
			err:          errors.New("attribute already exists"),
			expectedCode: codes.AlreadyExists,
		},
		{
			name:         "Generic error",
			err:          errors.New("some random error"),
			expectedCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcErr := convertServiceError(tt.err, "test operation")
			assert.Error(t, grpcErr)
			assert.Equal(t, tt.expectedCode, status.Code(grpcErr))
		})
	}
}
