package grpc

import (
	"encoding/json"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service/listings"
)

// mapDomainStorefrontToProto maps domain.Storefront to proto StorefrontFull
func mapDomainStorefrontToProto(s *domain.Storefront) *listingspb.StorefrontFull {
	if s == nil {
		return nil
	}

	storefront := &listingspb.StorefrontFull{
		Id:                    s.ID,
		UserId:                s.UserID,
		Slug:                  s.Slug,
		Name:                  s.Name,
		Description:           getOptionalString(s.Description),
		LogoUrl:               getOptionalString(s.LogoURL),
		BannerUrl:             getOptionalString(s.BannerURL),
		Theme:                 mapJSONBToProtoStruct(s.Theme),
		Phone:                 getOptionalString(s.Phone),
		Email:                 getOptionalString(s.Email),
		Website:               getOptionalString(s.Website),
		Address:               getOptionalString(s.Address),
		City:                  getOptionalString(s.City),
		PostalCode:            getOptionalString(s.PostalCode),
		Country:               getOptionalString(s.Country),
		Latitude:              getOptionalDouble(s.Latitude),
		Longitude:             getOptionalDouble(s.Longitude),
		FormattedAddress:      getOptionalString(s.FormattedAddress),
		GeoStrategy:           mapDomainGeoStrategyToProto(s.GeoStrategy),
		DefaultPrivacyLevel:   mapDomainPrivacyLevelToProto(s.DefaultPrivacyLevel),
		AddressVerified:       s.AddressVerified,
		Settings:              mapJSONBToProtoStruct(s.Settings),
		SeoMeta:               mapJSONBToProtoStruct(s.SeoMeta),
		IsActive:              s.IsActive,
		IsVerified:            s.IsVerified,
		VerificationDate:      mapTimeToProtoTimestamp(s.VerificationDate),
		Rating:                s.Rating,
		ReviewsCount:          s.ReviewsCount,
		ProductsCount:         s.ProductsCount,
		SalesCount:            s.SalesCount,
		ViewsCount:            s.ViewsCount,
		SubscriptionPlan:      mapDomainSubscriptionPlanToProto(s.SubscriptionPlan),
		SubscriptionExpiresAt: mapTimeToProtoTimestamp(s.SubscriptionExpiresAt),
		CommissionRate:        s.CommissionRate,
		SubscriptionId:        getOptionalInt64(s.SubscriptionID),
		IsSubscriptionActive:  s.IsSubscriptionActive,
		AiAgentEnabled:        s.AIAgentEnabled,
		AiAgentConfig:         mapJSONBToProtoStruct(s.AIAgentConfig),
		LiveShoppingEnabled:   s.LiveShoppingEnabled,
		GroupBuyingEnabled:    s.GroupBuyingEnabled,
		FollowersCount:        s.FollowersCount,
		CreatedAt:             timestamppb.New(s.CreatedAt),
		UpdatedAt:             timestamppb.New(s.UpdatedAt),
	}

	if len(s.Staff) > 0 {
		storefront.Staff = make([]*listingspb.StorefrontStaff, len(s.Staff))
		for i, staff := range s.Staff {
			storefront.Staff[i] = mapDomainStaffToProto(&staff)
		}
	}

	if len(s.Hours) > 0 {
		storefront.Hours = make([]*listingspb.StorefrontHours, len(s.Hours))
		for i, hour := range s.Hours {
			storefront.Hours[i] = mapDomainHoursToProto(&hour)
		}
	}

	if len(s.PaymentMethods) > 0 {
		storefront.PaymentMethods = make([]*listingspb.StorefrontPaymentMethod, len(s.PaymentMethods))
		for i, method := range s.PaymentMethods {
			storefront.PaymentMethods[i] = mapDomainPaymentMethodToProto(&method)
		}
	}

	if len(s.DeliveryOptions) > 0 {
		storefront.DeliveryOptions = make([]*listingspb.StorefrontDeliveryOption, len(s.DeliveryOptions))
		for i, option := range s.DeliveryOptions {
			storefront.DeliveryOptions[i] = mapDomainDeliveryOptionToProto(&option)
		}
	}

	return storefront
}

