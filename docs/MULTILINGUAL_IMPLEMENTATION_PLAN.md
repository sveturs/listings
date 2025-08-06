# План реализации улучшений мультиязычности Sve Tu Platform

**Дата создания**: 2025-08-03  
**Статус**: Готов к реализации  
**Основано на**: [MULTILINGUAL_AUDIT_REPORT.md](./MULTILINGUAL_AUDIT_REPORT.md)  
**Сложность**: Высокая  
**Время выполнения**: 6-9 недель  
**Команда**: 1-2 разработчика + DevOps

---

## 🎯 Цели проекта

### Основная цель
Устранить все критические проблемы мультиязычности и привести систему к enterprise-уровню поддержки 3 языков (ru, sr, en).

### Ключевые результаты (KPI)
- ✅ **Покрытие переводами**: с 95% до 99%+ для всех языков
- ✅ **Производительность поиска**: +40% релевантности для неанглийского контента  
- ✅ **Консистентность**: 100% согласованность конфигурации языков
- ✅ **Hardcoded строки**: 0 нелокализованных текстов в production коде
- ✅ **API локализация**: полная поддержка Accept-Language/Content-Language

---

## 📋 Структура проблем и решений

### 🔴 КРИТИЧЕСКИЕ (Блокирующие, решать первыми)

#### Проблема #1: OpenSearch без мультиязычности
**Влияние**: Главная страница маркетплейса показывает неточные результаты поиска  
**Файлы**: `backend/internal/proj/marketplace/storage/opensearch/`  
**Техническая суть**: Индекс не содержит переводы из БД, нет language analyzers

<details>
<summary>📁 Детальное решение для OpenSearch</summary>

**Шаг 1.1: Анализ текущего состояния**
```bash
# Проверить текущую структуру индекса
curl -X GET "http://localhost:9200/marketplace/_mapping" | jq '.marketplace.mappings.properties' > /tmp/current_mapping.json

# Проверить настройки анализаторов
curl -X GET "http://localhost:9200/marketplace/_settings" | jq '.marketplace.settings.index.analysis' > /tmp/current_analysis.json
```

**Шаг 1.2: Создание backup и новой конфигурации**
```bash
# Backup текущего индекса
curl -X POST "http://localhost:9200/_reindex" -H "Content-Type: application/json" -d '{
  "source": { "index": "marketplace" },
  "dest": { "index": "marketplace_backup_20250803" }
}'
```

**Шаг 1.3: Обновление маппинга индекса**
Создать файл: `backend/opensearch/marketplace_mapping_v2.json`
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "serbian_analyzer": {
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding", "serbian_stemmer"]
        },
        "russian_analyzer": {
          "tokenizer": "standard", 
          "filter": ["lowercase", "russian_stemmer"]
        },
        "english_analyzer": {
          "tokenizer": "standard",
          "filter": ["lowercase", "english_stemmer"] 
        },
        "multilingual_search": {
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding"]
        }
      },
      "filter": {
        "serbian_stemmer": {
          "type": "stemmer",
          "language": "light_nynorsk"
        },
        "russian_stemmer": {
          "type": "stemmer", 
          "language": "russian"
        },
        "english_stemmer": {
          "type": "stemmer",
          "language": "english"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "multilingual_search",
        "fields": {
          "ru": {"type": "text", "analyzer": "russian_analyzer"},
          "sr": {"type": "text", "analyzer": "serbian_analyzer"}, 
          "en": {"type": "text", "analyzer": "english_analyzer"}
        }
      },
      "description": {
        "type": "text",
        "analyzer": "multilingual_search",
        "fields": {
          "ru": {"type": "text", "analyzer": "russian_analyzer"},
          "sr": {"type": "text", "analyzer": "serbian_analyzer"},
          "en": {"type": "text", "analyzer": "english_analyzer"}
        }
      },
      "translations": {
        "type": "nested",
        "properties": {
          "language": {"type": "keyword"},
          "field_name": {"type": "keyword"},
          "translated_text": {
            "type": "text",
            "analyzer": "multilingual_search"
          }
        }
      },
      "original_language": {"type": "keyword"},
      "supported_languages": {"type": "keyword"}
    }
  }
}
```

**Шаг 1.4: Модификация Go кода для включения переводов**
Файл: `backend/internal/proj/marketplace/storage/opensearch/marketplace.go`

Добавить в структуру документа:
```go
type MarketplaceDocument struct {
    // ... существующие поля
    Translations []Translation `json:"translations"`
    SupportedLanguages []string `json:"supported_languages"`
}

type Translation struct {
    Language string `json:"language"`
    FieldName string `json:"field_name"`
    TranslatedText string `json:"translated_text"`
}
```

Модифицировать функцию индексирования:
```go
func (r *MarketplaceRepository) IndexListing(ctx context.Context, listing *domain.MarketplaceListing) error {
    // Загрузить переводы из БД
    translations, err := r.getListingTranslations(ctx, listing.ID)
    if err != nil {
        return err
    }
    
    doc := MarketplaceDocument{
        // ... заполнить существующие поля
        Translations: translations,
        SupportedLanguages: extractSupportedLanguages(translations),
    }
    
    // Индексировать с переводами
    return r.indexDocument(ctx, doc)
}

