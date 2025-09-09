package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/proj/postexpress/models"
)

// WSPManifestRequest представляет запрос для транзакции 73 (B2B Manifest)
type WSPManifestRequest struct {
	// Информация о клиенте/отправителе
	Posiljalac WSPPosiljalac `json:"Posiljalac"`

	// Список посылок для создания
	Posiljke []WSPPosiljka `json:"Posiljke"`

	// Дата приема (формат: YYYY-MM-DD)
	DatumPrijema string `json:"DatumPrijema"`

	// Время приема (формат: HH:MM)
	VremePrijema string `json:"VremePrijema,omitempty"`

	// ID отделения приема
	IdPostePrijema int `json:"IdPostePrijema,omitempty"`

	// ДОБАВЛЕНО: Обязательные поля для B2B
	IdPartnera     int    `json:"IdPartnera,omitempty"`     // ID партнера (10109 для svetu.rs)
	NazivManifesta string `json:"NazivManifesta,omitempty"` // Название манифеста
}

// WSPPosiljalac представляет данные отправителя
type WSPPosiljalac struct {
	// Обязательные поля
	Naziv         string `json:"Naziv"`         // Название компании/имя
	Adresa        string `json:"Adresa"`        // Адрес
	Mesto         string `json:"Mesto"`         // Город
	PostanskiBroj string `json:"PostanskiBroj"` // Почтовый индекс
	Telefon       string `json:"Telefon"`       // Телефон

	// Опциональные поля
	Email        string `json:"Email,omitempty"`
	PIB          string `json:"PIB,omitempty"`          // ИНН для юр.лиц
	MaticniBroj  string `json:"MaticniBroj,omitempty"`  // Регистрационный номер
	Kontakt      string `json:"Kontakt,omitempty"`      // Контактное лицо
	IdUgovor     int    `json:"IdUgovor,omitempty"`     // ДОБАВЛЕНО: ID договора (если есть)
	OznakaZemlje string `json:"OznakaZemlje,omitempty"` // Код страны (RS для Сербии)
}

// WSPPosiljka представляет одну посылку в манифесте
type WSPPosiljka struct {
	// Обязательные поля
	BrojPosiljke  string      `json:"BrojPosiljke"`  // Уникальный номер посылки
	IdRukovanje   int         `json:"IdRukovanje"`   // ID услуги (29, 30, 55, 58, 59, 71)
	IdTipPosiljke int         `json:"IdTipPosiljke"` // 1-обычная, 2-возврат
	Primalac      WSPPrimalac `json:"Primalac"`      // Данные получателя
	Masa          float64     `json:"Masa"`          // Вес в кг

	// Габариты (опционально, но рекомендуется)
	Duzina float64 `json:"Duzina,omitempty"` // Длина в см
	Sirina float64 `json:"Sirina,omitempty"` // Ширина в см
	Visina float64 `json:"Visina,omitempty"` // Высота в см

	// Дополнительные услуги
	Otkupnina          *WSPOtkupnina `json:"Otkupnina,omitempty"`          // COD
	ObjavljenaVrednost float64       `json:"ObjavljenaVrednost,omitempty"` // Страховая стоимость
	SMS                bool          `json:"SMS,omitempty"`                // SMS уведомления
	Povratnica         bool          `json:"Povratnica,omitempty"`         // Уведомление о доставке (AR)
	LicnoUrucenje      bool          `json:"LicnoUrucenje,omitempty"`      // Личное вручение
	PDK                bool          `json:"PDK,omitempty"`                // Возврат документов
	VD                 bool          `json:"VD,omitempty"`                 // Ценная посылка (обязательно при ObjavljenaVrednost > 0)

	// Содержимое и примечания
	Sadrzaj       string `json:"Sadrzaj,omitempty"`       // Описание содержимого
	Napomena      string `json:"Napomena,omitempty"`      // Примечания для курьера
	ReferencaBroj string `json:"ReferencaBroj,omitempty"` // Референсный номер клиента
}

// WSPPrimalac представляет данные получателя
type WSPPrimalac struct {
	// Тип адреса: S-стандарт, F-Fah (ящик), P-Post restant
	TipAdrese string `json:"TipAdrese"`

	// Общие поля
	Naziv   string `json:"Naziv"`   // Имя/название
	Telefon string `json:"Telefon"` // Телефон
	Email   string `json:"Email,omitempty"`

	// Для стандартного адреса (TipAdrese = "S")
	Adresa        string `json:"Adresa,omitempty"`        // Улица и номер дома
	Mesto         string `json:"Mesto,omitempty"`         // Город
	PostanskiBroj string `json:"PostanskiBroj,omitempty"` // Почтовый индекс
	PAK           string `json:"PAK,omitempty"`           // Почтовый адресный код (6 цифр)
	OznakaZemlje  string `json:"OznakaZemlje,omitempty"`  // Код страны (RS для Сербии)

	// Для Fah адреса (TipAdrese = "F")
	Fah      string `json:"Fah,omitempty"`      // Название Fah
	BrojFaha string `json:"BrojFaha,omitempty"` // Номер ящика

	// Для Post restant (TipAdrese = "P")
	IdPoste int `json:"IdPoste,omitempty"` // ID почтового отделения
}