// mapDomainStaffToProto maps domain.StorefrontStaff to proto StorefrontStaff
func mapDomainStaffToProto(s *domain.StorefrontStaff) *listingspb.StorefrontStaff {
	if s == nil {
		return nil
	}

	return &listingspb.StorefrontStaff{
		Id:           s.ID,
		StorefrontId: s.StorefrontID,
		UserId:       s.UserID,
		Role:         mapDomainStaffRoleToProto(s.Role),
		Permissions:  mapJSONBToProtoStruct(s.Permissions),
		LastActiveAt: mapTimeToProtoTimestamp(s.LastActiveAt),
		ActionsCount: s.ActionsCount,
		CreatedAt:    timestamppb.New(s.CreatedAt),
		UpdatedAt:    timestamppb.New(s.UpdatedAt),
	}
}

// mapDomainHoursToProto maps domain.StorefrontHours to proto StorefrontHours
func mapDomainHoursToProto(h *domain.StorefrontHours) *listingspb.StorefrontHours {
	if h == nil {
		return nil
	}

	return &listingspb.StorefrontHours{
		Id:           h.ID,
		StorefrontId: h.StorefrontID,
		DayOfWeek:    h.DayOfWeek,
		OpenTime:     getOptionalString(h.OpenTime),
		CloseTime:    getOptionalString(h.CloseTime),
		IsClosed:     h.IsClosed,
		SpecialDate:  getOptionalString(h.SpecialDate),
		SpecialNote:  getOptionalString(h.SpecialNote),
	}
}

// mapProtoHoursToDomain maps proto StorefrontHours to domain.StorefrontHours
func mapProtoHoursToDomain(h *listingspb.StorefrontHours) *domain.StorefrontHours {
	if h == nil {
		return nil
	}

	return &domain.StorefrontHours{
		ID:           h.Id,
		StorefrontID: h.StorefrontId,
		DayOfWeek:    h.DayOfWeek,
		OpenTime:     getOptionalStringPtr(h.OpenTime),
		CloseTime:    getOptionalStringPtr(h.CloseTime),
		IsClosed:     h.IsClosed,
		SpecialDate:  getOptionalStringPtr(h.SpecialDate),
		SpecialNote:  getOptionalStringPtr(h.SpecialNote),
	}
}

// mapDomainPaymentMethodToProto maps domain.PaymentMethod to proto StorefrontPaymentMethod
func mapDomainPaymentMethodToProto(m *domain.PaymentMethod) *listingspb.StorefrontPaymentMethod {
	if m == nil {
		return nil
	}

	return &listingspb.StorefrontPaymentMethod{
		Id:             m.ID,
		StorefrontId:   m.StorefrontID,
		MethodType:     mapDomainPaymentMethodTypeToProto(m.MethodType),
		IsEnabled:      m.IsEnabled,
		Provider:       getOptionalString(m.Provider),
		Settings:       mapJSONBToProtoStruct(m.Settings),
		TransactionFee: m.TransactionFee,
		MinAmount:      getOptionalDouble(m.MinAmount),
		MaxAmount:      getOptionalDouble(m.MaxAmount),
		CreatedAt:      timestamppb.New(m.CreatedAt),
	}
}

// mapProtoPaymentMethodToDomain maps proto StorefrontPaymentMethod to domain.PaymentMethod
func mapProtoPaymentMethodToDomain(m *listingspb.StorefrontPaymentMethod) *domain.PaymentMethod {
	if m == nil {
		return nil
	}

	return &domain.PaymentMethod{
		ID:             m.Id,
		StorefrontID:   m.StorefrontId,
		MethodType:     mapProtoPaymentMethodTypeToDomain(m.MethodType),
		IsEnabled:      m.IsEnabled,
		Provider:       getOptionalStringPtr(m.Provider),
		Settings:       mapProtoStructToJSONB(m.Settings),
		TransactionFee: m.TransactionFee,
		MinAmount:      getOptionalDoublePtr(m.MinAmount),
		MaxAmount:      getOptionalDoublePtr(m.MaxAmount),
	}
}

