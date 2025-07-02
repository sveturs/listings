# СИСТЕМНЫЙ ПАСПОРТ: Reviews Handler

## 📋 Обзор модуля

**Назначение**: Комплексная система управления отзывами и рейтингами  
**Расположение**: `/backend/internal/proj/reviews/`  
**Тип**: Backend handler  
**Статус**: ✅ Активный  

### 🎯 Основные функции
- Двухэтапное создание отзывов (черновик → публикация)
- Загрузка и управление фотографиями к отзывам
- Система голосования за полезность отзывов
- Ответы продавцов на отзывы
- Подтверждение/спор отзывов продавцами
- Агрегированные рейтинги для пользователей и витрин
- Мультиязычность и автопереводы
- Интеграция с поисковым индексом

## 🏗️ Архитектура модуля

### 📁 Структура файлов
```
backend/internal/proj/reviews/
├── handler/
│   ├── handler.go          # Регистрация маршрутов и фабрика
│   ├── reviews.go          # Основные HTTP handlers (1019 строк)
│   └── responses.go        # Структуры ответов API
├── service/
│   ├── interface.go        # Интерфейс ReviewServiceInterface
│   ├── service.go          # Фабрика сервисов
│   └── review.go          # Основная бизнес-логика
├── middleware/             # Дополнительные middleware
└── storage/
    ├── interface.go        # Интерфейс ReviewRepository
    └── postgres/
        ├── storage.go      # PostgreSQL Storage factory
        └── reviews.go      # Реализация для PostgreSQL
```

### 🔧 Основные компоненты

#### Handler (handler.go:19-21)
```go
type Handler struct {
    Review *ReviewHandler
}
```

#### ReviewHandler (reviews.go:19-36)
```go
type ReviewHandler struct {
    services      globalService.ServicesInterface
    reviewService service.ReviewServiceInterface
}
```

#### ReviewService (service/review.go:14-25)
```go
type ReviewService struct {
    storage storage.Storage
}
```

## 🛠️ API Endpoints

### 🌐 Публичные маршруты

| Метод | Путь | Функция | Описание |
|-------|------|---------|----------|
| GET | `/api/v1/reviews` | GetReviews | Список отзывов с фильтрами |
| GET | `/api/v1/reviews/:id` | GetReviewByID | Получить отзыв по ID |
| GET | `/api/v1/reviews/stats` | GetStats | Статистика отзывов |
| GET | `/api/v1/entity/:type/:id/rating` | GetEntityRating | Рейтинг сущности |
| GET | `/api/v1/entity/:type/:id/stats` | GetEntityStats | Статистика сущности |
| GET | `/api/v1/users/:id/aggregated-rating` | GetUserAggregatedRating | Агрегированный рейтинг пользователя |
| GET | `/api/v1/storefronts/:id/aggregated-rating` | GetStorefrontAggregatedRating | Агрегированный рейтинг витрины |
| GET | `/api/v1/public/storefronts/:id/reviews` | GetStorefrontReviews | Публичные отзывы витрины |
| GET | `/api/v1/public/storefronts/:id/rating` | GetStorefrontRatingSummary | Публичный рейтинг витрины |

### 🔐 Защищенные маршруты (JWT + CSRF)

