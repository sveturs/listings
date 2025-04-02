// backend/internal/proj/marketplace/handler/handler.go
package handler
import (
    globalService "backend/internal/proj/global/service"
)

type Handler struct {
    Marketplace *MarketplaceHandler
    Chat       *ChatHandler
    Translation *TranslationHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
    return &Handler{
        Marketplace: NewMarketplaceHandler(services),
        Chat:       NewChatHandler(services),
        Translation: NewTranslationHandler(services),
    }
}