package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/proj/postexpress"
	"backend/internal/proj/postexpress/models"
)

// NOTE: Эти структуры больше НЕ используются! См. types.go для правильных B2B структур.
// Оставлены для обратной совместимости со старым кодом, но будут удалены в следующей версии.

// WSPManifestResponse представляет ответ на создание манифеста
type WSPManifestResponse struct {
	Success      bool                `json:"Rezultat"`
	ErrorMessage string              `json:"Greska,omitempty"`
	Manifest     *WSPManifestInfo    `json:"Manifest,omitempty"`
	Posiljke     []WSPPosiljkaResult `json:"Posiljke,omitempty"`
}

// WSPManifestInfo содержит информацию о созданном манифесте
type WSPManifestInfo struct {
	BrojManifesta  string    `json:"BrojManifesta"`
	DatumKreiranja time.Time `json:"DatumKreiranja"`
	Status         string    `json:"Status"`
}

// WSPPosiljkaResult представляет результат создания посылки
type WSPPosiljkaResult struct {
	BrojPosiljke    string `json:"BrojPosiljke"`     // Наш номер
	PostExpressBroj string `json:"PostExpressBroj"`  // Номер Post Express
	Barkod          string `json:"Barkod"`           // Штрих-код
	Status          string `json:"Status"`           // OK или ERROR
	Greska          string `json:"Greska,omitempty"` // Описание ошибки
}

// CreateManifest создает манифест с посылками через транзакцию 73 (B2B API)
func (c *WSPClientImpl) CreateManifest(ctx context.Context, manifest *postexpress.ManifestRequest) (*postexpress.ManifestResponse, error) {
	// Валидация обязательных полей
	if len(manifest.Porudzbine) == 0 {
		return nil, fmt.Errorf("manifest must contain at least one order")
	}
	if manifest.ExtIDManifest == "" {
		return nil, fmt.Errorf("ExtIdManifest is required")
	}

	// Установка даты приема если не указана
	if manifest.DatumPrijema == "" {
		manifest.DatumPrijema = time.Now().Format("2006-01-02")
	}

	// Установка ID партнера если не указан
	if manifest.IDPartnera == 0 {
		manifest.IDPartnera = 10109 // svetu.rs partner ID
	}

	// Маршалинг данных манифеста
	inputData, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest request: %w", err)
	}

	// Выполнение транзакции 73 (B2B Manifest)
	req := &models.TransactionRequest{
		TransactionType: 73, // ID транзакции для B2B Manifest
		InputData:       string(inputData),
	}

	resp, err := c.Transaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("manifest transaction failed: %w", err)
	}

	// Проверка успешности транзакции
	if !resp.Success {
		errMsg := "manifest creation failed"
		if resp.ErrorMessage != nil {
			errMsg = *resp.ErrorMessage
		}
		return &postexpress.ManifestResponse{
			Rezultat: 1,
			Poruka:   errMsg,
		}, nil
	}

	// Парсинг результата
	var result postexpress.ManifestResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse manifest response: %w", err)
	}

	return &result, nil
}