func (r *MarketplaceRepository) getListingTranslations(ctx context.Context, listingID int) ([]Translation, error) {
    query := `
        SELECT language, field_name, translated_text 
        FROM translations 
        WHERE entity_type = 'listing' AND entity_id = $1
    `
    
    rows, err := r.db.QueryContext(ctx, query, listingID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var translations []Translation
    for rows.Next() {
        var t Translation
        err := rows.Scan(&t.Language, &t.FieldName, &t.TranslatedText)
        if err != nil {
            return nil, err
        }
        translations = append(translations, t)
    }
    
    return translations, nil
}
```

**Шаг 1.5: Обновление поискового запроса**
Файл: `backend/internal/proj/marketplace/storage/opensearch/search.go`

```go
func (r *MarketplaceRepository) Search(ctx context.Context, params SearchParams) (*SearchResult, error) {
    language := params.Language // получить из контекста запроса
    
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "should": []map[string]interface{}{
                    // Поиск по оригинальным полям
                    {
                        "multi_match": map[string]interface{}{
                            "query":  params.Query,
                            "fields": []string{"title^2", "description"},
                        },
                    },
                    // Поиск по переводам для конкретного языка
                    {
                        "nested": map[string]interface{}{
                            "path": "translations",
                            "query": map[string]interface{}{
                                "bool": map[string]interface{}{
                                    "must": []map[string]interface{}{
                                        {
                                            "term": map[string]interface{}{
                                                "translations.language": language,
                                            },
                                        },
                                        {
                                            "match": map[string]interface{}{
                                                "translations.translated_text": params.Query,
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    // Поиск с language-specific анализатором
                    {
                        "multi_match": map[string]interface{}{
                            "query":    params.Query,
                            "fields":   []string{fmt.Sprintf("title.%s^2", language), fmt.Sprintf("description.%s", language)},
                            "analyzer": fmt.Sprintf("%s_analyzer", language),
                        },
                    },
                },
                "minimum_should_match": 1,
            },
        },
    }
    
    // Выполнить поиск
    return r.executeSearch(ctx, query)
}
```

**Шаг 1.6: Скрипт реиндексации**
Создать файл: `backend/scripts/reindex_with_translations.go`
```go
package main

import (
    "context"
    "fmt"
    "log"
    // ... импорты
)

func main() {
    // Подключиться к БД и OpenSearch
    db := connectDB()
    es := connectOpenSearch()
    
    // Получить все листинги
    listings, err := getAllListings(db)
    if err != nil {
        log.Fatal(err)
    }
    
    // Реиндексировать с переводами
    for _, listing := range listings {
        translations := getListingTranslations(db, listing.ID)
        doc := buildDocumentWithTranslations(listing, translations)
        
        err := indexDocument(es, doc)
        if err != nil {
            log.Printf("Failed to index listing %d: %v", listing.ID, err)
            continue
        }
        
        fmt.Printf("Reindexed listing %d with %d translations\n", listing.ID, len(translations))
    }
}
```

**Команды для выполнения:**
```bash
# 1. Создать новый индекс с правильным маппингом
curl -X PUT "http://localhost:9200/marketplace_v2" -H "Content-Type: application/json" -d @backend/opensearch/marketplace_mapping_v2.json

# 2. Запустить скрипт реиндексации
cd backend && go run scripts/reindex_with_translations.go

# 3. Переключить алиас на новый индекс  
curl -X POST "http://localhost:9200/_aliases" -H "Content-Type: application/json" -d '{
  "actions": [
    {"remove": {"index": "marketplace", "alias": "marketplace_current"}},
    {"add": {"index": "marketplace_v2", "alias": "marketplace_current"}}
  ]
}'

# 4. Обновить код для использования нового алиаса
# Заменить "marketplace" на "marketplace_current" в коде
```

</details>

---

#### Проблема #2: Конфликт конфигурации языков по умолчанию
**Влияние**: Inconsistent UX, SEO проблемы  
**Файлы**: `frontend/svetu/src/i18n/config.ts`  
**Техническая суть**: Frontend использует `sr` по умолчанию, но `ru` наиболее полный

<details>
<summary>📁 Детальное решение для конфигурации языков</summary>

**Шаг 2.1: Анализ статистики использования**
```bash
# Проверить статистику переводов
PGPASSWORD=password psql -h localhost -U postgres -d svetubd -c "
SELECT language, COUNT(*) as translation_count,
       ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
FROM translations 
GROUP BY language 
ORDER BY translation_count DESC;
"

# Результат показывает: ru - 541 (35%), sr - 508 (33%), en - 504 (32%)
```

**Шаг 2.2: Обновление конфигурации Frontend**
Файл: `frontend/svetu/src/i18n/config.ts`
```typescript
// СТАРАЯ конфигурация
export const defaultLocale: Locale = 'sr'; // Проблема: sr неполный

// НОВАЯ конфигурация (обоснованная данными)
export const locales = ['ru', 'sr', 'en'] as const; // ru первый = приоритетный
export type Locale = (typeof locales)[number];

// Стратегия по доменам
export const getDefaultLocale = (): Locale => {
  if (typeof window !== 'undefined') {
    const hostname = window.location.hostname;
    
    // Для .rs домена приоритет: sr -> ru -> en  
    if (hostname.endsWith('.rs')) {
      return 'sr';
    }
    // Для .ru домена: ru -> sr -> en
    if (hostname.endsWith('.ru')) {
      return 'ru';  
    }
    // Для международных доменов: en -> ru -> sr
    return 'en';
  }
  
  // Server-side default (самый полный)
  return 'ru';
};

export const defaultLocale: Locale = getDefaultLocale();

// Fallback цепочка для отсутствующих переводов
export const localeFallbacks: Record<Locale, Locale[]> = {
  'sr': ['ru', 'en'], // sr -> ru -> en
  'ru': ['en', 'sr'], // ru -> en -> sr  
  'en': ['ru', 'sr'], // en -> ru -> sr
};

export const i18n = {
  locales,
  defaultLocale,
  localeFallbacks,
  localeDetection: {
    enabled: true,
    cookieName: 'locale-preference',
    cookieMaxAge: 365 * 24 * 60 * 60, // 1 год
    
    // Новые настройки детекции
    sources: [
      'cookie',           // 1. Сохранённые предпочтения  
      'header',           // 2. Accept-Language header
      'domain',           // 3. По домену (.rs, .ru)
      'default'           // 4. Fallback
    ],
  },
} as const;

// Обновленная функция загрузки сообщений с fallback
export async function getLocaleMessages(locale: Locale): Promise<Record<string, any>> {
  try {
    const messages = await import(`../messages/${locale}.json`).then(m => m.default);
    
    // Если это не основной язык, подгружаем fallback для недостающих ключей
    if (locale !== 'ru') {
      const fallbackLocale = localeFallbacks[locale][0];
      const fallbackMessages = await import(`../messages/${fallbackLocale}.json`).then(m => m.default);
      
      // Мерджим с fallback (приоритет у основного языка)
      return deepMerge(fallbackMessages, messages);
    }
    
    return messages;
  } catch (error) {
    console.error(`Failed to load messages for locale ${locale}:`, error);
    
    // Fallback к русскому если загрузка не удалась
    if (locale !== 'ru') {
      return import('../messages/ru.json').then(m => m.default);
    }
    
    throw error;
  }
}

// Вспомогательная функция deep merge
function deepMerge(target: any, source: any): any {
  const result = { ...target };
  
  for (const key in source) {
    if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
      result[key] = deepMerge(result[key] || {}, source[key]);
    } else {
      result[key] = source[key];
    }
  }
  
  return result;
}
```

**Шаг 2.3: Обновление роутинга**
Файл: `frontend/svetu/src/i18n/routing.ts`
```typescript
import { createNavigation } from 'next-intl/navigation';
import { defineRouting } from 'next-intl/routing';
import { i18n, getDefaultLocale } from './config';

export const routing = defineRouting({
  locales: i18n.locales,
  defaultLocale: getDefaultLocale(),
  
  // Новые настройки роутинга
  localePrefix: {
    mode: 'as-needed',
    prefixes: {
      // Не показывать префикс для языка по умолчанию на соответствующих доменах
      'sr': '', // svetu.rs/ вместо svetu.rs/sr/
      'ru': '/ru', // всегда показывать /ru
      'en': '/en', // всегда показывать /en
    }
  },
  
  // Альтернативные домены
  domains: [
    {
      domain: 'svetu.rs',
      defaultLocale: 'sr',
      locales: ['sr', 'ru', 'en']
    },
    {
      domain: 'svetu.ru', 
      defaultLocale: 'ru',
      locales: ['ru', 'sr', 'en']
    }
  ]
});

export const { Link, redirect, usePathname, useRouter, getPathname } =
  createNavigation(routing);
```

**Шаг 2.4: Обновление Layout**
Файл: `frontend/svetu/src/app/[locale]/layout.tsx`
```typescript
// Добавить валидацию locale
export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  
  // Валидация поддерживаемых локалей
  if (!i18n.locales.includes(locale as Locale)) {
    notFound();
  }
  
  const t = await getTranslations({ locale, namespace: 'metadata' });
  
  // Обновленные hreflang ссылки
  const alternateUrls = i18n.locales.reduce((acc, loc) => {
    acc[loc] = `https://svetu.rs/${loc === getDefaultLocale() ? '' : loc}`;
    return acc; 
  }, {} as Record<string, string>);

  return {
    title: t('title'),
    description: t('description'),
    alternates: {
      canonical: `https://svetu.rs/${locale === getDefaultLocale() ? '' : locale}`,
      languages: alternateUrls,
    },
    openGraph: {
      locale: locale,
      alternateLocale: i18n.locales.filter(l => l !== locale),
    },
  };
}
```

**Шаг 2.5: Middleware для детекции языка**
Файл: `frontend/svetu/src/middleware.ts`  
```typescript
import createMiddleware from 'next-intl/middleware';
import { routing } from './src/i18n/routing';

export default createMiddleware({
  ...routing,
  
  // Кастомная логика детекции
  localeDetection: {
    enabled: true,
    cookieName: 'NEXT_LOCALE',
  },
  
  // Перенаправления для SEO
  redirects: [
    // Перенаправить старые URL с /sr/ на корень для .rs домена
    {
      source: '/sr/:path*',
      destination: '/:path*',
      permanent: true,
      has: [{ type: 'host', value: 'svetu.rs' }],
    },
  ],
});

export const config = {
  matcher: [
    // Исключить API и статические файлы
    '/((?!api|_next|_vercel|.*\\..*).*)',
  ],
};
```

</details>

---

#### Проблема #3: Отсутствует раздел "storefront" в сербском языке
**Влияние**: Ошибки в UI витрин для сербских пользователей  
**Файлы**: `frontend/svetu/src/messages/sr.json`  
**Техническая суть**: В sr.json отсутствует весь раздел "storefront"

<details>
<summary>📁 Детальное решение для отсутствующего раздела</summary>

**Шаг 3.1: Анализ отсутствующих ключей**
```bash
# Сравнить ключи между языками
cd frontend/svetu/src/messages

# Найти различия в структуре
jq -r 'keys[]' en.json | sort > /tmp/en_keys.txt
jq -r 'keys[]' ru.json | sort > /tmp/ru_keys.txt  
jq -r 'keys[]' sr.json | sort > /tmp/sr_keys.txt

# Показать отсутствующие в sr.json
echo "Отсутствует в sr.json:"
comm -23 /tmp/en_keys.txt /tmp/sr_keys.txt

# Результат: storefront
```

**Шаг 3.2: Извлечение раздела storefront из en.json**
```bash
cd frontend/svetu/src/messages

# Извлечь раздел storefront
jq '.storefront' en.json > /tmp/storefront_en.json

# Проверить размер раздела
echo "Количество ключей в storefront:"
jq 'keys | length' /tmp/storefront_en.json
# Результат покажет количество ключей для перевода
```

**Шаг 3.3: Создание переводов на сербский**
Создать файл: `/tmp/storefront_sr_translations.json`
```json
{
  "storefront": {
    "title": "Продавница",
    "description": "Опис продавнице", 
    "create": "Направи продавницу",
    "edit": "Измени продавницу",
    "delete": "Обриши продавницу",
    "publish": "Објави",
    "unpublish": "Скини са објаве",
    "draft": "Нацрт",
    "active": "Активна",
    "inactive": "Неактивна",
    
    "dashboard": {
      "title": "Контролна табла",
      "overview": "Преглед",
      "statistics": "Статистике",
      "orders": "Поруџбине", 
      "products": "Производи",
      "settings": "Подешавања"
    },
    
    "products": {
      "title": "Производи",
      "add": "Додај производ", 
      "edit": "Измени производ",
      "delete": "Обриши производ",
      "publish": "Објави производ",
      "unpublish": "Скини са објаве",
      "stock": "Залиха",
      "price": "Цена",
      "description": "Опис",
      "images": "Слике",
      "category": "Категорија",
      "attributes": "Атрибути"
    },
    
    "orders": {
      "title": "Поруџбине",
      "new": "Нове поруџбине",
      "processing": "У обради", 
      "shipped": "Послато",
      "delivered": "Испоручено",
      "cancelled": "Отказано",
      "refunded": "Враћено",
      "total": "Укупно",
      "customer": "Купац",
      "date": "Датум",
      "status": "Статус"
    },
    
    "settings": {
      "title": "Подешавања",
      "general": "Опште",
      "appearance": "Изглед", 
      "payment": "Плаћање",
      "shipping": "Достава",
      "notifications": "Обавештења",
      "integrations": "Интеграције"
    },
    
    "public": {
      "welcome": "Добродошли у нашу продавницу",
      "featured": "Издвојени производи",
      "categories": "Категорије", 
      "search": "Претрага",
      "cart": "Корпа",
      "checkout": "Наплата",
      "account": "Налог"
    }
  }
}
```

**Шаг 3.4: Автоматическое объединение с sr.json**
```bash
cd frontend/svetu/src/messages

# Backup текущего файла
cp sr.json sr.json.backup

# Объединить файлы
jq -s '.[0] * .[1]' sr.json /tmp/storefront_sr_translations.json > sr.json.new

# Проверить результат
echo "Проверка объединения:"
jq 'has("storefront")' sr.json.new
# Должно вернуть: true

# Проверить количество ключей в новом разделе
jq '.storefront | keys | length' sr.json.new

# Заменить файл
mv sr.json.new sr.json
```

**Шаг 3.5: Валидация новых переводов**
Создать скрипт: `frontend/svetu/scripts/validate_translations.js`
```javascript
const fs = require('fs');
const path = require('path');

const messagesDir = path.join(__dirname, '../src/messages');
const languages = ['en', 'ru', 'sr'];

function loadMessages(lang) {
  const filePath = path.join(messagesDir, `${lang}.json`);
  return JSON.parse(fs.readFileSync(filePath, 'utf8'));
}

function getAllKeys(obj, prefix = '') {
  let keys = [];
  
  for (const [key, value] of Object.entries(obj)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    
    if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      keys.push(...getAllKeys(value, fullKey));
    } else {
      keys.push(fullKey);
    }
  }
  
  return keys;
}

function validateTranslations() {
  const messages = {};
  const allKeys = {};
  
  // Загрузить все языки
  for (const lang of languages) {
    messages[lang] = loadMessages(lang);
    allKeys[lang] = new Set(getAllKeys(messages[lang]));
  }
  
  // Найти отсутствующие ключи
  const baseKeys = allKeys['en']; // английский как база
  
  for (const lang of languages) {
    if (lang === 'en') continue;
    
    const missing = [...baseKeys].filter(key => !allKeys[lang].has(key));
    const extra = [...allKeys[lang]].filter(key => !baseKeys.has(key));
    
    console.log(`\n=== ${lang.toUpperCase()} ===`);
    console.log(`Всего ключей: ${allKeys[lang].size}`);
    console.log(`Отсутствует: ${missing.length}`);
    console.log(`Лишних: ${extra.length}`);
    
    if (missing.length > 0) {
      console.log('\nОтсутствующие ключи:');
      missing.slice(0, 10).forEach(key => console.log(`  - ${key}`));
      if (missing.length > 10) {
        console.log(`  ... и ещё ${missing.length - 10}`);
      }
    }
    
    if (extra.length > 0) {
      console.log('\nЛишние ключи:');
      extra.slice(0, 5).forEach(key => console.log(`  + ${key}`));
      if (extra.length > 5) {
        console.log(`  ... и ещё ${extra.length - 5}`);
      }
    }
  }
  
  // Проверить конкретно storefront
  console.log('\n=== STOREFRONT SECTION ===');
  for (const lang of languages) {
    const hasStorefront = 'storefront' in messages[lang];
    const storefrontKeys = hasStorefront ? getAllKeys(messages[lang].storefront, 'storefront').length : 0;
    
    console.log(`${lang}: ${hasStorefront ? '✅' : '❌'} (${storefrontKeys} ключей)`);
  }
}

validateTranslations();
```

**Команды для выполнения:**
```bash
# 1. Создать переводы и объединить
cd frontend/svetu
node scripts/validate_translations.js  # до изменений

# 2. Выполнить объединение файлов (команды выше)

# 3. Проверить результат
node scripts/validate_translations.js  # после изменений

# 4. Тест в браузере
yarn dev -p 3001
# Открыть http://localhost:3001/sr/storefronts и проверить отсутствие ошибок
```

</details>

---

### 🟡 ВЫСОКИЙ ПРИОРИТЕТ (Влияют на UX, решать после критических)

#### Проблема #4: 20+ файлов с hardcoded строками
**Влияние**: Нелокализованный контент для пользователей  
**Файлы**: Множество компонентов в `frontend/svetu/src/components/`  
**Техническая суть**: Строки захардкожены в коде вместо использования переводов

<details>
<summary>📁 Детальное решение для hardcoded строк</summary>

**Шаг 4.1: Полный аудит hardcoded строк**
```bash
cd frontend/svetu/src

# Найти все файлы с кириллическими строками
grep -r "\"[А-Яа-я][А-Яа-я ]{3,}\"" components/ --include="*.tsx" -n > /tmp/hardcoded_cyrillic.txt

# Найти файлы с латинскими строками (исключая className, testId и т.д.)
grep -r "\"[A-Z][a-zA-Z ]{4,}\"" components/ --include="*.tsx" -n | grep -v "className\|testId\|id=\|data-\|aria-" > /tmp/hardcoded_latin.txt

# Объединить результаты
cat /tmp/hardcoded_cyrillic.txt /tmp/hardcoded_latin.txt | sort | uniq > /tmp/all_hardcoded.txt

echo "Найдено hardcoded строк:"
wc -l /tmp/all_hardcoded.txt
```

**Шаг 4.2: Анализ найденных файлов**
Создать скрипт: `frontend/svetu/scripts/analyze_hardcoded.js`
```javascript
const fs = require('fs');
const path = require('path');

// Читаем результаты grep
const hardcodedFile = '/tmp/all_hardcoded.txt';
const lines = fs.readFileSync(hardcodedFile, 'utf8').split('\n').filter(Boolean);

const fileStats = {};
const translations = {};

lines.forEach(line => {
  const [filePath, lineNum, content] = line.split(':');
  
  if (!fileStats[filePath]) {
    fileStats[filePath] = { count: 0, issues: [] };
  }
  
  fileStats[filePath].count++;
  fileStats[filePath].issues.push({
    line: parseInt(lineNum),
    content: content.trim()
  });
  
  // Извлечь строку для перевода
  const match = content.match(/"([^"]+)"/);
  if (match) {
    const text = match[1];
    if (text.length > 3 && !text.includes('className') && !text.includes('test')) {
      translations[text] = text; // добавить в список для перевода
    }
  }
});

