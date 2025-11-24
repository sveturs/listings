package grpc

import (
	searchv1 "github.com/sveturs/listings/api/proto/search/v1"
	"github.com/sveturs/listings/internal/service/search"
)

// ============================================================================
// PHASE 21.2: Advanced Search Converters (Proto ↔ Domain)
// ============================================================================

// ProtoToFacetsRequest converts proto GetSearchFacetsRequest to domain FacetsRequest
func ProtoToFacetsRequest(req *searchv1.GetSearchFacetsRequest) *search.FacetsRequest {
	domainReq := &search.FacetsRequest{
		UseCache: true, // Default to using cache
	}

	// Optional query
	if req.Query != nil {
		domainReq.Query = *req.Query
	}

	// Optional category_id
	if req.CategoryId != nil {
		categoryID := *req.CategoryId
		domainReq.CategoryID = &categoryID
	}

	// Optional filters
	if req.Filters != nil {
		domainReq.Filters = ProtoToSearchFilters(req.Filters)
	}

	return domainReq
}

// FacetsResponseToProto converts domain FacetsResponse to proto GetSearchFacetsResponse
func FacetsResponseToProto(resp *search.FacetsResponse) *searchv1.GetSearchFacetsResponse {
	protoResp := &searchv1.GetSearchFacetsResponse{
		TookMs: resp.TookMs,
		Cached: resp.Cached,
	}

	// Convert categories
	if len(resp.Categories) > 0 {
		protoResp.Categories = make([]*searchv1.CategoryFacet, 0, len(resp.Categories))
		for _, cat := range resp.Categories {
			protoResp.Categories = append(protoResp.Categories, &searchv1.CategoryFacet{
				CategoryId: cat.CategoryID,
				Count:      cat.Count,
			})
		}
	}

	// Convert price ranges
	if len(resp.PriceRanges) > 0 {
		protoResp.PriceRanges = make([]*searchv1.PriceRangeFacet, 0, len(resp.PriceRanges))
		for _, pr := range resp.PriceRanges {
			protoResp.PriceRanges = append(protoResp.PriceRanges, &searchv1.PriceRangeFacet{
				Min:   pr.Min,
				Max:   pr.Max,
				Count: pr.Count,
			})
		}
	}

	// Convert attributes
	if len(resp.Attributes) > 0 {
		protoResp.Attributes = make(map[string]*searchv1.AttributeFacet)
		for key, attr := range resp.Attributes {
			protoAttr := &searchv1.AttributeFacet{
				Key: attr.Key,
			}

			if len(attr.Values) > 0 {
				protoAttr.Values = make([]*searchv1.AttributeValueCount, 0, len(attr.Values))
				for _, val := range attr.Values {
					protoAttr.Values = append(protoAttr.Values, &searchv1.AttributeValueCount{
						Value: val.Value,
						Count: val.Count,
					})
				}
			}

			protoResp.Attributes[key] = protoAttr
		}
	}

	// Convert source types
	if len(resp.SourceTypes) > 0 {
		protoResp.SourceTypes = make([]*searchv1.Facet, 0, len(resp.SourceTypes))
		for _, st := range resp.SourceTypes {
			protoResp.SourceTypes = append(protoResp.SourceTypes, &searchv1.Facet{
				Key:   st.Key,
				Count: st.Count,
			})
		}
	}

	// Convert stock statuses
	if len(resp.StockStatuses) > 0 {
		protoResp.StockStatuses = make([]*searchv1.Facet, 0, len(resp.StockStatuses))
		for _, ss := range resp.StockStatuses {
			protoResp.StockStatuses = append(protoResp.StockStatuses, &searchv1.Facet{
				Key:   ss.Key,
				Count: ss.Count,
			})
		}
	}

	return protoResp
}

