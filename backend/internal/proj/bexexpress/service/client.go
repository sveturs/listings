package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"backend/internal/proj/bexexpress/models"
)

// BEXClient представляет клиент для работы с BEX API
type BEXClient struct {
	authToken  string
	clientID   string
	baseURL    string
	httpClient *http.Client
}

// NewBEXClient создает новый клиент BEX
func NewBEXClient(authToken, clientID, baseURL string) *BEXClient {
	return &BEXClient{
		authToken: authToken,
		clientID:  clientID,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateShipment создает новое отправление в BEX
func (c *BEXClient) CreateShipment(shipment *models.BEXShipmentData) (*models.BEXShipmentResponse, error) {
	request := models.BEXShipmentRequest{
		ShipmentsList: []models.BEXShipmentData{*shipment},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/ship/api/Ship/postShipments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-AUTH-TOKEN", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response models.BEXShipmentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.RequestState {
		return nil, fmt.Errorf("request failed: %s", response.RequestError)
	}

	return &response, nil
}

// GetShipmentStatus получает статус отправления
func (c *BEXClient) GetShipmentStatus(shipmentID int, detailed bool) (*models.BEXStatusResponse, error) {
	params := url.Values{}
	params.Set("shipmentid", strconv.Itoa(shipmentID))
	if detailed {
		params.Set("mtype", "2")
	} else {
		params.Set("mtype", "1")
	}
	params.Set("lang", "1") // Serbian

	fullURL := fmt.Sprintf("%s/ship/api/Ship/getstate?%s", c.baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-AUTH-TOKEN", c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response models.BEXStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// If it's not JSON, parse the text response
		return &models.BEXStatusResponse{
			StatusText: string(body),
		}, nil
	}

	return &response, nil
}

// GetShipmentLabel получает этикетку для печати
func (c *BEXClient) GetShipmentLabel(shipmentID int, pageSize int, pagePosition int) ([]byte, error) {
	params := url.Values{}
	params.Set("shipmentId", strconv.Itoa(shipmentID))
	params.Set("pageSize", strconv.Itoa(pageSize))
	params.Set("pagePosition", strconv.Itoa(pagePosition))
	params.Set("parcelNo", "0")

	fullURL := fmt.Sprintf("%s/ShipDNF/Ship/getLabelWithProperties?%s", c.baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-AUTH-TOKEN", c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response models.BEXLabelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.State {
		return nil, fmt.Errorf("failed to get label: %s", response.Error)
	}

	// Decode base64 label
	labelData, err := base64.StdEncoding.DecodeString(response.ParcelLabel)
	if err != nil {
		return nil, fmt.Errorf("failed to decode label: %w", err)
	}

	return labelData, nil
}

// DeleteShipment удаляет отправление
func (c *BEXClient) DeleteShipment(shipmentID int) error {
	requestData := map[string]int{
		"shipmentId": shipmentID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/ship/api/Ship/delete", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-AUTH-TOKEN", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListParcelShops получает список пунктов выдачи
func (c *BEXClient) ListParcelShops(city string) ([]models.BEXParcelShop, error) {
	requestData := map[string]string{}
	if city != "" {
		requestData["city"] = city
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/ship/api/Ship/listParcelShops", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-AUTH-TOKEN", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var parcelShops []models.BEXParcelShop
	if err := json.Unmarshal(body, &parcelShops); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return parcelShops, nil
}

// Helper функции для создания отправления

// BuildShipmentData создает структуру данных отправления для API BEX
func BuildShipmentData(req *models.CreateShipmentRequest, senderClientID string) *models.BEXShipmentData {
	now := time.Now()

	// Pickup task (sender)
	pickupTask := models.BEXShipmentTask{
		Type:            1, // Pickup
		NameType:        3, // Client ID
		Name1:           senderClientID,
		Name2:           "",
		TaxID:           "",
		AddressType:     3,   // Using IDs
		Municipalities:  104, // Novi Sad
		Place:           "0",
		Street:          "0",
		HouseNumber:     18,
		Apartment:       "",
		ContactPerson:   "",
		Phone:           "",
		Date:            now.Format("2006-01-02"),
		TimeFrom:        "07:55",
		TimeTo:          "14:55",
		PreNotification: 0,
		Comment:         "",
		ParcelShop:      0,
	}

	// Delivery task (recipient)
	deliveryTask := models.BEXShipmentTask{
		Type:            2, // Delivery
		NameType:        1, // Person
		Name1:           req.RecipientName,
		Name2:           "",
		TaxID:           "",
		AddressType:     1, // Text address
		Municipalities:  0,
		Place:           req.RecipientCity,
		Street:          req.RecipientAddress,
		HouseNumber:     0,
		Apartment:       "",
		ContactPerson:   req.RecipientName,
		Phone:           req.RecipientPhone,
		Date:            "",
		TimeFrom:        "",
		TimeTo:          "",
		PreNotification: req.PreNotificationMinutes,
		Comment:         "",
		ParcelShop:      0,
	}

	if req.DeliveryInstructions != nil {
		deliveryTask.Comment = *req.DeliveryInstructions
	}

	// Build shipment data
	shipmentData := &models.BEXShipmentData{
		ShipmentID:               0, // 0 for new shipment
		ServiceSpeed:             1, // Always 1
		ShipmentType:             1, // Standard
		ShipmentCategory:         req.ShipmentCategory,
		ShipmentWeight:           req.WeightKg,
		TotalPackages:            req.TotalPackages,
		InvoiceAmount:            0,
		ShipmentContents:         req.ShipmentContents,
		CommentPublic:            "",
		CommentPrivate:           "",
		PersonalDelivery:         req.PersonalDelivery,
		ReturnSignedInvoices:     false,
		ReturnSignedConfirmation: false,
		ReturnPackage:            false,
		PayType:                  6, // Bank transfer by default
		InsuranceAmount:          0,
		PayToSender:              0,
		PayToSenderViaAccount:    false,
		SendersAccountNumber:     "",
		BankTransferComment:      "",
		Tasks: []models.BEXShipmentTask{
			pickupTask,
			deliveryTask,
		},
		Reports: []interface{}{},
	}

	// Set COD if provided
	if req.CODAmount != nil && *req.CODAmount > 0 {
		shipmentData.PayToSender = *req.CODAmount
		shipmentData.PayType = 2 // Receiver pays cash
	}

	// Set insurance if provided
	if req.InsuranceAmount != nil && *req.InsuranceAmount > 0 {
		shipmentData.InsuranceAmount = *req.InsuranceAmount
	}

	// Set comments if provided
	if req.Notes != nil {
		shipmentData.CommentPublic = *req.Notes
	}

	// Set order IDs as private comment
	if req.OrderID != nil {
		shipmentData.CommentPrivate = fmt.Sprintf("Order #%d", *req.OrderID)
	} else if req.StorefrontOrderID != nil {
		shipmentData.CommentPrivate = fmt.Sprintf("Storefront Order #%d", *req.StorefrontOrderID)
	}

	return shipmentData
}

// CalculateShipmentCategory определяет категорию отправления по весу
func CalculateShipmentCategory(weightKg float64) int {
	switch {
	case weightKg <= 0.5:
		return 1 // Documents up to 0.5kg
	case weightKg <= 1.0:
		return 2 // Package up to 1kg
	case weightKg <= 2.0:
		return 3 // Package up to 2kg
	default:
		return 31 // Package per kg
	}
}
