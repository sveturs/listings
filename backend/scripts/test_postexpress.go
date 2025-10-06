package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/postexpress"
)

func main() {
	// Настройка логирования
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Загрузка .env файла
	if err := godotenv.Load("../.env"); err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file, using environment variables")
	}

	fmt.Println("=================================================================")
	fmt.Println("Post Express API Integration Test")
	fmt.Println("=================================================================\n")

	// Создание сервиса
	service, err := postexpress.NewService(nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Post Express service")
	}

	config := service.GetConfig()
	fmt.Printf("API URL: %s\n", config.APIURL)
	fmt.Printf("Username: %s\n", config.Username)
	fmt.Printf("Brand: %s\n", config.Brand)
	fmt.Printf("Production: %v\n\n", config.IsProduction)

	ctx := context.Background()

	// Test 1: Получение списка офисов
	fmt.Println("Test 1: Получение списка офисов в Белграде")
	fmt.Println("-----------------------------------------------------------------")
	if err := testGetOffices(ctx, service); err != nil {
		log.Error().Err(err).Msg("Office list test failed")
	}
	fmt.Println()

	// Test 2: Расчет стоимости доставки
	fmt.Println("Test 2: Расчет стоимости доставки (Белград → Нови Сад)")
	fmt.Println("-----------------------------------------------------------------")
	if err := testCalculateRate(ctx, service); err != nil {
		log.Error().Err(err).Msg("Rate calculation test failed")
	}
	fmt.Println()

	// Test 3: Создание тестового отправления
	fmt.Println("Test 3: Создание тестового отправления")
	fmt.Println("-----------------------------------------------------------------")
	shipmentResp, err := testCreateShipment(ctx, service)
	if err != nil {
		log.Error().Err(err).Msg("Shipment creation test failed")
	}
	fmt.Println()

	// Test 4: Отслеживание отправления (если создали)
	if shipmentResp != nil && shipmentResp.TrackingNumber != "" {
		fmt.Println("Test 4: Отслеживание созданного отправления")
		fmt.Println("-----------------------------------------------------------------")
		if err := testTrackShipment(ctx, service, shipmentResp.TrackingNumber); err != nil {
			log.Error().Err(err).Msg("Tracking test failed")
		}
		fmt.Println()
	}

	fmt.Println("=================================================================")
	fmt.Println("Tests completed!")
	fmt.Println("=================================================================")
}

func testGetOffices(ctx context.Context, service *postexpress.Service) error {
	req := &postexpress.OfficeListRequest{
		City: "Београд", // Belgrade
	}

	resp, err := service.GetOffices(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get offices: %w", err)
	}

	fmt.Printf("✓ Found %d offices in Belgrade\n", len(resp.Offices))
	if len(resp.Offices) > 0 {
		fmt.Printf("  First office: %s - %s\n", resp.Offices[0].Name, resp.Offices[0].Address)
	}

	return nil
}

func testCalculateRate(ctx context.Context, service *postexpress.Service) error {
	req := &postexpress.RateRequest{
		FromCity:  "Београд",
		ToCity:    "Нови Сад",
		Weight:    2.5,  // 2.5 kg
		Value:     5000, // 5000 RSD
		CODAmount: 0,    // No COD
		Services:  []string{"tracking"},
	}

	resp, err := service.CalculateRate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to calculate rate: %w", err)
	}

	fmt.Printf("✓ Rate calculated successfully\n")
	fmt.Printf("  Available delivery options: %d\n", len(resp.DeliveryOptions))
	for i, option := range resp.DeliveryOptions {
		fmt.Printf("  %d. %s - %.2f RSD (estimated: %d days)\n",
			i+1, option.Name, option.TotalPrice, option.EstimatedDays)
	}

	return nil
}

