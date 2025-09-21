package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	globalService "backend/internal/proj/global/service"
	marketplaceService "backend/internal/proj/marketplace/service"
	storefrontService "backend/internal/proj/storefronts/service"
	"backend/internal/proj/viber/config"
	"backend/internal/proj/viber/service"
)

// MessageHandler обрабатывает сообщения от пользователей Viber
type MessageHandler struct {
	botService         *service.BotService
	infobipService     *service.InfobipBotService
	services           globalService.ServicesInterface
	marketplaceService marketplaceService.MarketplaceServiceInterface
	storefrontService  storefrontService.StorefrontService
	useInfobip         bool
	config             *config.ViberConfig
}

// NewMessageHandler создаёт новый обработчик сообщений
func NewMessageHandler(
	botService *service.BotService,
	infobipService *service.InfobipBotService,
	services globalService.ServicesInterface,
	marketplaceService marketplaceService.MarketplaceServiceInterface,
	storefrontService storefrontService.StorefrontService,
	useInfobip bool,
	cfg *config.ViberConfig,
) *MessageHandler {
	return &MessageHandler{
		botService:         botService,
		infobipService:     infobipService,
		services:           services,
		marketplaceService: marketplaceService,
		storefrontService:  storefrontService,
		useInfobip:         useInfobip,
		config:             cfg,
	}
}

// HandleSearch обрабатывает поиск товаров
func (m *MessageHandler) HandleSearch(ctx context.Context, viberID, query string) error {
	// Очищаем запрос
	cleanQuery := strings.TrimSpace(query)
	if len(cleanQuery) < 2 {
		msg := "Пожалуйста, введите поисковый запрос (минимум 2 символа).\n\n" +
			"Например:\n• iPhone\n• велосипед\n• ноутбук\n• мебель"
		return m.sendMessage(ctx, viberID, msg)
	}

	// Ищем товары в маркетплейсе
	searchParams := &search.ServiceParams{
		Query:    cleanQuery,
		Page:     1,
		Size:     5,
		Language: "ru",
	}

	results, err := m.marketplaceService.SearchListingsAdvanced(ctx, searchParams)
	if err != nil {
		msg := "Произошла ошибка при поиске. Попробуйте позже или измените запрос."
		return m.sendMessage(ctx, viberID, msg)
	}

	if len(results.Items) == 0 {
		// Предлагаем альтернативы
		suggestions, _ := m.marketplaceService.GetSuggestions(ctx, cleanQuery, 3)
		msg := fmt.Sprintf("По запросу \"%s\" ничего не найдено.\n\n", cleanQuery)

		if len(suggestions) > 0 {
			msg += "Возможно, вы искали:\n"
			for _, suggestion := range suggestions {
				msg += fmt.Sprintf("• %s\n", suggestion)
			}
		} else {
			msg += "Попробуйте:\n• Изменить запрос\n• Использовать синонимы\n• Проверить орфографию"
		}

		return m.sendMessage(ctx, viberID, msg)
	}

	// Формируем ответ с результатами
	msg := fmt.Sprintf("Найдено %d товаров по запросу \"%s\":\n\n", len(results.Items), cleanQuery)

	for i, listing := range results.Items {
		price := "Цена не указана"
		if listing.Price > 0 {
			price = fmt.Sprintf("%.0f RSD", listing.Price)
		}

		msg += fmt.Sprintf("%d. %s\n💰 %s\n📍 %s\n🔗 https://svetu.rs/marketplace/listing/%d\n\n",
			i+1, listing.Title, price, listing.Location, listing.ID)
	}

	if results.Total > len(results.Items) {
		msg += fmt.Sprintf("Показаны первые %d из %d результатов.\n", len(results.Items), results.Total)
		msg += "Больше результатов на сайте: https://svetu.rs/marketplace/search?q=" + cleanQuery
	}

	return m.sendMessage(ctx, viberID, msg)
}