| Метод | Путь | Функция | Описание |
|-------|------|---------|----------|
| GET | `/api/v1/reviews/can-review/:type/:id` | CanReview | Проверка возможности оставить отзыв |
| POST | `/api/v1/reviews/draft` | CreateDraftReview | Создать черновик отзыва |
| POST | `/api/v1/reviews/:id/photos` | UploadPhotos | Загрузить фото к отзыву |
| POST | `/api/v1/reviews/:id/publish` | PublishReview | Опубликовать черновик |
| PUT | `/api/v1/reviews/:id` | UpdateReview | Обновить отзыв |
| DELETE | `/api/v1/reviews/:id` | DeleteReview | Удалить отзыв |
| POST | `/api/v1/reviews/:id/vote` | VoteForReview | Голосовать за отзыв |
| POST | `/api/v1/reviews/:id/response` | AddResponse | Добавить ответ на отзыв |
| POST | `/api/v1/reviews/:id/confirm` | ConfirmReview | Подтвердить отзыв |
| POST | `/api/v1/reviews/:id/dispute` | DisputeReview | Создать спор по отзыву |
| POST | `/api/v1/reviews/upload-photos` | UploadPhotosForNewReview | Загрузка фото (legacy) |
| GET | `/api/v1/users/:id/reviews` | GetUserReviews | Отзывы пользователя |
| GET | `/api/v1/users/:id/rating` | GetUserRatingSummary | Рейтинг пользователя |
| GET | `/api/v1/storefronts/:id/reviews` | GetStorefrontReviews | Отзывы витрины |
| GET | `/api/v1/storefronts/:id/rating` | GetStorefrontRatingSummary | Рейтинг витрины |

## 🗄️ Модели данных

### Review (основная модель отзыва)
```go
type Review struct {
    ID                 int                          `json:"id"`
    UserID             int                          `json:"user_id"`
    EntityType         string                       `json:"entity_type"`
    EntityID           int                          `json:"entity_id"`
    EntityOriginType   string                       `json:"entity_origin_type,omitempty"`
    EntityOriginID     int                          `json:"entity_origin_id,omitempty"`
    Rating             int                          `json:"rating"`
    Comment            string                       `json:"comment,omitempty"`
    Pros               string                       `json:"pros,omitempty"`
    Cons               string                       `json:"cons,omitempty"`
    Photos             []string                     `json:"photos,omitempty"`
    LikesCount         int                          `json:"likes_count"`
    IsVerifiedPurchase bool                         `json:"is_verified_purchase"`
    Status             string                       `json:"status"`
    HelpfulVotes       int                          `json:"helpful_votes"`
    NotHelpfulVotes    int                          `json:"not_helpful_votes"`
    SellerConfirmed    bool                         `json:"seller_confirmed"`
    HasActiveDispute   bool                         `json:"has_active_dispute"`
    OriginalLanguage   string                       `json:"original_language"`
    Translations       map[string]map[string]string `json:"translations,omitempty"`
    User               *User                        `json:"user,omitempty"`
    Responses          []ReviewResponse             `json:"responses,omitempty"`
    CreatedAt          time.Time                    `json:"created_at"`
    UpdatedAt          time.Time                    `json:"updated_at"`
}
```

### CreateReviewRequest (запрос создания)
```go
type CreateReviewRequest struct {
    EntityType       string   `json:"entity_type" validate:"required,oneof=listing room car"`
    EntityID         int      `json:"entity_id" validate:"required"`
    Rating           int      `json:"rating" validate:"required,min=1,max=5"`
    StorefrontID     *int     `json:"storefront_id,omitempty"`
    Comment          string   `json:"comment"`
    Pros             string   `json:"pros,omitempty"`
    Cons             string   `json:"cons,omitempty"`
    Photos           []string `json:"photos"`
    OriginalLanguage string   `json:"original_language" validate:"required"`
}
```

### ReviewsFilter (фильтры поиска)
```go
type ReviewsFilter struct {
    EntityType string `query:"entity_type"`
    EntityID   int    `query:"entity_id"`
    UserID     int    `query:"user_id"`
    MinRating  int    `query:"min_rating"`
    MaxRating  int    `query:"max_rating"`
    Status     string `query:"status"`
    SortBy     string `query:"sort_by"`    // rating, date, likes
    SortOrder  string `query:"sort_order"` // asc, desc
    Page       int    `query:"page"`
    Limit      int    `query:"limit"`
}
```

