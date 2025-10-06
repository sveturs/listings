package service

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// VariantAttribute представляет атрибут варианта (цвет, размер и т.д.)
type VariantAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductVariant представляет один вариант товара
type ProductVariant struct {
	Name               string                      `json:"name"`
	SKU                string                      `json:"sku"`
	VariantAttributes  map[string]string           `json:"variant_attributes"`
	Price              float64                     `json:"price,omitempty"`
	StockQuantity      int                         `json:"stock_quantity,omitempty"`
	ImageURL           string                      `json:"image_url,omitempty"`
	OriginalAttributes map[string]interface{}      `json:"original_attributes,omitempty"`
}

// VariantGroup представляет группу вариантов товара
type VariantGroup struct {
	BaseName           string                 `json:"base_name"`
	BaseProduct        *ProductVariant        `json:"base_product"`
	Variants           []*ProductVariant      `json:"variants"`
	VariantCount       int                    `json:"variant_count"`
	VariantAttributes  []string               `json:"variant_attributes"`  // Названия атрибутов (color, size)
	Confidence         float64                `json:"confidence"`
	IsGrouped          bool                   `json:"is_grouped"`
}

// VariantDetector обнаруживает и группирует варианты товаров
type VariantDetector struct {
	// Паттерны для цветов
	colorPatterns []*regexp.Regexp

	// Паттерны для размеров
	sizePatterns []*regexp.Regexp

	// Паттерны для моделей
	modelPatterns []*regexp.Regexp

	// Минимальная уверенность для группировки
	minConfidence float64

	// Минимальный размер группы
	minGroupSize int
}

// NewVariantDetector создает новый детектор вариантов
func NewVariantDetector() *VariantDetector {
	return &VariantDetector{
		colorPatterns: compileColorPatterns(),
		sizePatterns:  compileSizePatterns(),
		modelPatterns: compileModelPatterns(),
		minConfidence: 0.7,
		minGroupSize:  2,
	}
}

