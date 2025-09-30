// backend/internal/storage/opensearch/client.go
package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// OpenSearchClient предоставляет методы для работы с OpenSearch
type OpenSearchClient struct {
	client *opensearch.Client
}

// Config содержит настройки подключения к OpenSearch
type Config struct {
	URL      string
	Username string
	Password string
}

// NewOpenSearchClient создает новый клиент OpenSearch
func NewOpenSearchClient(config Config) (*OpenSearchClient, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{config.URL},
		Username:  config.Username,
		Password:  config.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать клиент OpenSearch: %w", err)
	}

	return &OpenSearchClient{client: client}, nil
}

// CreateIndex создает индекс с указанной схемой
func (c *OpenSearchClient) CreateIndex(ctx context.Context, indexName string, mapping string) error {
	exists, err := c.IndexExists(ctx, indexName)
	if err != nil {
		return err
	}

	if exists {
		return nil // Индекс уже существует
	}

	createIndex := opensearchapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := createIndex.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("ошибка создания индекса: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		return fmt.Errorf("ошибка создания индекса: %s", res.String())
	}

	return nil
}

// IndexExists проверяет существование индекса
func (c *OpenSearchClient) IndexExists(ctx context.Context, indexName string) (bool, error) {
	exists := opensearchapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := exists.Do(ctx, c.client)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	return res.StatusCode == 200, nil
}

// IndexDocument индексирует документ
func (c *OpenSearchClient) IndexDocument(ctx context.Context, indexName, id string, document interface{}) error {
	docBytes, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("ошибка сериализации документа: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       strings.NewReader(string(docBytes)),
		Refresh:    "true", // Добавляем refresh для немедленной видимости изменений
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("ошибка индексации документа: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		return fmt.Errorf("ошибка индексации документа: %s", res.String())
	}

	return nil
}

// BulkIndex индексирует несколько документов за один запрос
func (c *OpenSearchClient) BulkIndex(ctx context.Context, indexName string, documents []map[string]interface{}) error {
	if len(documents) == 0 {
		return nil
	}

	fmt.Printf("BulkIndex: Starting to index %d documents to index %s\n", len(documents), indexName)

	var bulkBody strings.Builder

	for _, doc := range documents {
		// Получаем ID документа
		docID, ok := doc["id"].(string)
		if !ok {
			// Преобразуем число в строку, если ID - число
			if numID, ok := doc["id"].(int); ok {
				docID = fmt.Sprintf("%d", numID)
			} else {
				return fmt.Errorf("ID документа не найден или имеет неверный тип")
			}
		}

		// Создаем метаданные действия
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    docID,
			},
		}

		actionJSON, err := json.Marshal(action)
		if err != nil {
			return fmt.Errorf("ошибка сериализации action: %w", err)
		}

		// Добавляем метаданные действия
		bulkBody.WriteString(string(actionJSON))
		bulkBody.WriteString("\n")

		// Сохраняем документ как есть, включая ID
		docCopy := doc

		docJSON, err := json.Marshal(docCopy)
		if err != nil {
			return fmt.Errorf("ошибка сериализации документа: %w", err)
		}

		// Добавляем документ
		bulkBody.WriteString(string(docJSON))
		bulkBody.WriteString("\n")
	}

	req := opensearchapi.BulkRequest{
		Body: strings.NewReader(bulkBody.String()),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("ошибка bulk индексации: %w", err)
	}

	// Всегда закрываем body в конце
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	// Читаем ответ для проверки ошибок
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа bulk: %w", err)
	}

	// Проверяем HTTP статус
	if res.IsError() {
		return fmt.Errorf("ошибка bulk индексации HTTP %d: %s, body: %s", res.StatusCode, res.String(), string(body))
	}

	// Парсим ответ для проверки индивидуальных ошибок
	var bulkResponse map[string]interface{}
	if err := json.Unmarshal(body, &bulkResponse); err != nil {
		return fmt.Errorf("ошибка парсинга ответа bulk: %w, body: %s", err, string(body))
	}

	// Проверяем на наличие ошибок в ответе
	if hasErrors, ok := bulkResponse["errors"].(bool); ok && hasErrors {
		// Логируем детали ошибок
		if items, ok := bulkResponse["items"].([]interface{}); ok {
			for i, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					for action, details := range itemMap {
						if detailsMap, ok := details.(map[string]interface{}); ok {
							if errorInfo, hasError := detailsMap["error"]; hasError {
								fmt.Printf("Bulk error for doc %d (action: %s): %v\n", i, action, errorInfo)
							}
						}
					}
				}
			}
		}
		return fmt.Errorf("bulk индексация содержит ошибки, см. логи выше")
	}

	fmt.Printf("BulkIndex: Successfully completed bulk operation (took: %v, errors: %v)\n", bulkResponse["took"], bulkResponse["errors"])

	return nil
}

// UpdateDocument частично обновляет документ в индексе
func (c *OpenSearchClient) UpdateDocument(ctx context.Context, indexName, id string, doc map[string]interface{}) error {
	updateData := map[string]interface{}{
		"doc": doc,
	}

	docJSON, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("ошибка сериализации документа для обновления: %w", err)
	}

	req := opensearchapi.UpdateRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       strings.NewReader(string(docJSON)),
		Refresh:    "true", // Добавляем refresh для немедленной видимости изменений
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("ошибка обновления документа: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("ошибка обновления документа: %s, body: %s", res.String(), string(body))
	}

	return nil
}

