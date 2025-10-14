package postexpress

import "time"

// ========================================
// B2B Manifest API Structures (Transaction 73)
// ========================================

// ManifestRequest - запрос создания манифеста B2B (правильная структура)
type ManifestRequest struct {
	ExtIDManifest  string         `json:"ExtIdManifest"`            // ОБЯЗАТЕЛЬНО: Уникальный ID манифеста
	IDTipPosiljke  int            `json:"IdTipPosiljke"`            // ОБЯЗАТЕЛЬНО: Тип пошљке на уровне манифеста: 1-обычная, 2-возврат
	Posiljalac     SenderInfo     `json:"Posiljalac"`               // ОБЯЗАТЕЛЬНО: Отправитель (объект)
	Porudzbine     []OrderRequest `json:"Porudzbine"`               // ОБЯЗАТЕЛЬНО: Список заказов
	DatumPrijema   string         `json:"DatumPrijema"`             // ОБЯЗАТЕЛЬНО: Дата приема (YYYY-MM-DD)
	VremePrijema   string         `json:"VremePrijema,omitempty"`   // Время приема (HH:MM)
	IDPartnera     int            `json:"IdPartnera,omitempty"`     // ID партнера (10109 для svetu.rs)
	NazivManifesta string         `json:"NazivManifesta,omitempty"` // Название манифеста
}

// OrderRequest - данные заказа в манифесте (nested structure)
type OrderRequest struct {
	ExtIdPorudzbina      string            `json:"ExtIdPorudzbina,omitempty"`      // Внешний ID заказа
	ExtIdPorudzbinaKupca string            `json:"ExtIdPorudzbinaKupca,omitempty"` // ID заказа покупателя
	IndGrupnostUrucenja  *bool             `json:"IndGrupnostUrucenja,omitempty"`  // Индикатор групповой доставки
	Posiljke             []ShipmentRequest `json:"Posiljke"`                       // ОБЯЗАТЕЛЬНО: Список отправлений в заказе
}

// ShipmentRequest - данные отправления (B2B структура)
type ShipmentRequest struct {
	// === ОБЯЗАТЕЛЬНЫЕ B2B поля ===
	ExtBrend          string      `json:"ExtBrend"`                    // ОБЯЗАТЕЛЬНО: Бренд (напр. "SVETU")
	ExtMagacin        string      `json:"ExtMagacin"`                  // ОБЯЗАТЕЛЬНО: Склад (напр. "WAREHOUSE1")
	ExtReferenca      string      `json:"ExtReferenca"`                // ОБЯЗАТЕЛЬНО: Уникальная референция
	NacinPrijema      string      `json:"NacinPrijema"`                // ОБЯЗАТЕЛЬНО: K-курьер, O-офис
	ImaPrijemniBrojDN *bool       `json:"ImaPrijemniBrojDN,omitempty"` // ОБЯЗАТЕЛЬНО: false (указатель!)
	NacinPlacanja     string      `json:"NacinPlacanja"`               // ОБЯЗАТЕЛЬНО: POF, N, K
	Posiljalac        SenderInfo  `json:"Posiljalac"`                  // ОБЯЗАТЕЛЬНО: Отправитель в отправлении
	MestoPreuzimanja  *SenderInfo `json:"MestoPreuzimanja,omitempty"`  // Место забора (объект!)

	// === Основные данные ===
	BrojPosiljke string       `json:"BrojPosiljke"` // ОБЯЗАТЕЛЬНО: Номер отправления
	IDRukovanje  int          `json:"IdRukovanje"`  // ОБЯЗАТЕЛЬНО: ID услуги (29, 30, 55, 58, 59, 71, 85)
	Primalac     ReceiverInfo `json:"Primalac"`     // ОБЯЗАТЕЛЬНО: Получатель
	Masa         int          `json:"Masa"`         // ОБЯЗАТЕЛЬНО: Вес в ГРАММАХ (integer!)

	// === COD и ценности (в PARA - 1 RSD = 100 para) ===
	Otkupnina int `json:"Otkupnina,omitempty"` // COD в para (5000 RSD = 500000)
	Vrednost  int `json:"Vrednost,omitempty"`  // Объявленная ценность в para (ОБЯЗАТЕЛЬНО для COD!)

	// === Дополнительные услуги (строка через запятую!) ===
	PosebneUsluge string `json:"PosebneUsluge,omitempty"` // "PNA,OTK,VD" - НЕ массив!

	// === Опциональные поля ===
	Sadrzaj       string `json:"Sadrzaj,omitempty"`       // Описание содержимого
	ReferencaBroj string `json:"ReferencaBroj,omitempty"` // Референсный номер
	Napomena      string `json:"Napomena,omitempty"`      // Примечание
}

