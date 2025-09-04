package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Конфигурация тестирования
const (
	baseURL           = "http://localhost:3000/api/v2"
	numWorkers        = 10  // Количество параллельных воркеров
	requestsPerWorker = 100 // Запросов на воркер
	testDuration      = 30  // Длительность теста в секундах
)

// Метрики
var (
	totalRequests      int64
	successfulRequests int64
	failedRequests     int64
	totalResponseTime  int64
	minResponseTime    int64 = 999999
	maxResponseTime    int64
)

// Структуры для запросов
type AttributeRequest struct {
	CategoryID int  `json:"category_id,omitempty"`
	Required   bool `json:"required,omitempty"`
	Searchable bool `json:"searchable,omitempty"`
}

type AttributeValue struct {
	AttributeID  int     `json:"attribute_id"`
	TextValue    string  `json:"text_value,omitempty"`
	NumericValue float64 `json:"numeric_value,omitempty"`
	BooleanValue bool    `json:"boolean_value,omitempty"`
}

// JWT токен для авторизации (нужно получить из backend/scripts/create_test_jwt.go)
var authToken string

func init() {
	// Получаем токен из переменной окружения или генерируем
	authToken = os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		log.Println("WARNING: AUTH_TOKEN not set, некоторые запросы могут не работать")
	}
}

// Функция для выполнения HTTP запроса с метриками
func makeRequest(method, url string, body interface{}) (*http.Response, time.Duration) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		atomic.AddInt64(&failedRequests, 1)
		return nil, 0
	}

	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	start := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	duration := time.Since(start)

	atomic.AddInt64(&totalRequests, 1)
	atomic.AddInt64(&totalResponseTime, int64(duration.Milliseconds()))

	if err != nil {
		atomic.AddInt64(&failedRequests, 1)
		return nil, duration
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		atomic.AddInt64(&successfulRequests, 1)
	} else {
		atomic.AddInt64(&failedRequests, 1)
	}

	// Обновляем min/max время ответа
	ms := int64(duration.Milliseconds())
	for {
		oldMin := atomic.LoadInt64(&minResponseTime)
		if ms >= oldMin || atomic.CompareAndSwapInt64(&minResponseTime, oldMin, ms) {
			break
		}
	}
	for {
		oldMax := atomic.LoadInt64(&maxResponseTime)
		if ms <= oldMax || atomic.CompareAndSwapInt64(&maxResponseTime, oldMax, ms) {
			break
		}
	}

	return resp, duration
}

// Тест получения атрибутов категории
func testGetCategoryAttributes(categoryID int) {
	url := fmt.Sprintf("%s/attributes/category/%d", baseURL, categoryID)
	resp, duration := makeRequest("GET", url, nil)

	if resp != nil {
		defer resp.Body.Close()
		log.Printf("[GET] Категория %d - Статус: %d, Время: %dms",
			categoryID, resp.StatusCode, duration.Milliseconds())
	}
}

// Тест получения атрибутов с фильтрами
func testGetAttributesWithFilters() {
	// Случайные параметры фильтрации
	filters := []string{
		"?required=true",
		"?searchable=true",
		"?filterable=true",
		"?category_id=1",
		"?page=1&limit=20",
		"?required=true&searchable=true",
		"",
	}

	filter := filters[rand.Intn(len(filters))]
	url := fmt.Sprintf("%s/attributes%s", baseURL, filter)

	resp, duration := makeRequest("GET", url, nil)
	if resp != nil {
		defer resp.Body.Close()
		log.Printf("[GET] Атрибуты с фильтром '%s' - Статус: %d, Время: %dms",
			filter, resp.StatusCode, duration.Milliseconds())
	}
}

// Тест создания значений атрибутов
func testCreateAttributeValues(listingID int) {
	// Генерируем случайные значения атрибутов
	values := []AttributeValue{
		{
			AttributeID: rand.Intn(10) + 1,
			TextValue:   fmt.Sprintf("Test value %d", rand.Intn(1000)),
		},
		{
			AttributeID:  rand.Intn(10) + 10,
			NumericValue: rand.Float64() * 1000,
		},
		{
			AttributeID:  rand.Intn(10) + 20,
			BooleanValue: rand.Intn(2) == 1,
		},
	}

	url := fmt.Sprintf("%s/listings/%d/attributes/batch", baseURL, listingID)
	resp, duration := makeRequest("POST", url, map[string]interface{}{
		"attributes": values,
	})

	if resp != nil {
		defer resp.Body.Close()
		log.Printf("[POST] Создание атрибутов для листинга %d - Статус: %d, Время: %dms",
			listingID, resp.StatusCode, duration.Milliseconds())
	}
}

