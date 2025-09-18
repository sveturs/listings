package config

import "os"

// ViberConfig содержит настройки Viber Bot
type ViberConfig struct {
	// Direct Viber API (если используется напрямую)
	AuthToken   string
	BotName     string
	BotAvatar   string
	WebhookURL  string
	APIEndpoint string
	FrontendURL string // URL фронтенда для генерации ссылок

	// Infobip API (рекомендуемый способ)
	UseInfobip      bool
	InfobipAPIKey   string
	InfobipBaseURL  string
	InfobipSenderID string // ID отправителя в Infobip (Viber Service ID)
}

// LoadViberConfig загружает конфигурацию из переменных окружения
func LoadViberConfig() *ViberConfig {
	// Проверяем наличие Infobip конфигурации
	useInfobip := getEnv("INFOBIP_API_KEY", "") != ""

	return &ViberConfig{
		// Direct Viber API
		AuthToken:   getEnv("VIBER_AUTH_TOKEN", ""),
		BotName:     getEnv("VIBER_BOT_NAME", "SveTu Marketplace"),
		BotAvatar:   getEnv("VIBER_BOT_AVATAR", "https://svetu.rs/logo-720.png"),
		WebhookURL:  getEnv("VIBER_WEBHOOK_URL", "https://svetu.rs/api/viber/webhook"),
		APIEndpoint: getEnv("VIBER_API_ENDPOINT", "https://chatapi.viber.com/pa"),
		FrontendURL: getEnv("VIBER_PUBLIC_URL", getEnv("FRONTEND_URL", "https://svetu.rs")),

		// Infobip API
		UseInfobip:      useInfobip,
		InfobipAPIKey:   getEnv("INFOBIP_API_KEY", ""), // Требуется из переменной окружения
		InfobipBaseURL:  getEnv("INFOBIP_BASE_URL", ""),
		InfobipSenderID: getEnv("INFOBIP_SENDER_ID", ""), // Требуется из переменной окружения
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
