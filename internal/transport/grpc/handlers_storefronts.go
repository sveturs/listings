package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingsv1 "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/service/listings"
)

// CreateStorefront creates a new storefront
func (s *Server) CreateStorefront(ctx context.Context, req *listingsv1.CreateStorefrontRequest) (*listingsv1.StorefrontFull, error) {
	s.logger.Info().Int64("user_id", req.UserId).Msg("CreateStorefront called")

	if err := validateCreateStorefrontRequest(req); err != nil {
		s.logger.Warn().Err(err).Msg("Invalid CreateStorefront request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	serviceReq := &listings.CreateStorefrontRequest{
		UserID:      req.UserId,
		Name:        req.Name,
		Slug:        stringFromOptional(req.Slug),
		Description: getOptionalStringPtr(req.Description),
		Logo:        req.Logo,
		Banner:      req.Banner,
		Theme:       mapProtoStructToJSONB(req.Theme),
		Phone:       getOptionalStringPtr(req.Phone),
		Email:       getOptionalStringPtr(req.Email),
		Website:     getOptionalStringPtr(req.Website),
		Location:    mapProtoLocationToService(req.Location),
		Settings:    mapProtoStructToJSONB(req.Settings),
		SeoMeta:     mapProtoStructToJSONB(req.SeoMeta),
	}

	storefront, err := s.storefrontService.CreateStorefront(ctx, serviceReq)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create storefront")
		return nil, status.Errorf(codes.Internal, "failed to create storefront: %v", err)
	}

	s.logger.Info().Int64("storefront_id", storefront.ID).Msg("Storefront created successfully")

	return mapDomainStorefrontToProto(storefront), nil
}

// GetStorefront retrieves a storefront by ID or slug
func (s *Server) GetStorefront(ctx context.Context, req *listingsv1.GetStorefrontRequest) (*listingsv1.GetStorefrontResponse, error) {
	s.logger.Info().Msg("GetStorefront called")

	var id *int64
	var slug *string

	switch identifier := req.Identifier.(type) {
	case *listingsv1.GetStorefrontRequest_Id:
		id = &identifier.Id
	case *listingsv1.GetStorefrontRequest_Slug:
		slug = &identifier.Slug
	default:
		return nil, status.Error(codes.InvalidArgument, "either id or slug must be provided")
	}

	includes := &domain.Includes{
		Staff:           req.IncludeStaff,
		Hours:           req.IncludeHours,
		PaymentMethods:  req.IncludePaymentMethods,
		DeliveryOptions: req.IncludeDeliveryOptions,
	}

	storefront, err := s.storefrontService.GetStorefront(ctx, id, slug, includes)
	if err != nil {
		s.logger.Error().Err(err).Msg("Storefront not found")
		return nil, status.Errorf(codes.NotFound, "storefront not found: %v", err)
	}

	s.logger.Info().Int64("storefront_id", storefront.ID).Msg("Storefront retrieved successfully")

	return &listingsv1.GetStorefrontResponse{
		Storefront: mapDomainStorefrontToProto(storefront),
	}, nil
}

// GetStorefrontBySlug retrieves a storefront by slug
func (s *Server) GetStorefrontBySlug(ctx context.Context, req *listingsv1.GetStorefrontBySlugRequest) (*listingsv1.GetStorefrontResponse, error) {
	// Convert to GetStorefrontRequest
	getReq := &listingsv1.GetStorefrontRequest{
		Identifier: &listingsv1.GetStorefrontRequest_Slug{
			Slug: req.Slug,
		},
		IncludeStaff:           false,
		IncludeHours:           false,
		IncludePaymentMethods:  false,
		IncludeDeliveryOptions: false,
	}
	return s.GetStorefront(ctx, getReq)
}

