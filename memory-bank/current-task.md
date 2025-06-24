# Текущая задача: Массовое управление товарами и проблема безопасности изображений

## Статус: Активна
## Дата начала: 2025-06-22
## Последнее обновление: 2025-06-24 12:20 UTC

### История задач:

#### Предыдущая задача: Улучшение функционала витрин (Завершена)
- ✅ Создана таблица для хранения реальной аналитики (storefront_analytics)
- ✅ Создана таблица для событий (storefront_events)
- ✅ Реализован API endpoint для получения аналитики витрины
- ✅ Реализован сервис сбора событий аналитики
- ✅ Создан хук useAnalytics для отслеживания событий на frontend
- ✅ Обновлен дашборд для отображения реальных данных
- ✅ Добавлено отслеживание просмотров страниц витрины
- ✅ Создан background job для агрегации аналитики

#### Предыдущая задача: Исправление импорта товаров (Завершена)
- ✅ Исправлены селекторы Redux (storefronts вместо listings)
- ✅ Добавлена авторизация в запросы (Bearer token)
- ✅ Создан роут на основе slug витрины
- ✅ Исправлен парсинг XML файлов (images теперь массив)
- ✅ Успешно протестирован импорт товаров (ID: 19, 20)

### Текущие задачи:

#### 1. Массовое управление товарами (В процессе реализации Frontend)

**Предыстория задачи:**
- ✅ 2025-06-22: Исправлена ошибка загрузки изображений с Unsplash (async/await в getImageUrl)
- ✅ 2025-06-22: Исправлена проблема с async params в Next.js 13+ (использован useEffect для params)
- ✅ 2025-06-23: Закоммичены все изменения и синхронизирована ветка с main
- ✅ 2025-06-23: Исправлены замечания по PR #50 (убран хардкод localhost в MarketplaceCard, перегенерированы типы)
- ✅ 2025-06-24: PR #50 принят и влит в main

**ВАЖНОЕ ОТКРЫТИЕ (2025-06-24 11:00 UTC):**
При детальном изучении кода выяснилось, что функционал массовых операций, который был отмечен как реализованный в предыдущих сессиях, на самом деле НЕ СУЩЕСТВОВАЛ. Все отметки о создании UI компонентов были ошибочными.

**Реализация Backend (2025-06-24 11:00-11:40 UTC):**
✅ Добавлены модели в `/backend/internal/domain/models/storefront_product.go` (строки 168-224):
  - `BulkCreateProductsRequest/Response` - для массового создания
  - `BulkUpdateProductsRequest/Response` с `BulkUpdateItem` - для массового обновления
  - `BulkDeleteProductsRequest/Response` - для массового удаления
  - `BulkUpdateStatusRequest/Response` - для массового изменения статуса
  - `BulkOperationError` - для отчета об ошибках с указанием индекса/ID

✅ Реализованы storage методы в `/backend/internal/storage/postgres/storefront_product.go` (строки 543-805):
  - `BulkCreateProducts` - использует транзакцию, batch INSERT
  - `BulkUpdateProducts` - проверка принадлежности товаров, динамическое построение UPDATE
  - `BulkDeleteProducts` - каскадное удаление с RETURNING для отчета
  - `BulkUpdateStatus` - оптимизированное обновление через ANY($2)
  - Исправлена ошибка: заменен "NOW()" на CURRENT_TIMESTAMP (строки 679, 774)
  - Исправлена неиспользуемая переменная i (строка 626)

✅ Добавлены handlers в `/backend/internal/proj/storefronts/handler/product_handler.go` (строки 419-601):
  - `BulkCreateProducts` - POST /products/bulk/create
  - `BulkUpdateProducts` - PUT /products/bulk/update
  - `BulkDeleteProducts` - DELETE /products/bulk/delete
  - `BulkUpdateStatus` - PUT /products/bulk/status
  - Все handlers используют существующие helper функции для проверки доступа

✅ Обновлен service слой в `/backend/internal/proj/storefronts/service/product_service.go`:
  - Добавлены методы в интерфейс Storage (строки 28-32)
  - Реализованы service методы (строки 302-481) с:
    - Проверкой владельца витрины
    - Валидацией всех элементов запроса
    - Асинхронной индексацией в OpenSearch через goroutines
    - Конвертацией ошибок в формат ответа

✅ Зарегистрированы роуты в `/backend/internal/proj/storefronts/module.go`:
  - Добавлены 4 новых роута (строки 121-125)
  - Созданы wrapper функции для работы через slug (строки 489-544)

