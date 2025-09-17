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

		// Infobip API
		UseInfobip:      useInfobip,
		InfobipAPIKey:   getEnv("INFOBIP_API_KEY", "d225d6436a020569e4dee8468919de13-5d9211b6-0f93-4829-bd40-59855773e867"),
		InfobipBaseURL:  getEnv("INFOBIP_BASE_URL", "d9vgp1.api.infobip.com"),
		InfobipSenderID: getEnv("INFOBIP_SENDER_ID", "SveTuBot"), // Нужно будет зарегистрировать в Infobip
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}