// mapDomainDeliveryOptionToProto maps domain.StorefrontDeliveryOption to proto StorefrontDeliveryOption
func mapDomainDeliveryOptionToProto(o *domain.StorefrontDeliveryOption) *listingspb.StorefrontDeliveryOption {
	if o == nil {
		return nil
	}

	return &listingspb.StorefrontDeliveryOption{
		Id:               o.ID,
		StorefrontId:     o.StorefrontID,
		Name:             o.Name,
		Description:      getOptionalString(o.Description),
		BasePrice:        o.BasePrice,
		PricePerKm:       o.PricePerKm,
		PricePerKg:       o.PricePerKg,
		FreeAboveAmount:  getOptionalDouble(o.FreeAboveAmount),
		MinOrderAmount:   getOptionalDouble(o.MinOrderAmount),
		MaxWeightKg:      getOptionalDouble(o.MaxWeightKg),
		MaxDistanceKm:    getOptionalDouble(o.MaxDistanceKm),
		EstimatedDaysMin: o.EstimatedDaysMin,
		EstimatedDaysMax: o.EstimatedDaysMax,
		Zones:            mapJSONBToProtoStruct(o.Zones),
		AvailableDays:    mapJSONBToProtoStruct(o.AvailableDays),
		CutoffTime:       getOptionalString(o.CutoffTime),
		Provider:         getOptionalString(o.Provider),
		ProviderConfig:   mapJSONBToProtoStruct(o.ProviderConfig),
		IsActive:         o.IsActive,
		DisplayOrder:     o.DisplayOrder,
		CreatedAt:        timestamppb.New(o.CreatedAt),
		UpdatedAt:        timestamppb.New(o.UpdatedAt),
	}
}

// mapProtoDeliveryOptionToDomain maps proto StorefrontDeliveryOption to domain.StorefrontDeliveryOption
func mapProtoDeliveryOptionToDomain(o *listingspb.StorefrontDeliveryOption) *domain.StorefrontDeliveryOption {
	if o == nil {
		return nil
	}

	return &domain.StorefrontDeliveryOption{
		ID:               o.Id,
		StorefrontID:     o.StorefrontId,
		Name:             o.Name,
		Description:      getOptionalStringPtr(o.Description),
		BasePrice:        o.BasePrice,
		PricePerKm:       o.PricePerKm,
		PricePerKg:       o.PricePerKg,
		FreeAboveAmount:  getOptionalDoublePtr(o.FreeAboveAmount),
		MinOrderAmount:   getOptionalDoublePtr(o.MinOrderAmount),
		MaxWeightKg:      getOptionalDoublePtr(o.MaxWeightKg),
		MaxDistanceKm:    getOptionalDoublePtr(o.MaxDistanceKm),
		EstimatedDaysMin: o.EstimatedDaysMin,
		EstimatedDaysMax: o.EstimatedDaysMax,
		Zones:            mapProtoStructToJSONB(o.Zones),
		AvailableDays:    mapProtoStructToJSONB(o.AvailableDays),
		CutoffTime:       getOptionalStringPtr(o.CutoffTime),
		Provider:         getOptionalStringPtr(o.Provider),
		ProviderConfig:   mapProtoStructToJSONB(o.ProviderConfig),
		IsActive:         o.IsActive,
		DisplayOrder:     o.DisplayOrder,
	}
}

// mapDomainMapDataToProto maps domain.StorefrontMapData to proto StorefrontMapData
func mapDomainMapDataToProto(d *domain.StorefrontMapData) *listingspb.StorefrontMapData {
	if d == nil {
		return nil
	}

	return &listingspb.StorefrontMapData{
		Id:            d.ID,
		Slug:          d.Slug,
		Name:          d.Name,
		Latitude:      d.Latitude,
		Longitude:     d.Longitude,
		Rating:        d.Rating,
		LogoUrl:       stringFromOptional(d.LogoURL),
		Address:       stringFromOptional(d.Address),
		Phone:         stringFromOptional(d.Phone),
		WorkingNow:    d.WorkingNow,
		ProductsCount: d.ProductsCount,
		SupportsCod:   d.SupportsCOD,
		HasDelivery:   d.HasDelivery,
		HasSelfPickup: d.HasSelfPickup,
		AcceptsCards:  d.AcceptsCards,
	}
}

// mapDomainDashboardStatsToProto maps domain.StorefrontDashboardStats to proto DashboardStatsResponse
func mapDomainDashboardStatsToProto(s *domain.StorefrontDashboardStats) *listingspb.DashboardStatsResponse {
	if s == nil {
		return nil
	}

	return &listingspb.DashboardStatsResponse{
		TotalProducts:    s.TotalProducts,
		ActiveProducts:   s.ActiveProducts,
		OrdersCount:      s.OrdersCount,
		Revenue:          s.Revenue,
		AvgOrderValue:    s.AvgOrderValue,
		ViewsCount:       s.ViewsCount,
		UniqueVisitors:   s.UniqueVisitors,
		ConversionRate:   s.ConversionRate,
		PendingOrders:    s.PendingOrders,
		LowStockProducts: s.LowStockProducts,
	}
}