// ProtoToSearchFiltersRequest converts proto SearchWithFiltersRequest to domain SearchFiltersRequest
func ProtoToSearchFiltersRequest(req *searchv1.SearchWithFiltersRequest) *search.SearchFiltersRequest {
	domainReq := &search.SearchFiltersRequest{
		Query:         req.Query,
		Limit:         req.Limit,
		Offset:        req.Offset,
		UseCache:      req.UseCache,
		IncludeFacets: req.IncludeFacets,
	}

	// Optional category_id
	if req.CategoryId != nil {
		categoryID := *req.CategoryId
		domainReq.CategoryID = &categoryID
	}

	// Optional filters
	if req.Filters != nil {
		domainReq.Filters = ProtoToSearchFilters(req.Filters)
	}

	// Optional sort
	if req.Sort != nil {
		domainReq.Sort = ProtoToSortConfig(req.Sort)
	}

	// Set defaults
	if domainReq.Limit == 0 {
		domainReq.Limit = 20
	}

	return domainReq
}

// SearchFiltersResponseToProto converts domain SearchFiltersResponse to proto SearchWithFiltersResponse
func SearchFiltersResponseToProto(resp *search.SearchFiltersResponse) *searchv1.SearchWithFiltersResponse {
	protoResp := &searchv1.SearchWithFiltersResponse{
		Total:  resp.Total,
		TookMs: resp.TookMs,
		Cached: resp.Cached,
	}

	// Convert listings (reuse existing converter logic)
	if len(resp.Listings) > 0 {
		protoResp.Listings = make([]*searchv1.Listing, 0, len(resp.Listings))
		for _, listing := range resp.Listings {
			protoResp.Listings = append(protoResp.Listings, ListingSearchResultToProto(&listing))
		}
	}

	// Include facets if present
	if resp.Facets != nil {
		protoResp.Facets = FacetsResponseToProto(resp.Facets)
	}

	return protoResp
}

// ProtoToSearchFilters converts proto Filters to domain SearchFilters
func ProtoToSearchFilters(filters *searchv1.Filters) *search.SearchFilters {
	domainFilters := &search.SearchFilters{}

	// Price range
	if filters.Price != nil {
		domainFilters.Price = &search.PriceRange{}
		if filters.Price.Min != nil {
			min := *filters.Price.Min
			domainFilters.Price.Min = &min
		}
		if filters.Price.Max != nil {
			max := *filters.Price.Max
			domainFilters.Price.Max = &max
		}
	}

	// Attributes
	if len(filters.Attributes) > 0 {
		domainFilters.Attributes = make(map[string][]string)
		for key, attrValues := range filters.Attributes {
			if len(attrValues.Values) > 0 {
				domainFilters.Attributes[key] = attrValues.Values
			}
		}
	}

	// Location
	if filters.Location != nil {
		domainFilters.Location = &search.LocationFilter{
			Lat:      filters.Location.Lat,
			Lon:      filters.Location.Lon,
			RadiusKm: filters.Location.RadiusKm,
		}
	}

	// Source type
	if filters.SourceType != nil {
		sourceType := *filters.SourceType
		domainFilters.SourceType = &sourceType
	}

	// Stock status
	if filters.StockStatus != nil {
		stockStatus := *filters.StockStatus
		domainFilters.StockStatus = &stockStatus
	}

	return domainFilters
}

// ProtoToSortConfig converts proto SortConfig to domain SortConfig
func ProtoToSortConfig(sort *searchv1.SortConfig) *search.SortConfig {
	domainSort := &search.SortConfig{
		Field: sort.Field,
		Order: sort.Order,
	}

	// Set defaults
	if domainSort.Field == "" {
		domainSort.Field = "relevance"
	}
	if domainSort.Order == "" {
		domainSort.Order = "desc"
	}

	return domainSort
}

// ProtoToSuggestionsRequest converts proto GetSuggestionsRequest to domain SuggestionsRequest
func ProtoToSuggestionsRequest(req *searchv1.GetSuggestionsRequest) *search.SuggestionsRequest {
	domainReq := &search.SuggestionsRequest{
		Prefix:   req.Prefix,
		Limit:    req.Limit,
		UseCache: true, // Default to using cache
	}

	// Optional category_id
	if req.CategoryId != nil {
		categoryID := *req.CategoryId
		domainReq.CategoryID = &categoryID
	}

	// Set defaults
	if domainReq.Limit == 0 {
		domainReq.Limit = 10
	}

	return domainReq
}

