package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/internal/proj/postexpress"
	"backend/internal/proj/postexpress/models"
	"backend/pkg/logger"
)

// Constants
const (
	errMsgUnknown = "unknown error"
)

// Singleton instance
var (
	wspClientInstance WSPClient
	wspClientOnce     sync.Once
	wspClientMutex    sync.RWMutex
)

// WSPClientImpl представляет реализацию клиента WSP API
type WSPClientImpl struct {
	httpClient *http.Client
	config     *WSPConfig
	logger     logger.Logger
}

// WSPConfig представляет конфигурацию WSP клиента
type WSPConfig struct {
	Endpoint        string
	Username        string
	Password        string
	Language        string
	DeviceType      string // ИСПРАВЛЕНО: должен быть string ("2"), не int
	Timeout         time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
	TestMode        bool
	DeviceName      string
	ApplicationName string
	Version         string
	PartnerID       int    // Добавлено: ID партнера (10109 для svetu.rs)
	BankAccount     string // Банковский счёт для перевода откупнины
	PaymentCode     string // Шифра плаћања (обычно "189")
	PaymentModel    string // Модель платежа (обычно "97")
}

// NewWSPClient создает или возвращает singleton экземпляр WSP клиента
// Thread-safe singleton с lazy initialization
func NewWSPClient(config *WSPConfig, logger logger.Logger) WSPClient {
	wspClientOnce.Do(func() {
		// Настройка HTTP клиента с timeout и SSL
		httpClient := &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.TestMode, // #nosec G402 - только для тестового режима
				},
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}

		wspClientInstance = &WSPClientImpl{
			httpClient: httpClient,
			config:     config,
			logger:     logger,
		}

		logger.Info("WSPClient singleton initialized")
	})

	// Возвращаем singleton экземпляр
	wspClientMutex.RLock()
	defer wspClientMutex.RUnlock()
	return wspClientInstance
}

// ResetWSPClientSingleton сбрасывает singleton (только для тестов)
func ResetWSPClientSingleton() {
	wspClientMutex.Lock()
	defer wspClientMutex.Unlock()
	wspClientInstance = nil
	wspClientOnce = sync.Once{}
}

// Transaction выполняет базовую транзакцию к WSP API
func (c *WSPClientImpl) Transaction(ctx context.Context, req *models.TransactionRequest) (*models.TransactionResponse, error) {
	// Подготовка клиентских данных
	clientData := &models.ClientData{
		Username:          c.config.Username,
		Password:          c.config.Password,
		Jezik:             c.config.Language,
		IdTipUredjaja:     "2", // ВАЖНО: должна быть строка "2" для веб-приложения
		VerzijaOS:         "Linux",
		NazivUredjaja:     c.config.DeviceName,
		ModelUredjaja:     "API",
		VerzijaAplikacije: c.config.Version,
		IPAdresa:          c.getServerIP(),
		Geolokacija:       nil,
		IdPartnera:        c.config.PartnerID, // ДОБАВЛЕНО: ID партнера для B2B интеграции
	}

	clientJSON, err := json.Marshal(clientData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal client data: %w", err)
	}

	// Подготовка основного запроса
	transReq := &models.TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101,                 // ИСПРАВЛЕНО: 101 для B2B партнеров (не 3!)
		IdVrstaTransakcije: req.TransactionType, // ИСПРАВЛЕНО: IdVrstaTransakcije (не IdVrstaTranskacije!)
		TipSerijalizacije:  2,                   // JSON (2 для JSON согласно документации)
		IdTransakcija:      models.GenerateGUID(),
		StrIn:              req.InputData,
	}

	// Логирование запроса (без пароля)
	c.logger.Debug("WSP API Request - transaction_id: %s, type: %d, endpoint: %s",
		transReq.IdTransakcija, req.TransactionType, c.config.Endpoint)

	// Выполнение HTTP запроса с ретраями
	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err := c.executeRequest(ctx, transReq)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if attempt < c.config.MaxRetries {
			c.logger.Error("WSP API request failed, retrying attempt %d/%d: %v",
				attempt+1, c.config.MaxRetries, err)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.config.RetryDelay):
				// Продолжаем с следующей попыткой
			}
		}
	}

	return nil, fmt.Errorf("WSP API request failed after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

// executeRequest выполняет HTTP запрос
func (c *WSPClientImpl) executeRequest(ctx context.Context, transReq *models.TransakcijaIn) (*models.TransactionResponse, error) {
	// Сериализация запроса
	body, err := json.Marshal(transReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создание HTTP запроса
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Установка заголовков
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", fmt.Sprintf("%s/%s", c.config.ApplicationName, c.config.Version))

	// Выполнение запроса
	startTime := time.Now()
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		if closeErr := httpResp.Body.Close(); closeErr != nil {
			c.logger.Error("Failed to close HTTP response body: %v", closeErr)
		}
	}()

	executionTime := time.Since(startTime)

	// Логирование времени выполнения
	c.logger.Debug("WSP API Response - status_code: %d, execution_time_ms: %d",
		httpResp.StatusCode, executionTime.Milliseconds())

	// Проверка статуса ответа
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	// Читаем тело ответа в байты для логирования
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Логируем сырой ответ для диагностики
	c.logger.Debug("WSP API Raw Response Body (first 2000 chars): %s", truncateString(string(bodyBytes), 2000))

	// Декодирование ответа
	var wspResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &wspResp); err != nil {
		c.logger.Error("Failed to parse JSON response. Full raw response: %s", string(bodyBytes))
		return nil, fmt.Errorf("failed to parse manifest response: %w", err)
	}

	// Обработка ответа WSP
	return c.parseWSPResponse(wspResp)
}

