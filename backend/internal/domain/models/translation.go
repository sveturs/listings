// backend/internal/domain/models/translation.go
package models

import "time"

type Translation struct {
    ID                int       `json:"id"`
    EntityType        string    `json:"entity_type"`
    EntityID          int       `json:"entity_id"`
    FieldName         string    `json:"field_name"`
    Language          string    `json:"language"`
    TranslatedText    string    `json:"translated_text"`
    IsMachineTranslated bool    `json:"is_machine_translated"`
    IsVerified        bool      `json:"is_verified"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

type TranslatedFields struct {
    Translations map[string]map[string]string `json:"translations"` // language -> field -> text
    OriginalLanguage string                   `json:"original_language"`
}

// Интерфейс для сущностей, поддерживающих перевод
type Translatable interface {
    GetTranslationFields() []string
    GetOriginalLanguage() string
    SetTranslations(translations map[string]map[string]string)
}

// Добавляем поддержку переводов в существующие модели
func (l *MarketplaceListing) GetTranslationFields() []string {
    return []string{"title", "description"}
}

func (l *MarketplaceListing) GetOriginalLanguage() string {
    return l.OriginalLanguage
}

func (l *MarketplaceListing) SetTranslations(translations map[string]map[string]string) {
    if l.Translations == nil {
        l.Translations = make(map[string]map[string]string)
    }
    for lang, fields := range translations {
        if l.Translations[lang] == nil {
            l.Translations[lang] = make(map[string]string)
        }
        for field, text := range fields {
            l.Translations[lang][field] = text
        }
    }
}

// Аналогично для Review и MarketplaceMessage
func (r *Review) GetTranslationFields() []string {
    return []string{"comment", "pros", "cons"}
}

func (r *Review) GetOriginalLanguage() string {
    return r.OriginalLanguage
}

func (r *Review) SetTranslations(translations map[string]map[string]string) {
    if r.Translations == nil {
        r.Translations = make(map[string]map[string]string)
    }
    for lang, fields := range translations {
        if r.Translations[lang] == nil {
            r.Translations[lang] = make(map[string]string)
        }
        for field, text := range fields {
            r.Translations[lang][field] = text
        }
    }
}

func (m *MarketplaceMessage) GetTranslationFields() []string {
    return []string{"content"}
}

func (m *MarketplaceMessage) GetOriginalLanguage() string {
    return m.OriginalLanguage
}

func (m *MarketplaceMessage) SetTranslations(translations map[string]map[string]string) {
    if m.Translations == nil {
        m.Translations = make(map[string]map[string]string)
    }
    for lang, fields := range translations {
        if m.Translations[lang] == nil {
            m.Translations[lang] = make(map[string]string)
        }
        for field, text := range fields {
            m.Translations[lang][field] = text
        }
    }
}