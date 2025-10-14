package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/proj/postexpress/models"
	"backend/pkg/logger"
)

// Constants
const (
	errMsgUnknown = "unknown error"
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
	PartnerID       int // Добавлено: ID партнера (10109 для svetu.rs)
}

// NewWSPClient создает новый экземпляр WSP клиента
func NewWSPClient(config *WSPConfig, logger logger.Logger) WSPClient {
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

	return &WSPClientImpl{
		httpClient: httpClient,
		config:     config,
		logger:     logger,
	}
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

	// Декодирование ответа
	var wspResp map[string]interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&wspResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
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

	// Проверяем успешность
	success := true
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
		// Иногда данные приходят в поле StrOut
		outputData = json.RawMessage(fmt.Sprintf(`"%v"`, strOut))
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
		Barcode:         "",                       // Баркод не возвращается в B2B API
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

	// Выполнение запроса (тип транзакции 15 для отслеживания)
	req := &models.TransactionRequest{
		TransactionType: 15,
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

	// TODO: Декодировать base64 content в []byte
	// Пока возвращаем как есть
	return []byte(result.PDFContent), nil
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