// parseWSPResponse обрабатывает ответ от WSP API
func (c *WSPClientImpl) parseWSPResponse(resp map[string]interface{}) (*models.TransactionResponse, error) {
	// Проверяем наличие ошибки
	if errorMsg, exists := resp["Greska"]; exists && errorMsg != nil {
		return &models.TransactionResponse{
			Success:      false,
			ErrorMessage: stringPtr(fmt.Sprintf("%v", errorMsg)),
		}, nil
	}

	// Проверяем успешность транзакции
	success := true

	// Проверяем поле Rezultat (0 = успех, иначе ошибка)
	// ВАЖНО: Для B2B Manifest API результат может быть Rezultat!=0 на уровне транзакции,
	// но манифест может быть успешно создан (Rezultat=0 внутри StrOut)!
	// Поэтому сначала проверяем StrOut, и только если его нет - смотрим на внешний Rezultat
	if strOut, exists := resp["StrOut"]; exists && strOut != nil {
		if strOutStr, ok := strOut.(string); ok {
			// Логируем весь StrOut для диагностики
			c.logger.Debug("Full StrOut content (length %d): %s", len(strOutStr), strOutStr)

			// Парсим манифест из StrOut для проверки реального результата
			var manifestResp struct {
				Rezultat int    `json:"Rezultat"`
				Poruka   string `json:"Poruka"`
				Greske   []struct {
					ExtIDManifest   string `json:"ExtIdManifest"`
					ExtIDPorudzbina string `json:"ExtIdPorudzbina"`
					Rbr             int    `json:"Rbr"`
					PorukaGreske    string `json:"PorukaGreske"`
				} `json:"Greske"`
			}
			if err := json.Unmarshal([]byte(strOutStr), &manifestResp); err != nil {
				c.logger.Error("Failed to parse StrOut as manifest: %v", err)
			} else {
				c.logger.Debug("Parsed manifest - Rezultat: %d, Poruka: %s, Errors count: %d",
					manifestResp.Rezultat, manifestResp.Poruka, len(manifestResp.Greske))

				// РЕАЛЬНЫЙ результат берем из ВНУТРЕННЕГО Rezultat (в StrOut)
				if manifestResp.Rezultat != 0 {
					success = false
					c.logger.Error("Manifest creation failed - Rezultat: %d, Poruka: %s", manifestResp.Rezultat, manifestResp.Poruka)
				} else {
					success = true
					c.logger.Info("Manifest created successfully - Rezultat: 0")
				}

				// Логируем ошибки валидации (они могут быть даже при успехе - это warnings!)
				if len(manifestResp.Greske) > 0 {
					c.logger.Info("Post Express validation warnings (%d warnings):", len(manifestResp.Greske))
					for i, validErr := range manifestResp.Greske {
						c.logger.Info("  [%d] Manifest: %s, Order: %s, Rbr: %d, Message: %s",
							i+1, validErr.ExtIDManifest, validErr.ExtIDPorudzbina, validErr.Rbr, validErr.PorukaGreske)
					}
				}
			}
		}
	} else if rezultatField, exists := resp["Rezultat"]; exists {
		// Fallback: если нет StrOut, проверяем внешний Rezultat
		if rezultat, ok := rezultatField.(float64); ok {
			if rezultat != 0 {
				poruka := "unknown error"

				// Сначала пытаемся извлечь детальную ошибку из StrRezultat
				if strRezultat, exists := resp["StrRezultat"]; exists && strRezultat != nil {
					if strRezultatStr, ok := strRezultat.(string); ok {
						// Парсим StrRezultat как JSON для получения детальной ошибки
						var detailErr struct {
							Poruka         string `json:"Poruka"`
							PorukaKorisnik string `json:"PorukaKorisnik"`
							Info           string `json:"Info"`
						}
						if err := json.Unmarshal([]byte(strRezultatStr), &detailErr); err == nil {
							// Используем PorukaKorisnik если есть, иначе Poruka
							if detailErr.PorukaKorisnik != "" {
								poruka = detailErr.PorukaKorisnik
							} else if detailErr.Poruka != "" {
								poruka = detailErr.Poruka
							}
						}
					}
				}

				// Если не смогли извлечь из StrRezultat, используем обычное поле Poruka
				if poruka == "unknown error" {
					if porukaField, exists := resp["Poruka"]; exists && porukaField != nil {
						poruka = fmt.Sprintf("%v", porukaField)
					}
				}

				c.logger.Error("WSP transaction failed - Rezultat: %d, Poruka: %s", int(rezultat), poruka)

				// ВАЖНО: Возвращаем ошибку сразу с детальным сообщением
				return &models.TransactionResponse{
					Success:      false,
					ErrorMessage: &poruka,
				}, nil
			}
		}
	}

	// Также проверяем старое поле Uspesno (для обратной совместимости)
	if successField, exists := resp["Uspesno"]; exists {
		if b, ok := successField.(bool); ok {
			success = b
		}
	}

	// Извлекаем данные ответа
	var outputData json.RawMessage
	if data, exists := resp["Podaci"]; exists && data != nil {
		if dataBytes, err := json.Marshal(data); err == nil {
			outputData = dataBytes
		}
	} else if strOut, exists := resp["StrOut"]; exists && strOut != nil {
		// StrOut содержит JSON как строку - нужно распарсить его
		if strOutStr, ok := strOut.(string); ok {
			// Это уже JSON string, просто конвертируем в []byte
			outputData = json.RawMessage(strOutStr)
		} else {
			// Если не строка, пытаемся marshal
			if dataBytes, err := json.Marshal(strOut); err == nil {
				outputData = dataBytes
			}
		}
	}

	return &models.TransactionResponse{
		Success:    success,
		OutputData: outputData,
	}, nil
}

