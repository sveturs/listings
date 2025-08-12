package translation_admin

import (
	"context"
	"fmt"
	"os"
)

// AIProvider интерфейс для AI провайдеров переводов
type AIProvider interface {
	Translate(ctx context.Context, text, sourceLang, targetLang string) (string, float64, error)
	IsConfigured() bool
	GetName() string
}

// MockProvider мок-провайдер для тестирования
type MockProvider struct {
	name string
}

func NewMockProvider(name string) *MockProvider {
	if name == "" {
		name = "mock"
	}
	return &MockProvider{name: name}
}

func (p *MockProvider) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, float64, error) {
	// Простой мок - возвращаем текст с префиксом
	return fmt.Sprintf("[%s] %s", targetLang, text), 0.95, nil
}

func (p *MockProvider) IsConfigured() bool {
	return true
}

func (p *MockProvider) GetName() string {
	return p.name
}

// GetAIProvider возвращает провайдера по имени
func GetAIProvider(name string) AIProvider {
	// TODO: Реализовать реальных провайдеров
	// Пока возвращаем мок
	if name == "" {
		name = "openai"
	}
	
	// Проверяем наличие API ключей в env
	switch name {
	case "openai":
		if os.Getenv("OPENAI_API_KEY") != "" {
			// TODO: вернуть настоящий OpenAI провайдер
			return NewMockProvider("openai")
		}
	case "google":
		if os.Getenv("GOOGLE_API_KEY") != "" {
			// TODO: вернуть настоящий Google провайдер  
			return NewMockProvider("google")
		}
	case "deepl":
		if os.Getenv("DEEPL_API_KEY") != "" {
			// TODO: вернуть настоящий DeepL провайдер
			return NewMockProvider("deepl")
		}
	case "claude":
		if os.Getenv("CLAUDE_API_KEY") != "" {
			// TODO: вернуть настоящий Claude провайдер
			return NewMockProvider("claude")
		}
	}
	
	// По умолчанию возвращаем мок
	return NewMockProvider(name)
}