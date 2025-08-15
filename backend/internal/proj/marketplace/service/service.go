// backend/internal/proj/marketplace/service/service.go
package service

import (
	"context"

	"backend/internal/config"
	"backend/internal/proj/marketplace/repository"
	"backend/internal/proj/notifications/service"
	"backend/internal/storage"
)

type Service struct {
	Marketplace    MarketplaceServiceInterface
	Chat           ChatServiceInterface
	ChatAttachment ChatAttachmentServiceInterface
	Order          *OrderService
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

	return &Service{
		Marketplace:    NewMarketplaceService(storage, dummyTranslation, searchWeights, cache),
		Chat:           NewChatService(storage, notifService),
		ChatAttachment: nil, // Will be set by global service
		Order:          orderService,
	}
}

// SetChatAttachmentService sets the chat attachment service
// This is called by the global service after all dependencies are initialized
func (s *Service) SetChatAttachmentService(attachmentService ChatAttachmentServiceInterface) {
	s.ChatAttachment = attachmentService
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