// GetLocations ищет населенные пункты
func (c *WSPClientImpl) GetLocations(ctx context.Context, search string) ([]WSPLocation, error) {
	// Подготовка входных данных для поиска населенных пунктов
	searchReq := map[string]interface{}{
		"Naziv":           search,
		"BrojSlogova":     50, // Максимум результатов
		"NacinSortiranja": 0,  // Тип сортировки
	}

	inputData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	// Выполнение запроса (тип транзакции 3 для поиска населенных пунктов)
	req := &models.TransactionRequest{
		TransactionType: 3,
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetLocations transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("GetLocations failed: %s", errMsg)
	}

	// Парсинг результатов
	var result struct {
		Naselja []WSPLocation `json:"Naselja"`
	}

	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse locations response: %w", err)
	}

	return result.Naselja, nil
}

// GetOffices получает список отделений для населенного пункта
func (c *WSPClientImpl) GetOffices(ctx context.Context, locationID int) ([]WSPOffice, error) {
	// Подготовка запроса для получения отделений
	officesReq := map[string]interface{}{
		"IdNaselje": locationID,
	}

	inputData, err := json.Marshal(officesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal offices request: %w", err)
	}

	// Выполнение запроса (тип транзакции 10 для получения отделений)
	req := &models.TransactionRequest{
		TransactionType: 10,
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetOffices transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("GetOffices failed: %s", errMsg)
	}

	// Парсинг результатов
	var result struct {
		PostanskeJedinice []WSPOffice `json:"PostanskeJedinice"`
	}

	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse offices response: %w", err)
	}

	return result.PostanskeJedinice, nil
}

