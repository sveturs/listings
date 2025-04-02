// backend/internal/proj/marketplace/handler/translation.go
package handler

import (
    "backend/internal/proj/global/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    "backend/pkg/utils"
    "github.com/gofiber/fiber/v2"
    "log"
)

type TranslationHandler struct {
    services service.ServicesInterface
}

func NewTranslationHandler(services service.ServicesInterface) *TranslationHandler {
    return &TranslationHandler{
        services: services,
    }
}

// GetTranslationLimits возвращает информацию о лимитах перевода для разных провайдеров
func (h *TranslationHandler) GetTranslationLimits(c *fiber.Ctx) error {
    // Проверяем авторизацию
    _, ok := c.Locals("user_id").(int)
    if !ok {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
    }

    // Пытаемся получить доступ к фабрике перевода
    translationFactory, ok := h.services.Translation().(marketplaceService.TranslationFactoryInterface)
    if !ok {
        // Если интерфейс не реализован, возвращаем стандартные значения
        return utils.SuccessResponse(c, fiber.Map{
            "google": fiber.Map{
                "used": 0,
                "limit": 100,
            },
            "openai": fiber.Map{
                "used": 0,
                "limit": 0, // Нет жестких лимитов для OpenAI
            },
        })
    }

    // Получаем информацию о лимитах Google Translate
    googleUsed, googleLimit, err := translationFactory.GetTranslationCount(marketplaceService.GoogleTranslate)
    if err != nil {
        log.Printf("Ошибка получения лимитов Google Translate: %v", err)
        googleUsed = 0
        googleLimit = 100 // Значение по умолчанию
    }

    // Для OpenAI нет встроенных лимитов
    openaiUsed, openaiLimit := 0, 0

    // Получаем список доступных провайдеров
    availableProviders := translationFactory.GetAvailableProviders()
    providersStr := make([]string, len(availableProviders))
    for i, provider := range availableProviders {
        providersStr[i] = string(provider)
    }

    // Формируем ответ
    response := fiber.Map{
        "google": fiber.Map{
            "used": googleUsed,
            "limit": googleLimit,
            "available": googleLimit - googleUsed,
        },
        "openai": fiber.Map{
            "used": openaiUsed,
            "limit": openaiLimit,
        },
        "available_providers": providersStr,
        "default_provider": string(translationFactory.GetDefaultProvider()),
    }

    return utils.SuccessResponse(c, response)
}

// SetTranslationProvider устанавливает провайдер перевода по умолчанию
func (h *TranslationHandler) SetTranslationProvider(c *fiber.Ctx) error {
    // Проверяем авторизацию
    _, ok := c.Locals("user_id").(int)
    if !ok {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
    }

    var request struct {
        Provider string `json:"provider"`
    }

    if err := c.BodyParser(&request); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат запроса")
    }

    // Пытаемся получить доступ к фабрике перевода
    translationFactory, ok := h.services.Translation().(marketplaceService.TranslationFactoryInterface)
    if !ok {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Управление провайдером перевода недоступно")
    }

    // Преобразуем строку в enum
    var provider marketplaceService.TranslationProvider
    switch request.Provider {
    case "google":
        provider = marketplaceService.GoogleTranslate
    case "openai":
        provider = marketplaceService.OpenAI
    default:
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неизвестный провайдер перевода: "+request.Provider)
    }

    // Устанавливаем провайдер
    if err := translationFactory.SetDefaultProvider(provider); err != nil {
        log.Printf("Ошибка установки провайдера перевода: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка установки провайдера перевода: "+err.Error())
    }

    return utils.SuccessResponse(c, fiber.Map{
        "provider": string(provider),
        "message": "Провайдер перевода успешно установлен",
    })
}