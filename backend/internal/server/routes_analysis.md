# Анализ маршрутов Backend API

## Обзор

Данный документ содержит полный анализ всех URL маршрутов в `backend/internal/server/server.go` с группировкой по логическим блокам и предложениями по рефакторингу.

**Дата анализа:** 6 декабря 2025  
**Общее количество маршрутов:** ~150  
**Основная проблема:** Функция setupRoutes() содержит ~340 строк кода  

---

## 🗂️ Текущие URL маршруты по группам

### **1. Основные/Системные маршруты**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/` | - | inline | Главная страница API |
| GET | `/api/health` | - | inline | Health check |
| GET | `/swagger/*` | - | swagger.HandlerDefault | Swagger документация |
| GET | `/docs/*` | - | swagger.New | Документация API |

### **2. Статические файлы и медиа**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/listings/*` | - | inline | Редирект на MinIO для изображений |
| GET | `/uploads/*` | - | Static | Статические файлы загрузок |
| GET | `/public/*` | - | Static | Публичные статические файлы |
| GET | `/service-worker.js` | - | inline | Service Worker |

### **3. WebSocket**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/ws/chat` | AuthRequired | s.marketplace.Chat.HandleWebSocket | WebSocket для чата |

### **4. Публичные маршруты**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/reindex-ratings-public` | - | s.marketplace.Indexing.ReindexRatings | Публичная реиндексация рейтингов |
| POST | `/api/v1/public/reindex` | - | s.marketplace.Indexing.ReindexAll | Публичная реиндексация всех данных |
| POST | `/api/v1/public/send-email` | - | s.notifications.Notification.SendPublicEmail | Публичная отправка email |
| GET | `/api/v1/public/storefronts/:id` | - | s.storefront.Storefront.GetPublicStorefront | Публичная информация о витрине |
| GET | `/api/v1/public/storefronts/:id/reviews` | - | s.review.Review.GetStorefrontReviews | Отзывы о витрине |
| GET | `/api/v1/public/storefronts/:id/rating` | - | s.review.Review.GetStorefrontRatingSummary | Рейтинг витрины |
| GET | `/api/v1/admin-check/:email` | - | s.users.User.IsAdminPublic | Проверка статуса администратора |

### **5. Webhook'и**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/notifications/telegram/webhook` | - | s.notifications.Notification.HandleTelegramWebhook | Telegram webhook |
| POST | `/webhook/stripe` | - | inline | Stripe webhook для платежей |

### **6. Авторизация**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/users/register` | RegistrationRateLimit | s.users.User.Register | Регистрация пользователя |
| POST | `/api/v1/users/login` | AuthRateLimit | s.users.User.Login | Вход пользователя |
| GET | `/auth/session` | - | s.users.Auth.GetSession | Получение сессии |
| GET | `/auth/google` | RateLimitByIP(10, time.Minute) | s.users.Auth.GoogleAuth | Google авторизация |
| GET | `/auth/google/callback` | RateLimitByIP(10, time.Minute) | s.users.Auth.GoogleCallback | Google callback |
| GET | `/auth/logout` | - | s.users.Auth.Logout | Выход |
| GET | `/api/v1/csrf-token` | - | s.middleware.GetCSRFToken() | CSRF токен |

### **7. Marketplace (публичные)**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/marketplace/listings` | - | s.marketplace.Listings.GetListings | Список объявлений |
| GET | `/api/v1/marketplace/categories` | - | s.marketplace.Categories.GetCategories | Категории |
| GET | `/api/v1/marketplace/category-tree` | - | s.marketplace.Categories.GetCategoryTree | Дерево категорий |
| GET | `/api/v1/marketplace/listings/:id` | - | s.marketplace.Listings.GetListing | Конкретное объявление |
| GET | `/api/v1/marketplace/search` | - | s.marketplace.Search.SearchListingsAdvanced | Расширенный поиск |
| GET | `/api/v1/marketplace/suggestions` | - | s.marketplace.Search.GetSuggestions | Автодополнение |
| GET | `/api/v1/marketplace/category-suggestions` | - | s.marketplace.Search.GetCategorySuggestions | Подсказки категорий |
| GET | `/api/v1/marketplace/categories/:id/attributes` | - | s.marketplace.Categories.GetCategoryAttributes | Атрибуты категории |
| GET | `/api/v1/marketplace/listings/:id/price-history` | - | s.marketplace.Listings.GetPriceHistory | История цен |
| GET | `/api/v1/marketplace/listings/:id/similar` | - | s.marketplace.Search.GetSimilarListings | Похожие объявления |
| GET | `/api/v1/marketplace/categories/:id/attribute-ranges` | - | s.marketplace.Categories.GetAttributeRanges | Диапазоны атрибутов |
| GET | `/api/v1/marketplace/enhanced-suggestions` | - | s.marketplace.Search.GetEnhancedSuggestions | Улучшенные подсказки |
| GET | `/api/v1/marketplace/map/bounds` | - | s.marketplace.GetListingsInBounds | Объявления в границах |
| GET | `/api/v1/marketplace/map/clusters` | - | s.marketplace.GetMapClusters | Кластеры на карте |