✅ Перегенерированы OpenAPI типы:
  - Выполнен `make generate-types`
  - Обновлены TypeScript типы в frontend

**Реализация Frontend (2025-06-24 11:40-12:20 UTC):**
✅ Создан Redux slice `/frontend/svetu/src/store/slices/productSlice.ts`:
  - Комплексное состояние с фильтрами, пагинацией, UI настройками
  - Использование Set для selectedIds (оптимальная производительность)
  - Async thunks: bulkDeleteProducts, bulkUpdateStatus, exportProducts
  - Reducers для всех операций выбора (toggle, selectAll, clearSelection, selectByFilter)
  - Обработка состояний loading/success/error для bulk операций
  - Интеграция с toast для уведомлений

✅ Обновлен store в `/frontend/svetu/src/store/index.ts`:
  - Добавлен productReducer (строка 15)
  - Добавлено игнорирование сериализации для products.selectedIds (строка 35)

✅ Создан компонент BulkActions `/frontend/svetu/src/components/products/BulkActions.tsx`:
  - Современный дизайн с DaisyUI компонентами
  - Sticky позиционирование для удобства
  - Dropdown меню для выбора активации/деактивации
  - Двухэтапное подтверждение удаления с таймаутом
  - Прогресс-бар для отображения процесса
  - Поддержка отключения кнопок при обработке

✅ Создан компонент ProductCard `/frontend/svetu/src/components/products/ProductCard.tsx`:
  - Три режима отображения: grid, list, table
  - Оптимизация с React.memo
  - Чекбоксы с правильной обработкой событий (stopPropagation)
  - Визуальная индикация выбранных товаров (ring-2 ring-primary)
  - Отображение всех важных данных: цена, склад, SKU, статус
  - Цветовая индикация статуса склада

✅ Создан API сервис `/frontend/svetu/src/services/productApi.ts`:
  - Методы для всех bulk операций с правильной типизацией
  - Экспорт в CSV/XML с автоматическим скачиванием
  - Использование сгенерированных типов из OpenAPI

✅ Создан компонент ProductList `/frontend/svetu/src/components/products/ProductList.tsx`:
  - Панель управления с поиском (debounced) и фильтрами
  - Расширяемая панель фильтров (статус, склад, диапазон цен)
  - Переключение режимов отображения (grid/list/table)
  - Интеграция с Redux для состояния
  - Поддержка бесконечной прокрутки
  - Режим выбора с визуальной индикацией
  - Показ общего количества товаров

✅ Добавлены переводы в `/frontend/svetu/src/messages/ru.json` (строки 1763-1799):
  - Полный набор переводов для bulk операций
  - Pluralization для правильного склонения
  - Переводы для статусов склада
  - Все необходимые UI элементы

✅ Создана новая версия страницы `/frontend/svetu/src/app/[locale]/storefronts/[slug]/products/page-new.tsx`:
  - Интеграция с новыми компонентами
  - Упрощенная логика с использованием Redux
  - Сохранена совместимость с существующим API

#### 2. Критическая проблема безопасности (Требует решения)
**Проблема**: Система хранит только URL изображений из XML файлов. Это создает уязвимость - изображения могут быть подменены на источнике после модерации.

**Решение**: Необходимо реализовать загрузку изображений в MinIO при импорте товаров.

### План действий (обновлен 2025-06-24 12:20 UTC):
1. [x] Завершить реализацию backend для массовых операций
2. [x] Реализовать frontend для массовых операций
3. [ ] Протестировать функционал массового управления
   - [ ] Заменить старую страницу на новую (переименовать page-new.tsx в page.tsx)
   - [ ] Протестировать массовое удаление товаров
   - [ ] Протестировать массовое изменение статуса
   - [ ] Протестировать экспорт в CSV/XML
   - [ ] Проверить работу фильтров и поиска
   - [ ] Проверить производительность при большом количестве товаров
4. [ ] Реализовать загрузку изображений в MinIO при импорте (КРИТИЧНО)
5. [ ] Добавить прогресс-бар для массовых операций (частично реализовано)
6. [ ] Добавить фильтрацию по атрибутам товаров

### Корректировки к плану:
1. **Добавить экспорт товаров** - не было в изначальном плане, но реализовано в API сервисе. Нужно добавить endpoints на backend.
2. **Массовое редактирование** - кнопка добавлена в UI, но функционал не реализован. Решить, нужен ли.
3. **WebSocket для прогресса** - для больших операций может понадобиться real-time обновление прогресса.

