// backend/internal/proj/c2c/storage/postgres/storage.go
package postgres

import (
	"sync"
	"time"

	"backend/internal/domain/models"
	"backend/internal/proj/c2c/service"

	"github.com/jackc/pgx/v5/pgxpool"
	authservice "github.com/sveturs/auth/pkg/http/service"
)

const (
	// Attribute names
	attrNameArea           = "area"
	attrNameLandArea       = "land_area"
	attrNameMileage        = "mileage"
	attrNameEngineCapacity = "engine_capacity"
	attrNamePower          = "power"
	attrNameYear           = "year"
)

type Storage struct {
	pool               *pgxpool.Pool
	translationService service.TranslationServiceInterface
	userService        *authservice.UserService

	// Кэш для атрибутов категорий
	attributeCacheMutex sync.RWMutex
	attributeCache      map[int][]models.CategoryAttribute
	attributeCacheTime  map[int]time.Time

	// Кэш для ranges
	rangesCacheMutex sync.RWMutex
	rangesCache      map[int]map[string]map[string]interface{}
	rangesCacheTime  map[int]time.Time

	cacheTTL time.Duration
}

func NewStorage(pool *pgxpool.Pool, translationService service.TranslationServiceInterface, userService *authservice.UserService) *Storage {
	return &Storage{
		pool:               pool,
		translationService: translationService,
		userService:        userService,
		attributeCache:     make(map[int][]models.CategoryAttribute),
		attributeCacheTime: make(map[int]time.Time),
		rangesCache:        make(map[int]map[string]map[string]interface{}),
		rangesCacheTime:    make(map[int]time.Time),
		cacheTTL:           30 * time.Minute,
	}
}

// SetUserService устанавливает UserService для Storage
func (s *Storage) SetUserService(userService *authservice.UserService) {
	s.userService = userService
}
