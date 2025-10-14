// backend/internal/proj/c2c/storage/opensearch/repository_helpers.go
package opensearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

func (r *Repository) getBoostWeight(weightName string, defaultValue float64) float64 {
	if r.boostWeights == nil {
		return defaultValue
	}

	switch weightName {
	case "Title":
		return r.boostWeights.Title
	case "TitleNgram":
		return r.boostWeights.TitleNgram
	case "Description":
		return r.boostWeights.Description
	case "TranslationTitle":
		return r.boostWeights.TranslationTitle
	case "TranslationDesc":
		return r.boostWeights.TranslationDesc
	case "AttributeTextValue":
		return r.boostWeights.AttributeTextValue
	case "AttributeDisplayValue":
		return r.boostWeights.AttributeDisplayValue
	case "AttributeTextValueKeyword":
		return r.boostWeights.AttributeTextValueKeyword
	case "AttributeGeneralBoost":
		return r.boostWeights.AttributeGeneralBoost
	case "RealEstateAttributesCombined":
		return r.boostWeights.RealEstateAttributesCombined
	case "PropertyType":
		return r.boostWeights.PropertyType
	case "RoomsText":
		return r.boostWeights.RoomsText
	case "CarMake":
		return r.boostWeights.CarMake
	case "CarModel":
		return r.boostWeights.CarModel
	case "CarKeywords":
		return r.boostWeights.CarKeywords
	case "PerWordTitle":
		return r.boostWeights.PerWordTitle
	case "PerWordDescription":
		return r.boostWeights.PerWordDescription
	case "PerWordAllAttributes":
		return r.boostWeights.PerWordAllAttributes
	case "PerWordRealEstateAttributes":
		return r.boostWeights.PerWordRealEstateAttributes
	case "PerWordRoomsText":
		return r.boostWeights.PerWordRoomsText
	case "AutomotiveAttributePriority":
		return defaultValue // Automotive system removed
	case "SynonymBoost":
		return r.boostWeights.SynonymBoost
	default:
		return defaultValue
	}
}

func extractSuggestionsFromAgg(aggs map[string]interface{}, aggName string, suggestions map[string]bool) {
	if agg, ok := aggs[aggName].(map[string]interface{}); ok {
		extractBucketsFromAgg(agg, suggestions)
	}
}

// Вспомогательная функция для извлечения бакетов из агрегации
func extractBucketsFromAgg(agg map[string]interface{}, suggestions map[string]bool) {
	if buckets, ok := agg["buckets"].([]interface{}); ok {
		for _, bucket := range buckets {
			if bucketObj, ok := bucket.(map[string]interface{}); ok {
				if key, ok := bucketObj["key"].(string); ok && key != "" {
					suggestions[key] = true
				}
			}
		}
	}
}

// ReindexAll переиндексирует все объявления

func (r *Repository) getAttributeOptionTranslations(ctx context.Context, attrName, value string) (map[string]string, error) {
	query := `
        SELECT option_value, ru_translation, sr_translation
        FROM attribute_option_translations
        WHERE attribute_name = $1 AND option_value = $2
    `

	var optionValue, ruTranslation, srTranslation string
	err := r.storage.QueryRow(ctx, query, attrName, value).Scan(
		&optionValue, &ruTranslation, &srTranslation,
	)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logger.Info().Msgf("Ошибка получения переводов для атрибута %s, значение %s: %v", attrName, value, err)
		}
		return nil, err
	}

	translations := map[string]string{
		"ru": ruTranslation,
		"sr": srTranslation,
	}

	return translations, nil
}

// getListingTranslationsFromDB загружает переводы для объявления из таблицы translations
func (r *Repository) getListingTranslationsFromDB(ctx context.Context, listingID int) ([]DBTranslation, error) {
	query := `
		SELECT language, field_name, translated_text 
		FROM translations 
		WHERE entity_type = 'listing' AND entity_id = $1
		ORDER BY language, field_name
	`

	rows, err := r.storage.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса переводов: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var translations []DBTranslation
	for rows.Next() {
		var t DBTranslation
		err := rows.Scan(&t.Language, &t.FieldName, &t.TranslatedText)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования перевода: %w", err)
		}
		translations = append(translations, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по переводам: %w", err)
	}

	return translations, nil
}