// CreateShipment создает новое отправление через манифест (транзакция 73)
// ВАЖНО: Транзакция 63 предназначена только для отслеживания, не для создания!
func (c *WSPClientImpl) CreateShipment(ctx context.Context, shipment *WSPShipmentRequest) (*WSPShipmentResponse, error) {
	// Используем новый метод создания через манифест B2B (транзакция 73)
	manifestResp, err := c.CreateShipmentViaManifest(ctx, shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment via manifest: %w", err)
	}

	// Преобразуем ответ манифеста в формат WSPShipmentResponse
	// Новая структура: ManifestResponse с Rezultat (0=success), Poruka, Porudzbine[]
	if manifestResp.Rezultat != 0 {
		return &WSPShipmentResponse{
			Success:      false,
			ErrorMessage: manifestResp.Poruka,
		}, nil
	}

	// Проверяем результат создания посылки
	if len(manifestResp.Porudzbine) == 0 || len(manifestResp.Porudzbine[0].Posiljke) == 0 {
		return &WSPShipmentResponse{
			Success:      false,
			ErrorMessage: "no shipment result in manifest response",
		}, nil
	}

	posiljkaResult := manifestResp.Porudzbine[0].Posiljke[0]
	if posiljkaResult.Rezultat != 0 {
		return &WSPShipmentResponse{
			Success:      false,
			ErrorMessage: posiljkaResult.Poruka,
		}, nil
	}

	// Возвращаем успешный результат
	return &WSPShipmentResponse{
		Success:         true,
		TrackingNumber:  posiljkaResult.TrackingNumber,
		Barcode:         "", // Баркод не возвращается в B2B API
		ReferenceNumber: posiljkaResult.BrojPosiljke,
	}, nil
}

// GetShipmentStatus получает статус отправления
func (c *WSPClientImpl) GetShipmentStatus(ctx context.Context, trackingNumber string) (*WSPTrackingResponse, error) {
	trackingReq := map[string]interface{}{
		"BrojPosiljke": trackingNumber,
	}

	inputData, err := json.Marshal(trackingReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tracking request: %w", err)
	}

	// Выполнение запроса (тип транзакции 63 для отслеживания - как указал Nikola Dmitrašinović)
	req := &models.TransactionRequest{
		TransactionType: 63,
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetShipmentStatus transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("GetShipmentStatus failed: %s", errMsg)
	}

	// Парсинг результата
	var result WSPTrackingResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse tracking response: %w", err)
	}

	return &result, nil
}

