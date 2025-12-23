package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	categoriesv2 "github.com/vondi-global/listings/api/proto/categories/v2"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service"
)

// CategoryDetectionHandler - gRPC handler для детекции категорий
type CategoryDetectionHandler struct {
	categoriesv2.UnimplementedCategoryDetectionServiceServer
	service *service.CategoryDetectionService
	logger  zerolog.Logger
}

// NewCategoryDetectionHandler создаёт новый handler
func NewCategoryDetectionHandler(
	svc *service.CategoryDetectionService,
	logger zerolog.Logger,
) *CategoryDetectionHandler {
	return &CategoryDetectionHandler{
		service: svc,
		logger:  logger.With().Str("handler", "category_detection").Logger(),
	}
}

// DetectFromText определяет категорию по тексту
func (h *CategoryDetectionHandler) DetectFromText(
	ctx context.Context,
	req *categoriesv2.DetectFromTextRequest,
) (*categoriesv2.DetectCategoryResponse, error) {
	h.logger.Debug().
		Str("title", req.Title).
		Str("language", req.Language).
		Msg("DetectFromText called")

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	input := domain.DetectFromTextInput{
		Title:             req.Title,
		Description:       req.Description,
		Language:          req.Language,
		Hints:             protoHintsToDomain(req.Hints),
		SuggestedCategory: req.SuggestedCategory,
	}

	if input.Language == "" {
		input.Language = "sr" // default
	}

	detection, err := h.service.DetectFromText(ctx, input)
	if err != nil {
		h.logger.Error().Err(err).Msg("DetectFromText failed")
		return nil, status.Errorf(codes.Internal, "detection failed: %v", err)
	}

	return detectionToProto(detection), nil
}

// DetectFromKeywords определяет категорию по ключевым словам
func (h *CategoryDetectionHandler) DetectFromKeywords(
	ctx context.Context,
	req *categoriesv2.DetectFromKeywordsRequest,
) (*categoriesv2.DetectCategoryResponse, error) {
	h.logger.Debug().
		Strs("keywords", req.Keywords).
		Str("language", req.Language).
		Msg("DetectFromKeywords called")

	if len(req.Keywords) == 0 {
		return nil, status.Error(codes.InvalidArgument, "keywords are required")
	}

	language := req.Language
	if language == "" {
		language = "sr"
	}

	detection, err := h.service.DetectFromKeywords(ctx, req.Keywords, language)
	if err != nil {
		h.logger.Error().Err(err).Msg("DetectFromKeywords failed")
		return nil, status.Errorf(codes.Internal, "detection failed: %v", err)
	}

	return detectionToProto(detection), nil
}

// DetectBatch выполняет batch детекцию
func (h *CategoryDetectionHandler) DetectBatch(
	ctx context.Context,
	req *categoriesv2.DetectBatchRequest,
) (*categoriesv2.DetectBatchResponse, error) {
	h.logger.Debug().
		Int("items_count", len(req.Items)).
		Msg("DetectBatch called")

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items are required")
	}

	if len(req.Items) > 100 {
		return nil, status.Error(codes.InvalidArgument, "max 100 items per batch")
	}

	input := domain.DetectBatchInput{
		Items: make([]domain.DetectFromTextInput, len(req.Items)),
	}

	for i, item := range req.Items {
		language := item.Language
		if language == "" {
			language = "sr"
		}
		input.Items[i] = domain.DetectFromTextInput{
			Title:             item.Title,
			Description:       item.Description,
			Language:          language,
			Hints:             protoHintsToDomain(item.Hints),
			SuggestedCategory: item.SuggestedCategory,
		}
	}

	result, err := h.service.DetectBatch(ctx, input)
	if err != nil {
		h.logger.Error().Err(err).Msg("DetectBatch failed")
		return nil, status.Errorf(codes.Internal, "batch detection failed: %v", err)
	}

	response := &categoriesv2.DetectBatchResponse{
		Results:               make([]*categoriesv2.DetectCategoryResponse, len(result.Results)),
		TotalProcessingTimeMs: result.TotalProcessingTime,
	}

	for i, detection := range result.Results {
		response.Results[i] = detectionToProto(&detection)
	}

	return response, nil
}

// ConfirmSelection подтверждает выбор пользователя
func (h *CategoryDetectionHandler) ConfirmSelection(
	ctx context.Context,
	req *categoriesv2.ConfirmSelectionRequest,
) (*categoriesv2.ConfirmSelectionResponse, error) {
	h.logger.Debug().
		Str("detection_id", req.DetectionId).
		Str("selected_category_id", req.SelectedCategoryId).
		Msg("ConfirmSelection called")

	if req.DetectionId == "" || req.SelectedCategoryId == "" {
		return nil, status.Error(codes.InvalidArgument, "detection_id and selected_category_id are required")
	}

	detectionID, err := uuid.Parse(req.DetectionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid detection_id format")
	}

	selectedID, err := uuid.Parse(req.SelectedCategoryId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid selected_category_id format")
	}

	err = h.service.ConfirmSelection(ctx, detectionID, selectedID)
	if err != nil {
		h.logger.Error().Err(err).Msg("ConfirmSelection failed")
		return nil, status.Errorf(codes.Internal, "confirmation failed: %v", err)
	}

	return &categoriesv2.ConfirmSelectionResponse{Success: true}, nil
}

// protoHintsToDomain конвертирует proto hints в domain
func protoHintsToDomain(hints *categoriesv2.CategoryHints) *domain.CategoryHints {
	if hints == nil {
		return nil
	}
	return &domain.CategoryHints{
		Domain:      hints.Domain,
		ProductType: hints.ProductType,
		Keywords:    hints.Keywords,
	}
}

// detectionToProto конвертирует domain detection в proto
func detectionToProto(d *domain.CategoryDetection) *categoriesv2.DetectCategoryResponse {
	if d == nil {
		return &categoriesv2.DetectCategoryResponse{}
	}

	response := &categoriesv2.DetectCategoryResponse{
		DetectionId:      d.ID.String(),
		ProcessingTimeMs: d.ProcessingTimeMs,
	}

	if d.Primary != nil {
		response.Primary = categoryMatchToProto(d.Primary)
	}

	if len(d.Alternatives) > 0 {
		response.Alternatives = make([]*categoriesv2.CategoryMatch, len(d.Alternatives))
		for i, alt := range d.Alternatives {
			response.Alternatives[i] = categoryMatchToProto(&alt)
		}
	}

	return response
}

// categoryMatchToProto конвертирует domain match в proto
func categoryMatchToProto(m *domain.CategoryMatch) *categoriesv2.CategoryMatch {
	if m == nil {
		return nil
	}
	return &categoriesv2.CategoryMatch{
		CategoryId:      m.CategoryID.String(),
		CategoryName:    m.CategoryName,
		CategorySlug:    m.CategorySlug,
		CategoryPath:    m.CategoryPath,
		ConfidenceScore: m.ConfidenceScore,
		DetectionMethod: string(m.DetectionMethod),
		MatchedKeywords: m.MatchedKeywords,
	}
}
