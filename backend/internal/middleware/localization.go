package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// LocaleContextKey is the key for storing locale in context
type contextKey string

const LocaleContextKey contextKey = "locale"

// SupportedLocales defines the list of supported locales
var SupportedLocales = []string{"sr", "ru", "en"}

// DefaultLocale is the fallback locale
const DefaultLocale = "sr"

// LocalizationConfig holds configuration for localization middleware
type LocalizationConfig struct {
	// SupportedLocales is the list of supported locale codes
	SupportedLocales []string
	// DefaultLocale is the fallback locale when none is detected
	DefaultLocale string
	// CookieName is the name of the cookie to store locale preference
	CookieName string
	// HeaderName is the name of the header to read locale from
	HeaderName string
}

// DefaultLocalizationConfig returns default configuration
func DefaultLocalizationConfig() LocalizationConfig {
	return LocalizationConfig{
		SupportedLocales: SupportedLocales,
		DefaultLocale:    DefaultLocale,
		CookieName:       "locale-preference",
		HeaderName:       "Accept-Language",
	}
}

// Localization creates a new localization middleware
func Localization(config ...LocalizationConfig) fiber.Handler {
	cfg := DefaultLocalizationConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		locale := detectLocale(c, cfg)
		
		// Store locale in context
		c.Locals(string(LocaleContextKey), locale)
		
		// Add locale to response headers for debugging
		c.Set("X-Detected-Locale", locale)
		
		log.Debug().
			Str("locale", locale).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Msg("Detected locale for request")

		return c.Next()
	}
}

// detectLocale detects the best locale based on various sources
func detectLocale(c *fiber.Ctx, cfg LocalizationConfig) string {
	// 1. Check URL parameter (?lang=sr)
	if urlLang := c.Query("lang"); urlLang != "" {
		if isValidLocale(urlLang, cfg.SupportedLocales) {
			log.Debug().Str("source", "url_param").Str("locale", urlLang).Msg("Locale detected")
			return urlLang
		}
	}

	// 2. Check cookie
	if cookieLang := c.Cookies(cfg.CookieName); cookieLang != "" {
		if isValidLocale(cookieLang, cfg.SupportedLocales) {
			log.Debug().Str("source", "cookie").Str("locale", cookieLang).Msg("Locale detected")
			return cookieLang
		}
	}

	// 3. Check Accept-Language header
	if headerLang := parseAcceptLanguage(c.Get(cfg.HeaderName), cfg.SupportedLocales); headerLang != "" {
		log.Debug().Str("source", "accept_language").Str("locale", headerLang).Msg("Locale detected")
		return headerLang
	}

	// 4. Check custom X-Locale header (for API clients)
	if customLang := c.Get("X-Locale"); customLang != "" {
		if isValidLocale(customLang, cfg.SupportedLocales) {
			log.Debug().Str("source", "x_locale_header").Str("locale", customLang).Msg("Locale detected")
			return customLang
		}
	}

	// 5. Fallback to default
	log.Debug().Str("source", "default").Str("locale", cfg.DefaultLocale).Msg("Using default locale")
	return cfg.DefaultLocale
}

// parseAcceptLanguage parses Accept-Language header and returns best match
func parseAcceptLanguage(acceptLang string, supportedLocales []string) string {
	if acceptLang == "" {
		return ""
	}

	// Parse Accept-Language header: "en-US,en;q=0.9,ru;q=0.8,sr;q=0.7"
	languages := strings.Split(acceptLang, ",")
	
	type langWithQuality struct {
		lang    string
		quality float32
	}
	
	var parsed []langWithQuality
	
	for _, lang := range languages {
		lang = strings.TrimSpace(lang)
		parts := strings.Split(lang, ";")
		
		langCode := strings.TrimSpace(parts[0])
		quality := float32(1.0) // Default quality
		
		// Parse quality value if present
		if len(parts) > 1 {
			for _, part := range parts[1:] {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "q=") {
					if _, err := fmt.Sscanf(part, "q=%f", &quality); err != nil {
						quality = 1.0
					}
					break
				}
			}
		}
		
		// Extract main language code (en-US -> en)
		if idx := strings.Index(langCode, "-"); idx > 0 {
			langCode = langCode[:idx]
		}
		
		parsed = append(parsed, langWithQuality{lang: langCode, quality: quality})
	}
	
	// Sort by quality (highest first)
	for i := 0; i < len(parsed)-1; i++ {
		for j := i + 1; j < len(parsed); j++ {
			if parsed[j].quality > parsed[i].quality {
				parsed[i], parsed[j] = parsed[j], parsed[i]
			}
		}
	}
	
	// Find first supported language
	for _, item := range parsed {
		if isValidLocale(item.lang, supportedLocales) {
			return item.lang
		}
	}
	
	return ""
}

// isValidLocale checks if locale is supported
func isValidLocale(locale string, supportedLocales []string) bool {
	locale = strings.ToLower(locale)
	for _, supported := range supportedLocales {
		if strings.ToLower(supported) == locale {
			return true
		}
	}
	return false
}

// GetLocaleFromContext extracts locale from fiber context
func GetLocaleFromContext(c *fiber.Ctx) string {
	if locale, ok := c.Locals(string(LocaleContextKey)).(string); ok {
		return locale
	}
	return DefaultLocale
}

// GetLocaleFromGoContext extracts locale from Go context (for use in services)
func GetLocaleFromGoContext(ctx context.Context) string {
	if locale, ok := ctx.Value(LocaleContextKey).(string); ok {
		return locale
	}
	return DefaultLocale
}

// WithLocale adds locale to Go context
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, LocaleContextKey, locale)
}