// DeleteDocument удаляет документ из индекса
func (c *OpenSearchClient) DeleteDocument(ctx context.Context, indexName, id string) error {
	req := opensearchapi.DeleteRequest{
		Index:      indexName,
		DocumentID: id,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("ошибка удаления документа: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() && res.StatusCode != 404 { // Игнорируем ошибку, если документ не найден
		return fmt.Errorf("ошибка удаления документа: %s", res.String())
	}

	return nil
}

// Search выполняет поиск по индексу
func (c *OpenSearchClient) Search(ctx context.Context, indexName string, query map[string]interface{}) ([]byte, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	req := opensearchapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(string(queryJSON)),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		return nil, fmt.Errorf("ошибка выполнения поиска: %s", res.String())
	}

	// Читаем весь результат
	var result json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка декодирования результата: %w", err)
	}

	return result, nil
}

// Suggest выполняет запрос на автодополнение
func (c *OpenSearchClient) Suggest(ctx context.Context, indexName, field, prefix string, size int) ([]byte, error) {
	suggestQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			"completion": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field": field,
					"size":  size,
				},
			},
		},
	}

	return c.Search(ctx, indexName, suggestQuery)
}

// Execute выполняет прямой запрос к OpenSearch с указанным методом, путём и телом запроса
func (c *OpenSearchClient) Execute(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	// Для универсальных запросов используем существующие методы API
	// В зависимости от типа запроса выбираем соответствующий метод

	var res *opensearchapi.Response
	var err error

	// Обрабатываем запросы к служебным API
	if strings.HasPrefix(path, "/_") || strings.Contains(path, "/_stats") ||
		strings.Contains(path, "/_mapping") || strings.Contains(path, "/_settings") ||
		strings.Contains(path, "/_alias") || strings.Contains(path, "/_count") {
		// Для запросов к служебным API используем Cat API или Indices API
		switch {
		case strings.Contains(path, "/_stats"):
			// Индексная статистика
			indexName := strings.TrimSuffix(path, "/_stats")
			req := opensearchapi.IndicesStatsRequest{
				Index: []string{indexName},
			}
			res, err = req.Do(ctx, c.client)
		case strings.Contains(path, "/_mapping"):
			// Маппинг индекса
			indexName := strings.TrimSuffix(path, "/_mapping")
			req := opensearchapi.IndicesGetMappingRequest{
				Index: []string{indexName},
			}
			res, err = req.Do(ctx, c.client)
		case strings.Contains(path, "/_settings"):
			// Настройки индекса
			indexName := strings.TrimSuffix(path, "/_settings")
			req := opensearchapi.IndicesGetSettingsRequest{
				Index: []string{indexName},
			}
			res, err = req.Do(ctx, c.client)
		case strings.Contains(path, "/_alias"):
			// Алиасы индекса
			indexName := strings.TrimSuffix(path, "/_alias")
			req := opensearchapi.IndicesGetAliasRequest{
				Index: []string{indexName},
			}
			res, err = req.Do(ctx, c.client)
		case strings.Contains(path, "/_count"):
			// Количество документов
			indexName := strings.TrimSuffix(path, "/_count")
			req := opensearchapi.CountRequest{
				Index: []string{indexName},
			}
			res, err = req.Do(ctx, c.client)
		case strings.Contains(path, "/_cluster/health"):
			// Здоровье кластера
			parts := strings.Split(path, "/_cluster/health/")
			var req opensearchapi.ClusterHealthRequest
			if len(parts) > 1 {
				req = opensearchapi.ClusterHealthRequest{
					Index: []string{parts[1]},
				}
			} else {
				req = opensearchapi.ClusterHealthRequest{}
			}
			res, err = req.Do(ctx, c.client)
		default:
			// Для остальных служебных запросов используем Info
			req := opensearchapi.InfoRequest{}
			res, err = req.Do(ctx, c.client)
		}
	} else {
		// Обычные запросы к индексам (поиск и т.д.)
		var bodyReader io.Reader
		if body != nil {
			bodyReader = strings.NewReader(string(body))
		}

		switch method {
		case "GET", "POST":
			switch {
			case strings.Contains(path, "/_search"):
				// Поисковый запрос
				parts := strings.Split(path, "/_search")
				req := opensearchapi.SearchRequest{
					Index: []string{parts[0]},
					Body:  bodyReader,
				}
				res, err = req.Do(ctx, c.client)
			case strings.Contains(path, "/_delete_by_query"):
				// Удаление по запросу
				parts := strings.Split(path, "/_delete_by_query")
				indexName := strings.TrimPrefix(parts[0], "/")
				req := opensearchapi.DeleteByQueryRequest{
					Index: []string{indexName},
					Body:  bodyReader,
				}
				res, err = req.Do(ctx, c.client)
			case path == "/_bulk":
				// Bulk операции
				req := opensearchapi.BulkRequest{
					Body: bodyReader,
				}
				res, err = req.Do(ctx, c.client)
			default:
				// Обычный GET документа
				req := opensearchapi.GetRequest{
					Index:      path,
					DocumentID: "",
				}
				res, err = req.Do(ctx, c.client)
			}
		default:
			return nil, fmt.Errorf("unsupported method: %s", method)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	// Читаем ответ
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Проверяем статус ответа
	if res.IsError() {
		return nil, fmt.Errorf("opensearch error: %s", res.String())
	}

	return responseBody, nil
}

func (c *OpenSearchClient) DeleteIndex(ctx context.Context, indexName string) error {
	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("ошибка удаления индекса: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() && res.StatusCode != 404 { // Игнорируем, если индекс не существует
		return fmt.Errorf("ошибка удаления индекса: %s", res.String())
	}

	return nil
}
