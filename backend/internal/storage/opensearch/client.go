package opensearch

import (
    "context"
    "encoding/json"
    "fmt"
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
func (c *OpenSearchClient) CreateIndex(indexName string, mapping string) error {
    exists, err := c.IndexExists(indexName)
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
    
    res, err := createIndex.Do(context.Background(), c.client)
    if err != nil {
        return fmt.Errorf("ошибка создания индекса: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("ошибка создания индекса: %s", res.String())
    }
    
    return nil
}

// IndexExists проверяет существование индекса
func (c *OpenSearchClient) IndexExists(indexName string) (bool, error) {
    exists := opensearchapi.IndicesExistsRequest{
        Index: []string{indexName},
    }
    
    res, err := exists.Do(context.Background(), c.client)
    if err != nil {
        return false, fmt.Errorf("ошибка проверки индекса: %w", err)
    }
    
    return res.StatusCode == 200, nil
}

// IndexDocument индексирует документ
func (c *OpenSearchClient) IndexDocument(indexName, id string, document interface{}) error {
    docBytes, err := json.Marshal(document)
    if err != nil {
        return fmt.Errorf("ошибка сериализации документа: %w", err)
    }
    
    req := opensearchapi.IndexRequest{
        Index:      indexName,
        DocumentID: id,
        Body:       strings.NewReader(string(docBytes)),
    }
    
    res, err := req.Do(context.Background(), c.client)
    if err != nil {
        return fmt.Errorf("ошибка индексации документа: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("ошибка индексации документа: %s", res.String())
    }
    
    return nil
}

// BulkIndex индексирует несколько документов за один запрос
func (c *OpenSearchClient) BulkIndex(indexName string, documents []map[string]interface{}) error {
    if len(documents) == 0 {
        return nil
    }
    
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
        
        // Удаляем ID из документа для индексации
        docCopy := make(map[string]interface{})
        for k, v := range doc {
            if k != "id" {
                docCopy[k] = v
            }
        }
        
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
    
    res, err := req.Do(context.Background(), c.client)
    if err != nil {
        return fmt.Errorf("ошибка bulk индексации: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("ошибка bulk индексации: %s", res.String())
    }
    
    return nil
}

// DeleteDocument удаляет документ из индекса
func (c *OpenSearchClient) DeleteDocument(indexName, id string) error {
    req := opensearchapi.DeleteRequest{
        Index:      indexName,
        DocumentID: id,
    }
    
    res, err := req.Do(context.Background(), c.client)
    if err != nil {
        return fmt.Errorf("ошибка удаления документа: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() && res.StatusCode != 404 { // Игнорируем ошибку, если документ не найден
        return fmt.Errorf("ошибка удаления документа: %s", res.String())
    }
    
    return nil
}

// Search выполняет поиск по индексу
func (c *OpenSearchClient) Search(indexName string, query map[string]interface{}) ([]byte, error) {
    queryJSON, err := json.Marshal(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
    }
    
    req := opensearchapi.SearchRequest{
        Index: []string{indexName},
        Body:  strings.NewReader(string(queryJSON)),
    }
    
    res, err := req.Do(context.Background(), c.client)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения поиска: %w", err)
    }
    defer res.Body.Close()
    
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
func (c *OpenSearchClient) Suggest(indexName, field, prefix string, size int) ([]byte, error) {
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
    
    return c.Search(indexName, suggestQuery)
}
func (c *OpenSearchClient) DeleteIndex(indexName string) error {
    req := opensearchapi.IndicesDeleteRequest{
        Index: []string{indexName},
    }
    
    res, err := req.Do(context.Background(), c.client)
    if err != nil {
        return fmt.Errorf("ошибка удаления индекса: %w", err)
    }
    defer res.Body.Close()
    
    if res.IsError() && res.StatusCode != 404 { // Игнорируем, если индекс не существует
        return fmt.Errorf("ошибка удаления индекса: %s", res.String())
    }
    
    return nil
}