// CreateShipmentViaManifest создает отправление через манифест B2B API (правильная структура)
func (c *WSPClientImpl) CreateShipmentViaManifest(ctx context.Context, shipment *WSPShipmentRequest) (*postexpress.ManifestResponse, error) {
	timestamp := time.Now().Unix()
	boolFalse := false // Helper для *bool полей

	// Определяем ID услуги на основе типа сервиса
	idRukovanje := 29 // По умолчанию: PE_Danas_za_sutra_12
	switch shipment.ServiceType {
	case "PE_Danas_za_danas":
		idRukovanje = 30
	case "PE_Danas_za_odmah":
		idRukovanje = 55
	case "PE_Danas_za_sutra_19":
		idRukovanje = 58
	case "PE_Danas_za_odmah_Bg":
		idRukovanje = 59
	case "PE_Danas_za_sutra_isporuka":
		idRukovanje = 71
	case "PE_Klasicna":
		idRukovanje = 85
	}

	// Парсим адрес отправителя (улица и номер)
	senderStreet, senderNumber := parseAddress(shipment.SenderAddress)
	recipientStreet, recipientNumber := parseAddress(shipment.RecipientAddress)

	// Конвертируем вес из kg в граммы
	weightGrams := int(shipment.Weight * 1000)
	if weightGrams < 1 {
		weightGrams = 500 // минимум 500 грамм
	}

	// Формируем услуги
	services := "PNA" // Pickup notification - для курьера всегда нужна
	if shipment.CODAmount > 0 {
		services += ",OTK,VD" // Cash on Delivery + Valuable delivery
	}

	// Конвертируем COD в para (1 RSD = 100 para)
	codPara := int(shipment.CODAmount * 100)
	valuePara := int(shipment.InsuranceAmount * 100)
	if codPara > 0 && valuePara == 0 {
		valuePara = codPara // Vrednost должна быть минимум как Otkupnina
	}

	// Создаем правильную B2B структуру отправления
	posiljka := postexpress.ShipmentRequest{
		// Обязательные B2B поля
		ExtBrend:          "SVETU",
		ExtMagacin:        "WAREHOUSE1",
		ExtReferenca:      fmt.Sprintf("SVETU-REF-%d", timestamp),
		NacinPrijema:      "K", // K=courier, O=office
		ImaPrijemniBrojDN: &boolFalse,
		NacinPlacanja:     "POF", // POF=poslato od firme (sent by company)

		// Отправитель ВНУТРИ отправления
		Posiljalac: postexpress.SenderInfo{
			Naziv: shipment.SenderName,
			Adresa: &postexpress.AddressInfo{
				Ulica:         senderStreet,
				Broj:          senderNumber,
				Mesto:         shipment.SenderCity,
				PostanskiBroj: shipment.SenderPostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         shipment.SenderCity,
			PostanskiBroj: shipment.SenderPostalCode,
			Telefon:       shipment.SenderPhone,
			Email:         "b2b@svetu.rs",
			OznakaZemlje:  "RS",
		},

		// Место забора (для курьера - обязательно!)
		MestoPreuzimanja: &postexpress.SenderInfo{
			Naziv: shipment.SenderName,
			Adresa: &postexpress.AddressInfo{
				Ulica:         senderStreet,
				Broj:          senderNumber,
				Mesto:         shipment.SenderCity,
				PostanskiBroj: shipment.SenderPostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         shipment.SenderCity,
			PostanskiBroj: shipment.SenderPostalCode,
			Telefon:       shipment.SenderPhone,
			OznakaZemlje:  "RS",
		},

		// Основные данные
		BrojPosiljke: fmt.Sprintf("SVETU-%d", timestamp),
		IDRukovanje:  idRukovanje,
		Primalac: postexpress.ReceiverInfo{
			Naziv: shipment.RecipientName,
			Adresa: &postexpress.AddressInfo{
				Ulica:         recipientStreet,
				Broj:          recipientNumber,
				Mesto:         shipment.RecipientCity,
				PostanskiBroj: shipment.RecipientPostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         shipment.RecipientCity,
			PostanskiBroj: shipment.RecipientPostalCode,
			Telefon:       shipment.RecipientPhone,
			OznakaZemlje:  "RS",
		},
		Masa: weightGrams, // В граммах!

		// COD и ценность (в para!)
		Otkupnina: codPara,
		Vrednost:  valuePara,

		// Услуги (строка через запятую!)
		PosebneUsluge: services,

		// Опциональные поля
		Sadrzaj:       shipment.Content,
		ReferencaBroj: fmt.Sprintf("SVETU-%d", timestamp),
		Napomena:      shipment.Note,
	}

	// Создаем заказ с одной посылкой
	order := postexpress.OrderRequest{
		ExtIdPorudzbina: fmt.Sprintf("ORDER-%d", timestamp),
		Posiljke:        []postexpress.ShipmentRequest{posiljka},
	}

	// Создаем манифест с правильной структурой
	manifest := &postexpress.ManifestRequest{
		ExtIDManifest: fmt.Sprintf("MANIFEST-%d", timestamp),
		IDTipPosiljke: 1, // На уровне манифеста: 1=standard, 2=return
		Posiljalac: postexpress.SenderInfo{
			Naziv: shipment.SenderName,
			Adresa: &postexpress.AddressInfo{
				Ulica:         senderStreet,
				Broj:          senderNumber,
				Mesto:         shipment.SenderCity,
				PostanskiBroj: shipment.SenderPostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         shipment.SenderCity,
			PostanskiBroj: shipment.SenderPostalCode,
			Telefon:       shipment.SenderPhone,
			Email:         "b2b@svetu.rs",
			OznakaZemlje:  "RS",
		},
		Porudzbine:     []postexpress.OrderRequest{order},
		DatumPrijema:   time.Now().Format("2006-01-02"),
		VremePrijema:   time.Now().Format("15:04"),
		IDPartnera:     10109, // svetu.rs partner ID
		NazivManifesta: fmt.Sprintf("SVETU-%s", time.Now().Format("20060102-150405")),
	}

	// Создаем манифест через B2B API
	return c.CreateManifest(ctx, manifest)
}

// parseAddress разбирает адрес на улицу и номер дома
// Например: "Takovska 2" → ("Takovska", "2")
func parseAddress(fullAddress string) (street string, number string) {
	// Простое разбиение по последнему пробелу
	if fullAddress == "" {
		return "", ""
	}

	// Пробуем найти последний пробел
	lastSpace := -1
	for i := len(fullAddress) - 1; i >= 0; i-- {
		if fullAddress[i] == ' ' {
			lastSpace = i
			break
		}
	}

	if lastSpace == -1 {
		// Нет пробелов - весь адрес это улица
		return fullAddress, ""
	}

	street = fullAddress[:lastSpace]
	number = fullAddress[lastSpace+1:]
	return street, number
}