### **8. Marketplace (защищенные)**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/marketplace/listings` | AuthRequired | s.marketplace.Listings.CreateListing | Создание объявления |
| PUT | `/api/v1/marketplace/listings/:id` | AuthRequired | s.marketplace.Listings.UpdateListing | Обновление объявления |
| DELETE | `/api/v1/marketplace/listings/:id` | AuthRequired | s.marketplace.Listings.DeleteListing | Удаление объявления |
| POST | `/api/v1/marketplace/listings/:id/images` | AuthRequired | s.marketplace.Images.UploadImages | Загрузка изображений |
| POST | `/api/v1/marketplace/listings/:id/favorite` | AuthRequired | s.marketplace.Favorites.AddToFavorites | Добавить в избранное |
| DELETE | `/api/v1/marketplace/listings/:id/favorite` | AuthRequired | s.marketplace.Favorites.RemoveFromFavorites | Удалить из избранного |
| GET | `/api/v1/marketplace/favorites` | AuthRequired | s.marketplace.Favorites.GetFavorites | Получить избранное |

### **9. Переводы**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/translation/limits` | - | s.marketplace.Translations.GetTranslationLimits | Лимиты переводов |
| POST | `/api/v1/translation/provider` | - | s.marketplace.Translations.SetTranslationProvider | Установка провайдера |
| PUT | `/api/v1/marketplace/translations/:id` | AuthRequired | s.marketplace.Translations.UpdateTranslations | Обновление переводов |
| POST | `/api/v1/marketplace/translations/batch` | AuthRequired | s.marketplace.Translations.TranslateText | Пакетный перевод (дубликат) |
| POST | `/api/v1/marketplace/translations/batch-translate` | AuthRequired | s.marketplace.Translations.BatchTranslateListings | Пакетный перевод объявлений |
| POST | `/api/v1/marketplace/translations/translate` | AuthRequired | s.marketplace.Translations.TranslateText | Перевод текста |
| POST | `/api/v1/marketplace/translations/detect-language` | AuthRequired | s.marketplace.Translations.DetectLanguage | Определение языка |
| GET | `/api/v1/marketplace/translations/:id` | AuthRequired | s.marketplace.Translations.GetTranslations | Получение переводов |

### **10. Обработка изображений**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/marketplace/moderate-image` | AuthRequired | s.marketplace.Images.ModerateImage | Модерация изображений |
| POST | `/api/v1/marketplace/enhance-preview` | AuthRequired | s.marketplace.Images.EnhancePreview | Улучшение превью |
| POST | `/api/v1/marketplace/enhance-images` | AuthRequired | s.marketplace.Images.EnhanceImages | Улучшение изображений |

### **11. Отзывы (публичные)**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/reviews` | - | s.review.Review.GetReviews | Список отзывов |
| GET | `/api/v1/reviews/:id` | - | s.review.Review.GetReviewByID | Конкретный отзыв |
| GET | `/api/v1/reviews/stats` | - | s.review.Review.GetStats | Статистика отзывов |