// HandleMyOrders обрабатывает запрос "Мои заказы"
func (m *MessageHandler) HandleMyOrders(ctx context.Context, viberID string) error {
	// TODO: Нужно связать viberID с userID в системе
	msg := "Функция \"Мои заказы\" пока недоступна в боте.\n\n" +
		"Для просмотра заказов перейдите на сайт:\n" +
		"🔗 https://svetu.rs/profile/orders\n\n" +
		"Скоро эта функция будет добавлена в бота!"

	return m.sendMessage(ctx, viberID, msg)
}

// HandleCart обрабатывает запрос "Корзина"
func (m *MessageHandler) HandleCart(ctx context.Context, viberID string) error {
	// TODO: Нужно связать viberID с userID в системе
	msg := "Функция \"Корзина\" пока недоступна в боте.\n\n" +
		"Для работы с корзиной перейдите на сайт:\n" +
		"🔗 https://svetu.rs/cart\n\n" +
		"Скоро эта функция будет добавлена в бота!"

	return m.sendMessage(ctx, viberID, msg)
}

// HandleTrackDelivery обрабатывает запрос на отслеживание доставки
func (m *MessageHandler) HandleTrackDelivery(ctx context.Context, viberID, trackingToken string) error {
	// TODO: Получить реальные данные из БД через сервис
	// Пока используем тестовые данные

	// Генерируем динамическую карту для Viber
	delivery := &service.DeliveryInfo{
		TrackingToken:     trackingToken,
		CourierLatitude:   44.95, // Между Белградом и Нови-Садом
		CourierLongitude:  20.10,
		DeliveryLatitude:  45.2671, // Нови-Сад
		DeliveryLongitude: 19.8335,
		EstimatedTime:     time.Now().Add(2 * time.Hour),
	}

	// Используем InfobipService для отправки Rich Media с картой
	if m.useInfobip {
		return m.infobipService.SendTrackingNotification(ctx, viberID, delivery)
	}

	// Для обычного бота отправляем ссылку
	msg := fmt.Sprintf("📦 Отслеживание посылки: %s\n\n", trackingToken)
	msg += "📍 Статус: В пути\n"
	msg += "🚴 Курьер движется к вам\n\n"
	msg += fmt.Sprintf("🔗 Отследить на карте:\nhttps://svetu.rs/track/%s\n\n", trackingToken)
	msg += "💡 Совет: Откройте ссылку для просмотра карты с текущим положением курьера в реальном времени!"

	return m.sendMessage(ctx, viberID, msg)
}

// HandleStorefronts обрабатывает запрос "Витрины"
func (m *MessageHandler) HandleStorefronts(ctx context.Context, viberID string) error {
	// Ищем популярные витрины
	isActive := true
	filter := &models.StorefrontFilter{
		IsActive:  &isActive,
		Limit:     5,
		Offset:    0,
		SortBy:    "created_at",
		SortOrder: "DESC",
	}

	storefronts, total, err := m.storefrontService.Search(ctx, filter)
	if err != nil {
		msg := "Произошла ошибка при загрузке витрин. Попробуйте позже."
		return m.sendMessage(ctx, viberID, msg)
	}

	if len(storefronts) == 0 {
		msg := "Пока нет активных витрин.\n\n" +
			"Создайте свою витрину на сайте:\n" +
			"🔗 https://svetu.rs/storefronts/create"
		return m.sendMessage(ctx, viberID, msg)
	}

	msg := fmt.Sprintf("Найдено витрин: %d\n\n", total)

	for i, storefront := range storefronts {
		msg += fmt.Sprintf("%d. %s\n", i+1, storefront.Name)
		if storefront.Description != nil && *storefront.Description != "" {
			// Обрезаем описание до 100 символов
			desc := *storefront.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			msg += fmt.Sprintf("📝 %s\n", desc)
		}
		msg += fmt.Sprintf("🔗 https://svetu.rs/storefront/%s\n\n", storefront.Slug)
	}

	if total > len(storefronts) {
		msg += fmt.Sprintf("Показаны первые %d из %d витрин.\n", len(storefronts), total)
		msg += "Все витрины: https://svetu.rs/storefronts"
	}

	return m.sendMessage(ctx, viberID, msg)
}