// Сортировать файлы по количеству проблем
const sortedFiles = Object.entries(fileStats)
  .sort(([,a], [,b]) => b.count - a.count)
  .slice(0, 20); // топ 20

console.log('=== ТОП ФАЙЛОВ С HARDCODED СТРОКАМИ ===');
sortedFiles.forEach(([file, stats]) => {
  console.log(`${stats.count} проблем: ${file}`);
});

console.log('\n=== НАЙДЕННЫЕ СТРОКИ ДЛЯ ПЕРЕВОДА ===');
Object.keys(translations).slice(0, 30).forEach(text => {
  console.log(`"${text}"`);
});

// Сохранить для использования
fs.writeFileSync('/tmp/translation_candidates.json', JSON.stringify(translations, null, 2));
fs.writeFileSync('/tmp/files_to_fix.json', JSON.stringify(fileStats, null, 2));
```

**Шаг 4.3: Приоритизация файлов для исправления**
```bash
node scripts/analyze_hardcoded.js

# Выберем топ-5 файлов для первой итерации
echo "Приоритетные файлы для исправления:"
head -5 /tmp/files_to_fix.json
```

**Шаг 4.4: Пример исправления конкретного файла**
Возьмем файл с максимальным количеством проблем, например `components/GIS/demo/MapDemo.tsx`:

**Исходный код (проблемный):**
```typescript
// components/GIS/demo/MapDemo.tsx
const demoData = [
  {
    id: 1,
    title: 'Хостел "Центр"',  // ПРОБЛЕМА: hardcoded
    description: 'Уютный хостел в центре города', // ПРОБЛЕМА: hardcoded
    // ...
  }
];

return (
  <div>
    <h2>Демо карты</h2> {/* ПРОБЛЕМА: hardcoded */}
    <p>Это пример интерактивной карты</p> {/* ПРОБЛЕМА: hardcoded */}
  </div>
);
```

**Исправленный код:**
```typescript
// components/GIS/demo/MapDemo.tsx
import { useTranslations } from 'next-intl';