### ReviewStats (статистика)
```go
type ReviewStats struct {
    TotalReviews       int         `json:"total_reviews"`
    AverageRating      float64     `json:"average_rating"`
    VerifiedReviews    int         `json:"verified_reviews"`
    RatingDistribution map[int]int `json:"rating_distribution"`
    PhotoReviews       int         `json:"photo_reviews"`
}
```

## 🔄 Бизнес-процессы

### Двухэтапное создание отзыва
1. **Черновик** (POST `/reviews/draft`):
   - Создание отзыва со статусом `draft`
   - Определение языка текста
   - Санитизация от XSS
   - Проверка верифицированной покупки
   - Установка entity_origin для агрегации

2. **Загрузка фото** (POST `/reviews/:id/photos`):
   - Валидация форматов (JPEG/PNG/WebP)
   - Ограничение размера (5MB на файл, макс 5 файлов)
   - Загрузка в MinIO с уникальными именами

3. **Публикация** (POST `/reviews/:id/publish`):
   - Смена статуса на `published`
   - Обновление рейтинга в поисковом индексе
   - Отправка уведомлений через Notifications service

### Система голосования (reviews.go:243-295)
- Типы голосов: `helpful` / `not_helpful`
- Одно голосование на пользователя на отзыв
- Обновление счетчиков в режиме реального времени
- Уведомления автору отзыва

### Ответы продавцов (reviews.go:311-366)
- Только владелец объявления может отвечать
- Санитизация XSS для ответов
- Множественные ответы разрешены
- Уведомления автору отзыва

### Подтверждения и споры (reviews.go:959-1018)
- **Подтверждение**: продавец подтверждает отзыв как легитимный
- **Спор**: продавец оспаривает отзыв (причины: not_a_customer, false_information, deal_cancelled, spam, other)
- Система флагов `seller_confirmed` и `has_active_dispute`

## 🔒 Безопасность и валидация

### Input Validation
- XSS санитизация через `utils.SanitizeText()` (reviews.go:61, 327, 420)
- Валидация файлов: типы, размеры, количество
- Проверка авторства для операций редактирования
- CSRF защита для изменяющих операций

### File Upload Security (reviews.go:507-547)
```go
// Разрешенные форматы
allowedFormats := map[string]bool{
    "image/jpeg": true,
    "image/jpg":  true,
    "image/png":  true,
    "image/webp": true,
}

// Проверка размера (максимум 5MB)
if file.Size > 5*1024*1024 {
    return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.file_too_large")
}
```

### Permissions
- Авторы могут редактировать только свои отзывы
- Владельцы объявлений могут отвечать на отзывы
- Система проверки `CanUserReviewEntity`

## 🗃️ База данных

### Связанные таблицы
- `reviews` - основная таблица отзывов
- `review_responses` - ответы на отзывы
- `review_votes` - голоса за полезность
- `users` - связь с пользователями
- `marketplace_listings` - связь с объявлениями
- `storefronts` - связь с витринами

### Агрегация данных
```sql
-- Обновление рейтинга в поисковом индексе
SELECT COUNT(*), COALESCE(AVG(rating), 0)
FROM reviews
WHERE entity_type = $1 AND entity_id = $2 AND status = 'published'
```

## 🔗 Внешние интеграции

### OpenSearch Integration (service/review.go:27-71)
- Автоматическое обновление рейтингов в поисковом индексе
- Пересчет при публикации/удалении отзывов
- Обновление полей `average_rating` и `review_count` в листингах

### Translation Service
- Автоопределение языка отзывов
- Поддержка мультиязычных переводов
- Хранение оригинального языка

### MinIO File Storage
- Загрузка фотографий к отзывам
- Генерация уникальных имен файлов
- Bucket: `reviews/` для постоянных фото, `temp/` для временных

### Notification Service
- Уведомления о новых отзывах владельцам объявлений
- Уведомления о голосовании авторам отзывов
- Уведомления об ответах авторам отзывов

## 📊 Специальные возможности