// HandleHelp обрабатывает запрос помощи
func (m *MessageHandler) HandleHelp(ctx context.Context, viberID string) error {
	msg := "🤖 Помощь SveTu Bot\n\n" +
		"Я умею:\n" +
		"🔍 Искать товары - просто напишите что ищете\n" +
		"🏪 Показывать витрины - \"витрины\" или \"магазины\"\n" +
		"📦 Отслеживать доставку - \"отследить\" + номер\n" +
		"📱 Находить товары рядом - отправьте геолокацию\n\n" +
		"Команды:\n" +
		"• поиск, найти + запрос\n" +
		"• витрины, магазины\n" +
		"• заказы\n" +
		"• корзина\n" +
		"• помощь\n\n" +
		"🌐 Сайт: https://svetu.rs\n" +
		"📞 Поддержка: support@svetu.rs"

	return m.sendMessage(ctx, viberID, msg)
}

// HandleNearbyProducts обрабатывает поиск товаров рядом
func (m *MessageHandler) HandleNearbyProducts(ctx context.Context, viberID string, lat, lng float64) error {
	// Ищем товары в радиусе 5 км
	markers, err := m.marketplaceService.GetListingsInBounds(ctx,
		lat+0.045, // ~5км к северу
		lng+0.063, // ~5км к востоку
		lat-0.045, // ~5км к югу
		lng-0.063, // ~5км к западу
		14,        // zoom level
		"",        // все категории
		"",        // все состояния
		nil, nil,  // без фильтра по цене
		"") // без фильтра атрибутов
	if err != nil {
		msg := "Произошла ошибка при поиске товаров рядом с вами. Попробуйте позже."
		return m.sendMessage(ctx, viberID, msg)
	}

	if len(markers) == 0 {
		msg := "Рядом с вами пока нет товаров.\n\n" +
			"Попробуйте:\n" +
			"🔍 Поиск по всему сайту\n" +
			"🏪 Просмотр витрин\n" +
			"📍 Изменить местоположение\n\n" +
			"🌐 Перейти на сайт: https://svetu.rs"
		return m.sendMessage(ctx, viberID, msg)
	}

	msg := fmt.Sprintf("📍 Найдено %d товаров рядом с вами:\n\n", len(markers))

	// Ограничиваем до 5 товаров
	limit := len(markers)
	if limit > 5 {
		limit = 5
	}

	for i := 0; i < limit; i++ {
		marker := markers[i]
		distance := m.calculateDistance(lat, lng, marker.Latitude, marker.Longitude)

		msg += fmt.Sprintf("%d. %s\n", i+1, marker.Title)
		if marker.Price > 0 {
			msg += fmt.Sprintf("💰 %.0f RSD\n", marker.Price)
		}
		msg += fmt.Sprintf("📍 ~%.1f км от вас\n", distance)
		if marker.City != "" {
			msg += fmt.Sprintf("🏠 %s\n", marker.City)
		}
		msg += fmt.Sprintf("🔗 https://svetu.rs/marketplace/listing/%d\n\n", marker.ID)
	}

	if len(markers) > 5 {
		msg += fmt.Sprintf("Показаны первые 5 из %d товаров.\n", len(markers))
		msg += "Больше на карте: https://svetu.rs/map"
	}

	return m.sendMessage(ctx, viberID, msg)
}

// calculateDistance вычисляет расстояние между двумя точками в километрах
func (m *MessageHandler) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Простая формула для приблизительного расчёта расстояния
	const earthRadius = 6371.0 // km

	dlat := (lat2 - lat1) * 3.14159265359 / 180.0
	dlng := (lng2 - lng1) * 3.14159265359 / 180.0

	a := dlat*dlat + dlng*dlng
	return earthRadius * 2 * 0.7071067811865476 * 1.4142135623730951 * a // упрощённая формула
}

// sendMessage отправляет сообщение через подходящий сервис
func (m *MessageHandler) sendMessage(ctx context.Context, viberID, message string) error {
	if m.useInfobip && m.infobipService != nil {
		return m.infobipService.SendTextMessage(ctx, viberID, message)
	} else if m.botService != nil {
		return m.botService.SendTextMessage(ctx, viberID, message)
	}

	return fmt.Errorf("no bot service available")
}
