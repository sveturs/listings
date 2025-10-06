# План улучшенного импорта для Digital Vision (расширенная версия)

**Дата создания:** 2025-10-06
**Версия:** 2.0 (Enhanced)
**Клиент:** Digital Vision (https://digitalvision.rs)
**Статус:** 📋 Детальное планирование
**Приоритет:** 🔥 КРИТИЧЕСКИЙ (Premium клиент)

---

## 🎯 Расширенные требования

### Ключевые вопросы и ответы

#### 1. ✅ Фотографии → S3/MinIO
**Вопрос:** Будут ли фотографии загружены на S3?
**Ответ:** **ДА**, уже реализовано!
- Функция `importProductImages()` скачивает изображения из URL
- Загружает в MinIO (наш S3-совместимый storage)
- Bucket: `storefront-products`
- Генерирует thumbnails автоматически
- Graceful обработка ошибок (недоступные URL)

**Статус:** ✅ Работает (Спринт 2, задача 2.2)

---

#### 2. ❌ Preview с маппингом категорий ДО импорта
**Вопрос:** Будет ли возможность сопоставить категории клиента с нашими перед импортом?
**Ответ:** **НЕТ**, сейчас не реализовано!

**Что есть:**
- ✅ Preview показывает 10 товаров
- ✅ Валидация данных
- ❌ НЕТ маппинга категорий в preview

**Что нужно:**
```tsx
// В ImportPreviewTable показать:
<CategoryMappingStep>
  <ExternalCategory>
    OPREMA ZA MOBILNI > MASKE > SAMSUNG
  </ExternalCategory>

  <MappingSuggestion type="ai" confidence={0.95}>
    → Электроника > Аксессуары для телефонов > Чехлы Samsung
  </MappingSuggestion>

  <ManualOverride>
    <CategorySelector
      categories={marketplaceCategories}
      onSelect={updateMapping}
    />
  </ManualOverride>
</CategoryMappingStep>
```

**Приоритет:** 🔥 КРИТИЧЕСКИЙ

---

#### 3. ❌ AI автоматическое сопоставление категорий
**Вопрос:** Может ли AI предложить автоматическое сопоставление?
**Ответ:** **Частично реализовано**, но НЕ используется!

**Что есть:**
- ✅ `AICategoryDetector` - определяет категорию по названию/описанию товара
- ✅ `CategoryMappingService` - сохраняет маппинги в БД
- ❌ НЕТ AI маппинга внешних категорий на наши

**Что нужно:**
```go
// backend/internal/proj/storefronts/service/ai_category_mapper.go

type AICategoryMapper struct {
    aiDetector  *services.AICategoryDetector
    marketplaceCategories []models.MarketplaceCategory
}

func (m *AICategoryMapper) MapExternalCategory(
    externalCategory string,  // "OPREMA ZA MOBILNI > MASKE > SAMSUNG"
) (*CategoryMappingSuggestion, error) {
    // 1. Разбить external category на уровни
    levels := strings.Split(externalCategory, ">")

    // 2. Для каждого уровня найти похожую категорию в нашей БД
    suggestions := m.findSimilarCategories(levels)

    // 3. Использовать AI для финального выбора
    bestMatch := m.aiDetector.SelectBestMatch(suggestions, externalCategory)

    // 4. Вернуть с confidence score
    return &CategoryMappingSuggestion{
        ExternalCategory:      externalCategory,
        SuggestedCategoryID:   bestMatch.ID,
        SuggestedCategoryPath: bestMatch.Path,
        ConfidenceScore:       bestMatch.Confidence, // 0.0-1.0
        ReasoningSteps:        bestMatch.Reasoning,
    }, nil
}
```

**Приоритет:** 🔥 КРИТИЧЕСКИЙ

---

#### 4. ✅ Ручной маппинг только для проблемных категорий
**Вопрос:** Вручную сопоставлять только то, с чем AI не справился?
**Ответ:** **ДА**, именно так и должно работать!

**Workflow:**
```
1. AI анализирует ВСЕ категории Digital Vision (388 штук)
   ↓
2. Для каждой предлагает маппинг с confidence score
   ↓
3. Пользователь видит:
   ✅ High confidence (>0.90): 320 категорий - auto-approve
   ⚠️ Medium confidence (0.70-0.90): 50 категорий - review recommended
   ❌ Low confidence (<0.70): 18 категорий - manual required
   ↓
4. Пользователь проверяет только Medium + Low (68 категорий)
   ↓
5. Остальные 320 применяются автоматически
```

**UI:**
```tsx
<CategoryMappingReview>
  <AutoApprovedSection count={320} expanded={false}>
    ✅ Высокая уверенность - применено автоматически
  </AutoApprovedSection>

  <ReviewSection count={50} expanded={true}>
    ⚠️ Средняя уверенность - рекомендуется проверить
    {mediumConfidenceCategories.map(cat => (
      <MappingRow
        external={cat.external}
        suggested={cat.suggested}
        confidence={cat.confidence}
        onApprove={approve}
        onEdit={edit}
      />
    ))}
  </ReviewSection>

  <ManualSection count={18} expanded={true}>
    ❌ Низкая уверенность - требуется ручной маппинг
    {lowConfidenceCategories.map(cat => (
      <ManualMappingRow
        external={cat.external}
        onSelect={selectCategory}
      />
    ))}
  </ManualSection>
</CategoryMappingReview>
```

**Приоритет:** 🔥 КРИТИЧЕСКИЙ

---

#### 5. ❌ AI предложение новых категорий
**Вопрос:** Может ли AI найти важные категории у клиента и предложить добавить их нам?
**Ответ:** **НЕТ**, сейчас не реализовано, но ОТЛИЧНАЯ идея!

**Сценарий:**
```
Digital Vision имеет категорию:
"OPREMA ZA MOBILNI > BATERIJE > BATERIJE ECO GRADE" (188 товаров!)

У нас в маркетплейсе:
"Электроника > Аксессуары > Батареи" (без разделения на ECO/Outlet)

AI анализ:
1. Обнаруживает что "ECO GRADE" - это отдельный tier качества
2. Видит 188 товаров в этой категории (значимый объем!)
3. Проверяет что у нас нет такой подкатегории
4. Предлагает создать:
   "Электроника > Аксессуары > Батареи > Эко-класс (восстановленные)"
```

**Реализация:**
```go
// backend/internal/proj/storefronts/service/ai_category_analyzer.go

type CategoryInsight struct {
    ExternalCategory    string
    ProductCount        int
    Importance          float64  // 0-1, based on product count
    IsUnique            bool     // Нет аналога у нас
    SuggestedNewCategory *NewCategoryProposal
}

type NewCategoryProposal struct {
    ParentCategoryID  int
    Name              string
    Description       string
    Reasoning         string
    ExpectedProducts  int
    SimilarCategories []int  // Связанные категории
}

func (a *AICategoryAnalyzer) AnalyzeClientCategories(
    clientCategories []ClientCategory,
) []CategoryInsight {
    insights := []CategoryInsight{}

    for _, cat := range clientCategories {
        // 1. Проверить есть ли у нас похожая
        ourCategory := a.findSimilarCategory(cat.Path)

        // 2. Если нет и товаров много - это важная категория
        if ourCategory == nil && cat.ProductCount > 50 {
            proposal := a.generateNewCategoryProposal(cat)
            insights = append(insights, CategoryInsight{
                ExternalCategory: cat.Path,
                ProductCount: cat.ProductCount,
                Importance: a.calculateImportance(cat),
                IsUnique: true,
                SuggestedNewCategory: proposal,
            })
        }
    }

    return insights
}
```

**UI:**
```tsx
<NewCategoryProposals>
  <Proposal importance="high">
    <ExternalCategory>
      OPREMA ZA MOBILNI > BATERIJE > BATERIJE ECO GRADE (188 товаров)
    </ExternalCategory>

    <Analysis>
      ✨ AI обнаружил значимую категорию без аналога в системе

      Характеристики:
      - 188 товаров (1.1% от общего прайса)
      - Специализация: восстановленные/эко батареи
      - Средняя цена: 890 RSD (ниже новых на 30%)
    </Analysis>

    <Proposal>
      Предлагается создать:
      📁 Электроника > Аксессуары > Батареи
         └── ♻️ Эко-класс (восстановленные)

      Теги: eco, refurbished, economy
    </Proposal>

    <Actions>
      <Button onClick={createCategory}>✅ Создать категорию</Button>
      <Button onClick={mapToExisting}>🔗 Сопоставить с существующей</Button>
      <Button onClick={skip}>⏭️ Пропустить</Button>
    </Actions>
  </Proposal>
</NewCategoryProposals>
```

**Приоритет:** 🟡 ВАЖНО (но не блокирующее)

---

#### 6. ❌ Маппинг атрибутов клиента
**Вопрос:** Присутствуют ли дополнительные атрибуты? Нужно ли их сопоставлять?
**Ответ:** **ДА**, атрибуты есть, но сейчас просто складываются в JSONB!

**Что есть в Digital Vision XML:**
```xml
<artikal>
  <uvoznik>Digital Vision doo</uvoznik>         <!-- Импортер -->
  <godinaUvoza>2025.</godinaUvoza>              <!-- Год импорта -->
  <zemljaPorekla>Kina</zemljaPorekla>           <!-- Страна происхождения -->
  <dostupan>1</dostupan>                        <!-- В наличии -->
  <naAkciji>1</naAkciji>                        <!-- На акции -->
  <barKod>1234567890</barKod>                   <!-- Штрих-код -->
</artikal>
```

**Сейчас они просто складываются в JSONB:**
```go
product.Attributes = map[string]interface{}{
    "uvoznik":        dvProduct.Uvoznik,
    "godina_uvoza":   dvProduct.GodinaUvoza,
    "zemlja_porekla": dvProduct.ZemljaPorekla,
    // ...
}
```

**Что нужно:**

1. **Структурированные атрибуты в БД:**
```sql
-- У нас уже есть:
CREATE TABLE product_variant_attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    type VARCHAR(20),  -- text, number, boolean, select, multiselect
    is_required BOOLEAN DEFAULT false,
    is_variant_defining BOOLEAN DEFAULT false,  -- Для цвета, размера и т.д.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Нужно добавить атрибуты Digital Vision:
INSERT INTO product_variant_attributes (name, display_name, type) VALUES
('importer', 'Импортер', 'text'),
('import_year', 'Год импорта', 'text'),
('country_of_origin', 'Страна происхождения', 'select'),
('on_sale', 'На акции', 'boolean');
```

2. **Маппинг внешних атрибутов:**
```go
// backend/internal/proj/storefronts/service/attribute_mapper.go

type AttributeMapping struct {
    ExternalName    string  // "uvoznik"
    InternalAttrID  int     // ID в product_variant_attributes
    Transform       func(value interface{}) interface{}
}

var digitalVisionAttributeMap = map[string]AttributeMapping{
    "uvoznik": {
        InternalAttrID: 101,  // "importer"
        Transform: func(v interface{}) interface{} {
            return v  // Прямое копирование
        },
    },
    "godinaUvoza": {
        InternalAttrID: 102,  // "import_year"
        Transform: func(v interface{}) interface{} {
            year := strings.TrimSuffix(v.(string), ".")
            return year
        },
    },
    "zemljaPorekla": {
        InternalAttrID: 103,  // "country_of_origin"
        Transform: func(v interface{}) interface{} {
            // Маппинг названий стран на стандартные
            countryMap := map[string]string{
                "Kina": "China",
                "SAD": "USA",
                // ...
            }
            return countryMap[v.(string)]
        },
    },
    "naAkciji": {
        InternalAttrID: 104,  // "on_sale"
        Transform: func(v interface{}) interface{} {
            return v == "1"
        },
    },
}
```

3. **Preview атрибутов перед импортом:**
```tsx
<AttributeMappingPreview>
  <DetectedAttributes>
    Обнаружено атрибутов в файле: 6

    <AttributeList>
      <Attribute status="mapped">
        ✅ uvoznik → Импортер (text)
        Примеры: "Digital Vision doo", "Digital Vision doo"
      </Attribute>

      <Attribute status="mapped">
        ✅ godinaUvoza → Год импорта (text)
        Примеры: "2025.", "2024."
      </Attribute>

      <Attribute status="mapped">
        ✅ zemljaPorekla → Страна происхождения (select)
        Уникальных значений: 5
        Топ: Kina (15,234), Vietnam (1,234), Taiwan (456)
      </Attribute>

      <Attribute status="new">
        ⚠️ kategorija1 → Не сопоставлен
        <Suggestion>
          AI предлагает: создать атрибут "Категория производителя"
          или использовать для маппинга категорий
        </Suggestion>
      </Attribute>
    </AttributeList>
  </DetectedAttributes>

  <Actions>
    <Button onClick={autoApplyMappings}>
      ✅ Применить все сопоставления
    </Button>
    <Button onClick={customizeMapping}>
      🔧 Настроить вручную
    </Button>
  </Actions>
</AttributeMappingPreview>
```

**Приоритет:** 🟡 ВАЖНО

---

#### 7. ❌ Автоматическая группировка в варианты
**Вопрос:** Можно ли автоматически группировать товары в варианты (например, 10 цветов → 1 карточка с 10 вариантами)?
**Ответ:** **НЕТ**, сейчас не реализовано, но КРИТИЧЕСКИ ВАЖНО!

**Текущая проблема:**
Digital Vision прайс содержит множество товаров-вариантов:
```
✅ Система вариантов есть в БД (storefront_product_variants)
✅ Структура поддерживает variant_attributes (JSONB)
✅ Есть has_variants flag и is_default variant
❌ НЕТ автоматической группировки при импорте
```

**Примеры из Digital Vision:**
```
Fidget Spinner - 5 цветов (crni, crveni, plavi, beli, zeleni)
→ Сейчас: 5 отдельных карточек товаров
→ Нужно: 1 карточка "Fidget Spinner" + 5 вариантов цвета

Narukvica za Apple Watch Silicone Strap - 175+ вариантов!
→ Разные цвета: dark blue, light yellow, camellia red, black, white...
→ Разные размеры: S/M, M/L
→ Разные модели часов: 38/40/41mm, 42/44/45/49mm
→ Сейчас: 175 отдельных карточек
→ Нужно: 1 карточка + 175 вариантов (color × size × watch_model)
```

**Алгоритм определения вариантов:**

```go
// backend/internal/proj/storefronts/service/variant_detector.go

type VariantDetector struct {
    colorPatterns []string
    sizePatterns  []string
    modelPatterns []string
}

type ProductGroup struct {
    BaseName     string
    BaseProduct  *models.ImportProductRequest
    Variants     []*ProductVariant
}

type ProductVariant struct {
    Product        *models.ImportProductRequest
    VariantAttrs   map[string]string  // {"color": "black", "size": "S/M"}
}

func (d *VariantDetector) GroupProducts(
    products []models.ImportProductRequest,
) []ProductGroup {
    groups := make(map[string]*ProductGroup)

    for _, product := range products {
        // 1. Извлечь base name (без цвета, размера и т.д.)
        baseName := d.extractBaseName(product.Name)

        // 2. Извлечь атрибуты варианта
        variantAttrs := d.extractVariantAttributes(product.Name)

        // 3. Группировать
        if group, exists := groups[baseName]; exists {
            // Добавить как вариант
            group.Variants = append(group.Variants, &ProductVariant{
                Product: &product,
                VariantAttrs: variantAttrs,
            })
        } else {
            // Создать новую группу
            groups[baseName] = &ProductGroup{
                BaseName: baseName,
                BaseProduct: &product,
                Variants: []*ProductVariant{
                    {Product: &product, VariantAttrs: variantAttrs},
                },
            }
        }
    }

    // 4. Фильтровать - группы с 1 вариантом это обычные товары
    result := []ProductGroup{}
    for _, group := range groups {
        if len(group.Variants) > 1 {
            result = append(result, *group)
        }
    }

    return result
}

func (d *VariantDetector) extractBaseName(productName string) string {
    name := productName

    // Убираем цвета
    colorRegex := regexp.MustCompile(`\s(crn[iao]|bel[iao]|crveni?|zeleni?|plav[iao]|pink|black|white|red|blue|green)\s?$`)
    name = colorRegex.ReplaceAllString(name, "")

    // Убираем размеры
    sizeRegex := regexp.MustCompile(`\s[SML]\/\s?[ML]\s`)
    name = sizeRegex.ReplaceAllString(name, "")

    // Убираем модели часов
    watchModelRegex := regexp.MustCompile(`\s\d+\/\s?\d+\/\s?\d+\s?mm`)
    name = watchModelRegex.ReplaceAllString(name, "")

    return strings.TrimSpace(name)
}

func (d *VariantDetector) extractVariantAttributes(productName string) map[string]string {
    attrs := make(map[string]string)

    // Извлекаем цвет
    if color := d.extractColor(productName); color != "" {
        attrs["color"] = color
    }

    // Извлекаем размер
    if size := d.extractSize(productName); size != "" {
        attrs["size"] = size
    }

    // Извлекаем модель
    if model := d.extractModel(productName); model != "" {
        attrs["model"] = model
    }

    return attrs
}
```

**Preview вариантов перед импортом:**
```tsx
<VariantDetectionPreview>
  <Summary>
    🔍 Обнаружено потенциальных групп вариантов: 1,234
    📦 Из 17,353 товаров можно сгруппировать: 8,456 (48.7%)

    Экономия карточек: 17,353 → 10,131 (-41.6%)
  </Summary>

  <VariantGroupsList>
    <VariantGroup confidence={0.98} productCount={175}>
      <BaseName>
        Narukvica za Apple Watch Silicone Strap
      </BaseName>

      <VariantDimensions>
        - Цвета: 35 вариантов (dark blue, light yellow, black, ...)
        - Размеры: 2 варианта (S/M, M/L)
        - Модели: 2 варианта (38/40/41mm, 42/44/45/49mm)

        Всего комбинаций: 35 × 2 × 2 = 140 вариантов
        Обнаружено в прайсе: 175 вариантов ✅
      </VariantDimensions>

      <PreviewVariants>
        Variant 1: color=dark blue, size=S/M, model=38/40/41mm
        Variant 2: color=dark blue, size=M/L, model=42/44/45/49mm
        Variant 3: color=light yellow, size=S/M, model=38/40/41mm
        ... (показать еще 172)
      </PreviewVariants>

      <Actions>
        <Button primary onClick={groupAsVariants}>
          ✅ Создать 1 товар с 175 вариантами
        </Button>
        <Button onClick={keepSeparate}>
          ❌ Оставить 175 отдельных карточек
        </Button>
        <Button onClick={customize}>
          🔧 Настроить вручную
        </Button>
      </Actions>
    </VariantGroup>

    <VariantGroup confidence={0.95} productCount={5}>
      <BaseName>Fidget Spinner</BaseName>
      <VariantDimensions>
        - Цвета: 5 вариантов (crni, crveni, plavi, beli, zeleni)
      </VariantDimensions>
      <!-- ... -->
    </VariantGroup>

    <!-- ... еще 1,232 группы -->
  </VariantGroupsList>

  <GlobalActions>
    <Button onClick={autoApplyAll}>
      ⚡ Автоматически сгруппировать все (confidence > 0.90)
    </Button>
    <Button onClick={reviewAll}>
      👀 Проверить все группы вручную
    </Button>
  </GlobalActions>
</VariantDetectionPreview>
```

**Реализация импорта вариантов:**
```go
func (s *ImportService) importProductGroup(
    ctx context.Context,
    storefrontID int,
    group ProductGroup,
) error {
    // 1. Создать базовый товар (parent product)
    baseProduct := &models.StorefrontProduct{
        StorefrontID:  storefrontID,
        Name:          group.BaseName,
        Description:   group.BaseProduct.Description,
        CategoryID:    group.BaseProduct.CategoryID,
        HasVariants:   true,  // ВАЖНО!
        // Цена и остатки берутся из default варианта
    }

    if err := s.productService.CreateProduct(ctx, baseProduct); err != nil {
        return err
    }

    // 2. Создать варианты
    for i, variant := range group.Variants {
        variantProduct := &models.StorefrontProductVariant{
            ProductID:         baseProduct.ID,
            SKU:               variant.Product.SKU,
            Barcode:           variant.Product.Barcode,
            Price:             variant.Product.Price,
            StockQuantity:     variant.Product.StockQuantity,
            VariantAttributes: variant.VariantAttrs,  // {"color": "black", "size": "S/M"}
            IsDefault:         i == 0,  // Первый вариант - default
        }

        if err := s.productService.CreateVariant(ctx, variantProduct); err != nil {
            return err
        }

        // 3. Загрузить изображения для варианта
        if len(variant.Product.ImageURLs) > 0 {
            err := s.importVariantImages(ctx, variantProduct.ID, variant.Product.ImageURLs)
            if err != nil {
                log.Printf("Failed to import variant images: %v", err)
            }
        }
    }

    return nil
}
```

**Приоритет:** 🔥 КРИТИЧЕСКИЙ

---

## 🚀 Обновленный план реализации

### Фаза 0: Подготовка и анализ (1 неделя)
**Цель:** Понять все категории и атрибуты Digital Vision

#### Задача 0.1: Полный анализ прайса (1 день)
```bash
# Скрипт анализа
python3 analyze_digital_vision.py --file DigitalVision.xml --output analysis.json

# Результат:
{
  "categories": {
    "total": 388,
    "level1": 7,
    "level2": 56,
    "level3": 388,
    "top_categories": [...]
  },
  "attributes": {
    "detected": ["uvoznik", "godinaUvoza", "zemljaPorekla", ...],
    "unique_values": {...}
  },
  "variants": {
    "potential_groups": 1234,
    "products_affected": 8456,
    "top_variant_patterns": [...]
  },
  "images": {
    "total_products_with_images": 14205,
    "percentage": 81.8%
  }
}
```

**Deliverables:**
- [ ] Полный список категорий с количеством товаров
- [ ] Список всех атрибутов и их значений
- [ ] Список потенциальных групп вариантов
- [ ] Статистика изображений

---

### Фаза 1: Умный Preview (2 недели)
**Цель:** Preview с AI маппингом категорий, атрибутов и вариантов

#### Задача 1.1: Backend - AI Category Mapper (3 дня)
**Файлы:**
- `backend/internal/proj/storefronts/service/ai_category_mapper.go`
- `backend/internal/proj/storefronts/service/ai_category_analyzer.go`

**API Endpoints:**
```go
POST /api/v1/storefronts/import/analyze-categories
→ Анализирует категории в файле
→ Возвращает предложения AI по маппингу

POST /api/v1/storefronts/import/analyze-attributes
→ Анализирует атрибуты в файле
→ Возвращает предложения по маппингу атрибутов

POST /api/v1/storefronts/import/detect-variants
→ Находит потенциальные группы вариантов
→ Возвращает предложения по группировке
```

**Критерии:**
- [ ] AI корректно мапит 90%+ категорий
- [ ] Confidence score точно отражает качество маппинга
- [ ] Предлагает новые категории для импорта в систему
- [ ] Находит все атрибуты в файле
- [ ] Группирует варианты с accuracy >95%

#### Задача 1.2: Frontend - Enhanced Preview UI (4 дня)
**Компоненты:**
```tsx
// frontend/svetu/src/components/import/ImportAnalysisWizard.tsx
// Многошаговый wizard:

Step 1: Upload File
  └─ Drag & Drop или URL

Step 2: File Analysis (Auto)
  └─ Парсинг + AI анализ

Step 3: Category Mapping
  ├─ Auto-approved (high confidence)
  ├─ Review recommended (medium confidence)
  └─ Manual required (low confidence)

Step 4: Attribute Mapping
  ├─ Detected attributes
  ├─ Suggested mappings
  └─ Create new attributes

Step 5: Variant Detection
  ├─ Detected variant groups
  ├─ Auto-group suggestions
  └─ Manual grouping editor

Step 6: Preview & Confirm
  ├─ Summary statistics
  ├─ Sample products
  └─ Start import button
```

**Критерии:**
- [ ] Wizard интуитивен и удобен
- [ ] Показывает прогресс анализа
- [ ] Позволяет редактировать любой маппинг
- [ ] Сохраняет маппинги для будущих импортов

#### Задача 1.3: Category Proposals System (2 дня)
```go
// backend/internal/proj/marketplace/service/category_management_service.go

func (s *CategoryService) ProposeNewCategory(
    proposal *models.NewCategoryProposal,
) (*models.CategoryProposal, error) {
    // 1. Создать proposal в БД (статус: pending)
    // 2. Назначить на review админам
    // 3. Отправить уведомление
}

func (s *CategoryService) ApproveProposal(
    proposalID int,
    approverUserID int,
) (*models.MarketplaceCategory, error) {
    // 1. Создать категорию
    // 2. Обновить proposal (статус: approved)
    // 3. Уведомить пользователя
}
```

**Таблица:**
```sql
CREATE TABLE category_proposals (
    id SERIAL PRIMARY KEY,
    proposed_by_user_id INT NOT NULL,
    storefront_id INT,
    name VARCHAR(255) NOT NULL,
    parent_category_id INT,
    description TEXT,
    reasoning TEXT,
    expected_products INT,
    external_category_source VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending',  -- pending, approved, rejected
    reviewed_by_user_id INT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Критерии:**
- [ ] Proposals сохраняются в БД
- [ ] Admin panel для review proposals
- [ ] Уведомления при approve/reject

---

### Фаза 2: Variant Import Engine (2 недели)
**Цель:** Автоматическая группировка товаров в варианты

#### Задача 2.1: Variant Detector (4 дня)
**Файлы:**
- `backend/internal/proj/storefronts/service/variant_detector.go`
- `backend/internal/proj/storefronts/service/variant_grouper.go`

**Алгоритм:**
```go
1. Извлечь базовые названия (без цвета/размера)
2. Группировать товары по базовому названию
3. Определить variant-defining attributes (color, size, model)
4. Создать variant groups
5. Валидировать (все варианты имеют одинаковые attrs?)
6. Вернуть с confidence score
```

**Критерии:**
- [ ] Корректно определяет базовое название
- [ ] Извлекает variant attributes
- [ ] Группирует >95% вариантов правильно
- [ ] Не группирует разные товары

#### Задача 2.2: Import с вариантами (3 дня)
```go
func (s *ImportService) importWithVariants(
    ctx context.Context,
    req models.ImportRequest,
    variantGroups []ProductGroup,
) (*models.ImportJob, error) {
    for _, group := range variantGroups {
        if len(group.Variants) > 1 {
            // Импорт как группа вариантов
            s.importProductGroup(ctx, req.StorefrontID, group)
        } else {
            // Обычный импорт (один товар)
            s.createProduct(ctx, req.StorefrontID, group.BaseProduct)
        }
    }
}
```

**Критерии:**
- [ ] Создает parent product с has_variants=true
- [ ] Создает все варианты
- [ ] Загружает изображения для каждого варианта
- [ ] Корректно устанавливает is_default
- [ ] Сохраняет variant_attributes

#### Задача 2.3: Variant Preview UI (3 дня)
```tsx
<VariantGroupPreview group={group}>
  <GroupHeader>
    {group.baseName}
    <Badge>{group.variants.length} вариантов</Badge>
  </GroupHeader>

  <VariantTable>
    {group.variants.map(v => (
      <VariantRow>
        <Image src={v.image} />
        <Attributes>
          {Object.entries(v.variantAttrs).map(([k, v]) => (
            <Chip>{k}: {v}</Chip>
          ))}
        </Attributes>
        <Price>{v.price}</Price>
        <Stock>{v.stock}</Stock>
      </VariantRow>
    ))}
  </VariantTable>

  <Actions>
    <Button onClick={confirmGroup}>✅ Группировать</Button>
    <Button onClick={editGroup}>✏️ Редактировать</Button>
    <Button onClick={splitGroup}>❌ Разделить</Button>
  </Actions>
</VariantGroupPreview>
```

**Критерии:**
- [ ] Показывает все варианты
- [ ] Позволяет редактировать группировку
- [ ] Позволяет исключить отдельные варианты
- [ ] Preview итоговой карточки

---

### Фаза 3: Attribute System (1 неделя)
**Цель:** Полноценный маппинг и управление атрибутами

#### Задача 3.1: Attribute Mapper (2 дня)
```go
type AttributeMapper struct {
    attributeTemplates map[string]*AttributeTemplate
}

func (m *AttributeMapper) MapExternalAttribute(
    externalName string,
    externalValue interface{},
    productCategory int,
) (*MappedAttribute, error) {
    // 1. Найти подходящий attribute template
    template := m.findMatchingTemplate(externalName, productCategory)

    // 2. Трансформировать значение
    value := m.transformValue(externalValue, template.Type)

    // 3. Валидировать
    if err := m.validateValue(value, template); err != nil {
        return nil, err
    }

    return &MappedAttribute{
        AttributeID: template.ID,
        Value: value,
    }, nil
}
```

**Критерии:**
- [ ] Мапит стандартные атрибуты
- [ ] Трансформирует значения (типы, форматы)
- [ ] Валидирует значения
- [ ] Создает новые атрибуты при необходимости

#### Задача 3.2: Attribute Preview UI (2 дня)
```tsx
<AttributeMappingStep>
  {detectedAttributes.map(attr => (
    <AttributeMapping key={attr.name}>
      <External>
        {attr.name}: {attr.sampleValues.slice(0, 3).join(', ')}
        <Badge>{attr.uniqueValues} unique values</Badge>
      </External>

      <MappingArrow />

      <Internal>
        {attr.suggestedMapping ? (
          <Select
            value={attr.suggestedMapping.id}
            options={availableAttributes}
            onChange={updateMapping}
          />
        ) : (
          <CreateNewAttribute
            defaultName={attr.name}
            onCreate={createAndMap}
          />
        )}
      </Internal>

      <Confidence>{attr.mappingConfidence}</Confidence>
    </AttributeMapping>
  ))}
</AttributeMappingStep>
```

**Критерии:**
- [ ] Показывает все обнаруженные атрибуты
- [ ] Предлагает маппинг на существующие
- [ ] Позволяет создать новые атрибуты
- [ ] Показывает примеры значений

---

### Фаза 4: Production Ready (1 неделя)
**Цель:** Тестирование, оптимизация, документация

#### Задача 4.1: Полное тестирование (3 дня)
```bash
# 1. Импорт полного прайса Digital Vision (17,353 товаров)
# 2. Проверка результатов:
#    - Все категории смаппированы
#    - Все атрибуты сохранены
#    - Варианты сгруппированы
#    - Изображения загружены
# 3. Performance тесты
# 4. Stress тесты
```

**Критерии:**
- [ ] Импорт 17K товаров < 15 минут
- [ ] Accuracy категоризации >95%
- [ ] Accuracy группировки вариантов >95%
- [ ] Все изображения загружены
- [ ] Нет memory leaks

#### Задача 4.2: Документация (2 дня)
**Для Digital Vision:**
- [ ] Quick Start Guide (как сделать первый импорт)
- [ ] Category Mapping Guide
- [ ] Variant Grouping Guide
- [ ] Scheduled Import Setup

**Для других клиентов:**
- [ ] Generic Import Guide
- [ ] Supported Formats
- [ ] Troubleshooting

#### Задача 4.3: Мониторинг и алерты (1 день)
```go
// Метрики импорта
type ImportMetrics struct {
    TotalProducts      int
    ImportedProducts   int
    FailedProducts     int
    VariantGroups      int
    CategoriesMapped   int
    AttributesMapped   int
    ImagesDownloaded   int
    ProcessingTime     time.Duration
}

// Алерты
- Импорт завис (>30 минут для 10K товаров)
- Высокий процент ошибок (>5%)
- Низкая confidence категоризации (<0.8 average)
- Не удалось загрузить изображения (>20%)
```

---

## 📊 Итоговая статистика по плану

### Временные рамки
- **Фаза 0:** 1 неделя (анализ)
- **Фаза 1:** 2 недели (умный preview)
- **Фаза 2:** 2 недели (варианты)
- **Фаза 3:** 1 неделя (атрибуты)
- **Фаза 4:** 1 неделя (production ready)

**Итого:** 7 недель (~1.5 месяца)

### Приоритеты задач
🔥 **КРИТИЧЕСКИЕ (must-have для Digital Vision):**
1. AI маппинг категорий с preview
2. Автоматическая группировка вариантов
3. Маппинг атрибутов
4. Preview перед импортом

🟡 **ВАЖНЫЕ (nice-to-have):**
1. AI предложения новых категорий
2. Scheduled импорт
3. Webhook триггеры

🟢 **ЖЕЛАТЕЛЬНЫЕ (future):**
1. Инкрементальный импорт
2. Batch обработка изображений
3. Advanced analytics

### Ожидаемые результаты

**Для Digital Vision:**
- ✅ Импорт 17,353 товаров за 10-15 минут
- ✅ 95%+ автоматическая категоризация
- ✅ ~8,500 товаров сгруппированы в ~1,200 карточек с вариантами
- ✅ Экономия карточек: -41% (17,353 → 10,100)
- ✅ Все атрибуты сохранены и структурированы
- ✅ Все изображения загружены в S3

**Для платформы:**
- ✅ Универсальная система импорта для любых клиентов
- ✅ AI-powered категоризация
- ✅ Автоматическое определение вариантов
- ✅ Гибкий маппинг атрибутов
- ✅ Референс для других крупных клиентов

---

## 📝 Документация и примеры

### Пример полного импорта
```typescript
// 1. Загрузка файла
const file = await uploadFile('DigitalVision.zip');

// 2. Анализ
const analysis = await analyzeImportFile(file);
// {
//   categories: 388,
//   products: 17353,
//   variants_detected: 1234 groups,
//   attributes: 6,
//   confidence: {
//     high: 320 categories (82%),
//     medium: 50 categories (13%),
//     low: 18 categories (5%)
//   }
// }

// 3. Review и корректировка
const reviewed = await reviewCategoryMappings(analysis.categoryMappings);
const confirmedVariants = await reviewVariantGroups(analysis.variantGroups);

// 4. Импорт
const job = await startImport({
  file,
  categoryMappings: reviewed.categoryMappings,
  variantGroups: confirmedVariants,
  attributeMappings: reviewed.attributeMappings,
  updateMode: 'upsert'
});

// 5. Мониторинг
watchImportProgress(job.id);
```

---

**Статус:** 📋 Готов к реализации
**Next Step:** Начать Фазу 0 (анализ прайса Digital Vision)
**Дата последнего обновления:** 2025-10-06