export default function MapDemo() {
  const t = useTranslations('demo.map');
  
  // Данные теперь приходят из переводов
  const demoData = [
    {
      id: 1,
      title: t('locations.hostel_center.title'),
      description: t('locations.hostel_center.description'),
      // ...
    }
  ];

  return (
    <div>
      <h2>{t('title')}</h2>
      <p>{t('description')}</p>
    </div>
  );
}
```

**Добавить в файлы переводов:**
```json
// messages/ru.json
{
  "demo": {
    "map": {
      "title": "Демо карты",
      "description": "Это пример интерактивной карты",
      "locations": {
        "hostel_center": {
          "title": "Хостел \"Центр\"",
          "description": "Уютный хостел в центре города"
        }
      }
    }
  }
}

// messages/en.json  
{
  "demo": {
    "map": {
      "title": "Map Demo", 
      "description": "This is an example of interactive map",
      "locations": {
        "hostel_center": {
          "title": "Hostel \"Center\"",
          "description": "Cozy hostel in the city center"
        }
      }
    }
  }
}

// messages/sr.json
{
  "demo": {
    "map": {
      "title": "Демо мапе",
      "description": "Ово је пример интерактивне мапе", 
      "locations": {
        "hostel_center": {
          "title": "Хостел \"Центар\"",
          "description": "Уютан хостел у центру града"
        }
      }
    }
  }
}
```

**Шаг 4.5: Автоматизация исправлений**
Создать скрипт: `frontend/svetu/scripts/fix_hardcoded.js`
```javascript
const fs = require('fs');
const path = require('path');

class HardcodedFixer {
  constructor() {
    this.replacements = new Map();
    this.translationKeys = new Map();
  }
  