// AddressInfo - адрес (ВСЕГДА объект, НЕ строка!)
type AddressInfo struct {
	Ulica         string `json:"Ulica,omitempty"`         // Улица
	Broj          string `json:"Broj,omitempty"`          // Номер дома
	Mesto         string `json:"Mesto,omitempty"`         // Город
	PostanskiBroj string `json:"PostanskiBroj,omitempty"` // Почтовый индекс
	PAK           string `json:"PAK,omitempty"`           // Postal Address Code
	OznakaZemlje  string `json:"OznakaZemlje,omitempty"`  // Код страны (RS)
}

// SenderInfo - информация об отправителе/клиенте (Adresa - объект!)
type SenderInfo struct {
	Naziv         string       `json:"Naziv"`                   // ОБЯЗАТЕЛЬНО: Название/имя
	Adresa        *AddressInfo `json:"Adresa"`                  // ОБЯЗАТЕЛЬНО: Адрес (объект!)
	Mesto         string       `json:"Mesto"`                   // ОБЯЗАТЕЛЬНО: Город
	PostanskiBroj string       `json:"PostanskiBroj"`           // ОБЯЗАТЕЛЬНО: Почтовый индекс
	Telefon       string       `json:"Telefon"`                 // ОБЯЗАТЕЛЬНО: Телефон
	Email         string       `json:"Email,omitempty"`         // Email
	OznakaZemlje  string       `json:"OznakaZemlje,omitempty"`  // Код страны
	PIB           string       `json:"PIB,omitempty"`           // ИНН
	MaticniBroj   string       `json:"MaticniBroj,omitempty"`   // Регистрационный номер
	IDUgovor      int          `json:"IdUgovor,omitempty"`      // ID договора
	SifraKlijenta string       `json:"SifraKlijenta,omitempty"` // Код клиента
}

// ReceiverInfo - информация о получателе (Adresa - объект!)
type ReceiverInfo struct {
	Naziv         string       `json:"Naziv"`                   // ОБЯЗАТЕЛЬНО: Имя/название
	Adresa        *AddressInfo `json:"Adresa"`                  // ОБЯЗАТЕЛЬНО: Адрес (объект!)
	Mesto         string       `json:"Mesto"`                   // ОБЯЗАТЕЛЬНО: Город
	PostanskiBroj string       `json:"PostanskiBroj"`           // ОБЯЗАТЕЛЬНО: Почтовый индекс
	Telefon       string       `json:"Telefon"`                 // ОБЯЗАТЕЛЬНО: Телефон
	Email         string       `json:"Email,omitempty"`         // Email
	OznakaZemlje  string       `json:"OznakaZemlje,omitempty"`  // Код страны
	TipAdrese     string       `json:"TipAdrese,omitempty"`     // S-standard, F-fah, P-post restant
	PAK           string       `json:"PAK,omitempty"`           // Postal Address Code
}

// ServiceRequest - дополнительная услуга (DEPRECATED - используйте PosebneUsluge string)
type ServiceRequest struct {
	SifraUsluge string      `json:"SifraUsluge"`         // Код услуги (например, "SMS")
	Parametri   interface{} `json:"Parametri,omitempty"` // Параметры услуги
}