### **12. Отзывы (защищенные)**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/reviews` | AuthRequired | s.review.Review.CreateReview | Создание отзыва |
| PUT | `/api/v1/reviews/:id` | AuthRequired | s.review.Review.UpdateReview | Обновление отзыва |
| DELETE | `/api/v1/reviews/:id` | AuthRequired | s.review.Review.DeleteReview | Удаление отзыва |
| POST | `/api/v1/reviews/:id/vote` | AuthRequired | s.review.Review.VoteForReview | Голосование за отзыв |
| POST | `/api/v1/reviews/:id/response` | AuthRequired | s.review.Review.AddResponse | Ответ на отзыв |
| POST | `/api/v1/reviews/:id/photos` | AuthRequired | s.review.Review.UploadPhotos | Загрузка фото к отзыву |
| GET | `/api/v1/users/:id/reviews` | AuthRequired | s.review.Review.GetUserReviews | Отзывы пользователя |
| GET | `/api/v1/users/:id/rating` | AuthRequired | s.review.Review.GetUserRatingSummary | Рейтинг пользователя |

### **13. Рейтинги сущностей**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/entity/:type/:id/rating` | - | s.review.Review.GetEntityRating | Рейтинг сущности |
| GET | `/api/v1/entity/:type/:id/stats` | - | s.review.Review.GetEntityStats | Статистика сущности |

### **14. Витрины (защищенные)**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/storefronts` | AuthRequired | s.storefront.Storefront.GetUserStorefronts | Витрины пользователя |
| POST | `/api/v1/storefronts` | AuthRequired | s.storefront.Storefront.CreateStorefront | Создание витрины |
| GET | `/api/v1/storefronts/:id` | AuthRequired | s.storefront.Storefront.GetStorefront | Получение витрины |
| PUT | `/api/v1/storefronts/:id` | AuthRequired | s.storefront.Storefront.UpdateStorefront | Обновление витрины |
| DELETE | `/api/v1/storefronts/:id` | AuthRequired | s.storefront.Storefront.DeleteStorefront | Удаление витрины |

### **15. Импорт данных**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/storefronts/:id/import-sources` | AuthRequired | s.storefront.Storefront.GetImportSources | Источники импорта |
| POST | `/api/v1/storefronts/import-sources` | AuthRequired | s.storefront.Storefront.CreateImportSource | Создание источника |
| PUT | `/api/v1/storefronts/import-sources/:id` | AuthRequired | s.storefront.Storefront.UpdateImportSource | Обновление источника |
| DELETE | `/api/v1/storefronts/import-sources/:id` | AuthRequired | s.storefront.Storefront.DeleteImportSource | Удаление источника |
| POST | `/api/v1/storefronts/import-sources/:id/run` | AuthRequired | s.storefront.Storefront.RunImport | Запуск импорта |
| GET | `/api/v1/storefronts/import-sources/:id/history` | AuthRequired | s.storefront.Storefront.GetImportHistory | История импорта |
| GET | `/api/v1/storefronts/import-sources/:id/category-mappings` | AuthRequired | s.storefront.Storefront.GetCategoryMappings | Маппинг категорий |
| PUT | `/api/v1/storefronts/import-sources/:id/category-mappings` | AuthRequired | s.storefront.Storefront.UpdateCategoryMappings | Обновление маппинга |
| GET | `/api/v1/storefronts/import-sources/:id/imported-categories` | AuthRequired | s.storefront.Storefront.GetImportedCategories | Импортированные категории |
| POST | `/api/v1/storefronts/import-sources/:id/apply-category-mappings` | AuthRequired | s.storefront.Storefront.ApplyCategoryMappings | Применение маппинга |

### **16. Баланс и платежи**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/balance` | AuthRequired | s.balance.Balance.GetBalance | Получение баланса |
| GET | `/api/v1/balance/transactions` | AuthRequired | s.balance.Balance.GetTransactions | Транзакции |
| GET | `/api/v1/balance/payment-methods` | AuthRequired | s.balance.Balance.GetPaymentMethods | Способы оплаты |
| POST | `/api/v1/balance/deposit` | AuthRequired | s.balance.Balance.CreateDeposit | Пополнение баланса |

### **17. Пользователи (защищенные)**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/users/me` | AuthRequired | s.users.User.GetProfile | Профиль пользователя (deprecated) |
| PUT | `/api/v1/users/me` | AuthRequired | s.users.User.UpdateProfile | Обновление профиля (deprecated) |
| GET | `/api/v1/users/profile` | AuthRequired | s.users.User.GetProfile | Профиль пользователя |
| PUT | `/api/v1/users/profile` | AuthRequired | s.users.User.UpdateProfile | Обновление профиля |
| GET | `/api/v1/users/:id/profile` | AuthRequired | s.users.User.GetProfileByID | Профиль по ID |