// SuggestionsResponseToProto converts domain SuggestionsResponse to proto GetSuggestionsResponse
func SuggestionsResponseToProto(resp *search.SuggestionsResponse) *searchv1.GetSuggestionsResponse {
	protoResp := &searchv1.GetSuggestionsResponse{
		TookMs: resp.TookMs,
		Cached: resp.Cached,
	}

	// Convert suggestions
	if len(resp.Suggestions) > 0 {
		protoResp.Suggestions = make([]*searchv1.Suggestion, 0, len(resp.Suggestions))
		for _, sug := range resp.Suggestions {
			protoSug := &searchv1.Suggestion{
				Text:  sug.Text,
				Score: sug.Score,
			}

			// Optional listing_id
			if sug.ListingID != nil {
				protoSug.ListingId = sug.ListingID
			}

			protoResp.Suggestions = append(protoResp.Suggestions, protoSug)
		}
	}

	return protoResp
}

// ProtoToPopularSearchesRequest converts proto GetPopularSearchesRequest to domain PopularSearchesRequest
func ProtoToPopularSearchesRequest(req *searchv1.GetPopularSearchesRequest) *search.PopularSearchesRequest {
	domainReq := &search.PopularSearchesRequest{
		Limit: req.Limit,
	}

	// Optional category_id
	if req.CategoryId != nil {
		categoryID := *req.CategoryId
		domainReq.CategoryID = &categoryID
	}

	// Optional time_range
	if req.TimeRange != nil {
		domainReq.TimeRange = *req.TimeRange
	} else {
		domainReq.TimeRange = "24h" // Default
	}

	// Set defaults
	if domainReq.Limit == 0 {
		domainReq.Limit = 10
	}

	return domainReq
}

// PopularSearchesResponseToProto converts domain PopularSearchesResponse to proto GetPopularSearchesResponse
func PopularSearchesResponseToProto(resp *search.PopularSearchesResponse) *searchv1.GetPopularSearchesResponse {
	protoResp := &searchv1.GetPopularSearchesResponse{
		TookMs: resp.TookMs,
	}

	// Convert searches
	if len(resp.Searches) > 0 {
		protoResp.Searches = make([]*searchv1.PopularSearch, 0, len(resp.Searches))
		for _, ps := range resp.Searches {
			protoResp.Searches = append(protoResp.Searches, &searchv1.PopularSearch{
				Query:       ps.Query,
				SearchCount: ps.SearchCount,
				TrendScore:  ps.TrendScore,
			})
		}
	}

	return protoResp
}

// ListingSearchResultToProto converts domain ListingSearchResult to proto Listing
// This is a helper to avoid duplication with existing converter in handlers_search.go
func ListingSearchResultToProto(listing *search.ListingSearchResult) *searchv1.Listing {
	protoListing := &searchv1.Listing{
		Id:          listing.ID,
		Uuid:        listing.UUID,
		Title:       listing.Title,
		Price:       listing.Price,
		Currency:    listing.Currency,
		CategoryId:  listing.CategoryID,
		Status:      listing.Status,
		CreatedAt:   listing.CreatedAt,
		UserId:      listing.UserID,
		Quantity:    listing.Quantity,
		SourceType:  listing.SourceType,
		StockStatus: listing.StockStatus,
	}

	// Add optional fields
	if listing.Description != nil {
		protoListing.Description = listing.Description
	}
	if listing.StorefrontID != nil {
		protoListing.StorefrontId = listing.StorefrontID
	}
	if listing.SKU != nil {
		protoListing.Sku = listing.SKU
	}

	// Add images
	if len(listing.Images) > 0 {
		protoImages := make([]*searchv1.ListingImage, 0, len(listing.Images))
		for _, img := range listing.Images {
			protoImages = append(protoImages, &searchv1.ListingImage{
				Id:           img.ID,
				Url:          img.URL,
				IsPrimary:    img.IsPrimary,
				DisplayOrder: img.DisplayOrder,
			})
		}
		protoListing.Images = protoImages
	}

	return protoListing
}
