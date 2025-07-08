package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"

	"gopkg.in/yaml.v3"
)

// CategoriesFile представляет структуру YAML файла с категориями
type CategoriesFile struct {
	Categories []models.MarketplaceCategory `yaml:"categories"`
}

// Config структура для хранения конфигурации подключения к базе данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// SimplifiedCategory представляет упрощенную структуру категории для YAML
type SimplifiedCategory struct {
	ID           int                   `yaml:"id"`
	Name         string                `yaml:"name"`
	Slug         string                `yaml:"slug"`
	ParentID     *int                  `yaml:"parentid"`
	Icon         string                `yaml:"icon"`
	Translations map[string]string     `yaml:"translations"`
	Children     []*SimplifiedCategory `yaml:"children"`
}

// SimplifiedCategoriesFile представляет структуру YAML файла с упрощенными категориями
type SimplifiedCategoriesFile struct {
	Categories []*SimplifiedCategory `yaml:"categories"`
}

// Построение иерархической структуры категорий
func buildCategoryTree(categories []*SimplifiedCategory) []*SimplifiedCategory {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})

	// Создаем карту для быстрого доступа к категориям по ID
	categoryMap := make(map[int]*SimplifiedCategory)
	for i := range categories {
		categoryMap[categories[i].ID] = categories[i]
	}

	// Строим дерево
	var rootCategories []*SimplifiedCategory
	for _, cat := range categories {
		c := cat
		// Если нет родителя, это корневая категория
		if cat.ParentID == nil {
			rootCategories = append(rootCategories, c)
		} else {
			// Если есть родитель, добавляем текущую категорию в дочерние родителя
			if parent, exists := categoryMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, c)
			} else {
				// Если родитель не найден, считаем категорию корневой
				rootCategories = append(rootCategories, c)
			}
		}
	}

	return rootCategories
}

// convertToSimplified преобразует MarketplaceCategory в SimplifiedCategory
func convertToSimplified(cat models.MarketplaceCategory) *SimplifiedCategory {
	simplified := &SimplifiedCategory{
		// ID:           cat.ID * 10, // TODO: когда полностью перейдем на файл
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		ParentID:     cat.ParentID,
		Icon:         cat.Icon,
		Translations: cat.Translations,
	}

	return simplified
}

// Запись категорий в YAML файл
func writeCategoriesFile(filePath string, categories []*SimplifiedCategory) error {
	// Создаем структуру для файла
	categoriesFile := SimplifiedCategoriesFile{
		Categories: categories,
	}

	// Преобразуем в YAML с красивым форматированием
	yamlData, err := yaml.Marshal(categoriesFile)
	if err != nil {
		return fmt.Errorf("ошибка преобразования в YAML: %w", err)
	}

	// Записываем в файл
	if err := os.WriteFile(filePath, yamlData, 0o644); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

func run(cfg *config.Config) error {
	ctx := context.Background()

	// Инициализируем файловое хранилище
	fileStorage, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		log.Printf("Ошибка инициализации файлового хранилища: %v. Функции загрузки файлов могут быть недоступны.", err)

		// Не прерываем выполнение программы, так как сервер может работать и без файлового хранилища
	}

	// Инициализируем клиент OpenSearch
	var osClient *opensearch.OpenSearchClient
	if cfg.OpenSearch.URL != "" {
		var err error
		osClient, err = opensearch.NewOpenSearchClient(opensearch.Config{
			URL:      cfg.OpenSearch.URL,
			Username: cfg.OpenSearch.Username,
			Password: cfg.OpenSearch.Password,
		})
		if err != nil {
			log.Printf("Ошибка подключения к OpenSearch: %v", err)
		} else {
			log.Println("Успешное подключение к OpenSearch")
		}
	} else {
		log.Println("OpenSearch URL не указан, поиск будет отключен")
	}

	// Инициализируем базу данных с OpenSearch и файловым хранилищем
	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage, cfg.SearchWeights)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Получаем путь к YAML файлу из аргументов командной строки
	if len(os.Args) < 2 {
		log.Fatal("Укажите путь для сохранения YAML файла с категориями")
	}
	yamlFilePath := os.Args[1]

	// Получаем все категории из базы данных
	flatCategories, err := db.GetCategories(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Конвертируем категории в упрощенный формат
	simplifiedFlatCategories := make([]*SimplifiedCategory, 0, len(flatCategories))
	for _, cat := range flatCategories {
		simplifiedFlatCategories = append(simplifiedFlatCategories, convertToSimplified(cat))
	}

	// Строим иерархическую структуру
	rootCategories := buildCategoryTree(simplifiedFlatCategories)

	// Записываем категории в YAML файл
	if err := writeCategoriesFile(yamlFilePath, rootCategories); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Успешно экспортировано %d категорий в файл %s\n",
		len(flatCategories), yamlFilePath)

	return nil
}

func main() { // Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}
	err = run(cfg)
	if err != nil {
		log.Printf("Error running application: %v", err)
		os.Exit(1)
	}

	log.Println("Application started successfully")
}
