package factory

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/delivery/interfaces"
	"backend/internal/proj/postexpress"
)

// ProviderFactory - фабрика для создания провайдеров доставки
type ProviderFactory struct {
	db                 *sqlx.DB
	postExpressService *postexpress.Service
	// Добавим другие сервисы по мере их реализации
}

// NewProviderFactory создает новую фабрику провайдеров
func NewProviderFactory(db *sqlx.DB, postExpressService *postexpress.Service) *ProviderFactory {
	return &ProviderFactory{
		db:                 db,
		postExpressService: postExpressService,
	}
}

// NewProviderFactoryWithDefaults создает фабрику с автоинициализацией сервисов
func NewProviderFactoryWithDefaults(db *sqlx.DB) (*ProviderFactory, error) {
	// Инициализируем Post Express сервис
	postExpressSvc, err := postexpress.NewService(nil) // nil = загрузит config из ENV
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to initialize Post Express service, using mock provider")
		postExpressSvc = nil // Fallback to mock
	}

	return &ProviderFactory{
		db:                 db,
		postExpressService: postExpressSvc,
	}, nil
}

// CreateProvider создает провайдера по коду
func (f *ProviderFactory) CreateProvider(code string) (interfaces.DeliveryProvider, error) {
	switch code {
	case "post_express":
		if f.postExpressService != nil {
			log.Debug().Msg("Creating Post Express adapter with real service")
			return NewPostExpressAdapter(f.postExpressService), nil
		}
		// Fallback на mock если сервис не инициализирован
		log.Warn().Msg("Post Express service not available, using mock provider")
		return NewMockProvider("post_express", "Post Express"), nil
	case "bex_express":
		// TODO: реализовать адаптер для BEX Express
		return NewMockProvider("bex_express", "BEX Express"), nil
	case "aks_express":
		return NewMockProvider("aks_express", "AKS Express"), nil
	case "d_express":
		return NewMockProvider("d_express", "D Express"), nil
	case "city_express":
		return NewMockProvider("city_express", "City Express"), nil
	case "dhl_express":
		return NewMockProvider("dhl_express", "DHL Express"), nil
	default:
		return nil, fmt.Errorf("unknown provider code: %s", code)
	}
}

// GetAvailableProviders возвращает список доступных провайдеров
func (f *ProviderFactory) GetAvailableProviders() ([]string, error) {
	return []string{
		"post_express",
		"bex_express",
		"aks_express",
		"d_express",
		"city_express",
		"dhl_express",
	}, nil
}