  // Анализировать файл и найти кандидатов для замены
  analyzeFile(filePath) {
    const content = fs.readFileSync(filePath, 'utf8');
    const issues = [];
    
    // Найти все строки в кавычках
    const regex = /"([^"]{4,})"/g;
    let match;
    
    while ((match = regex.exec(content)) !== null) {
      const text = match[1];
      
      // Пропустить технические строки
      if (this.shouldIgnore(text)) continue;
      
      issues.push({
        text: text,
        start: match.index,
        end: match.index + match[0].length,
        line: this.getLineNumber(content, match.index)
      });
    }
    
    return issues;
  }
  
  shouldIgnore(text) {
    const ignorePatterns = [
      /^[a-z-]+$/, // CSS классы
      /^\w+Id$/, // ID атрибуты  
      /^data-/, // data атрибуты
      /^aria-/, // aria атрибуты
      /^\d+$/, // только цифры
      /^[#@]/, // хэштеги, упоминания
      /^\w+\.\w+$/, // файлы
      /^https?:\/\//, // URL
    ];
    
    return ignorePatterns.some(pattern => pattern.test(text));
  }
  
  getLineNumber(content, index) {
    return content.substring(0, index).split('\n').length;  
  }
  
  // Генерировать ключ перевода из текста
  generateTranslationKey(text, namespace = 'common') {
    const key = text
      .toLowerCase()
      .replace(/[^\w\s]/g, '') // убрать пунктуацию
      .replace(/\s+/g, '_') // пробелы в подчеркивания
      .substring(0, 50); // ограничить длину
      
    return `${namespace}.${key}`;
  }
  
  // Исправить конкретный файл
  fixFile(filePath, namespace) {
    const content = fs.readFileSync(filePath, 'utf8');
    const issues = this.analyzeFile(filePath);
    
    if (issues.length === 0) return null;
    
    let newContent = content;
    let offset = 0;
    const addedTranslations = {};
    
    // Добавить импорт useTranslations если его нет
    if (!content.includes('useTranslations')) {
      const importLine = "import { useTranslations } from 'next-intl';\n";
      newContent = importLine + newContent;
      offset += importLine.length;
    }
    
    // Заменить строки на вызовы t()
    issues.reverse(); // идем с конца чтобы не сбить индексы
    
    issues.forEach(issue => {
      const translationKey = this.generateTranslationKey(issue.text, namespace); 
      const replacement = `{t('${translationKey.replace(namespace + '.', '')}')}`;
      
      const start = issue.start + offset;
      const end = issue.end + offset;
      
      newContent = 
        newContent.substring(0, start) + 
        replacement + 
        newContent.substring(end);
        
      // Сохранить для добавления в файлы переводов
      addedTranslations[translationKey] = issue.text;
    });
    
    // Добавить хук useTranslations в компонент
    if (issues.length > 0 && !content.includes('useTranslations(')) {
      const hookLine = `  const t = useTranslations('${namespace}');\n`;
      
      // Найти где вставить (после объявления компонента)
      const functionMatch = newContent.match(/(export\s+(?:default\s+)?function\s+\w+[^{]*\{)/);
      if (functionMatch) {
        const insertPos = functionMatch.index + functionMatch[1].length;
        newContent = 
          newContent.substring(0, insertPos) + 
          '\n' + hookLine + 
          newContent.substring(insertPos);
      }
    }
    
    return {
      content: newContent,
      translations: addedTranslations,
      issuesFixed: issues.length
    };
  }
}

// Использование
const fixer = new HardcodedFixer();

// Исправить топ файлы
const topFiles = [
  'components/GIS/demo/MapDemo.tsx',
  // добавить другие файлы из анализа
];

topFiles.forEach(filePath => {
  const fullPath = path.join(__dirname, '../src', filePath);
  
  if (!fs.existsSync(fullPath)) {
    console.log(`Файл не найден: ${filePath}`);
    return;
  }
  
  const namespace = 'demo'; // или определить динамически
  const result = fixer.fixFile(fullPath, namespace);
  
  if (result) {
    // Сохранить исправленный файл
    fs.writeFileSync(fullPath + '.fixed', result.content);
    
    console.log(`Исправлен ${filePath}: ${result.issuesFixed} проблем`);
    console.log('Добавленные переводы:', result.translations);
  }
});
```

**Команды для выполнения (поэтапно):**
```bash
# 1. Проанализировать все файлы
cd frontend/svetu
node scripts/analyze_hardcoded.js

# 2. Исправить топ-5 файлов вручную или скриптом
node scripts/fix_hardcoded.js

# 3. Проверить результаты
for file in components/**/*.tsx.fixed; do
  echo "=== $file ==="
  diff "${file%.fixed}" "$file" | head -10
done

# 4. Если все ОК, заменить оригинальные файлы
for file in components/**/*.tsx.fixed; do
  mv "$file" "${file%.fixed}"
done

# 5. Добавить новые переводы в JSON файлы (вручную)
# 6. Тестировать
yarn dev -p 3001
```

</details>

---

#### Проблема #5: Backend API без локализации
**Влияние**: API не учитывает язык пользователя  
**Файлы**: `backend/internal/middleware/`, `backend/internal/proj/*/handler/`  
**Техническая суть**: Нет обработки Accept-Language, нет Content-Language headers

<details>
<summary>📁 Детальное решение для локализации Backend</summary>

**Шаг 5.1: Создание middleware для языка**
Файл: `backend/internal/middleware/locale.go`
```go
package middleware

import (
    "context"
    "strings"
    
    "github.com/gofiber/fiber/v2"
    "golang.org/x/text/language"
)

type LocaleConfig struct {
    SupportedLocales []string
    DefaultLocale    string
    Sources          []string // cookie, header, query, domain
}

func DefaultLocaleConfig() LocaleConfig {
    return LocaleConfig{
        SupportedLocales: []string{"ru", "sr", "en"},
        DefaultLocale:    "ru",
        Sources:          []string{"cookie", "header", "query", "domain"},
    }
}

func New(config ...LocaleConfig) fiber.Handler {
    cfg := DefaultLocaleConfig()
    if len(config) > 0 {
        cfg = config[0]
    }
    
    return func(c *fiber.Ctx) error {
        locale := detectLocale(c, cfg)
        
        // Set locale in context for handlers
        c.Locals("locale", locale)
        c.Locals("language", locale) // alias
        
        // Set Content-Language header in response
        c.Set("Content-Language", locale)
        
        return c.Next()
    }
}

func detectLocale(c *fiber.Ctx, cfg LocaleConfig) string {
    for _, source := range cfg.Sources {
        switch source {
        case "query":
            if locale := c.Query("lang"); locale != "" {
                if isSupported(locale, cfg.SupportedLocales) {
                    return locale
                }
            }
            
        case "cookie":
            if locale := c.Cookies("locale-preference"); locale != "" {
                if isSupported(locale, cfg.SupportedLocales) {
                    return locale
                }
            }
            
        case "header":
            if locale := parseAcceptLanguage(c.Get("Accept-Language"), cfg.SupportedLocales); locale != "" {
                return locale
            }
            
        case "domain":
            if locale := detectFromDomain(c.Hostname(), cfg.SupportedLocales); locale != "" {
                return locale
            }
        }
    }
    
    return cfg.DefaultLocale
}

func parseAcceptLanguage(header string, supported []string) string {
    if header == "" {
        return ""
    }
    
    // Parse Accept-Language header
    // "ru-RU,ru;q=0.9,en;q=0.8,sr;q=0.7"
    parts := strings.Split(header, ",")
    
    for _, part := range parts {
        // Remove quality factor
        locale := strings.Split(strings.TrimSpace(part), ";")[0]
        
        // Extract primary language
        primaryLang := strings.Split(locale, "-")[0]
        
        if isSupported(primaryLang, supported) {
            return primaryLang
        }
    }
    
    return ""
}

func detectFromDomain(hostname string, supported []string) string {
    switch {
    case strings.HasSuffix(hostname, ".rs"):
        if isSupported("sr", supported) {
            return "sr"
        }
    case strings.HasSuffix(hostname, ".ru"):
        if isSupported("ru", supported) {
            return "ru"
        }
    }
    
    return ""
}

func isSupported(locale string, supported []string) bool {
    for _, s := range supported {
        if s == locale {
            return true
        }
    }
    return false
}

// Helper function to get locale from context
func GetLocale(c *fiber.Ctx) string {
    if locale, ok := c.Locals("locale").(string); ok {
        return locale
    }
    return "ru" // fallback
}
```

**Шаг 5.2: Создание сервиса переводов**
Файл: `backend/internal/services/translation_service.go`
```go
package services

import (
    "context"
    "database/sql"
    "fmt"
    "sync"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type TranslationService struct {
    db    *sql.DB
    cache *redis.Client
    mu    sync.RWMutex
    
    // In-memory cache для часто используемых переводов
    localCache map[string]map[string]string // [language][key]value
}

type Translation struct {
    EntityType     string `json:"entity_type"`
    EntityID       int    `json:"entity_id"`
    Language       string `json:"language"`
    FieldName      string `json:"field_name"`
    TranslatedText string `json:"translated_text"`
    IsVerified     bool   `json:"is_verified"`
}

func NewTranslationService(db *sql.DB, cache *redis.Client) *TranslationService {
    return &TranslationService{
        db:         db,
        cache:      cache,
        localCache: make(map[string]map[string]string),
    }
}

// GetTranslation получает перевод для конкретной сущности
func (s *TranslationService) GetTranslation(ctx context.Context, entityType string, entityID int, fieldName, language string) (string, error) {
    cacheKey := fmt.Sprintf("translation:%s:%d:%s:%s", entityType, entityID, fieldName, language)
    
    // 1. Попробовать Redis cache
    if s.cache != nil {
        if cached := s.cache.Get(ctx, cacheKey); cached.Err() == nil {
            return cached.Val(), nil
        }
    }
    
    // 2. Попробовать БД
    query := `
        SELECT translated_text 
        FROM translations 
        WHERE entity_type = $1 AND entity_id = $2 AND field_name = $3 AND language = $4
        LIMIT 1
    `
    
    var translatedText string
    err := s.db.QueryRowContext(ctx, query, entityType, entityID, fieldName, language).Scan(&translatedText)
    
    if err == sql.ErrNoRows {
        // Fallback к оригинальному тексту или другому языку
        return s.getFallbackTranslation(ctx, entityType, entityID, fieldName, language)
    }
    
    if err != nil {
        return "", err
    }
    
    // Кэшировать результат
    if s.cache != nil {
        s.cache.Set(ctx, cacheKey, translatedText, time.Hour)
    }
    
    return translatedText, nil
}

func (s *TranslationService) getFallbackTranslation(ctx context.Context, entityType string, entityID int, fieldName, language string) (string, error) {
    // Порядок fallback: ru -> en -> sr
    fallbackOrder := []string{"ru", "en", "sr"}
    
    for _, fallbackLang := range fallbackOrder {
        if fallbackLang == language {
            continue // пропустить тот же язык
        }
        
        query := `
            SELECT translated_text 
            FROM translations 
            WHERE entity_type = $1 AND entity_id = $2 AND field_name = $3 AND language = $4
            LIMIT 1
        `
        
        var translatedText string
        err := s.db.QueryRowContext(ctx, query, entityType, entityID, fieldName, fallbackLang).Scan(&translatedText)
        
        if err == nil {
            return translatedText, nil
        }
    }
    
    // Если нет переводов, попробовать получить оригинальный текст
    return s.getOriginalText(ctx, entityType, entityID, fieldName)
}

func (s *TranslationService) getOriginalText(ctx context.Context, entityType string, entityID int, fieldName string) (string, error) {
    // Здесь должна быть логика получения оригинального текста из основной таблицы
    // В зависимости от entity_type
    
    switch entityType {
    case "listing":
        return s.getListingOriginalText(ctx, entityID, fieldName)
    case "category":
        return s.getCategoryOriginalText(ctx, entityID, fieldName)
    // добавить другие типы
    }
    
    return "", fmt.Errorf("original text not found for %s:%d:%s", entityType, entityID, fieldName)
}

func (s *TranslationService) getListingOriginalText(ctx context.Context, listingID int, fieldName string) (string, error) {
    var query string
    switch fieldName {
    case "title":
        query = "SELECT title FROM marketplace_listings WHERE id = $1"
    case "description":
        query = "SELECT description FROM marketplace_listings WHERE id = $1"
    default:
        return "", fmt.Errorf("unknown field: %s", fieldName)
    }
    
    var text string
    err := s.db.QueryRowContext(ctx, query, listingID).Scan(&text)
    return text, err
}

func (s *TranslationService) getCategoryOriginalText(ctx context.Context, categoryID int, fieldName string) (string, error) {
    var query string
    switch fieldName {
    case "name":
        query = "SELECT name FROM marketplace_categories WHERE id = $1"
    case "description":
        query = "SELECT description FROM marketplace_categories WHERE id = $1"
    default:
        return "", fmt.Errorf("unknown field: %s", fieldName)
    }
    
    var text string
    err := s.db.QueryRowContext(ctx, query, categoryID).Scan(&text)
    return text, err
}

// GetEntityTranslations получает все переводы для сущности
func (s *TranslationService) GetEntityTranslations(ctx context.Context, entityType string, entityID int, language string) (map[string]string, error) {
    query := `
        SELECT field_name, translated_text 
        FROM translations 
        WHERE entity_type = $1 AND entity_id = $2 AND language = $3
    `
    
    rows, err := s.db.QueryContext(ctx, query, entityType, entityID, language)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    result := make(map[string]string)
    for rows.Next() {
        var fieldName, translatedText string
        if err := rows.Scan(&fieldName, &translatedText); err != nil {
            return nil, err
        }
        result[fieldName] = translatedText
    }
    
    return result, nil
}

// Статические переводы для системных сообщений
var SystemMessages = map[string]map[string]string{
    "errors": {
        "ru": map[string]string{
            "validation.required": "Поле обязательно для заполнения",
            "validation.email":    "Неверный формат email",
            "auth.unauthorized":   "Необходима авторизация",
            "auth.forbidden":      "Доступ запрещен",
            "server.internal":     "Внутренняя ошибка сервера",
        },
        "en": map[string]string{
            "validation.required": "Field is required",
            "validation.email":    "Invalid email format",
            "auth.unauthorized":   "Authorization required",
            "auth.forbidden":      "Access forbidden",
            "server.internal":     "Internal server error",
        },
        "sr": map[string]string{
            "validation.required": "Поље је обавезно",
            "validation.email":    "Неисправан формат емаил-а",
            "auth.unauthorized":   "Потребна је ауторизација",
            "auth.forbidden":      "Приступ забрањен", 
            "server.internal":     "Унутрашња грешка сервера",
        },
    },
}

func (s *TranslationService) GetSystemMessage(key, language string) string {
    if msgs, ok := SystemMessages["errors"][language]; ok {
        if msg, ok := msgs[key]; ok {
            return msg
        }
    }
    
    // Fallback к английскому
    if msgs, ok := SystemMessages["errors"]["en"]; ok {
        if msg, ok := msgs[key]; ok {
            return msg
        }
    }
    
    return key // возвращаем ключ если нет перевода
}
```

**Шаг 5.3: Обновление обработчиков**
Файл: `backend/internal/proj/marketplace/handler/handler.go`
```go
package handler

import (
    "context"
    "net/http"
    
    "github.com/gofiber/fiber/v2"
    "your-project/internal/middleware"
    "your-project/internal/services"
)

type Handler struct {
    marketplaceService *services.MarketplaceService
    translationService *services.TranslationService
}

func New(ms *services.MarketplaceService, ts *services.TranslationService) *Handler {
    return &Handler{
        marketplaceService: ms,
        translationService: ts,
    }
}

// GetListings возвращает список объявлений с переводами
func (h *Handler) GetListings(c *fiber.Ctx) error {
    language := middleware.GetLocale(c)
    
    // Получить объявления из основного сервиса
    listings, err := h.marketplaceService.GetListings(c.Context(), /* параметры */)
    if err != nil {
        return h.errorResponse(c, "listings.getError", err, language)
    }
    
    // Добавить переводы
    for i, listing := range listings {
        translations, err := h.translationService.GetEntityTranslations(
            c.Context(), "listing", listing.ID, language,
        )
        if err == nil {
            // Заменить поля переводами если они есть
            if title, ok := translations["title"]; ok {
                listings[i].Title = title
            }
            if description, ok := translations["description"]; ok {
                listings[i].Description = description
            }
        }
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data":    listings,
        "meta": fiber.Map{
            "language": language,
            "total":    len(listings),
        },
    })
}

