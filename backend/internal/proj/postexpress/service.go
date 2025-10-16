package postexpress

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// Service - основной сервис для работы с Post Express API
type Service struct {
	client *Client
	config *Config
}

// NewService создает новый сервис Post Express
func NewService(config *Config) (*Service, error) {
	if config == nil {
		var err error
		config, err = LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	client := NewClient(config)

	log.Info().
		Str("api_url", config.APIURL).
		Str("username", config.Username).
		Str("brand", config.Brand).
		Bool("is_production", config.IsProduction).
		Msg("Post Express service initialized")

	return &Service{
		client: client,
		config: config,
	}, nil
}

// CreateManifest создает манифест с отправлениями
func (s *Service) CreateManifest(ctx context.Context, req *ManifestRequest) (*ManifestResponse, error) {
	log.Info().
		Str("manifest_id", req.ExtIDManifest).
		Int("orders_count", len(req.Porudzbine)).
		Msg("Creating Post Express manifest")

	var resp ManifestResponse
	if err := s.client.Post(ctx, "/manifest/create", req, &resp); err != nil {
		log.Error().
			Err(err).
			Str("manifest_id", req.ExtIDManifest).
			Msg("Failed to create manifest")
		return nil, err
	}

	log.Info().
		Str("manifest_id", req.ExtIDManifest).
		Int("internal_id", resp.IDManifesta).
		Int("result_code", resp.Rezultat).
		Msg("Manifest created successfully")

	return &resp, nil
}

// TrackShipments отслеживает отправления по трек-номерам
func (s *Service) TrackShipments(ctx context.Context, trackingNumbers []string) (*TrackingResponse, error) {
	log.Debug().
		Strs("tracking_numbers", trackingNumbers).
		Msg("Tracking Post Express shipments")

	req := &TrackingRequest{
		TrackingNumbers: trackingNumbers,
	}

	var resp TrackingResponse
	if err := s.client.Post(ctx, "/tracking/query", req, &resp); err != nil {
		log.Error().
			Err(err).
			Strs("tracking_numbers", trackingNumbers).
			Msg("Failed to track shipments")
		return nil, err
	}

	log.Debug().
		Int("shipments_count", len(resp.Posiljke)).
		Msg("Tracking info retrieved")

	return &resp, nil
}

// TrackShipment отслеживает одно отправление
func (s *Service) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	resp, err := s.TrackShipments(ctx, []string{trackingNumber})
	if err != nil {
		return nil, err
	}

	if len(resp.Posiljke) == 0 {
		return nil, fmt.Errorf("shipment not found: %s", trackingNumber)
	}

	return &resp.Posiljke[0], nil
}

// CancelShipments отменяет отправления
func (s *Service) CancelShipments(ctx context.Context, trackingNumbers []string, reason string) (*CancelResponse, error) {
	log.Info().
		Strs("tracking_numbers", trackingNumbers).
		Str("reason", reason).
		Msg("Canceling Post Express shipments")

	req := &CancelRequest{
		TrackingNumbers: trackingNumbers,
		Reason:          reason,
	}

	var resp CancelResponse
	if err := s.client.Post(ctx, "/shipment/cancel", req, &resp); err != nil {
		log.Error().
			Err(err).
			Strs("tracking_numbers", trackingNumbers).
			Msg("Failed to cancel shipments")
		return nil, err
	}

	log.Info().
		Int("canceled_count", len(resp.CanceledShipments)).
		Msg("Shipments canceled")

	return &resp, nil
}

// CancelShipment отменяет одно отправление
func (s *Service) CancelShipment(ctx context.Context, trackingNumber string, reason string) error {
	resp, err := s.CancelShipments(ctx, []string{trackingNumber}, reason)
	if err != nil {
		return err
	}

	if len(resp.CanceledShipments) == 0 {
		return fmt.Errorf("no shipments were canceled")
	}

	canceled := resp.CanceledShipments[0]
	if canceled.Rezultat != 0 {
		return fmt.Errorf("failed to cancel shipment: %s", canceled.Poruka)
	}

	return nil
}

// CalculateRate рассчитывает стоимость доставки
func (s *Service) CalculateRate(ctx context.Context, req *RateRequest) (*RateResponse, error) {
	log.Debug().
		Str("from", req.FromCity).
		Str("to", req.ToCity).
		Float64("weight", req.Weight).
		Msg("Calculating Post Express rate")

	var resp RateResponse
	if err := s.client.Post(ctx, "/rates/calculate", req, &resp); err != nil {
		log.Error().
			Err(err).
			Str("from", req.FromCity).
			Str("to", req.ToCity).
			Msg("Failed to calculate rate")
		return nil, err
	}

	log.Debug().
		Int("options_count", len(resp.DeliveryOptions)).
		Msg("Rate calculated")

	return &resp, nil
}

