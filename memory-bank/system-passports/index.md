# Паспорта подсистем Sve Tu Platform

## Назначение

Данная система паспортов создана для:
- Документирования структуры данных всех слоев системы
- Выявления и устранения несоответствий в наименованиях
- Обеспечения согласованности между слоями
- Упрощения разработки и поддержки

## Структура паспортов

### Общие паспорта
- [Архитектура системы](./architecture.md) - Общая схема взаимодействия слоев
- [Соглашения о наименованиях](./naming-conventions.md) - Правила именования во всех слоях
- [Матрица соответствий](./field-mapping-matrix.md) - Таблица маппинга полей между слоями
- [Несоответствия](./inconsistencies.md) - Выявленные проблемы и план их устранения

### Паспорта по слоям

#### Backend (Go)
- [Главный паспорт Backend](./backend/index.md)
- [API Endpoints](./backend/api-endpoints.md)
- Модели:
  - [User](./backend/models/user.md)
  - [Marketplace](./backend/models/marketplace.md)
  - [Payments](./backend/models/payments.md)
  - [Общие типы](./backend/models/common.md)
- Handlers:
  - [Auth](./backend/handlers/auth.md)
  - [Marketplace](./backend/handlers/marketplace.md)
  - [Payments](./backend/handlers/payments.md)
- [Swagger/OpenAPI типы](./backend/swagger-types.md)

#### Frontend (React/TypeScript)
- [Главный паспорт Frontend](./frontend/index.md)
- [API клиент](./frontend/api-client.md)
- Типы:
  - [Сгенерированные из OpenAPI](./frontend/types/generated.md)
  - [Redux Store](./frontend/types/redux-store.md)
  - [Компоненты](./frontend/types/components.md)
- [Поток данных](./frontend/data-flow.md)
- [Маппинг имен](./frontend/naming-mapping.md)

##### Frontend компоненты (96/96 завершено - 100%) ✅
- Auth (5/5): ✅
  - [AuthButton](./frontend/auth/AuthButton.md) ✅
  - [AuthContext](./frontend/auth/AuthContext.md) ✅
  - [LoginModal](./frontend/auth/LoginModal.md) ✅
  - [LoginForm](./frontend/auth/LoginForm.md) ✅
  - [RegisterForm](./frontend/auth/RegisterForm.md) ✅
- Marketplace (7/7): ✅
  - [MarketplaceList](./frontend/marketplace/MarketplaceList.md) ✅
  - [MarketplaceCard](./frontend/marketplace/MarketplaceCard.md) ✅
  - [MarketplaceFilters](./frontend/marketplace/MarketplaceFilters.md) ✅
  - [ImageGallery](./frontend/marketplace/ImageGallery.md) ✅
  - [ListingActions](./frontend/marketplace/ListingActions.md) ✅
  - [SimilarListings](./frontend/marketplace/SimilarListings.md) ✅
  - [SellerInfo](./frontend/marketplace/SellerInfo.md) ✅
- Chat (9/9): ✅
  - [ChatLayout](./frontend/chat/ChatLayout.md) ✅
  - [ChatList](./frontend/chat/ChatList.md) ✅
  - [ChatWindow](./frontend/chat/ChatWindow.md) ✅
  - [MessageItem](./frontend/chat/MessageItem.md) ✅
  - [MessageInput](./frontend/chat/MessageInput.md) ✅
  - [ChatAttachments](./frontend/chat/ChatAttachments.md) ✅
  - [AnimatedEmoji](./frontend/chat/AnimatedEmoji.md) ✅
  - [EmojiPicker](./frontend/chat/EmojiPicker.md) ✅
  - [FileUploadProgress](./frontend/chat/FileUploadProgress.md) ✅
- Payment (3/3): ✅
  - [EscrowStatus](./frontend/payment/EscrowStatus.md) ✅
  - [PaymentMethodSelector](./frontend/payment/PaymentMethodSelector.md) ✅
  - [PaymentProcessing](./frontend/payment/PaymentProcessing.md) ✅
- Reviews (9/9): ✅
  - [ImageGallery](./frontend/reviews/ImageGallery.md) ✅
  - [RatingDisplay](./frontend/reviews/RatingDisplay.md) ✅
  - [RatingInput](./frontend/reviews/RatingInput.md) ✅
  - [RatingStats](./frontend/reviews/RatingStats.md) ✅
  - [ReviewForm](./frontend/reviews/ReviewForm.md) ✅
  - [ReviewList](./frontend/reviews/ReviewList.md) ✅
  - [ReviewsSection](./frontend/reviews/ReviewsSection.md) ✅
  - [ReviewFormExample](./frontend/reviews/ReviewFormExample.md) ✅
  - [index](./frontend/reviews/index.md) ✅