// convertDBTranslationsToMap преобразует массив DBTranslation в структуру map[язык]map[поле]значение
func (r *Repository) convertDBTranslationsToMap(translations []DBTranslation) map[string]map[string]string {
	result := make(map[string]map[string]string)

	for _, t := range translations {
		if _, exists := result[t.Language]; !exists {
			result[t.Language] = make(map[string]string)
		}
		result[t.Language][t.FieldName] = t.TranslatedText
	}

	return result
}

// extractSupportedLanguages извлекает список поддерживаемых языков из переводов
func (r *Repository) extractSupportedLanguages(translations []DBTranslation) []string {
	langMap := make(map[string]bool)
	for _, t := range translations {
		langMap[t.Language] = true
	}

	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}

	return languages
}

func hasAttributeValue(attr models.ListingAttributeValue) bool {
	return (attr.TextValue != nil && *attr.TextValue != "") ||
		(attr.NumericValue != nil) ||
		(attr.BooleanValue != nil) ||
		(attr.JSONValue != nil && string(attr.JSONValue) != "{}" && string(attr.JSONValue) != "[]") ||
		attr.DisplayValue != ""
}

func isImportantTextAttribute(attrName string) bool {
	importantTextAttrs := map[string]bool{
		"brand":         true,
		"color":         true,
		"fuel_type":     true,
		"transmission":  true,
		"body_type":     true,
		"property_type": true,
	}
	return importantTextAttrs[attrName]
}

func formatAttributeDisplayValue(attr models.ListingAttributeValue) string {
	numVal := *attr.NumericValue
	unitStr := attr.Unit

	if unitStr == "" {
		switch attr.AttributeName {
		case "area":
			unitStr = "m²"
		case "land_area":
			unitStr = "ar"
		case "mileage":
			unitStr = "km"
		case "engine_capacity":
			unitStr = "l"
		case "power":
			unitStr = "ks"
		case "screen_size":
			unitStr = "inč"
		case "rooms":
			unitStr = "soba"
		case "floor", "total_floors":
			unitStr = "sprat"
		}
	}

	if attr.AttributeName == "year" {
		return fmt.Sprintf("%d", int(numVal))
	} else if unitStr != "" {
		return fmt.Sprintf("%g %s", numVal, unitStr)
	}
	return fmt.Sprintf("%g", numVal)
}

func addRangesForAttribute(doc map[string]interface{}, attr models.ListingAttributeValue) {
	if attr.NumericValue == nil {
		return
	}
	numVal := *attr.NumericValue

	switch attr.AttributeName {
	case fieldNamePrice:
		doc["price_range"] = getPriceRange(int(numVal))
	case "mileage":
		doc["mileage_range"] = getMileageRange(int(numVal))
	case "area":
		switch {
		case numVal <= 30:
			doc["area_range"] = "do 30 m²"
		case numVal <= 50:
			doc["area_range"] = "30-50 m²"
		case numVal <= 80:
			doc["area_range"] = "50-80 m²"
		case numVal <= 120:
			doc["area_range"] = "80-120 m²"
		default:
			doc["area_range"] = "od 120 m²"
		}
	}
}

func createAttributeDocument(attr models.ListingAttributeValue) map[string]interface{} {
	attrDoc := map[string]interface{}{
		"attribute_id":   attr.AttributeID,
		"attribute_name": attr.AttributeName,
		"display_name":   attr.DisplayName,
		"attribute_type": attr.AttributeType,
		"display_value":  attr.DisplayValue,
	}

	if attr.TextValue != nil && *attr.TextValue != "" {
		textValue := *attr.TextValue
		attrDoc["text_value"] = textValue
		attrDoc["text_value_lowercase"] = strings.ToLower(textValue)
	}

	if attr.NumericValue != nil && !math.IsNaN(*attr.NumericValue) && !math.IsInf(*attr.NumericValue, 0) {
		attrDoc["numeric_value"] = *attr.NumericValue
		if attr.Unit != "" {
			attrDoc["unit"] = attr.Unit
		}
	}

	if attr.BooleanValue != nil {
		attrDoc["boolean_value"] = *attr.BooleanValue
	}

	if attr.JSONValue != nil {
		jsonStr := string(attr.JSONValue)
		if jsonStr != "" && jsonStr != "{}" && jsonStr != "[]" {
			attrDoc["json_value"] = jsonStr
		}
	}

	return attrDoc
}

