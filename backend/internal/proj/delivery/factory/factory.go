package factory

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/delivery/interfaces"
	// "backend/internal/proj/postexpress" // TODO: исправить когда postexpress будет доступен
)

// ProviderFactory - фабрика для создания провайдеров доставки
type ProviderFactory struct {
	db                 *sqlx.DB
	postExpressService interface{} // *postexpress.Service // TODO: исправить тип когда postexpress будет доступен
	// Добавим другие сервисы по мере их реализации
}

// NewProviderFactory создает новую фабрику провайдеров
func NewProviderFactory(db *sqlx.DB, postExpressService interface{}) *ProviderFactory {
	return &ProviderFactory{
		db:                 db,
		postExpressService: postExpressService,
	}
}

// CreateProvider создает провайдера по коду
func (f *ProviderFactory) CreateProvider(code string) (interfaces.DeliveryProvider, error) {
	switch code {
	case "post_express":
		if f.postExpressService != nil {
			return NewPostExpressAdapter(f.postExpressService), nil
		}
		// Fallback на mock если сервис не инициализирован
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