### **18. Геокодирование**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/geocode/reverse` | - | s.geocode.ReverseGeocode | Обратное геокодирование |
| GET | `/api/v1/cities/suggest` | - | s.geocode.GetCitySuggestions | Подсказки городов |

### **19. Чат**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/marketplace/chat` | AuthRequired | s.marketplace.Chat.GetChats | Список чатов |
| GET | `/api/v1/marketplace/chat/messages` | AuthRequired | s.marketplace.Chat.GetMessages | Сообщения |
| POST | `/api/v1/marketplace/chat/messages` | AuthRequired, RateLimitMessages | s.marketplace.Chat.SendMessage | Отправка сообщения |
| PUT | `/api/v1/marketplace/chat/messages/read` | AuthRequired | s.marketplace.Chat.MarkAsRead | Отметить как прочитанное |
| POST | `/api/v1/marketplace/chat/:chat_id/archive` | AuthRequired | s.marketplace.Chat.ArchiveChat | Архивация чата |
| GET | `/api/v1/marketplace/chat/unread-count` | AuthRequired | s.marketplace.Chat.GetUnreadCount | Количество непрочитанных |
| POST | `/api/v1/marketplace/chat/messages/:id/attachments` | AuthRequired, RateLimitMessages | s.marketplace.Chat.UploadAttachments | Загрузка вложений |
| GET | `/api/v1/marketplace/chat/attachments/:id` | AuthRequired | s.marketplace.Chat.GetAttachment | Получение вложения |
| DELETE | `/api/v1/marketplace/chat/attachments/:id` | AuthRequired | s.marketplace.Chat.DeleteAttachment | Удаление вложения |

### **20. Контакты**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/contacts` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.GetContacts | Список контактов |
| POST | `/api/v1/contacts` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.AddContact | Добавление контакта |
| PUT | `/api/v1/contacts/:contact_user_id` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.UpdateContactStatus | Обновление статуса контакта |
| DELETE | `/api/v1/contacts/:contact_user_id` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.RemoveContact | Удаление контакта |
| GET | `/api/v1/contacts/privacy` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.GetPrivacySettings | Настройки приватности |
| PUT | `/api/v1/contacts/privacy` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.UpdatePrivacySettings | Обновление настроек приватности |
| GET | `/api/v1/contacts/status/:contact_user_id` | AuthRequired, RateLimitByUser(300, time.Minute) | s.contacts.GetContactStatus | Статус контакта |

### **21. Уведомления**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/v1/notifications/telegram` | - | s.notifications.Notification.GetTelegramStatus | Статус Telegram |
| GET | `/api/v1/notifications` | AuthRequired | s.notifications.Notification.GetNotifications | Список уведомлений |
| GET | `/api/v1/notifications/settings` | AuthRequired | s.notifications.Notification.GetSettings | Настройки уведомлений |
| PUT | `/api/v1/notifications/settings` | AuthRequired | s.notifications.Notification.UpdateSettings | Обновление настроек |
| GET | `/api/v1/notifications/telegram` | AuthRequired | s.notifications.Notification.GetTelegramStatus | Статус Telegram |
| GET | `/api/v1/notifications/telegram/token` | AuthRequired | s.notifications.Notification.GetTelegramToken | Telegram токен |
| PUT | `/api/v1/notifications/:id/read` | AuthRequired | s.notifications.Notification.MarkAsRead | Отметить как прочитанное |
| POST | `/api/v1/notifications/telegram/token` | AuthRequired | s.notifications.Notification.GetTelegramToken | Получить токен |

### **22. Администрирование - Категории**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/admin/categories` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.CreateCategory | Создание категории |
| GET | `/api/v1/admin/categories` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.GetCategories | Список категорий |
| GET | `/api/v1/admin/categories/:id` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.GetCategoryByID | Категория по ID |
| PUT | `/api/v1/admin/categories/:id` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.UpdateCategory | Обновление категории |
| DELETE | `/api/v1/admin/categories/:id` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.DeleteCategory | Удаление категории |
| POST | `/api/v1/admin/categories/:id/reorder` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.ReorderCategories | Сортировка категорий |
| PUT | `/api/v1/admin/categories/:id/move` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.MoveCategory | Перемещение категории |