func ensureImportantAttributes(doc map[string]interface{}, makeValue, modelValue string, listingID int) {
	if makeValue != "" && doc["make"] == nil {
		doc["make"] = makeValue
		doc["make_lowercase"] = strings.ToLower(makeValue)
		logger.Info().Msgf("FINAL CHECK: Добавлена марка '%s' в корень документа для объявления %d",
			makeValue, listingID)
	}

	if modelValue != "" && doc["model"] == nil {
		doc["model"] = modelValue
		doc["model_lowercase"] = strings.ToLower(modelValue)
		logger.Info().Msgf("FINAL CHECK: Добавлена модель '%s' в корень документа для объявления %d",
			modelValue, listingID)
	}
}

func getUniqueValues(values []string) []string {
	seen := make(map[string]bool)
	unique := make([]string, 0, len(values))

	for _, val := range values {
		if !seen[val] {
			seen[val] = true
			unique = append(unique, val)
		}
	}

	return unique
}

func flattenAttributeValues(attributeTextValues map[string][]string) []string {
	var result []string
	for _, values := range attributeTextValues {
		result = append(result, values...)
	}
	return result
}

func createRealEstateFieldsMap() map[string]bool {
	return map[string]bool{
		"rooms":            true,
		"floor":            true,
		"total_floors":     true,
		"area":             true,
		"land_area":        true,
		"property_type":    true,
		"year_built":       true,
		"bathrooms":        true,
		"condition":        true,
		"amenities":        true,
		"heating_type":     true,
		"parking":          true,
		"balcony":          true,
		"furnished":        true,
		"air_conditioning": true,
	}
}

func createCarFieldsMap() map[string]bool {
	return map[string]bool{
		"make":            true,
		"model":           true,
		"year":            true,
		"mileage":         true,
		"engine_capacity": true,
		"fuel_type":       true,
		"transmission":    true,
		"body_type":       true,
	}
}

func createImportantAttributesMap() map[string]bool {
	return map[string]bool{
		"make":            true,
		"model":           true,
		"brand":           true,
		"year":            true,
		"color":           true,
		"rooms":           true,
		"property_type":   true,
		"body_type":       true,
		"engine_capacity": true,
		"fuel_type":       true,
		"transmission":    true,
		"cpu":             true,
		"gpu":             true,
		"memory":          true,
		"ram":             true,
		"storage_type":    true,
		"screen_size":     true,
	}
}

func getMileageRange(mileage int) string {
	switch {
	case mileage <= 5000:
		return "0-5000"
	case mileage <= 10000:
		return "5001-10000"
	case mileage <= 50000:
		return "10001-50000"
	case mileage <= 100000:
		return "50001-100000"
	case mileage <= 150000:
		return "100001-150000"
	case mileage <= 200000:
		return "150001-200000"
	default:
		return "200001+"
	}
}

func getPriceRange(price int) string {
	switch {
	case price <= 5000:
		return "0-5000"
	case price <= 10000:
		return "5001-10000"
	case price <= 20000:
		return "10001-20000"
	case price <= 50000:
		return "20001-50000"
	case price <= 100000:
		return "50001-100000"
	case price <= 500000:
		return "100001-500000"
	default:
		return "500001+"
	}
}

func (r *Repository) geocodeCity(city, country string) (*struct{ Lat, Lon float64 }, error) {
	query := city
	if country != "" {
		query += ", " + country
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "HostelBookingSystem/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный статус ответа: %d", resp.StatusCode)
	}

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("город не найден")
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return nil, err
	}

	return &struct{ Lat, Lon float64 }{lat, lon}, nil
}

func getDocKeys(doc map[string]interface{}) []string {
	keys := make([]string, 0, len(doc))
	for k := range doc {
		keys = append(keys, k)
	}
	return keys
}