### Агрегированные рейтинги (reviews.go:857-907)
- **Пользователи**: рейтинг по всем отзывам, полученным пользователем
- **Витрины**: совокупный рейтинг по всем товарам витрины
- Разбивка по источникам (прямые отзывы, через товары)
- Процент верифицированных покупок

### Верифицированные покупки (service/review.go:96)
- Проверка факта покупки через платежную систему
- Специальный маркер `is_verified_purchase`
- Повышенный вес в агрегированных рейтингах

### Мультиязычность
- Определение языка оригинального текста
- Хранение переводов в JSON формате
- API возвращает переводы на запрашиваемом языке

## 🏭 Фабричные методы

### Service Factory (service/service.go:11-15)
```go
func NewService(storage storage.Storage) *Service {
    return &Service{
        Review: NewReviewService(storage),
    }
}
```

### Handler Factory (handler.go:12-17)
```go
func NewHandler(services globalService.ServicesInterface) *Handler {
    return &Handler{
        Review: NewReviewHandler(services),
    }
}
```

## 📝 Структуры ответов

### ReviewsListResponse (responses.go:18-22)
```go
type ReviewsListResponse struct {
    Success bool            `json:"success"`
    Data    []models.Review `json:"data"`
    Meta    ReviewsMeta     `json:"meta"`
}
```

### PhotosResponse (responses.go:48-52)
```go
type PhotosResponse struct {
    Success bool     `json:"success"`
    Message string   `json:"message"`
    Photos  []string `json:"photos"`
}
```

### RatingResponse (responses.go:55-58)
```go
type RatingResponse struct {
    Success bool    `json:"success"`
    Rating  float64 `json:"rating"`
}
```

## ⚠️ Особенности реализации

### Двухэтапное создание
- Решает проблему потери данных при загрузке фото
- Позволяет предварительный просмотр перед публикацией
- Автоматическая очистка неопубликованных черновиков

### Поисковая интеграция
- Синхронное обновление рейтингов в OpenSearch
- Fallback при ошибках индексации
- Логирование всех операций индексации

### Entity Origin System
- Агрегация отзывов по источникам (пользователь/витрина)
- Автоматическое определение origin при создании
- Поддержка сложных схем агрегации

## 🔄 Связи с другими модулями

### Входящие зависимости
- `marketplace` handler - информация об объявлениях
- `users` handler - данные пользователей
- `storefronts` handler - информация о витринах
- `payments` handler - верификация покупок

### Исходящие зависимости
- PostgreSQL storage для всех операций с БД
- OpenSearch для обновления рейтингов в поиске
- MinIO для хранения фотографий
- Translation service для мультиязычности
- Notification service для уведомлений

## 🚀 TODO и улучшения

### Технические улучшения
- [ ] Batch обновления рейтингов в OpenSearch
- [ ] Кэширование агрегированных рейтингов
- [ ] Асинхронная обработка уведомлений
- [ ] Автоочистка старых черновиков

### Функциональные улучшения
- [ ] Система модерации отзывов
- [ ] AI-анализ тональности отзывов
- [ ] Рекомендации похожих отзывов
- [ ] Экспорт отзывов в различные форматы

### Безопасность
- [ ] Rate limiting для создания отзывов
- [ ] Система антиспама
- [ ] Детекция накрутки рейтингов
- [ ] Аудит логи всех операций

## 📊 Метрики и мониторинг

### Логируемые события
- Создание, публикация, обновление отзывов
- Загрузка фотографий и ошибки файлов
- Ошибки интеграции с OpenSearch
- Операции голосования и ответов

### Рекомендуемые метрики
- Количество отзывов по типам сущностей
- Процент верифицированных отзывов
- Среднее время от черновика до публикации
- Активность голосования и ответов
- Статистика споров и подтверждений

---

**Дата создания**: $(date)  
**Версия**: 1.0  
**Статус**: ✅ Активный модуль  
**Последнее обновление**: Двухэтапная система создания отзывов с интеграцией поиска