// ManifestResponse - ответ на создание манифеста
type ManifestResponse struct {
	Rezultat       int               `json:"Rezultat"`                   // 0 - успех, иначе ошибка
	Poruka         string            `json:"Poruka"`                     // Сообщение об ошибке
	IDManifesta    int               `json:"IdManifesta"`                // ID созданного манифеста
	ExtIDManifest  string            `json:"ExtIdManifest"`              // Внешний ID манифеста
	Porudzbine     []OrderResponse   `json:"Porudzbine"`                 // Обработанные заказы
	GreskeValidaci []ValidationError `json:"GreskeValidacije,omitempty"` // Ошибки валидации
}

// OrderResponse - ответ по заказу
type OrderResponse struct {
	BrojPorudzbine string             `json:"BrojPorudzbine"`
	Posiljke       []ShipmentResponse `json:"Posiljke"`
}

// ShipmentResponse - ответ по отправлению
type ShipmentResponse struct {
	BrojPosiljke   string `json:"BrojPosiljke"`       // Номер отправления
	IDPosiljke     int    `json:"IdPosiljke"`         // ID отправления в системе Pošta
	TrackingNumber string `json:"TrackingNumber"`     // Трек-номер
	Status         string `json:"Status"`             // Статус отправления
	Rezultat       int    `json:"Rezultat"`           // 0 - успех
	Poruka         string `json:"Poruka,omitempty"`   // Сообщение об ошибке
	LabelURL       string `json:"LabelUrl,omitempty"` // URL этикетки для печати
}

// ValidationError - ошибка валидации
type ValidationError struct {
	Polje    string `json:"Polje"`    // Поле с ошибкой
	Vrednost string `json:"Vrednost"` // Значение поля
	Poruka   string `json:"Poruka"`   // Сообщение об ошибке
}

// TrackingRequest - запрос отслеживания
type TrackingRequest struct {
	TrackingNumbers []string `json:"TrackingNumbers"` // Список трек-номеров для отслеживания
}

// TrackingResponse - ответ отслеживания
type TrackingResponse struct {
	Rezultat int            `json:"Rezultat"`
	Poruka   string         `json:"Poruka,omitempty"`
	Posiljke []TrackingInfo `json:"Posiljke"`
}

// TrackingInfo - информация об отслеживании отправления
type TrackingInfo struct {
	TrackingNumber  string           `json:"TrackingNumber"`
	Status          string           `json:"Status"`     // Текущий статус
	StatusText      string           `json:"StatusText"` // Описание статуса
	CurrentLocation string           `json:"CurrentLocation,omitempty"`
	EstimatedDate   *time.Time       `json:"EstimatedDate,omitempty"`
	DeliveredDate   *time.Time       `json:"DeliveredDate,omitempty"`
	Events          []TrackingEvent  `json:"Events"`
	ProofOfDelivery *ProofOfDelivery `json:"ProofOfDelivery,omitempty"`
}

// TrackingEvent - событие отслеживания
type TrackingEvent struct {
	Timestamp   time.Time `json:"Timestamp"`
	Status      string    `json:"Status"`
	Description string    `json:"Description"`
	Location    string    `json:"Location,omitempty"`
	Details     string    `json:"Details,omitempty"`
}

// ProofOfDelivery - подтверждение доставки
type ProofOfDelivery struct {
	RecipientName string    `json:"RecipientName"`
	SignatureURL  string    `json:"SignatureUrl,omitempty"`
	PhotoURL      string    `json:"PhotoUrl,omitempty"`
	DeliveredAt   time.Time `json:"DeliveredAt"`
	Notes         string    `json:"Notes,omitempty"`
}

// CancelRequest - запрос отмены отправления
type CancelRequest struct {
	TrackingNumbers []string `json:"TrackingNumbers"` // Трек-номера для отмены
	Reason          string   `json:"Reason"`          // Причина отмены
}

// CancelResponse - ответ на отмену
type CancelResponse struct {
	Rezultat          int                `json:"Rezultat"`
	Poruka            string             `json:"Poruka,omitempty"`
	CanceledShipments []CanceledShipment `json:"CanceledShipments"`
}

// CanceledShipment - отмененное отправление
type CanceledShipment struct {
	TrackingNumber string `json:"TrackingNumber"`
	Status         string `json:"Status"`           // "canceled" или ошибка
	Rezultat       int    `json:"Rezultat"`         // 0 - успех
	Poruka         string `json:"Poruka,omitempty"` // Сообщение
}