func testCreateShipment(ctx context.Context, service *postexpress.Service) (*postexpress.ShipmentResponse, error) {
	// Генерируем уникальный номер отправления
	shipmentNumber := fmt.Sprintf("SVETU-TEST-%d", time.Now().Unix())

	req := &postexpress.ShipmentRequest{
		BrojPosiljke: shipmentNumber,
		Tezina:       2.5,
		VrednostRSD:  5000,
		Otkupnina:    0,
		NacinPlacanj: postexpress.PaymentCash,

		// Получатель (тестовые данные)
		PrijemnoLice:       "Test Recipient",
		PrijemnoLiceAdresa: "Bulevar Kralja Aleksandra 73",
		PrijemnoLiceGrad:   "Београд",
		PrijemnoLicePosbr:  "11000",
		PrijemnoLiceTel:    "+381111234567",
		PrijemnoLiceEmail:  "test@example.com",

		// Отправитель (используем данные из конфига)
		PosaljalacNaziv:  service.GetConfig().Brand,
		PosaljalacAdresa: "Test Street 1",
		PosaljalacGrad:   "Београд",
		PosaljalacPosbr:  "11000",
		PosaljalacTel:    "+381111111111",
		PosaljalacEmail:  service.GetConfig().Username,

		// Добавляем SMS уведомление
		Usluge: []postexpress.ServiceRequest{
			{
				SifraUsluge: postexpress.ServiceSMS,
				Parametri:   "+381111234567",
			},
		},
	}

	// Валидация перед отправкой
	if err := service.ValidateShipment(req); err != nil {
		return nil, fmt.Errorf("shipment validation failed: %w", err)
	}
	fmt.Println("✓ Shipment data validated")

	// Создание отправления
	resp, err := service.CreateShipment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	fmt.Printf("✓ Shipment created successfully!\n")
	fmt.Printf("  Shipment ID: %d\n", resp.IDPosiljke)
	fmt.Printf("  Tracking Number: %s\n", resp.TrackingNumber)
	fmt.Printf("  External ID: %s\n", resp.BrojPosiljke)
	fmt.Printf("  Status: %s\n", resp.Status)
	if resp.LabelURL != "" {
		fmt.Printf("  Label URL: %s\n", resp.LabelURL)
	}

	// Сохраняем tracking number в файл для последующего использования
	if err := os.WriteFile("/tmp/postexpress_tracking.txt", []byte(resp.TrackingNumber), 0644); err != nil {
		log.Warn().Err(err).Msg("Failed to save tracking number")
	}

	return resp, nil
}

func testTrackShipment(ctx context.Context, service *postexpress.Service, trackingNumber string) error {
	tracking, err := service.TrackShipment(ctx, trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to track shipment: %w", err)
	}

	fmt.Printf("✓ Tracking info retrieved\n")
	fmt.Printf("  Tracking Number: %s\n", tracking.TrackingNumber)
	fmt.Printf("  Status: %s - %s\n", tracking.Status, tracking.StatusText)
	fmt.Printf("  Current Location: %s\n", tracking.CurrentLocation)

	if tracking.EstimatedDate != nil {
		fmt.Printf("  Estimated Delivery: %s\n", tracking.EstimatedDate.Format("2006-01-02"))
	}

	if len(tracking.Events) > 0 {
		fmt.Printf("  Events: %d\n", len(tracking.Events))
		for i, event := range tracking.Events {
			fmt.Printf("    %d. [%s] %s - %s\n",
				i+1,
				event.Timestamp.Format("2006-01-02 15:04"),
				event.Status,
				event.Description)
		}
	}

	if tracking.ProofOfDelivery != nil {
		fmt.Printf("\n  Proof of Delivery:\n")
		fmt.Printf("    Recipient: %s\n", tracking.ProofOfDelivery.RecipientName)
		fmt.Printf("    Delivered At: %s\n", tracking.ProofOfDelivery.DeliveredAt.Format("2006-01-02 15:04"))
		if tracking.ProofOfDelivery.SignatureURL != "" {
			fmt.Printf("    Signature URL: %s\n", tracking.ProofOfDelivery.SignatureURL)
		}
	}

	// Выводим полный JSON для анализа
	jsonData, _ := json.MarshalIndent(tracking, "  ", "  ")
	fmt.Printf("\n  Full tracking data (JSON):\n  %s\n", string(jsonData))

	return nil
}