// WSPOtkupnina представляет данные для наложенного платежа (COD)
type WSPOtkupnina struct {
	Iznos         float64 `json:"Iznos"`                 // Сумма COD
	NacinPlacanja string  `json:"NacinPlacanja"`         // N-наличные, E-евро, U-уплатница
	BrojRacuna    string  `json:"BrojRacuna,omitempty"`  // Номер счета для перевода
	PozivNaBroj   string  `json:"PozivNaBroj,omitempty"` // Референсный номер платежа
}

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

// CreateManifest создает манифест с посылками через транзакцию 73
func (c *WSPClientImpl) CreateManifest(ctx context.Context, manifest *WSPManifestRequest) (*WSPManifestResponse, error) {
	// Валидация обязательных полей
	if len(manifest.Posiljke) == 0 {
		return nil, fmt.Errorf("manifest must contain at least one shipment")
	}

	// Установка даты приема если не указана
	if manifest.DatumPrijema == "" {
		manifest.DatumPrijema = time.Now().Format("2006-01-02")
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
		return &WSPManifestResponse{
			Success:      false,
			ErrorMessage: errMsg,
		}, nil
	}

	// Парсинг результата
	var result WSPManifestResponse
	if err := json.Unmarshal(resp.OutputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse manifest response: %w", err)
	}

	return &result, nil
}

// CreateShipmentViaManifest создает отправление через манифест (правильный способ)
func (c *WSPClientImpl) CreateShipmentViaManifest(ctx context.Context, shipment *WSPShipmentRequest) (*WSPManifestResponse, error) {
	// Преобразуем старый формат запроса в новый для манифеста

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
	}

	// Создаем посылку для манифеста
	posiljka := WSPPosiljka{
		BrojPosiljke:  fmt.Sprintf("SVT-%d", time.Now().Unix()),
		IdRukovanje:   idRukovanje,
		IdTipPosiljke: 1, // Обычная посылка
		Masa:          shipment.Weight,
		// ДОБАВЛЕНО: Стандартные габариты если не указаны
		Duzina: 30, // 30 см по умолчанию
		Sirina: 20, // 20 см по умолчанию
		Visina: 10, // 10 см по умолчанию
		Primalac: WSPPrimalac{
			TipAdrese:     "S", // Стандартный адрес
			Naziv:         shipment.RecipientName,
			Telefon:       shipment.RecipientPhone,
			Adresa:        shipment.RecipientAddress,
			Mesto:         shipment.RecipientCity,
			PostanskiBroj: shipment.RecipientPostalCode,
			OznakaZemlje:  "RS", // ДОБАВЛЕНО: код страны для Сербии
		},
		Sadrzaj:       shipment.Content,
		Napomena:      shipment.Note,
		ReferencaBroj: fmt.Sprintf("SVETU-%d", time.Now().Unix()), // ДОБАВЛЕНО: референсный номер
	}

	// Добавляем COD если указан
	if shipment.CODAmount > 0 {
		posiljka.Otkupnina = &WSPOtkupnina{
			Iznos:         shipment.CODAmount,
			NacinPlacanja: "N", // Наличные по умолчанию
		}
	}

	// Добавляем страховку если указана
	if shipment.InsuranceAmount > 0 {
		posiljka.ObjavljenaVrednost = shipment.InsuranceAmount
		// ДОБАВЛЕНО: для ценных посылок требуется услуга VD
		posiljka.VD = true
	}

	// Создаем манифест с одной посылкой
	manifest := &WSPManifestRequest{
		Posiljalac: WSPPosiljalac{
			Naziv:         shipment.SenderName,
			Adresa:        shipment.SenderAddress,
			Mesto:         shipment.SenderCity,
			PostanskiBroj: shipment.SenderPostalCode,
			Telefon:       shipment.SenderPhone,
			Email:         "b2b@svetu.rs", // ДОБАВЛЕНО: email для уведомлений
			OznakaZemlje:  "RS",           // ДОБАВЛЕНО: код страны для Сербии
		},
		Posiljke:       []WSPPosiljka{posiljka},
		DatumPrijema:   time.Now().Format("2006-01-02"),
		VremePrijema:   time.Now().Format("15:04"),
		IdPartnera:     10109,                                                         // ДОБАВЛЕНО: Partner ID для svetu.rs
		NazivManifesta: fmt.Sprintf("SVETU-%s", time.Now().Format("20060102-150405")), // ДОБАВЛЕНО: уникальное имя
	}

	// Создаем манифест
	return c.CreateManifest(ctx, manifest)
}