// compileColorPatterns компилирует регулярные выражения для определения цветов
func compileColorPatterns() []*regexp.Regexp {
	patterns := []string{
		// Русские цвета (все падежи и формы) - используем (?:^|\s) и (?:\s|$) вместо \b
		`(?:^|\s)(черн(?:ый|ая|ое|ые|ого|ому|ым|ом)|бел(?:ый|ая|ое|ые|ого|ому|ым|ом)|красн(?:ый|ая|ое|ые|ого|ому|ым|ом)|син(?:ий|яя|ее|ие|его|ему|им|ем)|зелен(?:ый|ая|ое|ые|ого|ому|ым|ом)|желт(?:ый|ая|ое|ые|ого|ому|ым|ом)|оранжев(?:ый|ая|ое|ые|ого|ому|ым|ом)|фиолетов(?:ый|ая|ое|ые|ого|ому|ым|ом)|розов(?:ый|ая|ое|ые|ого|ому|ым|ом)|коричнев(?:ый|ая|ое|ые|ого|ому|ым|ом)|сер(?:ый|ая|ое|ые|ого|ому|ым|ом)|голуб(?:ой|ая|ое|ые|ого|ому|ым|ом)|бежев(?:ый|ая|ое|ые|ого|ому|ым|ом)|бордов(?:ый|ая|ое|ые|ого|ому|ым|ом))(?:\s|$)`,
		// Английские цвета
		`\b(black|white|red|blue|green|yellow|orange|purple|pink|brown|gray|grey|beige|burgundy|navy|cream|ivory)\b`,
		// Сербские цвета
		`(?:^|\s)(crn(?:a|i|e|o)|bel(?:a|i|e|o)|crven(?:a|i|e|o)|plav(?:a|i|e|o)|zelen(?:a|i|e|o)|žut(?:a|i|e|o)|narandžast(?:a|i|e|o)|ljubičast(?:a|i|e|o)|roze|braon|siv(?:a|i|e|o)|svetloplav(?:a|i|e|o)|bež|bordo)(?:\s|$)`,
		// Варианты написания
		`(?:^|\s)(чёрн(?:ый|ая|ое|ые)|тёмн(?:ый|ая|ое|ые)|светл(?:ый|ая|ое|ые)|тёмно[-\s]?син(?:ий|яя|ее|ие)|светло[-\s]?син(?:ий|яя|ее|ие))(?:\s|$)`,
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile("(?i)" + p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

// compileSizePatterns компилирует регулярные выражения для определения размеров
func compileSizePatterns() []*regexp.Regexp {
	patterns := []string{
		// Стандартные размеры одежды
		`\b(XXS|XS|S|M|L|XL|XXL|XXXL)\b`,
		// Числовые размеры
		`\b\d+\s?(см|mm|ml|л|kg|кг|г|g)\b`,
		// Европейские размеры одежды
		`\b(размер|size)?\s?\d{2,3}\b`,
		// Размеры обуви
		`\b(размер|size)?\s?\d{1,2}(\.5)?\b`,
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile("(?i)" + p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

// compileModelPatterns компилирует регулярные выражения для определения моделей/версий
func compileModelPatterns() []*regexp.Regexp {
	patterns := []string{
		// Версии и модели
		`\b(v\d+|ver\d+|версия\s?\d+|модель\s?\d+)\b`,
		`\b(model|модель)\s+[A-Z0-9\-]+\b`,
		// Годы
		`\b(20\d{2})\b`,
		// Поколения
		`\b(\d+(st|nd|rd|th)\s+generation|поколение)\b`,
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile("(?i)" + p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

// ExtractBaseName извлекает базовое название товара, убирая вариантные атрибуты
func (vd *VariantDetector) ExtractBaseName(productName string) string {
	name := strings.TrimSpace(productName)

	// Убираем цвета (заменяем match + окружение пробелами на один пробел)
	for _, pattern := range vd.colorPatterns {
		name = pattern.ReplaceAllString(name, " ")
	}

	// Убираем размеры
	for _, pattern := range vd.sizePatterns {
		name = pattern.ReplaceAllString(name, " ")
	}

	// Убираем модели
	for _, pattern := range vd.modelPatterns {
		name = pattern.ReplaceAllString(name, " ")
	}

	// Очищаем лишние пробелы
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")

	// Убираем trailing знаки пунктуации и пробелы
	name = strings.TrimRight(name, " ,-/")
	name = strings.TrimSpace(name)

	return name
}

// ExtractVariantAttributes извлекает вариантные атрибуты из названия товара
func (vd *VariantDetector) ExtractVariantAttributes(productName string) map[string]string {
	attributes := make(map[string]string)

	// Проверяем цвета (используем оригинальное название для матчинга)
	for _, pattern := range vd.colorPatterns {
		if submatches := pattern.FindStringSubmatch(productName); len(submatches) > 1 {
			// submatches[0] - полное совпадение (со spaces)
			// submatches[1] - группа захвата (цвет без spaces)
			attributes["color"] = strings.ToLower(strings.TrimSpace(submatches[1]))
			break
		}
	}

	// Проверяем размеры
	for _, pattern := range vd.sizePatterns {
		if submatches := pattern.FindStringSubmatch(productName); len(submatches) > 1 {
			attributes["size"] = strings.ToLower(strings.TrimSpace(submatches[1]))
			break
		}
	}

	// Проверяем модели
	for _, pattern := range vd.modelPatterns {
		if submatches := pattern.FindStringSubmatch(productName); len(submatches) > 1 {
			attributes["model"] = strings.ToLower(strings.TrimSpace(submatches[1]))
			break
		}
	}

	return attributes
}

// GroupProducts группирует товары в варианты на основе базового названия
func (vd *VariantDetector) GroupProducts(products []*ProductVariant) []*VariantGroup {
	// Группируем по базовому названию
	groupMap := make(map[string][]*ProductVariant)

	for _, product := range products {
		baseName := vd.ExtractBaseName(product.Name)
		if baseName == "" {
			baseName = product.Name
		}
		groupMap[baseName] = append(groupMap[baseName], product)
	}

	// Создаем группы вариантов
	groups := make([]*VariantGroup, 0)

	for baseName, variants := range groupMap {
		// Пропускаем слишком маленькие группы
		if len(variants) < vd.minGroupSize {
			// Создаем одиночный товар (не вариант)
			groups = append(groups, &VariantGroup{
				BaseName:     baseName,
				BaseProduct:  variants[0],
				Variants:     variants,
				VariantCount: 1,
				Confidence:   1.0,
				IsGrouped:    false,
			})
			continue
		}

		// Извлекаем атрибуты для всех вариантов
		for _, variant := range variants {
			variant.VariantAttributes = vd.ExtractVariantAttributes(variant.Name)
		}

		// Собираем уникальные типы атрибутов
		attrTypes := make(map[string]bool)
		attrCount := 0
		for _, variant := range variants {
			if len(variant.VariantAttributes) > 0 {
				attrCount++
				for attrName := range variant.VariantAttributes {
					attrTypes[attrName] = true
				}
			}
		}

		// Вычисляем уверенность: процент товаров с вариантными атрибутами
		confidence := float64(attrCount) / float64(len(variants))

		// Проверяем минимальную уверенность
		if confidence < vd.minConfidence {
			// Недостаточная уверенность - создаем одиночные товары
			for _, variant := range variants {
				groups = append(groups, &VariantGroup{
					BaseName:     variant.Name,
					BaseProduct:  variant,
					Variants:     []*ProductVariant{variant},
					VariantCount: 1,
					Confidence:   0.0,
					IsGrouped:    false,
				})
			}
			continue
		}

		// Преобразуем map в slice
		attrTypesList := make([]string, 0, len(attrTypes))
		for attrName := range attrTypes {
			attrTypesList = append(attrTypesList, attrName)
		}
		sort.Strings(attrTypesList)

		// Выбираем базовый товар (первый с изображением или просто первый)
		baseProduct := variants[0]
		for _, v := range variants {
			if v.ImageURL != "" {
				baseProduct = v
				break
			}
		}

		// Создаем группу вариантов
		group := &VariantGroup{
			BaseName:          baseName,
			BaseProduct:       baseProduct,
			Variants:          variants,
			VariantCount:      len(variants),
			VariantAttributes: attrTypesList,
			Confidence:        confidence,
			IsGrouped:         true,
		}

		groups = append(groups, group)
	}

	// Сортируем группы по количеству вариантов (по убыванию)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].VariantCount > groups[j].VariantCount
	})

	return groups
}

// ValidateVariantGroup проверяет валидность группы вариантов
func (vd *VariantDetector) ValidateVariantGroup(group *VariantGroup) []string {
	warnings := make([]string, 0)

	if group.VariantCount < 2 {
		return warnings
	}

	// Проверяем, что все варианты имеют одинаковый набор типов атрибутов
	expectedAttrs := make(map[string]bool)
	for _, attrType := range group.VariantAttributes {
		expectedAttrs[attrType] = true
	}

	for i, variant := range group.Variants {
		// Проверка наличия всех атрибутов
		for attrType := range expectedAttrs {
			if _, exists := variant.VariantAttributes[attrType]; !exists {
				warnings = append(warnings,
					fmt.Sprintf("Variant %d (%s) missing attribute '%s'", i+1, variant.SKU, attrType))
			}
		}

		// Проверка дублирования атрибутов
		for j := i + 1; j < len(group.Variants); j++ {
			other := group.Variants[j]
			if vd.hasSameAttributes(variant, other) {
				warnings = append(warnings,
					fmt.Sprintf("Variants %s and %s have identical attributes", variant.SKU, other.SKU))
			}
		}
	}

	return warnings
}

// hasSameAttributes проверяет, имеют ли два варианта идентичные атрибуты
func (vd *VariantDetector) hasSameAttributes(v1, v2 *ProductVariant) bool {
	if len(v1.VariantAttributes) != len(v2.VariantAttributes) {
		return false
	}

	for k, v := range v1.VariantAttributes {
		if v2.VariantAttributes[k] != v {
			return false
		}
	}

	return true
}

// SetMinConfidence устанавливает минимальную уверенность для группировки
func (vd *VariantDetector) SetMinConfidence(confidence float64) {
	if confidence >= 0 && confidence <= 1 {
		vd.minConfidence = confidence
	}
}

// SetMinGroupSize устанавливает минимальный размер группы
func (vd *VariantDetector) SetMinGroupSize(size int) {
	if size >= 1 {
		vd.minGroupSize = size
	}
}