// GetOffices получает список офисов/отделений
func (s *Service) GetOffices(ctx context.Context, req *OfficeListRequest) (*OfficeListResponse, error) {
	log.Debug().
		Str("city", req.City).
		Str("postal_code", req.PostalCode).
		Msg("Getting Post Express offices")

	var resp OfficeListResponse
	if err := s.client.Post(ctx, "/offices/list", req, &resp); err != nil {
		log.Error().
			Err(err).
			Str("city", req.City).
			Msg("Failed to get offices")
		return nil, err
	}

	log.Debug().
		Int("offices_count", len(resp.Offices)).
		Msg("Offices retrieved")

	return &resp, nil
}

// CreateShipment - вспомогательный метод для создания одного отправления
// Обертка над CreateManifest для простоты использования (UPDATED для B2B API)
func (s *Service) CreateShipment(ctx context.Context, shipment *ShipmentRequest) (*ShipmentResponse, error) {
	// Валидация перед созданием
	if err := s.ValidateShipment(shipment); err != nil {
		return nil, fmt.Errorf("shipment validation failed: %w", err)
	}

	// Генерируем уникальный ID манифеста
	manifestID := fmt.Sprintf("SVETU-M-%d", time.Now().Unix())
	orderID := fmt.Sprintf("SVETU-O-%d", time.Now().Unix())

	manifest := &ManifestRequest{
		ExtIDManifest: manifestID,
		IDTipPosiljke: 1, // Обычная отправка
		Posiljalac:    shipment.Posiljalac, // Отправитель на уровне манифеста
		Porudzbine: []OrderRequest{
			{
				ExtIdPorudzbina: orderID,
				Posiljke:        []ShipmentRequest{*shipment},
			},
		},
		DatumPrijema: time.Now().Format("2006-01-02"),
		IDPartnera:   10109, // svetu.rs
	}

	resp, err := s.CreateManifest(ctx, manifest)
	if err != nil {
		return nil, err
	}

	if len(resp.Porudzbine) == 0 || len(resp.Porudzbine[0].Posiljke) == 0 {
		return nil, fmt.Errorf("no shipment in manifest response")
	}

	return &resp.Porudzbine[0].Posiljke[0], nil
}

// ValidateShipment проверяет корректность данных отправления перед отправкой (UPDATED для B2B)
func (s *Service) ValidateShipment(shipment *ShipmentRequest) error {
	// Обязательные B2B поля
	if shipment.ExtBrend == "" {
		return fmt.Errorf("ExtBrend is required")
	}
	if shipment.ExtMagacin == "" {
		return fmt.Errorf("ExtMagacin is required")
	}
	if shipment.ExtReferenca == "" {
		return fmt.Errorf("ExtReferenca is required")
	}
	if shipment.NacinPrijema == "" {
		return fmt.Errorf("NacinPrijema is required (K or O)")
	}
	if shipment.NacinPlacanja == "" {
		return fmt.Errorf("NacinPlacanja is required (POF, N, K)")
	}

	// Основные поля
	if shipment.BrojPosiljke == "" {
		return fmt.Errorf("shipment number (BrojPosiljke) is required")
	}
	if shipment.Masa <= 0 {
		return fmt.Errorf("weight (Masa) must be greater than 0 grams")
	}
	if shipment.IDRukovanje == 0 {
		return fmt.Errorf("service type (IdRukovanje) is required")
	}

	// Получатель
	if shipment.Primalac.Naziv == "" {
		return fmt.Errorf("recipient name is required")
	}
	if shipment.Primalac.Adresa == nil {
		return fmt.Errorf("recipient address object is required")
	}
	if shipment.Primalac.Mesto == "" {
		return fmt.Errorf("recipient city is required")
	}
	if shipment.Primalac.PostanskiBroj == "" {
		return fmt.Errorf("recipient postal code is required")
	}
	if shipment.Primalac.Telefon == "" {
		return fmt.Errorf("recipient phone is required")
	}

	// Отправитель (внутри отправления для B2B)
	if shipment.Posiljalac.Naziv == "" {
		return fmt.Errorf("sender name is required")
	}
	if shipment.Posiljalac.Adresa == nil {
		return fmt.Errorf("sender address object is required")
	}
	if shipment.Posiljalac.Mesto == "" {
		return fmt.Errorf("sender city is required")
	}
	if shipment.Posiljalac.Telefon == "" {
		return fmt.Errorf("sender phone is required")
	}

	// COD валидация
	if shipment.Otkupnina != nil && shipment.Otkupnina.Iznos > 0 {
		if shipment.Vrednost == 0 {
			return fmt.Errorf("Vrednost is required when Otkupnina is set")
		}
		// Проверка на услуги OTK и VD
		if shipment.PosebneUsluge == "" {
			return fmt.Errorf("PosebneUsluge must include OTK and VD for COD shipments")
		}
		// Проверка обязательных полей откупнины
		if shipment.Otkupnina.TekuciRacun == "" {
			return fmt.Errorf("TekuciRacun is required for COD shipments")
		}
		if shipment.Otkupnina.VrstaDokumenta == "" {
			return fmt.Errorf("VrstaDokumenta is required for COD shipments")
		}
	}

	return nil
}

// GetConfig возвращает текущую конфигурацию
func (s *Service) GetConfig() *Config {
	return s.config
}

// IsProduction проверяет, работает ли сервис в production окружении
func (s *Service) IsProduction() bool {
	return s.config.IsProduction
}
