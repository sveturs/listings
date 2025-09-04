package config

import (
	"os"
	"strconv"
)

// FeatureFlags управляет флагами для постепенного перехода на новую систему
type FeatureFlags struct {
	// Флаги для унифицированной системы атрибутов
	UseUnifiedAttributes      bool `yaml:"use_unified_attributes"`      // Использовать новую систему атрибутов
	UnifiedAttributesFallback bool `yaml:"unified_attributes_fallback"` // Откат на старую систему при ошибках
	DualWriteAttributes       bool `yaml:"dual_write_attributes"`       // Дублировать запись в обе системы
	UnifiedAttributesPercent  int  `yaml:"unified_attributes_percent"`  // Процент трафика на новую систему (0-100)

	// Дополнительные флаги для контроля
	LogAttributeSystemCalls  bool `yaml:"log_attribute_system_calls"`  // Логировать все вызовы
	AttributeCacheEnabled    bool `yaml:"attribute_cache_enabled"`     // Включить кеширование атрибутов
	AttributeCacheTTLMinutes int  `yaml:"attribute_cache_ttl_minutes"` // TTL кеша в минутах
}

// LoadFeatureFlags загружает флаги из переменных окружения
func LoadFeatureFlags() *FeatureFlags {
	flags := &FeatureFlags{
		// Значения по умолчанию (безопасные для продакшена)
		UseUnifiedAttributes:      false,
		UnifiedAttributesFallback: true,
		DualWriteAttributes:       true,
		UnifiedAttributesPercent:  0,
		LogAttributeSystemCalls:   false,
		AttributeCacheEnabled:     true,
		AttributeCacheTTLMinutes:  30,
	}

	// Загрузка из переменных окружения
	if val := os.Getenv("USE_UNIFIED_ATTRIBUTES"); val != "" {
		flags.UseUnifiedAttributes = val == "true"
	}

	if val := os.Getenv("UNIFIED_ATTRIBUTES_FALLBACK"); val != "" {
		flags.UnifiedAttributesFallback = val == "true"
	}

	if val := os.Getenv("DUAL_WRITE_ATTRIBUTES"); val != "" {
		flags.DualWriteAttributes = val == "true"
	}

	if val := os.Getenv("UNIFIED_ATTRIBUTES_PERCENT"); val != "" {
		if percent, err := strconv.Atoi(val); err == nil {
			if percent >= 0 && percent <= 100 {
				flags.UnifiedAttributesPercent = percent
			}
		}
	}

	if val := os.Getenv("LOG_ATTRIBUTE_SYSTEM_CALLS"); val != "" {
		flags.LogAttributeSystemCalls = val == "true"
	}

	if val := os.Getenv("ATTRIBUTE_CACHE_ENABLED"); val != "" {
		flags.AttributeCacheEnabled = val == "true"
	}

	if val := os.Getenv("ATTRIBUTE_CACHE_TTL_MINUTES"); val != "" {
		if ttl, err := strconv.Atoi(val); err == nil && ttl > 0 {
			flags.AttributeCacheTTLMinutes = ttl
		}
	}

	return flags
}

// ShouldUseUnifiedAttributes определяет, должен ли конкретный запрос использовать новую систему
func (ff *FeatureFlags) ShouldUseUnifiedAttributes(userID int) bool {
	if !ff.UseUnifiedAttributes {
		return false
	}

	// A/B тестирование по проценту пользователей
	if ff.UnifiedAttributesPercent < 100 {
		// Используем простой hash от userID для определения группы
		userGroup := userID % 100
		return userGroup < ff.UnifiedAttributesPercent
	}

	return true
}

// IsFeatureEnabled проверяет, включена ли функция
func (ff *FeatureFlags) IsFeatureEnabled(feature string) bool {
	switch feature {
	case "unified_attributes":
		return ff.UseUnifiedAttributes
	case "unified_attributes_fallback":
		return ff.UnifiedAttributesFallback
	case "dual_write_attributes":
		return ff.DualWriteAttributes
	case "attribute_cache":
		return ff.AttributeCacheEnabled
	case "attribute_logging":
		return ff.LogAttributeSystemCalls
	default:
		return false
	}
}

// GetFeaturePercentage возвращает процент включения функции
func (ff *FeatureFlags) GetFeaturePercentage(feature string) int {
	switch feature {
	case "unified_attributes":
		return ff.UnifiedAttributesPercent
	default:
		return 0
	}
}

// UpdateFeatureFlag обновляет значение флага (для динамического управления)
func (ff *FeatureFlags) UpdateFeatureFlag(feature string, value interface{}) error {
	switch feature {
	case "use_unified_attributes":
		if boolVal, ok := value.(bool); ok {
			ff.UseUnifiedAttributes = boolVal
			return nil
		}
	case "unified_attributes_fallback":
		if boolVal, ok := value.(bool); ok {
			ff.UnifiedAttributesFallback = boolVal
			return nil
		}
	case "dual_write_attributes":
		if boolVal, ok := value.(bool); ok {
			ff.DualWriteAttributes = boolVal
			return nil
		}
	case "unified_attributes_percent":
		if intVal, ok := value.(int); ok && intVal >= 0 && intVal <= 100 {
			ff.UnifiedAttributesPercent = intVal
			return nil
		}
	case "log_attribute_system_calls":
		if boolVal, ok := value.(bool); ok {
			ff.LogAttributeSystemCalls = boolVal
			return nil
		}
	case "attribute_cache_enabled":
		if boolVal, ok := value.(bool); ok {
			ff.AttributeCacheEnabled = boolVal
			return nil
		}
	case "attribute_cache_ttl_minutes":
		if intVal, ok := value.(int); ok && intVal > 0 {
			ff.AttributeCacheTTLMinutes = intVal
			return nil
		}
	}
	return ErrInvalidFeatureFlag
}

// GetCurrentConfiguration возвращает текущую конфигурацию флагов
func (ff *FeatureFlags) GetCurrentConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"use_unified_attributes":      ff.UseUnifiedAttributes,
		"unified_attributes_fallback": ff.UnifiedAttributesFallback,
		"dual_write_attributes":       ff.DualWriteAttributes,
		"unified_attributes_percent":  ff.UnifiedAttributesPercent,
		"log_attribute_system_calls":  ff.LogAttributeSystemCalls,
		"attribute_cache_enabled":     ff.AttributeCacheEnabled,
		"attribute_cache_ttl_minutes": ff.AttributeCacheTTLMinutes,
	}
}

// План постепенного включения
// Неделя 1: Тестовое окружение
// - USE_UNIFIED_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_FALLBACK=true
// - DUAL_WRITE_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_PERCENT=100

// Неделя 2: 10% продакшн трафика
// - USE_UNIFIED_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_FALLBACK=true
// - DUAL_WRITE_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_PERCENT=10

// Неделя 3: 50% продакшн трафика
// - USE_UNIFIED_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_FALLBACK=true
// - DUAL_WRITE_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_PERCENT=50

// Неделя 4: 100% продакшн трафика
// - USE_UNIFIED_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_FALLBACK=true
// - DUAL_WRITE_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_PERCENT=100

// После стабилизации: Отключение старой системы
// - USE_UNIFIED_ATTRIBUTES=true
// - UNIFIED_ATTRIBUTES_FALLBACK=false
// - DUAL_WRITE_ATTRIBUTES=false
// - UNIFIED_ATTRIBUTES_PERCENT=100
