package main

import (
	"context"
	"database/sql"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
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

// Вспомогательная функция для получения значения переменной окружения или дефолтного значения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Вспомогательная функция для получения целочисленного значения из переменной окружения
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value := 0; valueStr != "" {
		fmt.Sscanf(valueStr, "%d", &value)
		return value
	}
	return defaultValue
}

// Подключение к базе данных
func connectDB(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с БД: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с БД: %w", err)
	}

	return db, nil
}

// Построение иерархической структуры категорий
func buildCategoryTree(categories []models.MarketplaceCategory) []models.MarketplaceCategory {
	// Создаем карту для быстрого доступа к категориям по ID
	categoryMap := make(map[int]*models.MarketplaceCategory)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	// Строим дерево
	var rootCategories []models.MarketplaceCategory
	for _, cat := range categories {
		// Если нет родителя, это корневая категория
		if cat.ParentID == nil {
			rootCategories = append(rootCategories, cat)
		} else {
			// Если есть родитель, добавляем текущую категорию в дочерние родителя
			if parent, exists := categoryMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, cat)
			} else {
				// Если родитель не найден, считаем категорию корневой
				rootCategories = append(rootCategories, cat)
			}
		}
	}

	return rootCategories
}

// Запись категорий в YAML файл
func writeCategoriesFile(filePath string, categories []models.MarketplaceCategory) error {
	// Создаем структуру для файла
	categoriesFile := CategoriesFile{
		Categories: categories,
	}

	// Преобразуем в YAML с красивым форматированием
	yamlData, err := yaml.Marshal(categoriesFile)
	if err != nil {
		return fmt.Errorf("ошибка преобразования в YAML: %w", err)
	}

	// Записываем в файл
	if err := ioutil.WriteFile(filePath, yamlData, 0644); err != nil {
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
	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage)
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

	// Строим иерархическую структуру
	rootCategories := buildCategoryTree(flatCategories)

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