### Результаты изучения кода:

#### Backend архитектура:
- **Handlers**: `/backend/internal/proj/storefronts/handler/product_handler.go`
- **Services**: `/backend/internal/proj/storefronts/service/product_service.go`
- **Storage**: `/backend/internal/storage/postgres/storefront_product.go`
- **Models**: `/backend/internal/domain/models/storefront_product.go`

**Текущие endpoints**:
- Есть CRUD операции для единичных товаров
- Есть импорт товаров через XML/CSV/ZIP
- **НЕТ** bulk API endpoints для массовых операций
- **НЕТ** оптимизированных bulk методов в storage

#### Проблема безопасности изображений:
- **Marketplace (объявления)**: БЕЗОПАСНО - изображения загружаются в MinIO
- **Storefronts (товары)**: НЕБЕЗОПАСНО - хранятся только внешние URL
- Риски: утечка IP, подмена контента, зависимость от внешних серверов

#### Frontend состояние:
- **ОТСУТСТВУЕТ** функционал массовых операций
- Нет чекбоксов для выбора товаров
- Нет компонента BulkActions
- Нет Redux slice для управления товарами
- Все операции только индивидуальные

#### База данных:
**Таблицы**:
- `storefront_products` - основная таблица товаров
- `storefront_product_images` - изображения (хранят только URL!)
- `storefront_product_variants` - варианты товаров
- `storefront_inventory_movements` - движение товаров

**Особенности**:
- Уникальные индексы: (storefront_id, sku), (storefront_id, barcode)
- JSON поле attributes для доп. атрибутов
- Полнотекстовый поиск через GIN индекс

#### OpenSearch:
- При создании/обновлении товара - индексация
- Нет bulk API для массовой индексации

#### MinIO Storage:
- Есть метод `UploadImageFromURL` в minio storage
- Используется только для marketplace объявлений
- **НЕ используется** для товаров витрин

### Примечания и важные детали реализации:

#### Статус веток и PR:
- Ветка feature/improve-product-attributes-selection полностью синхронизирована с main
- PR #50 был принят и влит в main 2025-06-24 08:23 UTC
- После merge ветка стала идентичной main (0 коммитов впереди)
- Все новые изменения пока не закоммичены

#### Технические детали реализации:
1. **Redux Set для selectedIds**: Использован Set вместо массива для O(1) производительности при проверке выбранных товаров
2. **Debounced поиск**: 300ms задержка для оптимизации запросов при вводе в поиск
3. **Транзакции в PostgreSQL**: Все bulk операции используют транзакции для атомарности
4. **Асинхронная индексация**: OpenSearch индексация выполняется в горутинах после основной операции
5. **Двухэтапное удаление**: Требует подтверждения с автосбросом через 5 секунд

#### Окружение:
- Backend работает на порту 3000 (проверено, процесс активен)
- Frontend должен работать на порту 3001
- База данных PostgreSQL, поисковый движок OpenSearch, хранилище MinIO

#### Что НЕ реализовано:
1. **Экспорт товаров на backend** - frontend готов, но нет endpoints
2. **Массовое редактирование товаров** - есть кнопка в UI, но нет функционала
3. **Real-time прогресс** - есть UI для прогресс-бара, но нет WebSocket интеграции
4. **Фильтрация по атрибутам** - нужна интеграция с системой атрибутов категорий

#### Файлы для новой сессии:
- **Основной анализ**: [@memory-bank/analysis/bulk-operations-and-security-analysis.md](../analysis/bulk-operations-and-security-analysis.md)
- **Новые компоненты**: 
  - `/frontend/svetu/src/store/slices/productSlice.ts`
  - `/frontend/svetu/src/components/products/BulkActions.tsx`
  - `/frontend/svetu/src/components/products/ProductCard.tsx`
  - `/frontend/svetu/src/components/products/ProductList.tsx`
  - `/frontend/svetu/src/services/productApi.ts`
- **Измененные файлы backend**:
  - `/backend/internal/domain/models/storefront_product.go`
  - `/backend/internal/storage/postgres/storefront_product.go`
  - `/backend/internal/proj/storefronts/handler/product_handler.go`
  - `/backend/internal/proj/storefronts/service/product_service.go`
  - `/backend/internal/proj/storefronts/module.go`