// UpdateStorefront updates a storefront
func (s *Server) UpdateStorefront(ctx context.Context, req *listingsv1.UpdateStorefrontRequest) (*listingsv1.StorefrontFull, error) {
	s.logger.Info().Int64("storefront_id", req.Id).Msg("UpdateStorefront called")

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	serviceReq := &listings.UpdateStorefrontRequest{
		Name:                getOptionalStringPtr(req.Name),
		Description:         getOptionalStringPtr(req.Description),
		IsActive:            getOptionalBoolPtr(req.IsActive),
		LogoURL:             getOptionalStringPtr(req.LogoUrl),
		BannerURL:           getOptionalStringPtr(req.BannerUrl),
		Theme:               mapProtoStructToJSONB(req.Theme),
		Phone:               getOptionalStringPtr(req.Phone),
		Email:               getOptionalStringPtr(req.Email),
		Website:             getOptionalStringPtr(req.Website),
		Location:            mapProtoLocationToServicePtr(req.Location),
		Settings:            mapProtoStructToJSONB(req.Settings),
		SeoMeta:             mapProtoStructToJSONB(req.SeoMeta),
		AIAgentEnabled:      getOptionalBoolPtr(req.AiAgentEnabled),
		LiveShoppingEnabled: getOptionalBoolPtr(req.LiveShoppingEnabled),
		GroupBuyingEnabled:  getOptionalBoolPtr(req.GroupBuyingEnabled),
	}

	storefront, err := s.storefrontService.UpdateStorefront(ctx, req.Id, serviceReq)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to update storefront")
		return nil, status.Errorf(codes.Internal, "failed to update storefront: %v", err)
	}

	s.logger.Info().Int64("storefront_id", req.Id).Msg("Storefront updated successfully")

	return mapDomainStorefrontToProto(storefront), nil
}

// DeleteStorefront deletes a storefront
func (s *Server) DeleteStorefront(ctx context.Context, req *listingsv1.DeleteStorefrontRequest) (*listingsv1.DeleteStorefrontResponse, error) {
	s.logger.Info().Int64("storefront_id", req.Id).Msg("DeleteStorefront called")

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.storefrontService.DeleteStorefront(ctx, req.Id, req.HardDelete)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to delete storefront")
		return nil, status.Errorf(codes.Internal, "failed to delete storefront: %v", err)
	}

	s.logger.Info().Int64("storefront_id", req.Id).Msg("Storefront deleted successfully")

	return &listingsv1.DeleteStorefrontResponse{
		Success: true,
		Message: "Storefront deleted successfully",
	}, nil
}

// ListStorefronts lists storefronts with filters
func (s *Server) ListStorefronts(ctx context.Context, req *listingsv1.ListStorefrontsRequest) (*listingsv1.ListStorefrontsResponse, error) {
	s.logger.Info().Msg("ListStorefronts called")

	filter := mapProtoFilterToDomain(req)

	storefronts, total, err := s.storefrontService.ListStorefronts(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list storefronts")
		return nil, status.Errorf(codes.Internal, "failed to list storefronts: %v", err)
	}

	protoStorefronts := make([]*listingsv1.StorefrontFull, len(storefronts))
	for i, sf := range storefronts {
		protoStorefronts[i] = mapDomainStorefrontToProto(&sf)
	}

	s.logger.Info().Int("count", len(storefronts)).Int("total", total).Msg("Storefronts listed successfully")

	return &listingsv1.ListStorefrontsResponse{
		Storefronts: protoStorefronts,
		Total:       int32(total),
	}, nil
}

// GetMyStorefronts retrieves storefronts for a specific user
func (s *Server) GetMyStorefronts(ctx context.Context, req *listingsv1.ListStorefrontsRequest) (*listingsv1.ListStorefrontsResponse, error) {
	return s.ListStorefronts(ctx, req)
}

// AddStaff adds a staff member to a storefront
func (s *Server) AddStaff(ctx context.Context, req *listingsv1.AddStaffRequest) (*listingsv1.StorefrontStaff, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Int64("user_id", req.UserId).Msg("AddStaff called")

	if req.StorefrontId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id and user_id are required")
	}

	role := mapProtoStaffRoleToDomain(req.Role)
	permissions := mapProtoStructToJSONB(req.Permissions)

	staff, err := s.storefrontService.AddStaff(ctx, req.StorefrontId, req.UserId, role, permissions)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to add staff")
		return nil, status.Errorf(codes.Internal, "failed to add staff: %v", err)
	}

	s.logger.Info().Int64("staff_id", staff.ID).Msg("Staff added successfully")

	return mapDomainStaffToProto(staff), nil
}

