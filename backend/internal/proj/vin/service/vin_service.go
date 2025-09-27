package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/proj/vin/models"
	"backend/internal/proj/vin/storage/postgres"
)

// VINService предоставляет бизнес-логику для работы с VIN
type VINService struct {
	storage *postgres.VINStorage
	decoder *VINDecoder
}

// NewVINService создает новый сервис VIN
func NewVINService(storage *postgres.VINStorage) *VINService {
	return &VINService{
		storage: storage,
		decoder: NewVINDecoder(),
	}
}

// DecodeVIN декодирует VIN номер с кэшированием
func (s *VINService) DecodeVIN(ctx context.Context, req *models.VINDecodeRequest, userID *int64) (*models.VINDecodeResponse, error) {
	// Валидация VIN
	if err := s.decoder.ValidateVIN(req.VIN); err != nil {
		return &models.VINDecodeResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Проверяем кэш
	cached, err := s.storage.GetCachedVIN(ctx, req.VIN)
	if err == nil && cached != nil {
		// Записываем в историю
		history := &models.VINCheckHistory{
			UserID:        userID,
			VIN:           req.VIN,
			ListingID:     req.ListingID,
			DecodeSuccess: true,
			DecodeCacheID: &cached.ID,
			CheckType:     "manual",
		}

		if err := s.storage.SaveCheckHistory(ctx, history); err != nil {
			// Не критично, просто логируем
			fmt.Printf("Failed to save check history: %v\n", err)
		}

		return &models.VINDecodeResponse{
			Success: true,
			Data:    cached,
			Source:  "cache",
		}, nil
	}

	// Декодируем локально
	decoded, err := s.decoder.DecodeLocal(req.VIN)
	if err != nil {
		// Сохраняем неудачную попытку
		history := &models.VINCheckHistory{
			UserID:        userID,
			VIN:           req.VIN,
			ListingID:     req.ListingID,
			DecodeSuccess: false,
			CheckType:     "manual",
		}
		s.storage.SaveCheckHistory(ctx, history)

		return &models.VINDecodeResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Сохраняем в кэш
	decoded.CreatedAt = time.Now()
	decoded.UpdatedAt = time.Now()

	cacheID, err := s.storage.SaveToCache(ctx, decoded)
	if err != nil {
		// Не критично, возвращаем результат без кэша
		fmt.Printf("Failed to save to cache: %v\n", err)
	} else {
		decoded.ID = cacheID
	}

	// Записываем в историю
	history := &models.VINCheckHistory{
		UserID:        userID,
		VIN:           req.VIN,
		ListingID:     req.ListingID,
		DecodeSuccess: true,
		DecodeCacheID: &cacheID,
		CheckType:     "manual",
	}

	if err := s.storage.SaveCheckHistory(ctx, history); err != nil {
		fmt.Printf("Failed to save check history: %v\n", err)
	}

	return &models.VINDecodeResponse{
		Success: true,
		Data:    decoded,
		Source:  "api",
	}, nil
}

// GetHistory получает историю проверок VIN
func (s *VINService) GetHistory(ctx context.Context, req *models.VINHistoryRequest) ([]*models.VINCheckHistory, error) {
	return s.storage.GetHistory(ctx, req)
}

// AutoFillFromVIN заполняет данные объявления из VIN
func (s *VINService) AutoFillFromVIN(ctx context.Context, vin string, userID *int64) (map[string]interface{}, error) {
	// Валидация VIN
	if err := s.decoder.ValidateVIN(vin); err != nil {
		return nil, err
	}

	// Декодируем VIN
	req := &models.VINDecodeRequest{VIN: vin}
	resp, err := s.DecodeVIN(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("не удалось декодировать VIN")
	}

	// Преобразуем в формат для автозаполнения
	autoFillData := make(map[string]interface{})

	if resp.Data.Make != nil {
		autoFillData["make"] = *resp.Data.Make
	}
	if resp.Data.Model != nil {
		autoFillData["model"] = *resp.Data.Model
	}
	if resp.Data.Year != nil {
		autoFillData["year"] = *resp.Data.Year
	}
	if resp.Data.EngineType != nil {
		autoFillData["engine_type"] = *resp.Data.EngineType
	}
	if resp.Data.EngineDisplacement != nil {
		autoFillData["engine_displacement"] = *resp.Data.EngineDisplacement
	}
	if resp.Data.TransmissionType != nil {
		autoFillData["transmission"] = *resp.Data.TransmissionType
	}
	if resp.Data.Drivetrain != nil {
		autoFillData["drivetrain"] = *resp.Data.Drivetrain
	}
	if resp.Data.BodyType != nil {
		autoFillData["body_type"] = *resp.Data.BodyType
	}
	if resp.Data.FuelType != nil {
		autoFillData["fuel_type"] = *resp.Data.FuelType
	}
	if resp.Data.Doors != nil {
		autoFillData["doors"] = *resp.Data.Doors
	}
	if resp.Data.Seats != nil {
		autoFillData["seats"] = *resp.Data.Seats
	}

	// Добавляем VIN в результат
	autoFillData["vin"] = vin
	autoFillData["vin_verified"] = true

	return autoFillData, nil
}

// GetVINStats получает статистику использования VIN декодера
func (s *VINService) GetVINStats(ctx context.Context, userID *int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Общее количество проверок
	totalChecks, err := s.storage.GetTotalChecks(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats["total_checks"] = totalChecks

	// Количество уникальных VIN
	uniqueVINs, err := s.storage.GetUniqueVINCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats["unique_vins"] = uniqueVINs

	// Последние проверки
	recentHistory := &models.VINHistoryRequest{
		UserID: userID,
		Limit:  5,
	}
	recent, err := s.storage.GetHistory(ctx, recentHistory)
	if err != nil {
		return nil, err
	}
	stats["recent_checks"] = recent

	// Статистика по производителям (если есть данные)
	manufacturerStats, err := s.storage.GetManufacturerStats(ctx, userID)
	if err == nil {
		stats["manufacturer_stats"] = manufacturerStats
	}

	return stats, nil
}

// ValidateAndEnrichListing проверяет и обогащает данные объявления с помощью VIN
func (s *VINService) ValidateAndEnrichListing(ctx context.Context, listingData map[string]interface{}, userID *int64) (map[string]interface{}, error) {
	// Проверяем наличие VIN в данных объявления
	vinInterface, ok := listingData["vin"]
	if !ok {
		return listingData, nil // VIN не указан, возвращаем как есть
	}

	vin, ok := vinInterface.(string)
	if !ok || vin == "" {
		return listingData, nil
	}

	// Декодируем VIN
	req := &models.VINDecodeRequest{VIN: vin}
	resp, err := s.DecodeVIN(ctx, req, userID)
	if err != nil || !resp.Success || resp.Data == nil {
		// Не критично, возвращаем данные как есть
		return listingData, nil
	}

	// Обогащаем данные объявления
	enriched := make(map[string]interface{})
	for k, v := range listingData {
		enriched[k] = v
	}

	// Добавляем или проверяем данные из VIN
	vinData := resp.Data

	// Функция для безопасного добавления данных
	addIfMissing := func(key string, vinValue interface{}) {
		if vinValue != nil {
			if existingValue, exists := enriched[key]; !exists || existingValue == nil || existingValue == "" {
				// Добавляем значение из VIN если оно отсутствует
				enriched[key] = vinValue
			} else {
				// Добавляем для валидации
				enriched[key+"_vin_validated"] = vinValue
			}
		}
	}

	// Добавляем/проверяем данные
	if vinData.Make != nil {
		addIfMissing("make", *vinData.Make)
	}
	if vinData.Model != nil {
		addIfMissing("model", *vinData.Model)
	}
	if vinData.Year != nil {
		addIfMissing("year", *vinData.Year)
	}
	if vinData.BodyType != nil {
		addIfMissing("body_type", *vinData.BodyType)
	}
	if vinData.FuelType != nil {
		addIfMissing("fuel_type", *vinData.FuelType)
	}
	if vinData.TransmissionType != nil {
		addIfMissing("transmission", *vinData.TransmissionType)
	}
	if vinData.EngineDisplacement != nil {
		addIfMissing("engine_displacement", *vinData.EngineDisplacement)
	}

	// Добавляем флаг проверки VIN
	enriched["vin_verified"] = true
	enriched["vin_decode_status"] = vinData.DecodeStatus

	// Добавляем метаданные о проверке
	vinMetadata := map[string]interface{}{
		"checked_at": time.Now().Format(time.RFC3339),
		"status":     vinData.DecodeStatus,
	}

	if existingMeta, ok := enriched["vin_metadata"]; ok {
		if metaMap, ok := existingMeta.(map[string]interface{}); ok {
			for k, v := range vinMetadata {
				metaMap[k] = v
			}
			enriched["vin_metadata"] = metaMap
		} else {
			enriched["vin_metadata"] = vinMetadata
		}
	} else {
		enriched["vin_metadata"] = vinMetadata
	}

	return enriched, nil
}

// ExportVINReport экспортирует отчет о VIN в JSON формате
func (s *VINService) ExportVINReport(ctx context.Context, vin string, userID *int64) ([]byte, error) {
	// Валидация VIN
	if err := s.decoder.ValidateVIN(vin); err != nil {
		return nil, err
	}

	// Получаем данные из кэша или декодируем
	req := &models.VINDecodeRequest{VIN: vin}
	resp, err := s.DecodeVIN(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("не удалось получить данные о VIN")
	}

	// Формируем отчет
	report := map[string]interface{}{
		"vin":         vin,
		"report_date": time.Now().Format(time.RFC3339),
		"vehicle":     resp.Data,
		"verification_status": map[string]interface{}{
			"vin_valid":     true,
			"decode_status": resp.Data.DecodeStatus,
			"source":        resp.Source,
		},
	}

	// Добавляем историю проверок для этого VIN
	historyReq := &models.VINHistoryRequest{
		VIN:   vin,
		Limit: 10,
	}
	history, err := s.storage.GetHistory(ctx, historyReq)
	if err == nil && len(history) > 0 {
		report["check_history"] = history
	}

	// Преобразуем в JSON
	return json.MarshalIndent(report, "", "  ")
}