// errorResponse возвращает локализованную ошибку
func (h *Handler) errorResponse(c *fiber.Ctx, errorKey string, err error, language string) error {
    message := h.translationService.GetSystemMessage(errorKey, language)
    
    // Логировать оригинальную ошибку
    // log.Error("API Error", "key", errorKey, "error", err, "language", language)
    
    return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
        "success": false,
        "error": fiber.Map{
            "message": message,
            "code":    errorKey,
        },
        "meta": fiber.Map{
            "language": language,
        },
    })
}
```

**Шаг 5.4: Интеграция в основное приложение**
Файл: `backend/internal/server/server.go`
```go
func (s *Server) setupRoutes() {
    // Добавить middleware для локализации
    s.app.Use(localeMiddleware.New(localeMiddleware.LocaleConfig{
        SupportedLocales: []string{"ru", "sr", "en"},
        DefaultLocale:    "ru",
        Sources:          []string{"cookie", "header", "query", "domain"},
    }))
    
    // Создать сервис переводов
    translationService := services.NewTranslationService(s.db, s.redis)
    
    // Передать в хендлеры
    marketplaceHandler := marketplaceHandler.New(s.marketplaceService, translationService)
    
    // Настроить маршруты
    api := s.app.Group("/api/v1")
    api.Get("/marketplace/listings", marketplaceHandler.GetListings)
    // ... другие маршруты
}
```

**Команды для тестирования:**
```bash
# 1. Тест с разными языками
curl -H "Accept-Language: ru-RU,ru;q=0.9" http://localhost:3000/api/v1/marketplace/listings
curl -H "Accept-Language: sr-RS,sr;q=0.9" http://localhost:3000/api/v1/marketplace/listings
curl -H "Accept-Language: en-US,en;q=0.9" http://localhost:3000/api/v1/marketplace/listings

# 2. Тест с cookie
curl -H "Cookie: locale-preference=sr" http://localhost:3000/api/v1/marketplace/listings

# 3. Тест с query параметром
curl http://localhost:3000/api/v1/marketplace/listings?lang=en

# 4. Проверить Content-Language header в ответе
curl -I -H "Accept-Language: ru" http://localhost:3000/api/v1/marketplace/listings
```

</details>

---

### 🟢 СРЕДНИЙ ПРИОРИТЕТ (Оптимизация, делать в последнюю очередь)

#### Проблема #6: Неоптимальная загрузка переводов
**Влияние**: Медленная загрузка страниц из-за больших файлов переводов  
**Файлы**: `frontend/svetu/src/messages/*.json`, `frontend/svetu/src/i18n/`  
**Техническая суть**: Все переводы загружаются сразу, нет lazy loading

<details>
<summary>📁 Детальное решение для оптимизации загрузки</summary>

**Шаг 6.1: Анализ размера файлов переводов**
```bash
cd frontend/svetu/src/messages

echo "Размеры файлов переводов:"
ls -lh *.json

echo -e "\nТоп разделов по размеру:"
for file in *.json; do
    echo "=== $file ==="
    jq -r 'to_entries[] | "\(.key): \(.value | tostring | length) символов"' "$file" | sort -n -k2 | tail -5
done
```

**Шаг 6.2: Разбиение переводов на модули**
Создать структуру: `frontend/svetu/src/messages/modules/`
```bash
mkdir -p frontend/svetu/src/messages/modules/{common,marketplace,storefronts,admin,auth}

# Разбить переводы по модулям
cd frontend/svetu/src/messages

# Общие переводы (всегда загружаются)
jq '{common: .common, errors: .errors, navigation: .navigation, header: .header}' ru.json > modules/common/ru.json
jq '{common: .common, errors: .errors, navigation: .navigation, header: .header}' en.json > modules/common/en.json  
jq '{common: .common, errors: .errors, navigation: .navigation, header: .header}' sr.json > modules/common/sr.json

# Модуль маркетплейса
jq '{marketplace: .marketplace, listing: .listing, search: .search, map: .map}' ru.json > modules/marketplace/ru.json
jq '{marketplace: .marketplace, listing: .listing, search: .search, map: .map}' en.json > modules/marketplace/en.json
jq '{marketplace: .marketplace, listing: .listing, search: .search, map: .map}' sr.json > modules/marketplace/sr.json

# Модуль витрин  
jq '{storefront: .storefront, products: .products, orders: .orders}' ru.json > modules/storefronts/ru.json
jq '{storefront: .storefront, products: .products, orders: .orders}' en.json > modules/storefronts/en.json
jq '{storefront: .storefront, products: .products, orders: .orders}' sr.json > modules/storefronts/sr.json

# Админка
jq '{admin: .admin, analytics: .analytics, permissions: .permissions, roles: .roles}' ru.json > modules/admin/ru.json
jq '{admin: .admin, analytics: .analytics, permissions: .permissions, roles: .roles}' en.json > modules/admin/en.json  
jq '{admin: .admin, analytics: .analytics, permissions: .permissions, roles: .roles}' sr.json > modules/admin/sr.json

# Авторизация
jq '{auth: .auth, profile: .profile, validation: .validation}' ru.json > modules/auth/ru.json
jq '{auth: .auth, profile: .profile, validation: .validation}' en.json > modules/auth/en.json
jq '{auth: .auth, profile: .profile, validation: .validation}' sr.json > modules/auth/sr.json
```

**Шаг 6.3: Создание системы динамической загрузки**
Файл: `frontend/svetu/src/i18n/dynamic-loader.ts`
```typescript
import { Locale } from './config';

export type MessageModule = 'common' | 'marketplace' | 'storefronts' | 'admin' | 'auth';

interface LoadedModules {
  [locale: string]: {
    [module: string]: Record<string, any>;
  };
}

class DynamicTranslationLoader {
  private loadedModules: LoadedModules = {};
  private loadingPromises: Map<string, Promise<Record<string, any>>> = new Map();

  // Загрузить общие переводы (синхронно при старте приложения)
  async loadCommonMessages(locale: Locale): Promise<Record<string, any>> {
    const cacheKey = `${locale}-common`;
    
    if (this.loadedModules[locale]?.common) {
      return this.loadedModules[locale].common;
    }

    try {
      const messages = await import(`../messages/modules/common/${locale}.json`).then(m => m.default);
      
      if (!this.loadedModules[locale]) {
        this.loadedModules[locale] = {};
      }
      
      this.loadedModules[locale].common = messages;
      return messages;
    } catch (error) {
      console.error(`Failed to load common messages for ${locale}:`, error);
      
      // Fallback к основному файлу
      const fallbackMessages = await import(`../messages/${locale}.json`).then(m => m.default);
      return {
        common: fallbackMessages.common || {},
        errors: fallbackMessages.errors || {},
        navigation: fallbackMessages.navigation || {},
        header: fallbackMessages.header || {},
      };
    }
  }

  // Загрузить модуль по требованию
  async loadModule(locale: Locale, module: MessageModule): Promise<Record<string, any>> {
    const cacheKey = `${locale}-${module}`;
    
    // Если модуль уже загружен
    if (this.loadedModules[locale]?.[module]) {
      return this.loadedModules[locale][module];
    }

    // Если модуль уже загружается
    if (this.loadingPromises.has(cacheKey)) {
      return this.loadingPromises.get(cacheKey)!;
    }

    // Загрузить модуль
    const loadingPromise = this.doLoadModule(locale, module);
    this.loadingPromises.set(cacheKey, loadingPromise);

    try {
      const messages = await loadingPromise;
      
      if (!this.loadedModules[locale]) {
        this.loadedModules[locale] = {};
      }
      
      this.loadedModules[locale][module] = messages;
      return messages;
    } finally {
      this.loadingPromises.delete(cacheKey);
    }
  }

  private async doLoadModule(locale: Locale, module: MessageModule): Promise<Record<string, any>> {
    try {
      return await import(`../messages/modules/${module}/${locale}.json`).then(m => m.default);
    } catch (error) {
      console.error(`Failed to load module ${module} for ${locale}:`, error);
      
      // Fallback к основному файлу переводов
      const fullMessages = await import(`../messages/${locale}.json`).then(m => m.default);
      
      // Извлечь нужные секции для модуля
      switch (module) {
        case 'marketplace':
          return {
            marketplace: fullMessages.marketplace || {},
            listing: fullMessages.listing || {},
            search: fullMessages.search || {},
            map: fullMessages.map || {},
          };
        case 'storefronts':
          return {
            storefront: fullMessages.storefront || {},
            products: fullMessages.products || {},
            orders: fullMessages.orders || {},
          };
        case 'admin':
          return {
            admin: fullMessages.admin || {},
            analytics: fullMessages.analytics || {},
            permissions: fullMessages.permissions || {},
            roles: fullMessages.roles || {},
          };
        case 'auth':
          return {
            auth: fullMessages.auth || {},
            profile: fullMessages.profile || {},
            validation: fullMessages.validation || {},
          };
        default:
          return {};
      }
    }
  }

  // Получить все загруженные сообщения для локали
  getAllLoadedMessages(locale: Locale): Record<string, any> {
    const localeModules = this.loadedModules[locale] || {};
    
    // Объединить все загруженные модули
    let allMessages = {};
    for (const moduleMessages of Object.values(localeModules)) {
      allMessages = { ...allMessages, ...moduleMessages };
    }
    
    return allMessages;
  }

  // Предзагрузить модули для определенного маршрута
  async preloadModulesForRoute(locale: Locale, route: string): Promise<void> {
    const routeModules = this.getModulesForRoute(route);
    
    await Promise.all(
      routeModules.map(module => this.loadModule(locale, module))
    );
  }

  private getModulesForRoute(route: string): MessageModule[] {
    if (route.startsWith('/marketplace') || route.startsWith('/listings')) {
      return ['marketplace'];
    }
    
    if (route.startsWith('/storefronts') || route.startsWith('/products')) {
      return ['storefronts'];
    }
    
    if (route.startsWith('/admin')) {
      return ['admin'];
    }
    
    if (route.startsWith('/auth') || route.startsWith('/profile')) {
      return ['auth'];
    }
    
    return [];
  }
}

export const dynamicLoader = new DynamicTranslationLoader();
```

**Шаг 6.4: Обновление конфигурации i18n**
Файл: `frontend/svetu/src/i18n/config.ts`
```typescript
import { dynamicLoader, MessageModule } from './dynamic-loader';

// Обновленная функция загрузки с поддержкой модулей
export async function getLocaleMessages(locale: Locale): Promise<Record<string, any>> {
  try {
    // Всегда загружаем общие переводы
    const commonMessages = await dynamicLoader.loadCommonMessages(locale);
    
    // Если это SSR или первая загрузка, загружаем только common
    if (typeof window === 'undefined') {
      return commonMessages;
    }
    
    // На клиенте можем загрузить дополнительные модули на основе маршрута
    const currentPath = window.location.pathname;
    await dynamicLoader.preloadModulesForRoute(locale, currentPath);
    
    // Получить все загруженные сообщения
    return dynamicLoader.getAllLoadedMessages(locale);
  } catch (error) {
    console.error(`Failed to load messages for locale ${locale}:`, error);
    
    // Fallback к полному файлу
    return import(`../messages/${locale}.json`).then(m => m.default);
  }
}

// Хук для динамической загрузки модулей в компонентах
export function useDynamicTranslations(module: MessageModule) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const locale = useLocale();
  
  useEffect(() => {
    let mounted = true;
    
    async function loadModule() {
      if (dynamicLoader.getAllLoadedMessages(locale)[module]) {
        return; // уже загружен
      }
      
      setIsLoading(true);
      setError(null);
      
      try {
        await dynamicLoader.loadModule(locale, module);
        
        if (mounted) {
          setIsLoading(false);
          // Trigger re-render если нужно
          // forceUpdate();
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err : new Error('Loading failed'));
          setIsLoading(false);
        }
      }
    }
    
    loadModule();
    
    return () => {
      mounted = false;
    };
  }, [locale, module]);
  
  return { isLoading, error };
}
```

**Шаг 6.5: Обновление компонентов для динамической загрузки**
Пример использования в компоненте:
```typescript
// components/marketplace/MarketplacePage.tsx
import { useDynamicTranslations } from '@/i18n/config';

export default function MarketplacePage() {
  const t = useTranslations();
  const { isLoading, error } = useDynamicTranslations('marketplace');
  
  if (isLoading) {
    return <div>Loading translations...</div>;
  }
  
  if (error) {
    console.error('Translation loading error:', error);
    // Продолжить с доступными переводами
  }
  
  return (
    <div>
      <h1>{t('marketplace.title')}</h1>
      {/* остальной контент */}
    </div>
  );
}
```

**Шаг 6.6: Предзагрузка модулей на основе навигации**
Файл: `frontend/svetu/src/components/NavigationPreloader.tsx`
```typescript
'use client';