- CreateListing (9/9): ✅
  - [StepWizard](./frontend/create-listing/StepWizard.md) ✅
  - [BasicInfoStep](./frontend/create-listing/BasicInfoStep.md) ✅
  - [CategorySelectionStep](./frontend/create-listing/CategorySelectionStep.md) ✅
  - [AttributesStep](./frontend/create-listing/AttributesStep.md) ✅
  - [PhotosStep](./frontend/create-listing/PhotosStep.md) ✅
  - [LocationStep](./frontend/create-listing/LocationStep.md) ✅
  - [PaymentDeliveryStep](./frontend/create-listing/PaymentDeliveryStep.md) ✅
  - [TrustSetupStep](./frontend/create-listing/TrustSetupStep.md) ✅
  - [PreviewPublishStep](./frontend/create-listing/PreviewPublishStep.md) ✅
- Storefronts (18/18): ✅
  - [StorefrontHeader](./frontend/storefronts/StorefrontHeader.md) ✅
  - [StorefrontInfo](./frontend/storefronts/StorefrontInfo.md) ✅
  - [StorefrontActions](./frontend/storefronts/StorefrontActions.md) ✅
  - [StorefrontProducts](./frontend/storefronts/StorefrontProducts.md) ✅
  - [StorefrontMap](./frontend/storefronts/StorefrontMap.md) ✅
  - [MapFilters](./frontend/storefronts/MapFilters.md) ✅
  - [AddressSearch](./frontend/storefronts/AddressSearch.md) ✅
  - [StorefrontMapContainer](./frontend/storefronts/StorefrontMapContainer.md) ✅
  - [BasicInfoStep](./frontend/storefronts/BasicInfoStep.md) ✅
  - [BusinessDetailsStep](./frontend/storefronts/BusinessDetailsStep.md) ✅
  - [BusinessHoursStep](./frontend/storefronts/BusinessHoursStep.md) ✅
  - [LocationStep](./frontend/storefronts/LocationStep.md) ✅
  - [PaymentDeliveryStep](./frontend/storefronts/PaymentDeliveryStep.md) ✅
  - [StaffSetupStep](./frontend/storefronts/StaffSetupStep.md) ✅
  - [PreviewPublishStep](./frontend/storefronts/PreviewPublishStep.md) ✅
  - [LocationMap](./frontend/storefronts/LocationMap.md) ✅
  - [CreateStorefrontContext](./frontend/storefronts/CreateStorefrontContext.md) ✅
  - [StorefrontLocationMap](./frontend/storefronts/StorefrontLocationMap.md) ✅
- Products (10/10): ✅
  - [ProductList](./frontend/products/ProductList.md) ✅
  - [ProductCard](./frontend/products/ProductCard.md) ✅
  - [BulkActions](./frontend/products/BulkActions.md) ✅
  - [ProductWizard](./frontend/products/ProductWizard.md) ✅
  - [CreateProductContext](./frontend/products/CreateProductContext.md) ✅
  - [BasicInfoStep](./frontend/products/BasicInfoStep.md) ✅
  - [CategoryStep](./frontend/products/CategoryStep.md) ✅
  - [PhotosStep](./frontend/products/PhotosStep.md) ✅
  - [PreviewStep](./frontend/products/PreviewStep.md) ✅
  - [productSlice](./frontend/products/productSlice.md) ✅
- Import (5/5): ✅
  - [ImportWizard](./frontend/import/ImportWizard.md) ✅
  - [ImportManager](./frontend/import/ImportManager.md) ✅
  - [ImportJobsList](./frontend/import/ImportJobsList.md) ✅
  - [ImportJobDetails](./frontend/import/ImportJobDetails.md) ✅
  - [ImportErrorsModal](./frontend/import/ImportErrorsModal.md) ✅
- Search (3/3): ✅
  - [SearchBar](./frontend/search/SearchBar.md) ✅
  - [SearchPage](./frontend/search/SearchPage.md) ✅
  - [UnifiedSearchService](./frontend/search/UnifiedSearchService.md) ✅