### **23. Администрирование - Атрибуты категорий**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/admin/categories/:id/attributes` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.AddAttributeToCategory | Добавить атрибут |
| DELETE | `/api/v1/admin/categories/:id/attributes/:attr_id` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.RemoveAttributeFromCategory | Удалить атрибут |
| PUT | `/api/v1/admin/categories/:id/attributes/:attr_id` | AuthRequired, AdminRequired | s.marketplace.AdminCategories.UpdateAttributeCategory | Обновить атрибут |

### **24. Администрирование - Атрибуты**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/admin/attributes` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.CreateAttribute | Создание атрибута |
| GET | `/api/v1/admin/attributes` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.GetAttributes | Список атрибутов |
| GET | `/api/v1/admin/attributes/:id` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.GetAttributeByID | Атрибут по ID |
| PUT | `/api/v1/admin/attributes/:id` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.UpdateAttribute | Обновление атрибута |
| DELETE | `/api/v1/admin/attributes/:id` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.DeleteAttribute | Удаление атрибута |
| POST | `/api/v1/admin/attributes/bulk-update` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.BulkUpdateAttributes | Массовое обновление |

### **25. Администрирование - Экспорт/импорт атрибутов**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/admin/categories/:categoryId/attributes/export` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.ExportCategoryAttributes | Экспорт атрибутов |
| POST | `/api/v1/admin/categories/:categoryId/attributes/import` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.ImportCategoryAttributes | Импорт атрибутов |
| POST | `/api/v1/admin/categories/:targetCategoryId/attributes/copy` | AuthRequired, AdminRequired | s.marketplace.AdminAttributes.CopyAttributesSettings | Копирование настроек |

### **26. Администрирование - Кастомные компоненты**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/admin/custom-components/templates` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.ListTemplates | Шаблоны компонентов |
| POST | `/api/v1/admin/custom-components/templates` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.CreateTemplate | Создание шаблона |
| GET | `/api/v1/admin/custom-components/usage` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.GetComponentUsages | Использование компонентов |
| POST | `/api/v1/admin/custom-components/usage` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.AddComponentUsage | Добавить использование |
| DELETE | `/api/v1/admin/custom-components/usage/:id` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.RemoveComponentUsage | Удалить использование |
| POST | `/api/v1/admin/custom-components` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.CreateComponent | Создание компонента |
| GET | `/api/v1/admin/custom-components` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.ListComponents | Список компонентов |
| GET | `/api/v1/admin/custom-components/:id` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.GetComponent | Компонент по ID |
| PUT | `/api/v1/admin/custom-components/:id` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.UpdateComponent | Обновление компонента |
| DELETE | `/api/v1/admin/custom-components/:id` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.DeleteComponent | Удаление компонента |
| GET | `/api/v1/admin/categories/:category_id/components` | AuthRequired, AdminRequired | s.marketplace.CustomComponents.GetCategoryComponents | Компоненты категории |

### **27. Администрирование - Группы атрибутов**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/admin/attribute-groups` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.ListAttributeGroups | Список групп |
| POST | `/api/v1/admin/attribute-groups` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.CreateAttributeGroup | Создание группы |
| GET | `/api/v1/admin/attribute-groups/:id` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.GetAttributeGroup | Группа по ID |
| PUT | `/api/v1/admin/attribute-groups/:id` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.UpdateAttributeGroup | Обновление группы |
| DELETE | `/api/v1/admin/attribute-groups/:id` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.DeleteAttributeGroup | Удаление группы |
| GET | `/api/v1/admin/attribute-groups/:id/items` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.GetAttributeGroupWithItems | Элементы группы |
| POST | `/api/v1/admin/attribute-groups/:id/items` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.AddItemToGroup | Добавить в группу |
| DELETE | `/api/v1/admin/attribute-groups/:id/items/:attributeId` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.RemoveItemFromGroup | Удалить из группы |
| GET | `/api/v1/admin/categories/:id/attribute-groups` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.GetCategoryGroups | Группы категории |
| POST | `/api/v1/admin/categories/:id/attribute-groups` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.AttachGroupToCategory | Привязать группу |
| DELETE | `/api/v1/admin/categories/:id/attribute-groups/:groupId` | AuthRequired, AdminRequired | s.marketplace.MarketplaceHandler.DetachGroupFromCategory | Отвязать группу |