// mapProtoFilterToDomain maps proto ListStorefrontsRequest to domain.ListStorefrontsFilter
func mapProtoFilterToDomain(req *listingspb.ListStorefrontsRequest) *domain.ListStorefrontsFilter {
	// Initialize with default values for required fields
	sortBy := ""
	if req.SortBy != nil {
		sortBy = *req.SortBy
	}
	sortOrder := ""
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	filter := &domain.ListStorefrontsFilter{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      req.Page,
		Limit:     req.Limit,
	}

	if req.UserId != nil {
		filter.UserID = req.UserId
	}
	if req.IsActive != nil {
		filter.IsActive = req.IsActive
	}
	if req.IsVerified != nil {
		filter.IsVerified = req.IsVerified
	}
	if req.City != nil {
		filter.City = req.City
	}
	if req.Country != nil {
		filter.Country = req.Country
	}
	if req.Latitude != nil && req.Longitude != nil && req.RadiusKm != nil {
		filter.Latitude = req.Latitude
		filter.Longitude = req.Longitude
		filter.RadiusKm = req.RadiusKm
	}
	if req.MinRating != nil {
		filter.MinRating = req.MinRating
	}
	if len(req.SubscriptionPlans) > 0 {
		filter.SubscriptionPlans = make([]string, len(req.SubscriptionPlans))
		for i, plan := range req.SubscriptionPlans {
			filter.SubscriptionPlans[i] = mapProtoSubscriptionPlanToDomain(plan)
		}
	}
	if len(req.PaymentMethods) > 0 {
		filter.PaymentMethods = make([]string, len(req.PaymentMethods))
		for i, method := range req.PaymentMethods {
			filter.PaymentMethods[i] = mapProtoPaymentMethodTypeToDomain(method)
		}
	}
	if req.SupportsCod != nil {
		filter.SupportsCOD = req.SupportsCod
	}
	if req.HasAiAgent != nil {
		filter.HasAIAgent = req.HasAiAgent
	}
	if req.HasLiveShopping != nil {
		filter.HasLiveShopping = req.HasLiveShopping
	}
	if req.HasGroupBuying != nil {
		filter.HasGroupBuying = req.HasGroupBuying
	}
	if req.Search != nil {
		filter.Search = req.Search
	}

	return filter
}

// mapProtoLocationToService maps proto Location to listings.StorefrontLocation
func mapProtoLocationToService(loc *listingspb.Location) listings.StorefrontLocation {
	return listings.StorefrontLocation{
		UserLat:     loc.UserLat,
		UserLng:     loc.UserLng,
		FullAddress: loc.FullAddress,
		City:        loc.City,
		PostalCode:  getOptionalStringPtr(loc.PostalCode),
		Country:     loc.Country,
	}
}

// mapProtoLocationToServicePtr maps proto Location to *listings.StorefrontLocation
func mapProtoLocationToServicePtr(loc *listingspb.Location) *listings.StorefrontLocation {
	if loc == nil {
		return nil
	}
	result := mapProtoLocationToService(loc)
	return &result
}

// Enum mapping functions

func mapDomainGeoStrategyToProto(strategy string) listingspb.StorefrontGeoStrategy {
	switch strategy {
	case "storefront_location":
		return listingspb.StorefrontGeoStrategy_STOREFRONT_GEO_STRATEGY_STOREFRONT_LOCATION
	case "individual_location":
		return listingspb.StorefrontGeoStrategy_STOREFRONT_GEO_STRATEGY_INDIVIDUAL_LOCATION
	default:
		return listingspb.StorefrontGeoStrategy_STOREFRONT_GEO_STRATEGY_UNSPECIFIED
	}
}

func mapDomainPrivacyLevelToProto(level string) listingspb.LocationPrivacyLevel {
	switch level {
	case "exact":
		return listingspb.LocationPrivacyLevel_LOCATION_PRIVACY_LEVEL_EXACT
	case "street":
		return listingspb.LocationPrivacyLevel_LOCATION_PRIVACY_LEVEL_STREET
	case "district":
		return listingspb.LocationPrivacyLevel_LOCATION_PRIVACY_LEVEL_DISTRICT
	case "city":
		return listingspb.LocationPrivacyLevel_LOCATION_PRIVACY_LEVEL_CITY
	default:
		return listingspb.LocationPrivacyLevel_LOCATION_PRIVACY_LEVEL_UNSPECIFIED
	}
}

