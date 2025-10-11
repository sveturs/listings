⏺ 📊 Детальный анализ полей: MarketplaceListing vs StorefrontProduct

✅ ОБЩИЕ поля (совпадают по смыслу):

| №   | Поле                | MarketplaceListing                    | StorefrontProduct                         | Совместимость                        |
  |-----|---------------------|---------------------------------------|-------------------------------------------|--------------------------------------|
| 1   | ID                  | ID int                                | ID int                                    | ✅ 100%                               |
| 2   | Название            | Title string                          | Name string                               | ✅ 95% (разные имена)                 |
| 3   | Описание            | Description string                    | Description string                        | ✅ 100%                               |
| 4   | Цена                | Price float64                         | Price float64                             | ✅ 100%                               |
| 5   | Категория           | CategoryID int                        | CategoryID int                            | ✅ 100%                               |
| 6   | Изображения         | Images []MarketplaceImage             | Images []StorefrontProductImage           | ✅ 90% (оба реализуют ImageInterface) |
| 7   | Категория (объект)  | Category *MarketplaceCategory         | Category *MarketplaceCategory             | ✅ 100%                               |
| 8   | Создано             | CreatedAt time.Time                   | CreatedAt time.Time                       | ✅ 100%                               |
| 9   | Обновлено           | UpdatedAt time.Time                   | UpdatedAt time.Time                       | ✅ 100%                               |
| 10  | Переводы            | Translations TranslationMap           | Translations map[string]map[string]string | ⚠️ 70% (разные типы)                 |
| 11  | Адрес переводы      | AddressMultilingual map[string]string | AddressTranslations map[string]string     | ✅ 95% (разные имена)                 |
| 12  | Широта              | Latitude *float64                     | IndividualLatitude *float64               | ✅ 100% (по смыслу)                   |
| 13  | Долгота             | Longitude *float64                    | IndividualLongitude *float64              | ✅ 100% (по смыслу)                   |
| 14  | Показать на карте   | ShowOnMap bool                        | ShowOnMap bool                            | ✅ 100%                               |
| 15  | Приватность локации | LocationPrivacy string                | LocationPrivacy *string                   | ✅ 95%                                |
| 16  | Адрес               | Location string + City + Country      | IndividualAddress *string                 | ✅ 90%                                |
| 17  | Атрибуты            | Attributes []ListingAttributeValue    | Attributes JSONB                          | ⚠️ 60% (разная структура)            |
| 18  | Активность          | Status string (active/...)            | IsActive bool                             | ⚠️ 70% (разная семантика)            |
| 19  | Просмотры           | ViewsCount int                        | ViewCount int                             | ✅ 100%                               |
| 20  | Варианты            | Variants []MarketplaceListingVariant  | Variants []StorefrontProductVariant       | ⚠️ 60% (разные структуры)            |
| 21  | Остаток             | StockQuantity *int                    | StockQuantity *int                        | ✅ 100%                               |
| 22  | Статус остатка      | StockStatus *string                   | StockStatus string                        | ✅ 95%                                |

Итого общих полей: 22 поля (с разной степенью совместимости)

  ---

⏺ 🔴 УНИКАЛЬНЫЕ поля MarketplaceListing (ТОЛЬКО для маркетплейса):

| №   | Поле                | Тип                    | Зачем нужно                           | Можно ли в NULL? |
  |-----|---------------------|------------------------|---------------------------------------|------------------|
| 1   | UserID              | int                    | Владелец объявления (P2P)             | ❌ NOT NULL       |
| 2   | Condition           | string                 | Состояние товара (новый/б/у)          | ❌ NOT NULL       |
| 3   | Status              | string                 | draft/active/sold/archived            | ❌ NOT NULL       |
| 4   | HelpfulVotes        | int                    | Голоса "полезно"                      | ✅ DEFAULT 0      |
| 5   | NotHelpfulVotes     | int                    | Голоса "не полезно"                   | ✅ DEFAULT 0      |
| 6   | IsFavorite          | bool                   | В избранном у текущего пользователя   | ✅ DEFAULT false  |
| 7   | OldPrice            | *float64               | Старая цена (для скидок)              | ✅ NULL           |
| 8   | HasDiscount         | bool                   | Есть ли скидка                        | ✅ DEFAULT false  |
| 9   | DiscountPercentage  | *int                   | Процент скидки                        | ✅ NULL           |
| 10  | Metadata            | map[string]interface{} | Дополнительные данные                 | ✅ NULL           |
| 11  | AverageRating       | float64                | Средняя оценка                        | ✅ DEFAULT 0      |
| 12  | ReviewCount         | int                    | Количество отзывов                    | ✅ DEFAULT 0      |
| 13  | StorefrontID        | *int                   | Связь с витриной (если товар витрины) | ✅ NULL           |
| 14  | Storefront          | *Storefront            | Данные витрины                        | ✅ NULL           |
| 15  | ExternalID          | string                 | ID из внешней системы                 | ✅ NULL           |
| 16  | IsStorefrontProduct | bool                   | Флаг что это товар витрины            | ✅ DEFAULT false  |
| 17  | OriginalLanguage    | string                 | Оригинальный язык                     | ✅ NULL           |
| 18  | RawTranslations     | interface{}            | Сырые данные переводов                | ✅ NULL           |
| 19  | CategoryPathNames   | []string               | Путь категорий (названия)             | ✅ NULL           |
| 20  | CategoryPathIds     | []int                  | Путь категорий (ID)                   | ✅ NULL           |
| 21  | CategoryPathSlugs   | []string               | Путь категорий (slugs)                | ✅ NULL           |
| 22  | CategoryPath        | []string               | Путь категорий                        | ✅ NULL           |
| 23  | User                | *User                  | Данные пользователя-владельца         | ✅ NULL (join)    |

Итого уникальных полей MarketplaceListing: 23 поля

  ---

⏺ 🟢 УНИКАЛЬНЫЕ поля StorefrontProduct (ТОЛЬКО для витрин):

| №   | Поле                  | Тип     | Зачем нужно            | Можно ли в NULL? |
  |-----|-----------------------|---------|------------------------|------------------|
| 1   | StorefrontID          | int     | Витрина-владелец (B2C) | ❌ NOT NULL       |
| 2   | Currency              | string  | Валюта товара          | ❌ NOT NULL       |
| 3   | SKU                   | *string | Артикул продавца       | ✅ NULL           |
| 4   | Barcode               | *string | Штрихкод (EAN/UPC)     | ✅ NULL           |
| 5   | IsActive              | bool    | Активность товара      | ❌ NOT NULL       |
| 6   | SoldCount             | int     | Количество продаж      | ✅ DEFAULT 0      |
| 7   | HasIndividualLocation | bool    | Есть ли своя локация   | ✅ DEFAULT false  |
| 8   | HasVariants           | bool    | Есть ли варианты       | ✅ DEFAULT false  |

Итого уникальных полей StorefrontProduct: 8 полей

  ---

⏺ 📈 ИТОГОВАЯ СТАТИСТИКА:

MarketplaceListing:
├─ Общие поля:      22 (49%)
├─ Уникальные поля: 23 (51%)
└─ ВСЕГО:           45 полей

StorefrontProduct:
├─ Общие поля:      22 (73%)
├─ Уникальные поля:  8 (27%)
└─ ВСЕГО:           30 полей

Объединенная сущность Product:
├─ Общие поля:      22
├─ Уникальные ML:   23
├─ Уникальные SP:    8
└─ ВСЕГО:           53 поля (!!)