### **28. Администрирование - Управление пользователями**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/admin/users` | AuthRequired, AdminRequired | s.users.User.GetAllUsers | Все пользователи |
| GET | `/api/v1/admin/users/:id` | AuthRequired, AdminRequired | s.users.User.GetUserByIDAdmin | Пользователь по ID |
| PUT | `/api/v1/admin/users/:id` | AuthRequired, AdminRequired | s.users.User.UpdateUserAdmin | Обновление пользователя |
| PUT | `/api/v1/admin/users/:id/status` | AuthRequired, AdminRequired | s.users.User.UpdateUserStatus | Обновление статуса |
| DELETE | `/api/v1/admin/users/:id` | AuthRequired, AdminRequired | s.users.User.DeleteUser | Удаление пользователя |
| GET | `/api/v1/admin/users/:id/balance` | AuthRequired, AdminRequired | s.users.User.GetUserBalance | Баланс пользователя |
| GET | `/api/v1/admin/users/:id/transactions` | AuthRequired, AdminRequired | s.users.User.GetUserTransactions | Транзакции пользователя |

### **29. Администрирование - Управление администраторами**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/api/v1/admin/admins` | AuthRequired, AdminRequired | s.users.User.GetAllAdmins | Все администраторы |
| POST | `/api/v1/admin/admins` | AuthRequired, AdminRequired | s.users.User.AddAdmin | Добавить администратора |
| DELETE | `/api/v1/admin/admins/:email` | AuthRequired, AdminRequired | s.users.User.RemoveAdmin | Удалить администратора |
| GET | `/api/v1/admin/admins/check/:email` | AuthRequired, AdminRequired | s.users.User.IsAdmin | Проверка администратора |

### **30. Администрирование - Системные операции**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| POST | `/api/v1/admin/reindex-listings` | AuthRequired, AdminRequired | s.marketplace.Indexing.ReindexAll | Реиндексация объявлений |
| POST | `/api/v1/admin/reindex-listings-with-translations` | AuthRequired, AdminRequired | s.marketplace.Indexing.ReindexAllWithTranslations | Реиндексация с переводами |
| POST | `/api/v1/admin/sync-discounts` | AuthRequired, AdminRequired | s.marketplace.Listings.SynchronizeDiscounts | Синхронизация скидок |
| POST | `/api/v1/admin/reindex-ratings` | AuthRequired, AdminRequired | s.marketplace.Indexing.ReindexRatings | Реиндексация рейтингов |

### **31. Legacy административные маршруты (/api/admin)**

**Примечание:** Полная копия всех современных админских маршрутов под префиксом `/api/admin` для обратной совместимости. Включает все маршруты из разделов 22-30, но с префиксом `/api/admin` вместо `/api/v1/admin`.

### **32. Временные тестовые маршруты**

| HTTP метод | URL путь | Middleware | Handler | Описание |
|------------|----------|------------|---------|-----------|
| GET | `/admin-categories-test` | - | s.marketplace.AdminCategories.GetCategories | Тестовый маршрут категорий |

---

## 🚨 Выявленные проблемы

### **1. Архитектурные проблемы**
- **Огромная функция setupRoutes()** - ~340 строк кода
- **Нарушение принципа единственной ответственности**
- **Сложность в поддержке и тестировании**

### **2. Дублирование кода**
- **Полное дублирование admin routes** (v1 и legacy)
- **Дублирующие handler'ы** для переводов
- **Повторяющиеся middleware паттерны**

### **3. Проблемы безопасности**
- **WebSocket без аутентификации** ❌ (строка 197-202)
- **Inconsistent rate limiting** применение
- **Отсутствует contacts handler** в структуре Server

### **4. Конфигурационные проблемы**
- **Хардкод MinIO URL** (`http://localhost:9000`)
- **Отсутствие конфигурируемых параметров**