// PrintLabel получает этикетку отправления
func (c *WSPClientImpl) PrintLabel(ctx context.Context, shipmentID string) ([]byte, error) {
	labelReq := map[string]interface{}{
		"IdPosiljke": shipmentID,
		"TipFajla":   "PDF", // PDF формат
	}

	inputData, err := json.Marshal(labelReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal label request: %w", err)
	}

	// Выполнение запроса (тип транзакции 20 для печати этикетки)
	req := &models.TransactionRequest{
		TransactionType: 20,
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("PrintLabel transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("PrintLabel failed: %s", errMsg)
	}

	// Результат должен содержать base64 encoded PDF
	var result struct {
		PDFContent string `json:"SadrzajPDF"`
	}

	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse label response: %w", err)
	}

	// Декодируем base64 content в []byte
	decodedPDF, err := base64.StdEncoding.DecodeString(result.PDFContent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 PDF content: %w", err)
	}

	return decodedPDF, nil
}

// CancelShipment отменяет отправление
func (c *WSPClientImpl) CancelShipment(ctx context.Context, shipmentID string) error {
	cancelReq := map[string]interface{}{
		"IdPosiljke": shipmentID,
		"Razlog":     "Отмена по требованию клиента",
	}

	inputData, err := json.Marshal(cancelReq)
	if err != nil {
		return fmt.Errorf("failed to marshal cancel request: %w", err)
	}

	// Выполнение запроса (тип транзакции 25 для отмены)
	req := &models.TransactionRequest{
		TransactionType: 25,
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return fmt.Errorf("CancelShipment transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return fmt.Errorf("CancelShipment failed: %s", errMsg)
	}

	return nil
}

// getServerIP получает IP адрес сервера
func (c *WSPClientImpl) getServerIP() string {
	// Пытаемся получить внешний IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			c.logger.Error("Failed to close connection: %v", closeErr)
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// stringPtr возвращает указатель на строку
func stringPtr(s string) *string {
	return &s
}

// Вспомогательные функции для конвертации типов

// ParseFloat безопасно парсит float из interface{}
func ParseFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0
}

// ParseInt безопасно парсит int из interface{}
func ParseInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// ParseString безопасно парсит string из interface{}
func ParseString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// ParseBool безопасно парсит bool из interface{}
func ParseBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1"
	case int:
		return v != 0
	}
	return false
}

// truncateString обрезает строку до максимальной длины
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ========================================
// TX 3-11: Новые методы для полной интеграции Post Express WSP API
// ========================================

// GetSettlements - TX 3: Поиск населённых пунктов по названию
func (c *WSPClientImpl) GetSettlements(ctx context.Context, query string) (*postexpress.GetSettlementsResponse, error) {
	searchReq := &postexpress.GetSettlementsRequest{
		Naziv: query,
	}

	inputData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GetSettlements request: %w", err)
	}

	req := &models.TransactionRequest{
		TransactionType: 3, // GetNaselje
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetSettlements transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("GetSettlements failed: %s", errMsg)
	}

	// Парсинг результатов
	var result postexpress.GetSettlementsResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse settlements response: %w", err)
	}

	c.logger.Debug("GetSettlements completed - Query: %s, Found: %d settlements, Rezultat: %d",
		query, len(result.Naselja), result.Rezultat)

	return &result, nil
}

// GetStreets - TX 4: Поиск улиц в населённом пункте
func (c *WSPClientImpl) GetStreets(ctx context.Context, settlementID int, query string) (*postexpress.GetStreetsResponse, error) {
	searchReq := &postexpress.GetStreetsRequest{
		IdNaselje: settlementID,
		Naziv:     query,
	}

	inputData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GetStreets request: %w", err)
	}

	req := &models.TransactionRequest{
		TransactionType: 4, // GetUlica
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetStreets transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("GetStreets failed: %s", errMsg)
	}

	// Парсинг результатов
	var result postexpress.GetStreetsResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse streets response: %w", err)
	}

	c.logger.Debug("GetStreets completed - SettlementID: %d, Query: %s, Found: %d streets, Rezultat: %d",
		settlementID, query, len(result.Ulice), result.Rezultat)

	return &result, nil
}