import { usePathname } from 'next/navigation';
import { useLocale } from 'next-intl';
import { useEffect } from 'react';
import { dynamicLoader } from '@/i18n/dynamic-loader';

export default function NavigationPreloader() {
  const pathname = usePathname();
  const locale = useLocale();
  
  useEffect(() => {
    // Предзагрузить модули для текущего маршрута
    dynamicLoader.preloadModulesForRoute(locale, pathname);
  }, [pathname, locale]);
  
  return null; // компонент невидимый
}
```

**Добавить в layout:**
```typescript
// app/[locale]/layout.tsx
import NavigationPreloader from '@/components/NavigationPreloader';

export default function RootLayout({ children, params }: {
  children: React.ReactNode;
  params: { locale: string };
}) {
  return (
    <html lang={params.locale}>
      <body>
        <NavigationPreloader />
        {children}
      </body>
    </html>
  );
}
```

**Команды для тестирования производительности:**
```bash
# 1. Сравнить размеры до и после
cd frontend/svetu/src/messages

echo "Размеры до разбиения:"
ls -lh *.json

echo -e "\nРазмеры после разбиения:"
find modules/ -name "*.json" -exec ls -lh {} \;

# 2. Проверить загрузку в браузере
yarn dev -p 3001
# Открыть DevTools -> Network -> отфильтровать по .json
# Проверить что загружаются только нужные модули

# 3. Тест производительности
yarn build
yarn start
# Lighthouse audit для проверки улучшения производительности
```

</details>

---

## ⏱️ Детальный календарный план

### Неделя 1: Критические исправления (OpenSearch)
**Понедельник-Вторник**: Анализ и backup
- Анализ текущего состояния OpenSearch
- Создание backup индекса
- Дизайн новой структуры с мультиязычностью

**Среда-Четверг**: Реализация мультиязычного индекса
- Создание нового маппинга с language analyzers
- Модификация Go кода для включения переводов
- Тестирование индексирования

**Пятница**: Переключение и тестирование
- Реиндексация с переводами
- Переключение алиаса на новый индекс
- Тестирование поиска на всех языках

### Неделя 2: Конфигурация языков и Backend локализация
**Понедельник**: Исправление конфигурации языков
- Анализ статистики использования языков
- Обновление конфигурации Frontend
- Обновление роутинга и метаданных

**Вторник-Среда**: Backend локализация
- Создание middleware для определения языка
- Реализация сервиса переводов
- Обновление API handlers

**Четверг-Пятница**: Интеграция и тестирование
- Интеграция middleware в приложение
- Тестирование API с разными языками
- Отладка и исправление ошибок

### Неделя 3: Отсутствующие переводы и hardcoded строки
**Понедельник**: Исправление отсутствующих переводов
- Добавление раздела "storefront" в сербский
- Создание скриптов валидации
- Проверка консистентности переводов

**Вторник-Четверг**: Устранение hardcoded строк
- Полный аудит hardcoded строк
- Исправление топ-10 файлов
- Добавление недостающих переводов в JSON

**Пятница**: Тестирование и валидация
- Проверка всех исправленных компонентов
- Валидация отсутствия нелокализованного контента
- Исправление найденных проблем

### Неделя 4: Оптимизация и финальное тестирование
**Понедельник-Вторник**: Оптимизация загрузки переводов
- Разбиение переводов на модули
- Реализация динамической загрузки
- Настройка предзагрузки по маршрутам

**Среда**: Комплексное тестирование
- Тестирование всех компонентов на 3 языках
- Проверка производительности
- Валидация SEO метаданных

**Четверг**: Исправления и доработки
- Исправление найденных в тестировании проблем
- Оптимизация производительности
- Финальная проверка

**Пятница**: Развертывание и мониторинг
- Развертывание на staging
- Проверка в production-подобной среде
- Настройка мониторинга

---

## 🧪 План тестирования

### Модульные тесты
```typescript
// frontend/svetu/src/__tests__/i18n/translations.test.ts
describe('Translation System', () => {
  test('should load all required languages', async () => {
    for (const locale of ['ru', 'sr', 'en']) {
      const messages = await getLocaleMessages(locale);
      expect(messages).toBeDefined();
      expect(Object.keys(messages).length).toBeGreaterThan(0);
    }
  });
  
  test('should have storefront section in all languages', async () => {
    for (const locale of ['ru', 'sr', 'en']) {
      const messages = await getLocaleMessages(locale);
      expect(messages.storefront).toBeDefined();
    }
  });
  
  test('should fallback to default language', async () => {
    // Тест fallback механизма
    const fallbackMessage = await getFallbackTranslation('nonexistent.key', 'invalid-lang');
    expect(fallbackMessage).toBeTruthy();
  });
});
```

### Интеграционные тесты
```bash
#!/bin/bash
# tests/integration/multilingual.sh

echo "Testing multilingual API endpoints..."