// Воркер для выполнения тестов
func worker(id int, wg *sync.WaitGroup, stop chan bool) {
	defer wg.Done()

	log.Printf("Воркер %d запущен", id)

	for {
		select {
		case <-stop:
			log.Printf("Воркер %d остановлен", id)
			return
		default:
			// Выполняем случайный тест
			testType := rand.Intn(3)
			switch testType {
			case 0:
				// Тест получения атрибутов категории
				categoryID := rand.Intn(100) + 1
				testGetCategoryAttributes(categoryID)
			case 1:
				// Тест получения атрибутов с фильтрами
				testGetAttributesWithFilters()
			case 2:
				// Тест создания значений атрибутов
				listingID := rand.Intn(1000) + 1
				testCreateAttributeValues(listingID)
			}

			// Небольшая задержка между запросами
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		}
	}
}

// Функция для вывода статистики в реальном времени
func printStats(stop chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			total := atomic.LoadInt64(&totalRequests)
			successful := atomic.LoadInt64(&successfulRequests)
			failed := atomic.LoadInt64(&failedRequests)
			totalTime := atomic.LoadInt64(&totalResponseTime)
			min := atomic.LoadInt64(&minResponseTime)
			max := atomic.LoadInt64(&maxResponseTime)

			elapsed := time.Since(startTime).Seconds()
			rps := float64(total) / elapsed

			avgTime := int64(0)
			if total > 0 {
				avgTime = totalTime / total
			}

			successRate := float64(0)
			if total > 0 {
				successRate = float64(successful) * 100 / float64(total)
			}

			fmt.Println("\n========== СТАТИСТИКА ==========")
			fmt.Printf("Время работы: %.0f сек\n", elapsed)
			fmt.Printf("Всего запросов: %d\n", total)
			fmt.Printf("Успешных: %d (%.2f%%)\n", successful, successRate)
			fmt.Printf("Неудачных: %d\n", failed)
			fmt.Printf("RPS: %.2f\n", rps)
			fmt.Printf("Время ответа: Min=%dms, Avg=%dms, Max=%dms\n", min, avgTime, max)
			fmt.Println("================================")
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("===========================================")
	fmt.Println("   НАГРУЗОЧНОЕ ТЕСТИРОВАНИЕ")
	fmt.Println("   UNIFIED ATTRIBUTES API V2")
	fmt.Println("===========================================")
	fmt.Printf("URL: %s\n", baseURL)
	fmt.Printf("Воркеров: %d\n", numWorkers)
	fmt.Printf("Длительность: %d сек\n", testDuration)
	fmt.Println("===========================================\n")

	// Каналы для управления
	stopWorkers := make(chan bool)
	stopStats := make(chan bool)

	// Запускаем вывод статистики
	go printStats(stopStats)

	// Запускаем воркеров
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, stopWorkers)
	}

	// Ждем заданное время
	time.Sleep(time.Duration(testDuration) * time.Second)

	// Останавливаем воркеров
	fmt.Println("\nОстановка тестирования...")
	close(stopWorkers)
	wg.Wait()

	// Останавливаем вывод статистики
	close(stopStats)

	// Финальная статистика
	total := atomic.LoadInt64(&totalRequests)
	successful := atomic.LoadInt64(&successfulRequests)
	failed := atomic.LoadInt64(&failedRequests)
	totalTime := atomic.LoadInt64(&totalResponseTime)
	min := atomic.LoadInt64(&minResponseTime)
	max := atomic.LoadInt64(&maxResponseTime)

	avgTime := int64(0)
	if total > 0 {
		avgTime = totalTime / total
	}

	successRate := float64(0)
	if total > 0 {
		successRate = float64(successful) * 100 / float64(total)
	}

	fmt.Println("\n===========================================")
	fmt.Println("         ФИНАЛЬНАЯ СТАТИСТИКА")
	fmt.Println("===========================================")
	fmt.Printf("Всего запросов: %d\n", total)
	fmt.Printf("Успешных: %d (%.2f%%)\n", successful, successRate)
	fmt.Printf("Неудачных: %d\n", failed)
	fmt.Printf("RPS: %.2f\n", float64(total)/float64(testDuration))
	fmt.Printf("Время ответа:\n")
	fmt.Printf("  Минимальное: %d мс\n", min)
	fmt.Printf("  Среднее: %d мс\n", avgTime)
	fmt.Printf("  Максимальное: %d мс\n", max)
	fmt.Println("===========================================")

	// Выводим результат
	if successRate >= 95 {
		fmt.Println("✅ ТЕСТ ПРОЙДЕН УСПЕШНО!")
	} else if successRate >= 80 {
		fmt.Println("⚠️  ТЕСТ ПРОЙДЕН С ПРЕДУПРЕЖДЕНИЯМИ")
	} else {
		fmt.Println("❌ ТЕСТ НЕ ПРОЙДЕН")
		os.Exit(1)
	}
}