// UpdateStaff updates a staff member
func (s *Server) UpdateStaff(ctx context.Context, req *listingsv1.UpdateStaffRequest) (*listingsv1.StorefrontStaff, error) {
	s.logger.Info().Int64("staff_id", req.Id).Msg("UpdateStaff called")

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	update := &domain.StaffUpdate{}
	if req.Role != nil {
		role := mapProtoStaffRoleToDomain(*req.Role)
		update.Role = &role
	}
	if req.Permissions != nil {
		update.Permissions = mapProtoStructToJSONB(req.Permissions)
	}

	_, err := s.storefrontService.UpdateStaff(ctx, req.Id, update)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to update staff")
		return nil, status.Errorf(codes.Internal, "failed to update staff: %v", err)
	}

	s.logger.Info().Int64("staff_id", req.Id).Msg("Staff updated successfully")

	return &listingsv1.StorefrontStaff{Id: req.Id}, nil
}

// RemoveStaff removes a staff member from a storefront
func (s *Server) RemoveStaff(ctx context.Context, req *listingsv1.RemoveStaffRequest) (*listingsv1.DeleteStorefrontResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Int64("user_id", req.UserId).Msg("RemoveStaff called")

	if req.StorefrontId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id and user_id are required")
	}

	err := s.storefrontService.RemoveStaff(ctx, req.StorefrontId, req.UserId)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to remove staff")
		return nil, status.Errorf(codes.Internal, "failed to remove staff: %v", err)
	}

	s.logger.Info().Msg("Staff removed successfully")

	return &listingsv1.DeleteStorefrontResponse{
		Success: true,
		Message: "Staff removed successfully",
	}, nil
}

// GetStaff retrieves all staff members for a storefront
func (s *Server) GetStaff(ctx context.Context, req *listingsv1.GetStaffRequest) (*listingsv1.GetStaffResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("GetStaff called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	staff, err := s.storefrontService.GetStaff(ctx, req.StorefrontId)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get staff")
		return nil, status.Errorf(codes.Internal, "failed to get staff: %v", err)
	}

	protoStaff := make([]*listingsv1.StorefrontStaff, len(staff))
	for i, member := range staff {
		protoStaff[i] = mapDomainStaffToProto(&member)
	}

	s.logger.Info().Int("count", len(staff)).Msg("Staff retrieved successfully")

	return &listingsv1.GetStaffResponse{
		Staff: protoStaff,
	}, nil
}

// SetWorkingHours sets working hours for a storefront
func (s *Server) SetWorkingHours(ctx context.Context, req *listingsv1.SetWorkingHoursRequest) (*listingsv1.GetWorkingHoursResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("SetWorkingHours called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	hours := make([]domain.StorefrontHours, len(req.Hours))
	for i, h := range req.Hours {
		hours[i] = *mapProtoHoursToDomain(h)
	}

	err := s.storefrontService.SetWorkingHours(ctx, req.StorefrontId, hours)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to set working hours")
		return nil, status.Errorf(codes.Internal, "failed to set working hours: %v", err)
	}

	savedHours, err := s.storefrontService.GetWorkingHours(ctx, req.StorefrontId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get working hours: %v", err)
	}

	protoHours := make([]*listingsv1.StorefrontHours, len(savedHours))
	for i, h := range savedHours {
		protoHours[i] = mapDomainHoursToProto(&h)
	}

	s.logger.Info().Int("count", len(savedHours)).Msg("Working hours set successfully")

	return &listingsv1.GetWorkingHoursResponse{
		Hours: protoHours,
	}, nil
}