# Test Russian
echo "Testing Russian API..."
RESPONSE=$(curl -s -H "Accept-Language: ru-RU" http://localhost:3000/api/v1/marketplace/listings)
echo $RESPONSE | jq -r '.meta.language' | grep -q "ru" || echo "FAIL: Russian API"

# Test Serbian  
echo "Testing Serbian API..."
RESPONSE=$(curl -s -H "Accept-Language: sr-RS" http://localhost:3000/api/v1/marketplace/listings)
echo $RESPONSE | jq -r '.meta.language' | grep -q "sr" || echo "FAIL: Serbian API"

# Test English
echo "Testing English API..."
RESPONSE=$(curl -s -H "Accept-Language: en-US" http://localhost:3000/api/v1/marketplace/listings)
echo $RESPONSE | jq -r '.meta.language' | grep -q "en" || echo "FAIL: English API"

echo "Multilingual API tests completed"
```

### E2E тесты
```typescript
// frontend/svetu/e2e/multilingual.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Multilingual functionality', () => {
  test('should switch languages correctly', async ({ page }) => {
    await page.goto('/');
    
    // Test Russian
    await page.click('[data-testid="language-switcher"]');
    await page.click('[data-testid="language-ru"]');
    await expect(page.locator('h1')).toContainText('Добро пожаловать');
    
    // Test Serbian
    await page.click('[data-testid="language-switcher"]');
    await page.click('[data-testid="language-sr"]');
    await expect(page.locator('h1')).toContainText('Добродошли');
    
    // Test English
    await page.click('[data-testid="language-switcher"]');
    await page.click('[data-testid="language-en"]');
    await expect(page.locator('h1')).toContainText('Welcome');
  });
  
  test('should show correct search results for different languages', async ({ page }) => {
    // Test search in Russian
    await page.goto('/ru/marketplace');
    await page.fill('[data-testid="search-input"]', 'автомобиль');
    await page.press('[data-testid="search-input"]', 'Enter');
    await expect(page.locator('[data-testid="search-results"]')).toBeVisible();
    
    // Test search in Serbian
    await page.goto('/sr/marketplace');
    await page.fill('[data-testid="search-input"]', 'аутомобил');
    await page.press('[data-testid="search-input"]', 'Enter');
    await expect(page.locator('[data-testid="search-results"]')).toBeVisible();
  });
});
```

---

## 📊 Метрики и мониторинг

### KPI для отслеживания прогресса
```yaml
# Покрытие переводами
translation_coverage:
  current: 95%
  target: 99%
  measurement: "% ключей с переводами для всех языков"

# Производительность поиска  
search_performance:
  current: "60% релевантности для неанглийского контента"
  target: "90% релевантности для всех языков"
  measurement: "средняя точность поиска по языкам"

# Размер переводов
translation_bundle_size:
  current: "4.2MB общий размер JSON файлов"
  target: "1.5MB при первой загрузке"
  measurement: "размер загружаемых переводов"

# Hardcoded строки
hardcoded_strings:
  current: "20+ файлов с нелокализованным контентом"
  target: "0 hardcoded строк в production коде"
  measurement: "количество файлов с нелокализаванным контентом"
```

### Инструменты мониторинга
```bash
# Скрипт проверки покрытия переводами
#!/bin/bash
# scripts/check_translation_coverage.sh

echo "Checking translation coverage..."

cd frontend/svetu/src/messages

echo "=== Coverage by language ==="
for lang in ru en sr; do
    total_keys=$(jq -r 'paths(scalars) as $p | $p | join(".")' ${lang}.json | wc -l)
    echo "$lang: $total_keys keys"
done

echo -e "\n=== Missing translations ==="
# Сравнить ключи между языками и найти отсутствующие
./scripts/find_missing_translations.js

echo -e "\n=== Hardcoded strings check ==="
grep -r "\"[А-Яа-я][А-Яа-я ]{3,}\"" ../components/ --include="*.tsx" | wc -l | xargs echo "Cyrillic hardcoded strings found:"
```

---

## 🚀 Критерии готовности к продакшену

### Обязательные требования (MUST HAVE)
- [ ] **OpenSearch индекс** содержит переводы из БД для всех языков
- [ ] **Поиск работает** корректно на русском, сербском и английском  
- [ ] **API возвращает** правильный Content-Language header
- [ ] **Frontend определяет** язык из Accept-Language, cookie и домена
- [ ] **Все файлы переводов** содержат одинаковый набор ключей
- [ ] **0 hardcoded строк** в production компонентах
- [ ] **Fallback механизм** работает при отсутствии перевода

### Желательные требования (SHOULD HAVE)  
- [ ] **Lazy loading** переводов по модулям работает
- [ ] **Предзагрузка** переводов по маршрутам настроена
- [ ] **Кэширование** переводов в Redis функционирует
- [ ] **Автоматическая валидация** переводов в CI/CD
- [ ] **Мониторинг** качества переводов настроен

### Критерии производительности
- [ ] **Первая загрузка**: < 1.5MB переводов
- [ ] **Переключение языка**: < 200ms
- [ ] **Поиск на любом языке**: < 500ms ответ
- [ ] **API с переводами**: < 300ms ответ

---

## 🛠️ Готовые скрипты и команды

### Быстрый старт для разработчика
```bash
#!/bin/bash
# scripts/quick_start_multilingual.sh

echo "🌍 Быстрый старт мультиязычности"

# 1. Backup текущего состояния
echo "📁 Создание backup..."
mkdir -p backups/$(date +%Y%m%d_%H%M%S)
cp -r frontend/svetu/src/messages backups/$(date +%Y%m%d_%H%M%S)/
curl -X POST "http://localhost:9200/_reindex" -H "Content-Type: application/json" -d '{
  "source": {"index": "marketplace"},
  "dest": {"index": "marketplace_backup_$(date +%Y%m%d)"
}}'

# 2. Проверка текущего состояния
echo "🔍 Анализ текущего состояния..."
./scripts/analyze_multilingual_status.sh

# 3. Быстрые исправления
echo "⚡ Применение быстрых исправлений..."

# Добавить отсутствующий storefront в sr.json
jq '.storefront = (.storefront // {})' frontend/svetu/src/messages/en.json > /tmp/storefront_template.json
echo "Добавьте переводы в /tmp/storefront_template.json и выполните:"
echo "jq -s '.[0] * .[1]' frontend/svetu/src/messages/sr.json /tmp/storefront_template.json > frontend/svetu/src/messages/sr.json.new"

echo "✅ Готово! Следуйте детальному плану для полной реализации"
```

### Валидация переводов
```bash
#!/bin/bash
# scripts/validate_translations.sh

echo "🔍 Валидация переводов..."

cd frontend/svetu/src/messages

# Проверить консистентность ключей
echo "Проверка консистентности ключей между языками:"
jq -r 'paths(scalars) as $p | $p | join(".")' ru.json | sort > /tmp/ru_keys_flat.txt
jq -r 'paths(scalars) as $p | $p | join(".")' en.json | sort > /tmp/en_keys_flat.txt  
jq -r 'paths(scalars) as $p | $p | join(".")' sr.json | sort > /tmp/sr_keys_flat.txt

echo "Отсутствует в русском:"
comm -23 /tmp/en_keys_flat.txt /tmp/ru_keys_flat.txt | head -5

echo "Отсутствует в сербском:"
comm -23 /tmp/en_keys_flat.txt /tmp/sr_keys_flat.txt | head -5

# Проверить пустые значения
echo -e "\nПроверка пустых переводов:"
for lang in ru en sr; do
    empty_count=$(jq -r 'paths(scalars) as $p | {"key": ($p | join(".")), "value": getpath($p)} | select(.value == "") | .key' ${lang}.json | wc -l)
    echo "$lang: $empty_count пустых переводов"
done

echo "✅ Валидация завершена"
```

---

## 📞 Поддержка и контакты

### При возникновении проблем
1. **Проверьте логи**: `docker logs hostel_backend` и `yarn dev` вывод
2. **Проверьте кэш**: `docker exec hostel_redis redis-cli FLUSHALL` 
3. **Переиндексируйте**: `cd backend && ./reindex`
4. **Проверьте переводы**: `node scripts/validate_translations.js`

### Полезные команды для отладки
```bash
# Проверить статус OpenSearch
curl -X GET "http://localhost:9200/_cluster/health?pretty"

# Проверить индекс marketplace
curl -X GET "http://localhost:9200/marketplace/_stats?pretty"

# Проверить переводы в БД
PGPASSWORD=password psql -h localhost -U postgres -d svetubd -c "SELECT language, COUNT(*) FROM translations GROUP BY language;"

# Проверить размер переводов
du -sh frontend/svetu/src/messages/

# Найти hardcoded строки
grep -r "\"[А-Яа-я][А-Яа-я ]{3,}\"" frontend/svetu/src/components/ --include="*.tsx" | wc -l
```

---

*Этот план содержит все необходимые детали для реализации улучшений мультиязычности без дополнительного изучения проекта. Следуйте плану поэтапно для достижения максимального результата.*