// RateRequest - запрос расчета тарифа
type RateRequest struct {
	FromCity  string   `json:"FromCity"`            // Город отправления
	ToCity    string   `json:"ToCity"`              // Город получения
	Weight    float64  `json:"Weight"`              // Вес в кг
	Value     float64  `json:"Value"`               // Объявленная ценность
	CODAmount float64  `json:"CodAmount,omitempty"` // Наложенный платеж
	Services  []string `json:"Services,omitempty"`  // Дополнительные услуги
}

// RateResponse - ответ с тарифом
type RateResponse struct {
	Rezultat        int              `json:"Rezultat"`
	Poruka          string           `json:"Poruka,omitempty"`
	DeliveryOptions []DeliveryOption `json:"DeliveryOptions"`
}

// DeliveryOption - вариант доставки
type DeliveryOption struct {
	Type          string  `json:"Type"`          // standard, express
	Name          string  `json:"Name"`          // Название
	BasePrice     float64 `json:"BasePrice"`     // Базовая стоимость
	CODFee        float64 `json:"CodFee"`        // Комиссия за COD
	InsuranceFee  float64 `json:"InsuranceFee"`  // Страхование
	FuelSurcharge float64 `json:"FuelSurcharge"` // Топливная надбавка
	TotalPrice    float64 `json:"TotalPrice"`    // Итого
	EstimatedDays int     `json:"EstimatedDays"` // Оценочное время доставки (дни)
	Currency      string  `json:"Currency"`      // Валюта (RSD)
}

// OfficeListRequest - запрос списка офисов
type OfficeListRequest struct {
	City       string `json:"City,omitempty"`       // Фильтр по городу
	PostalCode string `json:"PostalCode,omitempty"` // Фильтр по индексу
}

// OfficeListResponse - ответ со списком офисов
type OfficeListResponse struct {
	Rezultat int      `json:"Rezultat"`
	Poruka   string   `json:"Poruka,omitempty"`
	Offices  []Office `json:"Offices"`
}

// Office - информация об офисе/отделении
type Office struct {
	Code         string  `json:"Code"`                   // Код отделения
	Name         string  `json:"Name"`                   // Название
	Address      string  `json:"Address"`                // Адрес
	City         string  `json:"City"`                   // Город
	PostalCode   string  `json:"PostalCode"`             // Почтовый индекс
	Phone        string  `json:"Phone,omitempty"`        // Телефон
	WorkingHours string  `json:"WorkingHours,omitempty"` // Часы работы
	Latitude     float64 `json:"Latitude,omitempty"`     // Широта
	Longitude    float64 `json:"Longitude,omitempty"`    // Долгота
}

// Константы для статусов отправлений
const (
	StatusCreated           = "created"
	StatusPickupScheduled   = "pickup_scheduled"
	StatusPickedUp          = "picked_up"
	StatusInTransit         = "in_transit"
	StatusArrived           = "arrived"
	StatusOutForDelivery    = "out_for_delivery"
	StatusDelivered         = "delivered"
	StatusDeliveryAttempted = "delivery_attempted"
	StatusReturning         = "returning"
	StatusReturned          = "returned"
	StatusCancelled         = "canceled"
	StatusLost              = "lost"
	StatusDamaged           = "damaged"
)

// Константы для способов оплаты
const (
	PaymentCash    = "cash"    // Наличные
	PaymentCard    = "card"    // Карта
	PaymentAccount = "account" // На счет
)

// Константы для типов доставки
const (
	DeliveryTypeStandard = "standard" // Обычная доставка
	DeliveryTypeExpress  = "express"  // Экспресс доставка
)

// Константы для дополнительных услуг
const (
	ServiceSMS       = "SMS"        // SMS уведомление
	ServiceInsurance = "INSURANCE"  // Дополнительное страхование
	ServiceCOD       = "COD"        // Наложенный платеж
	ServiceReturnDoc = "RETURN_DOC" // Возврат документов
)