### **5. API Design проблемы**
- **Inconsistent URL patterns**
- **Смешение публичных и защищенных маршрутов в одних группах**
- **Deprecated endpoints** не помечены для удаления

---

## 🎯 Предлагаемый план рефакторинга

### **Фаза 1: Логическая группировка маршрутов**

```
📂 ГРУППА 1: Core & Infrastructure (10 маршрутов)
   ├── Health & Status checks
   ├── Static files & documentation  
   ├── CSRF tokens
   └── Basic utilities

📂 ГРУППА 2: Authentication & Security (7 маршрутов)
   ├── Login/Register + Rate limiting
   ├── Google OAuth flow
   ├── Session management
   └── Public security checks

📂 ГРУППА 3: Public Marketplace (15 маршрутов)
   ├── Browse listings & categories
   ├── Search & suggestions
   ├── Maps & geo features
   └── Public storefront info

📂 ГРУППА 4: Protected User Operations (25 маршрутов)
   ├── User profiles & reviews
   ├── Protected marketplace CRUD
   ├── Favorites & image operations
   └── User-specific data

📂 ГРУППА 5: Business Operations (20 маршрутов)
   ├── Storefronts management
   ├── Balance & payments
   ├── Import/export workflows
   └── Translation services

📂 ГРУППА 6: Communication (15 маршрутов)
   ├── Chat & WebSocket
   ├── Notifications system
   ├── Contacts management
   └── External webhooks

📂 ГРУППА 7: Administration (40+ маршрутов)
   ├── Categories & Attributes CRUD
   ├── User & admin management
   ├── System operations
   └── Custom components

📂 ГРУППА 8: Geocoding & Utilities (5 маршрутов)
   ├── Reverse geocoding
   ├── City suggestions
   └── Location utilities
```

### **Фаза 2: Создание отдельных функций**

1. **setupCoreRoutes()** - базовые системные маршруты
2. **setupAuthenticationRoutes()** - авторизация и безопасность
3. **setupPublicMarketplaceRoutes()** - публичный API маркетплейса
4. **setupProtectedUserRoutes()** - защищенные пользовательские операции
5. **setupBusinessRoutes()** - бизнес-логика и операции
6. **setupCommunicationRoutes()** - чаты, уведомления, контакты
7. **setupAdminRoutes()** - администрирование
8. **setupUtilityRoutes()** - геокодирование и утилиты

### **Фаза 3: Исправление проблем безопасности**

1. **Добавить аутентификацию для WebSocket**
2. **Унифицировать rate limiting применение**
3. **Добавить недостающий contacts handler**
4. **Конфигурируемые URLs вместо хардкода**

### **Фаза 4: Убрать дублирование**

1. **Консолидировать admin routes** (убрать legacy через deprecation)
2. **Унифицировать translation endpoints**
3. **Создать reusable middleware chains**

### **Фаза 5: Улучшение API Design**

1. **Consistent URL naming conventions**
2. **Proper HTTP status codes**
3. **OpenAPI documentation**
4. **Deprecation headers для старых endpoints**

---

## 📋 Рекомендации по внедрению

### **Приоритет 1 (Критично)**
- ✅ Разбить setupRoutes() на логические функции
- ✅ Исправить WebSocket безопасность
- ✅ Добавить отсутствующий contacts handler

### **Приоритет 2 (Важно)**
- 🔄 Убрать хардкод MinIO URL
- 🔄 Унифицировать rate limiting
- 🔄 Создать план deprecation для legacy routes

### **Приоритет 3 (Желательно)**
- 📋 Создать OpenAPI спецификацию
- 📋 Добавить middleware documentation
- 📋 Implement health checks для всех зависимостей

### **Метрики успеха**
- Функция setupRoutes() < 50 строк
- Каждая группа маршрутов в отдельной функции < 30 строк
- 100% endpoints имеют proper middleware
- 0 дублирующих маршрутов

---

## 🔗 Связанные файлы

- `backend/internal/middleware/` - middleware функции
- `backend/internal/proj/*/handler/` - handlers для маршрутов
- `backend/internal/config/config.go` - конфигурация
- `backend/docs/` - документация API

---

**Последнее обновление:** 6 декабря 2025  
**Статус:** Требует рефакторинг  
**Reviewer:** Backend Team  