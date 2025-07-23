package service

import (
	"context"
	"testing"
)

// mockTranslationServiceWithAddresses - мок сервиса перевода для тестирования адресов
type mockTranslationServiceWithAddresses struct {
	translateEntityFieldsCalled bool
	translatedFields            map[string]map[string]string
}

func (m *mockTranslationServiceWithAddresses) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	// Простой мок - возвращает тот же текст с префиксом языка
	return targetLanguage + "_" + text, nil
}

func (m *mockTranslationServiceWithAddresses) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	return "ru", 0.99, nil
}

func (m *mockTranslationServiceWithAddresses) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	return map[string]string{
		"en": "en_" + text,
		"ru": text,
		"sr": "sr_" + text,
	}, nil
}

func (m *mockTranslationServiceWithAddresses) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	m.translateEntityFieldsCalled = true

	// Создаем переводы для всех целевых языков
	result := make(map[string]map[string]string)

	// Добавляем исходный язык
	result[sourceLanguage] = fields

	// Добавляем переводы для целевых языков
	for _, lang := range targetLanguages {
		langFields := make(map[string]string)
		for fieldName, fieldValue := range fields {
			langFields[fieldName] = lang + "_" + fieldValue
		}
		result[lang] = langFields
	}

	m.translatedFields = result
	return result, nil
}

func (m *mockTranslationServiceWithAddresses) ModerateText(ctx context.Context, text string, language string) (string, error) {
	return text, nil
}

func (m *mockTranslationServiceWithAddresses) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	return targetLanguage + "_" + text, nil
}

// TestTranslateEntityFields тестирует правильность вызова метода перевода полей
func TestTranslateEntityFields(t *testing.T) {
	// Создаем мок сервис перевода
	mockTranslation := &mockTranslationServiceWithAddresses{}

	// Тестовые данные
	addressFields := map[string]string{
		"location": "Улица Пушкина, дом 10",
		"city":     "Белград",
		"country":  "Сербия",
	}
	sourceLanguage := "ru"
	targetLanguages := []string{"en", "sr"}

	// Выполняем тест
	translations, err := mockTranslation.TranslateEntityFields(context.Background(), sourceLanguage, targetLanguages, addressFields)
	// Проверяем результат
	if err != nil {
		t.Errorf("TranslateEntityFields returned error: %v", err)
	}

	// Проверяем, что метод был вызван
	if !mockTranslation.translateEntityFieldsCalled {
		t.Error("TranslateEntityFields was not called")
	}

	// Проверяем структуру переводов
	if len(translations) != 3 { // ru, en, sr
		t.Errorf("Expected 3 languages in translations, got %d", len(translations))
	}

	// Проверяем английские переводы
	if enFields, ok := translations["en"]; ok {
		if enFields["location"] != "en_Улица Пушкина, дом 10" {
			t.Errorf("Expected English location translation, got %s", enFields["location"])
		}
		if enFields["city"] != "en_Белград" {
			t.Errorf("Expected English city translation, got %s", enFields["city"])
		}
		if enFields["country"] != "en_Сербия" {
			t.Errorf("Expected English country translation, got %s", enFields["country"])
		}
	} else {
		t.Error("English translations not found")
	}

	// Проверяем сербские переводы
	if srFields, ok := translations["sr"]; ok {
		if srFields["location"] != "sr_Улица Пушкина, дом 10" {
			t.Errorf("Expected Serbian location translation, got %s", srFields["location"])
		}
		if srFields["city"] != "sr_Белград" {
			t.Errorf("Expected Serbian city translation, got %s", srFields["city"])
		}
		if srFields["country"] != "sr_Сербия" {
			t.Errorf("Expected Serbian country translation, got %s", srFields["country"])
		}
	} else {
		t.Error("Serbian translations not found")
	}

	// Проверяем исходный язык
	if ruFields, ok := translations["ru"]; ok {
		if ruFields["location"] != "Улица Пушкина, дом 10" {
			t.Errorf("Expected Russian location unchanged, got %s", ruFields["location"])
		}
		if ruFields["city"] != "Белград" {
			t.Errorf("Expected Russian city unchanged, got %s", ruFields["city"])
		}
		if ruFields["country"] != "Сербия" {
			t.Errorf("Expected Russian country unchanged, got %s", ruFields["country"])
		}
	} else {
		t.Error("Russian translations not found")
	}
}

// TestCreateListingWithAddressTranslation тестирует создание объявления с переводом адресов
func TestCreateListingWithAddressTranslation(t *testing.T) {
	// Этот тест требует полной настройки всех зависимостей
	// и доступа к базе данных, поэтому пропускаем его в unit тестах
	t.Skip("Integration test requires database connection")
}

// TestUpdateListingWithAddressTranslation тестирует обновление объявления с переводом адресов
func TestUpdateListingWithAddressTranslation(t *testing.T) {
	// Этот тест требует полной настройки всех зависимостей
	// и доступа к базе данных, поэтому пропускаем его в unit тестах
	t.Skip("Integration test requires database connection")
}
