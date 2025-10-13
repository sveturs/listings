// backend/internal/config/currencies.go
package config

// CurrencyConfig содержит конфигурацию валют для платформы
type CurrencyConfig struct {
	// DefaultCurrency - валюта по умолчанию для всех операций
	DefaultCurrency string `mapstructure:"default_currency"`

	// SupportedCurrencies - список поддерживаемых валют
	SupportedCurrencies []string `mapstructure:"supported_currencies"`

	// ExchangeRatesEnabled - включены ли автоматические курсы обмена
	ExchangeRatesEnabled bool `mapstructure:"exchange_rates_enabled"`

	// ExchangeRatesProvider - провайдер курсов валют (например, "ecb", "cbr", "manual")
	ExchangeRatesProvider string `mapstructure:"exchange_rates_provider"`
}

// GetDefaultCurrencyConfig возвращает конфигурацию валют по умолчанию
func GetDefaultCurrencyConfig() CurrencyConfig {
	return CurrencyConfig{
		DefaultCurrency:       "RSD", // Сербский динар
		SupportedCurrencies:   []string{"RSD", "EUR", "USD"},
		ExchangeRatesEnabled:  false, // По умолчанию отключено
		ExchangeRatesProvider: "manual",
	}
}

// IsCurrencySupported проверяет, поддерживается ли валюта
func (cc *CurrencyConfig) IsCurrencySupported(currency string) bool {
	for _, c := range cc.SupportedCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}

// GetCurrencySymbol возвращает символ валюты
func (cc *CurrencyConfig) GetCurrencySymbol(currency string) string {
	symbols := map[string]string{
		"RSD": "РСД",
		"EUR": "€",
		"USD": "$",
		"GBP": "£",
		"RUB": "₽",
	}

	if symbol, ok := symbols[currency]; ok {
		return symbol
	}
	return currency
}

// globalDefaultCurrency хранит дефолтную валюту (для быстрого доступа без config)
var globalDefaultCurrency = "RSD"

// SetGlobalDefaultCurrency устанавливает глобальную валюту по умолчанию
// Должна вызываться при инициализации приложения
func SetGlobalDefaultCurrency(currency string) {
	if currency != "" {
		globalDefaultCurrency = currency
	}
}

// GetGlobalDefaultCurrency возвращает глобальную валюту по умолчанию
func GetGlobalDefaultCurrency() string {
	return globalDefaultCurrency
}
