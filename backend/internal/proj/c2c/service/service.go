// backend/internal/proj/c2c/service/service.go
package service

import (
	"context"
	"time"

	"backend/internal/config"
	"backend/internal/proj/c2c/repository"
	"backend/internal/proj/notifications/service"
	"backend/internal/storage"
)

type Service struct {
	Marketplace     MarketplaceServiceInterface
	Chat            ChatServiceInterface
	ChatAttachment  ChatAttachmentServiceInterface
	ChatTranslation *ChatTranslationService
	Order           *OrderService
	UnifiedCar      *UnifiedCarService
}

func NewService(storage storage.Storage, notifService service.NotificationServiceInterface, searchWeights *config.SearchWeights, cache CacheInterface) *Service {
	// Create a minimal translation service for internal use
	// Note: The actual translation service will be injected later by the global service
	// This is a temporary service to satisfy the interface requirement
	dummyTranslation := &dummyTranslationService{}

	// Создаем OrderService если есть доступ к MarketplaceOrderRepository
	var orderService *OrderService
	if marketplaceOrderRepo := storage.MarketplaceOrder(); marketplaceOrderRepo != nil {
		// Создаем адаптеры для репозиториев
		orderRepoAdapter := repository.NewPostgresOrderAdapter(marketplaceOrderRepo)
		listingRepoAdapter := repository.NewPostgresMarketplaceAdapter(storage) // storage сам реализует GetListingByID

		userRepoAdapter := &SimpleUserRepository{storage: storage}

		orderService = NewOrderService(
			orderRepoAdapter,         // orderRepo
			listingRepoAdapter,       // listingRepo
			userRepoAdapter,          // userRepo
			NewPaymentService(),      // paymentService
			NewNotificationAdapter(), // notificationService
			5.0,                      // platformFeeRate (5%)
		)
	}

	// Создаем UnifiedCarService с включенным VIN декодером
	carServiceConfig := &CarServiceConfig{
		VINDecoderEnabled: true,
		CacheEnabled:      true,
		CacheTTL:          24 * time.Duration(time.Hour),
	}
	unifiedCarService := NewUnifiedCarService(storage, nil, carServiceConfig)

	return &Service{
		Marketplace:     NewMarketplaceService(storage, dummyTranslation, searchWeights, cache),
		Chat:            NewChatService(storage, notifService),
		ChatAttachment:  nil, // Will be set by global service
		ChatTranslation: nil, // Will be set by global service
		Order:           orderService,
		UnifiedCar:      unifiedCarService,
	}
}

// SetChatAttachmentService sets the chat attachment service
// This is called by the global service after all dependencies are initialized
func (s *Service) SetChatAttachmentService(attachmentService ChatAttachmentServiceInterface) {
	s.ChatAttachment = attachmentService
}

// SetChatTranslationService sets the chat translation service
// This is called by the global service after all dependencies are initialized
func (s *Service) SetChatTranslationService(translationService *ChatTranslationService) {
	s.ChatTranslation = translationService
}

// dummyTranslationService is a minimal implementation that does nothing
// It serves as a placeholder until the real translation service is injected
type dummyTranslationService struct{}

func (d *dummyTranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	return languageAuto, 1.0, nil
}

func (d *dummyTranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	return map[string]string{languageAuto: text}, nil
}

func (d *dummyTranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	result[sourceLanguage] = fields
	return result, nil
}

func (d *dummyTranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) TranslateWithToneModeration(ctx context.Context, text string, sourceLanguage string, targetLanguage string, moderateTone bool) (string, error) {
	return text, nil
}
