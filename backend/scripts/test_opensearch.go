//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go/v2"
)

func main() {
	fmt.Println("=== OpenSearch Integration Test ===")
	fmt.Println()

	// Загружаем .env файл
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("⚠️  Warning: .env file not found, using defaults\n\n")
	}

	// Получаем конфигурацию
	opensearchURL := getEnv("OPENSEARCH_URL", "http://localhost:9200")
	username := getEnv("OPENSEARCH_USERNAME", "admin")
	password := getEnv("OPENSEARCH_PASSWORD", "admin")

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   OpenSearch URL: %s\n", opensearchURL)
	fmt.Printf("   Username: %s\n", username)
	fmt.Println()

	// Создаем клиент
	cfg := opensearch.Config{
		Addresses: []string{opensearchURL},
		Username:  username,
		Password:  password,
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("❌ Failed to create OpenSearch client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Test 1: Cluster health
	fmt.Println("🏥 Test 1: Cluster Health Check")
	fmt.Println("--------------------------------")

	healthRes, err := client.Cluster.Health(
		client.Cluster.Health.WithContext(ctx),
	)
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		os.Exit(1)
	}
	defer healthRes.Body.Close()

	if healthRes.IsError() {
		fmt.Printf("❌ Health check error: %s\n", healthRes.Status())
		os.Exit(1)
	}

	var health map[string]interface{}
	if err := json.NewDecoder(healthRes.Body).Decode(&health); err != nil {
		fmt.Printf("❌ Failed to parse health response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Cluster Status: %v\n", health["status"])
	fmt.Printf("   Cluster Name: %v\n", health["cluster_name"])
	fmt.Printf("   Number of Nodes: %v\n", health["number_of_nodes"])
	fmt.Printf("   Active Shards: %v\n", health["active_shards"])
	fmt.Println()

	// Test 2: List indices
	fmt.Println("📚 Test 2: List Indices")
	fmt.Println("-----------------------")

	indicesRes, err := client.Cat.Indices(
		client.Cat.Indices.WithContext(ctx),
		client.Cat.Indices.WithFormat("json"),
	)
	if err != nil {
		fmt.Printf("❌ Failed to list indices: %v\n", err)
		os.Exit(1)
	}
	defer indicesRes.Body.Close()

	if indicesRes.IsError() {
		fmt.Printf("❌ List indices error: %s\n", indicesRes.Status())
		os.Exit(1)
	}

	var indices []map[string]interface{}
	if err := json.NewDecoder(indicesRes.Body).Decode(&indices); err != nil {
		fmt.Printf("❌ Failed to parse indices: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Found %d indices\n", len(indices))

	// Показываем только наши индексы (не системные)
	fmt.Println("   Application indices:")
	for _, index := range indices {
		indexName := index["index"].(string)
		if indexName[0] != '.' { // Skip system indices
			docsCount := index["docs.count"]
			storeSize := index["store.size"]
			fmt.Printf("   - %s (docs: %v, size: %v)\n", indexName, docsCount, storeSize)
		}
	}
	fmt.Println()

	// Test 3: Create test index and document
	testIndexName := "test-integration-" + fmt.Sprintf("%d", time.Now().Unix())

	fmt.Println("📝 Test 3: Create Index and Document")
	fmt.Println("-------------------------------------")
	fmt.Printf("Creating test index: %s\n", testIndexName)

	// Create index
	createRes, err := client.Indices.Create(testIndexName)
	if err != nil {
		fmt.Printf("❌ Failed to create index: %v\n", err)
		os.Exit(1)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		fmt.Printf("❌ Create index error: %s\n", createRes.Status())
		os.Exit(1)
	}
	fmt.Println("✅ Index created")

	// Index a test document
	testDoc := map[string]interface{}{
		"title":       "Test Document",
		"description": "This is a test document for integration testing",
		"timestamp":   time.Now().Format(time.RFC3339),
		"count":       42,
	}

	docJSON, _ := json.Marshal(testDoc)
	indexRes, err := client.Index(
		testIndexName,
		bytes.NewReader(docJSON),
		client.Index.WithContext(ctx),
		client.Index.WithDocumentID("test-1"),
		client.Index.WithRefresh("true"), // Force refresh for immediate search
	)
	if err != nil {
		fmt.Printf("❌ Failed to index document: %v\n", err)
		os.Exit(1)
	}
	defer indexRes.Body.Close()

	if indexRes.IsError() {
		fmt.Printf("❌ Index document error: %s\n", indexRes.Status())
		os.Exit(1)
	}
	fmt.Println("✅ Document indexed")
	fmt.Println()

	// Test 4: Search
	fmt.Println("🔍 Test 4: Search Document")
	fmt.Println("--------------------------")

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "Test",
			},
		},
	}

	queryJSON, _ := json.Marshal(query)
	searchRes, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(testIndexName),
		client.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		fmt.Printf("❌ Search failed: %v\n", err)
		os.Exit(1)
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		fmt.Printf("❌ Search error: %s\n", searchRes.Status())
		os.Exit(1)
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(searchRes.Body).Decode(&searchResult); err != nil {
		fmt.Printf("❌ Failed to parse search results: %v\n", err)
		os.Exit(1)
	}

	hits := searchResult["hits"].(map[string]interface{})
	totalHits := hits["total"].(map[string]interface{})["value"].(float64)

	fmt.Printf("✅ Search successful: found %.0f documents\n", totalHits)
	if totalHits > 0 {
		firstHit := hits["hits"].([]interface{})[0].(map[string]interface{})
		source := firstHit["_source"].(map[string]interface{})
		fmt.Printf("   First result: %s\n", source["title"])
	}
	fmt.Println()

	// Test 5: Delete test index
	fmt.Println("🗑️  Test 5: Cleanup")
	fmt.Println("------------------")

	deleteRes, err := client.Indices.Delete([]string{testIndexName})
	if err != nil {
		fmt.Printf("❌ Failed to delete index: %v\n", err)
		os.Exit(1)
	}
	defer deleteRes.Body.Close()

	if deleteRes.IsError() {
		fmt.Printf("❌ Delete index error: %s\n", deleteRes.Status())
		os.Exit(1)
	}
	fmt.Printf("✅ Test index deleted: %s\n", testIndexName)
	fmt.Println()

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  ✅ All OpenSearch tests passed! 🎉   ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
