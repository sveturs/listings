# Исправление переводов категорий и локализации интерфейса

## Проблема
На главной странице маркетплейса категории не реагировали на смену языка страницы и отображались только на сербском языке. Также некоторые элементы интерфейса (заголовок "Категории", фильтр маркетплейса) не переводились.

## Причины
1. В базе данных отсутствовали переводы для категорий
2. Backend не передавал язык из query параметра в контекст при получении категорий
3. В SQL запросах использовался неправильный entity_type ('category' вместо 'marketplace_category')
4. Frontend компоненты использовали захардкоженные строки вместо системы интернационализации

## Решение

### 1. Добавление переводов в базу данных
Создан и выполнен SQL скрипт для добавления переводов всех категорий на три языка (sr, ru, en):
```sql
-- Использование правильного entity_type
WHERE t.entity_type = 'marketplace_category'
AND t.field_name = 'name'
```

### 2. Обновление Backend

#### Файл: `backend/internal/proj/marketplace/handler/categories.go`
```go
// Добавлена передача языка в контекст
func (h *CategoriesHandler) GetCategories(c *fiber.Ctx) error {
    lang := c.Query("lang", "en")
    ctx := context.WithValue(c.UserContext(), "locale", lang)
    categories, err := h.marketplaceService.GetCategories(ctx)
    // ...
}
```

#### Файл: `backend/internal/proj/marketplace/storage/postgres/marketplace.go`
Исправлен entity_type во всех SQL запросах:
```sql
-- Было: WHERE t.entity_type = 'category'
-- Стало: WHERE t.entity_type = 'marketplace_category'
```

### 3. Обновление Frontend

#### Файл: `frontend/svetu/src/components/categories/CategorySidebar.tsx`
```tsx
// Добавлено использование интернационализации
import { useLocale, useTranslations } from 'next-intl';

export default function CategorySidebar() {
  const t = useTranslations('marketplace');
  // Заменены все захардкоженные строки на t('key')
}
```

#### Файл: `frontend/svetu/src/components/marketplace/HomePage.tsx`
```tsx
// Добавлена интернационализация для фильтра маркетплейса
{t('privateListingsOnly')}
{t('privateListingsOnlyDescription')}
```

### 4. Обновление файлов локализации

Добавлены недостающие переводы во все языковые файлы:

#### `frontend/svetu/src/messages/en.json`:
```json
{
  "categories": "Categories",
  "categoriesLoadError": "Error loading categories",
  "clearFilter": "Clear filter",
  "noCategoriesFound": "No categories found",
  "privateListingsOnly": "Private listings only (marketplace)",
  "privateListingsOnlyDescription": "When enabled - shows only marketplace products. When disabled - shows everything (marketplace + storefronts)"
}
```

#### `frontend/svetu/src/messages/sr.json`:
```json
{
  "categories": "Kategorije",
  "categoriesLoadError": "Greška pri učitavanju kategorija",
  "clearFilter": "Obriši filter",
  "noCategoriesFound": "Nema pronađenih kategorija",
  "privateListingsOnly": "Samo privatni oglasi (tržnica)",
  "privateListingsOnlyDescription": "Kada je uključeno - prikazuje samo proizvode tržnice. Kada je isključeno - prikazuje sve (tržnica + izlozi)",
  "placeholder": "Pretražite proizvode i usluge..."
}
```

#### `frontend/svetu/src/messages/ru.json`:
```json
{
  "categories": "Категории",
  "categoriesLoadError": "Ошибка загрузки категорий",
  "clearFilter": "Очистить фильтр",
  "noCategoriesFound": "Категории не найдены",
  "privateListingsOnly": "Только частные объявления (маркетплейс)",
  "privateListingsOnlyDescription": "Когда включено - показывает только товары маркетплейса. Когда выключено - показывает всё (маркетплейс + витрины)"
}
```

## Важные замечания

1. **Entity Type**: Всегда используйте `marketplace_category` для категорий маркетплейса в таблице переводов
2. **Контекст языка**: При добавлении новых эндпоинтов, работающих с переводами, не забывайте передавать язык через контекст
3. **Интернационализация**: Всегда используйте систему интернационализации next-intl вместо захардкоженных строк в компонентах
4. **Кеш категорий**: После изменения переводов в БД ОБЯЗАТЕЛЬНО очистите кеш Redis для обновления данных:
   ```bash
   # Очистка всего кеша Redis
   echo "FLUSHALL" | nc localhost 6379
   
   # Или если Redis в Docker:
   docker exec -it redis redis-cli FLUSHALL
   ```
   Категории кешируются на 6 часов с ключом, зависящим от языка

## Проблема в админке

В админке также была проблема с отображением переводов категорий.

### Исправления для админки

#### 1. Backend - добавление передачи языка в контекст

**Файл**: `backend/internal/proj/marketplace/handler/admin_categories.go`

```go
func (h *AdminCategoriesHandler) GetAllCategories(c *fiber.Ctx) error {
    logger.Info().Str("method", c.Method()).Str("path", c.Path()).Msg("GetAllCategories handler called")

    // Получаем язык из query параметра
    lang := c.Query("lang", "en")
    
    // Создаем контекст с языком
    ctx := context.WithValue(c.UserContext(), "locale", lang)
    
    categories, err := h.marketplaceService.GetAllCategories(ctx)
    // ...
}
```

#### 2. Frontend - исправление отображения текущего языка

**Файл**: `frontend/svetu/src/components/attributes/InlineTranslationEditor.tsx`

```tsx
// Добавлен импорт useLocale
import { useTranslations, useLocale } from 'next-intl';

// В компоненте добавлено получение текущей локали
const locale = useLocale();

// Исправлено отображение текста в соответствии с текущим языком
<span className="text-sm">
  {translations[locale] || translations[LANGUAGES[0]] || t('notTranslated')}
</span>
```

**Файл**: `frontend/svetu/src/app/[locale]/admin/categories/components/CategoryTree.tsx`

```tsx
// Исправлена инициализация переводов из объекта категории
const [translations, setTranslations] = useState<Record<string, string>>({
  en: category.translations?.en || category.name,
  ru: category.translations?.ru || category.name,
  sr: category.translations?.sr || category.name,
});
```

## Тестирование

1. Перезапустить backend сервер
2. Проверить смену языка на главной странице (ru/en/sr)
3. Убедиться, что категории отображаются на выбранном языке
4. Проверить, что все элементы интерфейса (заголовки, фильтры, placeholder'ы) переведены

## Команды для проверки

```bash
# Проверка переводов в БД
psql -d svetubd -c "SELECT id, name, translations FROM marketplace_categories LIMIT 5;"

# Проверка API
curl 'http://localhost:3000/api/v1/marketplace/categories?lang=ru' | jq '.data[0]'
curl 'http://localhost:3000/api/v1/marketplace/categories?lang=en' | jq '.data[0]'
curl 'http://localhost:3000/api/v1/marketplace/categories?lang=sr' | jq '.data[0]'
```