- Common (15/15): ✅
  - [Обзор Common компонентов](./frontend/common-components.md) ✅
  - [Индекс Common компонентов](./frontend/common-components-index.md) ✅
  - [InfiniteScrollTrigger](./frontend/common-components.md#1-infinitescrolltrigger) ✅
  - [ViewToggle](./frontend/common-components.md#2-viewtoggle) ✅
  - [ErrorBoundary](./frontend/common-components.md#3-errorboundary-autherrorboundary) ✅
  - [FormField](./frontend/common-components.md#4-formfield) ✅
  - [OptimizedImage](./frontend/common-components.md#5-optimizedimage) ✅
  - [SafeImage](./frontend/common-components.md#6-safeimage) ✅
  - [LanguageSwitcher](./frontend/common-components.md#7-languageswitcher) ✅
  - [DraftStatus](./frontend/common-components.md#8-draftstatus) ✅
  - [DraftsModal](./frontend/common-components.md#9-draftsmodal) ✅
  - [IconPicker](./frontend/common-components.md#10-iconpicker) ✅
  - [GoogleIcon](./frontend/common-components.md#11-googleicon) ✅
  - [WebSocketManager](./frontend/common-components.md#12-websocketmanager) ✅
  - [ReduxProvider](./frontend/common-components.md#13-reduxprovider) ✅
  - [AdminGuard](./frontend/common-components.md#14-adminguard) ✅
  - [AuthStateManager](./frontend/common-components.md#15-authstatemanager) ✅
- Base (11/11): ✅
  - [Паспорта базовых компонентов](./frontend/base-components.md) ✅
  - [RootLayout](./frontend/base-components.md#1-rootlayout-main-layout) ✅
  - [AdminLayout](./frontend/base-components.md#2-adminlayout) ✅
  - [Header](./frontend/base-components.md#3-header) ✅
  - [HomePage](./frontend/base-components.md#4-homepage-main-page) ✅
  - [ReduxProvider](./frontend/base-components.md#5-reduxprovider) ✅
  - [AuthStateManager](./frontend/base-components.md#6-authstatemanager) ✅
  - [WebSocketManager](./frontend/base-components.md#7-websocketmanager) ✅
  - [SearchBar](./frontend/base-components.md#8-searchbar) ✅
  - [LanguageSwitcher](./frontend/base-components.md#9-languageswitcher) ✅
  - [AuthButton](./frontend/base-components.md#10-authbutton) ✅
  - [ErrorBoundary](./frontend/base-components.md#11-errorboundary) ✅
- Redux Store (1/1): ✅
  - [ReduxStore](./frontend/redux/ReduxStore.md) ✅

#### Database (PostgreSQL)
- [Главный паспорт Database](./database/index.md) ✅
- Таблицы (5/38 завершено):
  - [users](./database/tables/users.md) ✅
  - [marketplace_categories](./database/tables/marketplace_categories.md) ✅
  - [marketplace_listings](./database/tables/marketplace_listings.md) ✅
  - [marketplace_images](./database/tables/marketplace_images.md) ✅
  - [marketplace_favorites](./database/tables/marketplace_favorites.md) ✅
  - [Все таблицы...](./database/index.md#таблицы-по-группам)
- [Связи таблиц](./database/relationships.md)
- [Правила именования](./database/naming-rules.md)

#### OpenSearch
- [Главный паспорт OpenSearch](./opensearch/index.md)
- Индексы:
  - [listings](./opensearch/indices/listings.md)
  - [users](./opensearch/indices/users.md)
- [Маппинги полей](./opensearch/mappings.md)
- [Процесс синхронизации](./opensearch/sync-process.md)

#### MinIO
- [Главный паспорт MinIO](./minio/index.md)
- [Структура buckets](./minio/buckets.md)
- [Структура путей](./minio/paths.md)
- [Метаданные](./minio/metadata.md)

## Как использовать паспорта

1. **При разработке новой функциональности**:
   - Проверьте соответствующие паспорта слоев
   - Используйте правильные наименования из матрицы соответствий
   - Обновите паспорта при добавлении новых полей/типов

2. **При исправлении багов**:
   - Проверьте раздел "Несоответствия"
   - Убедитесь в правильном маппинге между слоями

3. **При code review**:
   - Проверяйте соответствие кода паспортам
   - Отмечайте расхождения с документацией

## Поддержка актуальности

- Паспорта должны обновляться при любых изменениях структуры данных
- Раз в месяц проводить аудит соответствия кода паспортам
- Все найденные несоответствия документировать в [inconsistencies.md](./inconsistencies.md)