// GetWorkingHours retrieves working hours for a storefront
func (s *Server) GetWorkingHours(ctx context.Context, req *listingsv1.GetWorkingHoursRequest) (*listingsv1.GetWorkingHoursResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("GetWorkingHours called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	hours, err := s.storefrontService.GetWorkingHours(ctx, req.StorefrontId)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get working hours")
		return nil, status.Errorf(codes.Internal, "failed to get working hours: %v", err)
	}

	protoHours := make([]*listingsv1.StorefrontHours, len(hours))
	for i, h := range hours {
		protoHours[i] = mapDomainHoursToProto(&h)
	}

	s.logger.Info().Int("count", len(hours)).Msg("Working hours retrieved successfully")

	return &listingsv1.GetWorkingHoursResponse{
		Hours: protoHours,
	}, nil
}

// IsOpenNow checks if a storefront is currently open
func (s *Server) IsOpenNow(ctx context.Context, req *listingsv1.IsOpenNowRequest) (*listingsv1.IsOpenNowResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("IsOpenNow called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	isOpen, nextOpen, nextClose, err := s.storefrontService.IsOpenNow(ctx, req.StorefrontId)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to check if open")
		return nil, status.Errorf(codes.Internal, "failed to check if open: %v", err)
	}

	response := &listingsv1.IsOpenNowResponse{
		IsOpen: isOpen,
	}

	if nextOpen != nil {
		response.NextOpenTime = getStringPtr(nextOpen.Format("15:04:05"))
	}
	if nextClose != nil {
		response.NextCloseTime = getStringPtr(nextClose.Format("15:04:05"))
	}

	s.logger.Info().Bool("is_open", isOpen).Msg("Open status checked successfully")

	return response, nil
}

// SetPaymentMethods sets payment methods for a storefront
func (s *Server) SetPaymentMethods(ctx context.Context, req *listingsv1.SetPaymentMethodsRequest) (*listingsv1.GetPaymentMethodsResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("SetPaymentMethods called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	methods := make([]domain.PaymentMethod, len(req.Methods))
	for i, m := range req.Methods {
		methods[i] = *mapProtoPaymentMethodToDomain(m)
	}

	err := s.storefrontService.SetPaymentMethods(ctx, req.StorefrontId, methods)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to set payment methods")
		return nil, status.Errorf(codes.Internal, "failed to set payment methods: %v", err)
	}

	savedMethods, err := s.storefrontService.GetPaymentMethods(ctx, req.StorefrontId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment methods: %v", err)
	}

	protoMethods := make([]*listingsv1.StorefrontPaymentMethod, len(savedMethods))
	for i, m := range savedMethods {
		protoMethods[i] = mapDomainPaymentMethodToProto(&m)
	}

	s.logger.Info().Int("count", len(savedMethods)).Msg("Payment methods set successfully")

	return &listingsv1.GetPaymentMethodsResponse{
		Methods: protoMethods,
	}, nil
}

// GetPaymentMethods retrieves payment methods for a storefront
func (s *Server) GetPaymentMethods(ctx context.Context, req *listingsv1.GetPaymentMethodsRequest) (*listingsv1.GetPaymentMethodsResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("GetPaymentMethods called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	methods, err := s.storefrontService.GetPaymentMethods(ctx, req.StorefrontId)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get payment methods")
		return nil, status.Errorf(codes.Internal, "failed to get payment methods: %v", err)
	}

	protoMethods := make([]*listingsv1.StorefrontPaymentMethod, len(methods))
	for i, m := range methods {
		protoMethods[i] = mapDomainPaymentMethodToProto(&m)
	}

	s.logger.Info().Int("count", len(methods)).Msg("Payment methods retrieved successfully")

	return &listingsv1.GetPaymentMethodsResponse{
		Methods: protoMethods,
	}, nil
}