// ValidateAddress - TX 6: Валидация адреса
func (c *WSPClientImpl) ValidateAddress(ctx context.Context, req *postexpress.AddressValidationRequest) (*postexpress.AddressValidationResponse, error) {
	inputData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ValidateAddress request: %w", err)
	}

	txReq := &models.TransactionRequest{
		TransactionType: 6, // ProveraAdrese
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, txReq)
	if err != nil {
		return nil, fmt.Errorf("ValidateAddress transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("ValidateAddress failed: %s", errMsg)
	}

	// Парсинг результатов
	var result postexpress.AddressValidationResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse address validation response: %w", err)
	}

	c.logger.Debug("ValidateAddress completed - SettlementID: %d, PostojiAdresa: %t, Rezultat: %d",
		req.IdNaselje, result.PostojiAdresa, result.Rezultat)

	return &result, nil
}

// CheckServiceAvailability - TX 9: Проверка доступности услуги доставки
func (c *WSPClientImpl) CheckServiceAvailability(ctx context.Context, req *postexpress.ServiceAvailabilityRequest) (*postexpress.ServiceAvailabilityResponse, error) {
	inputData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CheckServiceAvailability request: %w", err)
	}

	// Логируем JSON запроса для диагностики
	c.logger.Debug("CheckServiceAvailability request JSON: %s", string(inputData))

	txReq := &models.TransactionRequest{
		TransactionType: 9, // ProveraDostupnostiUsluge
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, txReq)
	if err != nil {
		return nil, fmt.Errorf("CheckServiceAvailability transaction failed: %w", err)
	}

	if !resp.Success {
		errMsg := errMsgUnknown
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return nil, fmt.Errorf("CheckServiceAvailability failed: %s", errMsg)
	}

	// Парсинг результатов
	var result postexpress.ServiceAvailabilityResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse service availability response: %w", err)
	}

	c.logger.Debug("CheckServiceAvailability completed - ServiceID: %d, Dostupna: %t, OcekivanoDana: %d, Rezultat: %d",
		req.IdRukovanje, result.Dostupna, result.OcekivanoDana, result.Rezultat)

	return &result, nil
}

// CalculatePostage - TX 11: Расчёт стоимости доставки
func (c *WSPClientImpl) CalculatePostage(ctx context.Context, req *postexpress.PostageCalculationRequest) (*postexpress.PostageCalculationResponse, error) {
	inputData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CalculatePostage request: %w", err)
	}

	txReq := &models.TransactionRequest{
		TransactionType: 11, // PostarinaPosiljke
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, txReq)
	if err != nil {
		return nil, fmt.Errorf("CalculatePostage transaction failed: %w", err)
	}

	// ВАЖНО: Не проверяем resp.Success здесь!
	// Вместо этого парсим response и возвращаем структуру с Rezultat и Poruka
	// Это позволяет handler получить детальное сообщение об ошибке от Post Express

	// Если resp.Success == false, ErrorMessage содержит детальную ошибку из StrRezultat
	if !resp.Success && resp.ErrorMessage != nil {
		// Возвращаем структуру response с ошибкой для handler
		result := &postexpress.PostageCalculationResponse{
			Rezultat: 3, // Код ошибки от Post Express
			Poruka:   *resp.ErrorMessage,
		}
		c.logger.Debug("CalculatePostage failed - Rezultat: %d, Poruka: %s", result.Rezultat, result.Poruka)
		return result, nil
	}

	// Парсинг результатов
	var result postexpress.PostageCalculationResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse postage calculation response: %w", err)
	}

	c.logger.Debug("CalculatePostage completed - ServiceID: %d, Weight: %dg, Postage: %d para (%.2f RSD), Rezultat: %d",
		req.IdRukovanje, req.Masa, result.Postarina, float64(result.Postarina)/100, result.Rezultat)

	return &result, nil
}