func mapDomainSubscriptionPlanToProto(plan string) listingspb.SubscriptionPlanType {
	switch plan {
	case "starter":
		return listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_STARTER
	case "professional":
		return listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_PROFESSIONAL
	case "business":
		return listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_BUSINESS
	case "enterprise":
		return listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_ENTERPRISE
	default:
		return listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_UNSPECIFIED
	}
}

func mapProtoSubscriptionPlanToDomain(plan listingspb.SubscriptionPlanType) string {
	switch plan {
	case listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_STARTER:
		return "starter"
	case listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_PROFESSIONAL:
		return "professional"
	case listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_BUSINESS:
		return "business"
	case listingspb.SubscriptionPlanType_SUBSCRIPTION_PLAN_TYPE_ENTERPRISE:
		return "enterprise"
	default:
		return "starter"
	}
}

func mapDomainStaffRoleToProto(role string) listingspb.StaffRole {
	switch role {
	case "owner":
		return listingspb.StaffRole_STAFF_ROLE_OWNER
	case "manager":
		return listingspb.StaffRole_STAFF_ROLE_MANAGER
	case "cashier":
		return listingspb.StaffRole_STAFF_ROLE_CASHIER
	case "support":
		return listingspb.StaffRole_STAFF_ROLE_SUPPORT
	case "moderator":
		return listingspb.StaffRole_STAFF_ROLE_MODERATOR
	default:
		return listingspb.StaffRole_STAFF_ROLE_UNSPECIFIED
	}
}

func mapProtoStaffRoleToDomain(role listingspb.StaffRole) string {
	switch role {
	case listingspb.StaffRole_STAFF_ROLE_OWNER:
		return "owner"
	case listingspb.StaffRole_STAFF_ROLE_MANAGER:
		return "manager"
	case listingspb.StaffRole_STAFF_ROLE_CASHIER:
		return "cashier"
	case listingspb.StaffRole_STAFF_ROLE_SUPPORT:
		return "support"
	case listingspb.StaffRole_STAFF_ROLE_MODERATOR:
		return "moderator"
	default:
		return "owner"
	}
}

func mapDomainPaymentMethodTypeToProto(methodType string) listingspb.PaymentMethodType {
	switch methodType {
	case "cash":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_CASH
	case "cod":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_COD
	case "card":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_CARD
	case "bank_transfer":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_BANK_TRANSFER
	case "paypal":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAL
	case "crypto":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_CRYPTO
	case "postanska":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_POSTANSKA
	case "keks_pay":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_KEKS_PAY
	case "ips":
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_IPS
	default:
		return listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED
	}
}

func mapProtoPaymentMethodTypeToDomain(methodType listingspb.PaymentMethodType) string {
	switch methodType {
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_CASH:
		return "cash"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_COD:
		return "cod"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_CARD:
		return "card"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_BANK_TRANSFER:
		return "bank_transfer"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAL:
		return "paypal"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_CRYPTO:
		return "crypto"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_POSTANSKA:
		return "postanska"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_KEKS_PAY:
		return "keks_pay"
	case listingspb.PaymentMethodType_PAYMENT_METHOD_TYPE_IPS:
		return "ips"
	default:
		return "cash"
	}
}

// JSONB and Struct mapping helpers

func mapProtoStructToJSONB(s *structpb.Struct) domain.JSONB {
	if s == nil {
		return nil
	}
	data, err := s.MarshalJSON()
	if err != nil {
		return nil
	}
	return domain.JSONB(data)
}

func mapJSONBToProtoStruct(j domain.JSONB) *structpb.Struct {
	if len(j) == 0 {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(j, &result); err != nil {
		return nil
	}
	s, err := structpb.NewStruct(result)
	if err != nil {
		return nil
	}
	return s
}

// Helper functions for optional values

func getOptionalString(s *string) *string {
	if s == nil {
		return nil
	}
	return s
}

func getOptionalStringPtr(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func getOptionalDouble(d *float64) *float64 {
	if d == nil {
		return nil
	}
	return d
}

func getOptionalDoublePtr(d *float64) *float64 {
	if d == nil {
		return nil
	}
	return d
}

func getOptionalInt64(i *int64) *int64 {
	if i == nil {
		return nil
	}
	return i
}

func getOptionalBoolPtr(b *bool) *bool {
	if b == nil {
		return nil
	}
	return b
}

func getStringPtr(s string) *string {
	return &s
}

func mapTimeToProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// stringFromOptional converts *string to string, returning "" if nil
func stringFromOptional(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