// SetDeliveryOptions sets delivery options for a storefront
func (s *Server) SetDeliveryOptions(ctx context.Context, req *listingsv1.SetDeliveryOptionsRequest) (*listingsv1.GetDeliveryOptionsResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("SetDeliveryOptions called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	options := make([]domain.StorefrontDeliveryOption, len(req.Options))
	for i, o := range req.Options {
		options[i] = *mapProtoDeliveryOptionToDomain(o)
	}

	err := s.storefrontService.SetDeliveryOptions(ctx, req.StorefrontId, options)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to set delivery options")
		return nil, status.Errorf(codes.Internal, "failed to set delivery options: %v", err)
	}

	savedOptions, err := s.storefrontService.GetDeliveryOptions(ctx, req.StorefrontId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get delivery options: %v", err)
	}

	protoOptions := make([]*listingsv1.StorefrontDeliveryOption, len(savedOptions))
	for i, o := range savedOptions {
		protoOptions[i] = mapDomainDeliveryOptionToProto(&o)
	}

	s.logger.Info().Int("count", len(savedOptions)).Msg("Delivery options set successfully")

	return &listingsv1.GetDeliveryOptionsResponse{
		Options: protoOptions,
	}, nil
}

// GetDeliveryOptions retrieves delivery options for a storefront
func (s *Server) GetDeliveryOptions(ctx context.Context, req *listingsv1.GetDeliveryOptionsRequest) (*listingsv1.GetDeliveryOptionsResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("GetDeliveryOptions called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	options, err := s.storefrontService.GetDeliveryOptions(ctx, req.StorefrontId)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get delivery options")
		return nil, status.Errorf(codes.Internal, "failed to get delivery options: %v", err)
	}

	protoOptions := make([]*listingsv1.StorefrontDeliveryOption, len(options))
	for i, o := range options {
		protoOptions[i] = mapDomainDeliveryOptionToProto(&o)
	}

	s.logger.Info().Int("count", len(options)).Msg("Delivery options retrieved successfully")

	return &listingsv1.GetDeliveryOptionsResponse{
		Options: protoOptions,
	}, nil
}

// GetMapData retrieves storefronts for map display
func (s *Server) GetMapData(ctx context.Context, req *listingsv1.GetMapDataRequest) (*listingsv1.GetMapDataResponse, error) {
	s.logger.Info().Msg("GetMapData called")

	bounds := &domain.MapBounds{
		North: req.North,
		South: req.South,
		East:  req.East,
		West:  req.West,
	}

	var filter *domain.ListStorefrontsFilter
	if req.Filter != nil {
		filter = mapProtoFilterToDomain(req.Filter)
	}

	mapData, err := s.storefrontService.GetMapData(ctx, bounds, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get map data")
		return nil, status.Errorf(codes.Internal, "failed to get map data: %v", err)
	}

	protoMapData := make([]*listingsv1.StorefrontMapData, len(mapData))
	for i, data := range mapData {
		protoMapData[i] = mapDomainMapDataToProto(&data)
	}

	s.logger.Info().Int("count", len(mapData)).Msg("Map data retrieved successfully")

	return &listingsv1.GetMapDataResponse{
		Storefronts: protoMapData,
	}, nil
}

// GetDashboardStats retrieves dashboard statistics for a storefront
func (s *Server) GetDashboardStats(ctx context.Context, req *listingsv1.DashboardStatsRequest) (*listingsv1.DashboardStatsResponse, error) {
	s.logger.Info().Int64("storefront_id", req.StorefrontId).Msg("GetDashboardStats called")

	if req.StorefrontId == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	var fromTime, toTime *time.Time
	if req.DateFrom != nil {
		t := req.DateFrom.AsTime()
		fromTime = &t
	}
	if req.DateTo != nil {
		t := req.DateTo.AsTime()
		toTime = &t
	}

	stats, err := s.storefrontService.GetDashboardStats(ctx, req.StorefrontId, fromTime, toTime)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get dashboard stats")
		return nil, status.Errorf(codes.Internal, "failed to get dashboard stats: %v", err)
	}

	s.logger.Info().Msg("Dashboard stats retrieved successfully")

	return mapDomainDashboardStatsToProto(stats), nil
}

// Validation helper function
func validateCreateStorefrontRequest(req *listingsv1.CreateStorefrontRequest) error {
	if req.UserId == 0 {
		return fmt.Errorf("user_id is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Location == nil {
		return fmt.Errorf("location is required")
	}
	if req.Location.City == "" {
		return fmt.Errorf("location.city is required")
	}
	if req.Location.Country == "" {
		return fmt.Errorf("location.country is required")
	}
	return nil
}
