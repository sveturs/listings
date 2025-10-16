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

// OtkupninaData - структура откупной пошилки (COD) с полными банковскими данными
type OtkupninaData struct {
	Iznos          int    `json:"Iznos"`          // Сумма в para (1 RSD = 100 para)
	VrstaDokumenta string `json:"VrstaDokumenta"` // Тип документа: N=налогни документ
	TekuciRacun    string `json:"TekuciRacun"`    // Банковский счёт для перевода откупнины
	ModelPNB       string `json:"ModelPNB"`       // Модель платежа (обычно "97")
	PNB            string `json:"PNB"`            // Позив на број (payment reference)
	SifraPlacanja  string `json:"SifraPlacanja"`  // Шифра плаћања (обычно "189")
}

// ShipmentRequest - данные отправления (B2B структура)
type ShipmentRequest struct {
	// === ОБЯЗАТЕЛЬНЫЕ B2B поля ===
	ExtBrend          string      `json:"ExtBrend"`            // ОБЯЗАТЕЛЬНО: Бренд (напр. "SVETU")
	ExtMagacin        string      `json:"ExtMagacin"`          // ОБЯЗАТЕЛЬНО: Склад (напр. "WAREHOUSE1")
	ExtReferenca      string      `json:"ExtReferenca"`        // ОБЯЗАТЕЛЬНО: Уникальная референция
	NacinPrijema      string      `json:"NacinPrijema"`        // ОБЯЗАТЕЛЬНО: K-курьер, O-офис
	ImaPrijemniBrojDN string      `json:"ImaPrijemniBrojDN"`   // ОБЯЗАТЕЛЬНО: "D" если есть PrijemniBroj, "N" если нет
	NacinPlacanja     string      `json:"NacinPlacanja"`       // ОБЯЗАТЕЛЬНО: POF, N, K
	Posiljalac        SenderInfo  `json:"Posiljalac"`                  // ОБЯЗАТЕЛЬНО: Отправитель в отправлении
	MestoPreuzimanja  *SenderInfo `json:"MestoPreuzimanja,omitempty"`  // Место забора (объект!)

	// === Основные данные ===
	BrojPosiljke string       `json:"BrojPosiljke"` // ОБЯЗАТЕЛЬНО: Номер отправления
	IDRukovanje  int          `json:"IdRukovanje"`  // ОБЯЗАТЕЛЬНО: ID услуги (29, 30, 55, 58, 59, 71, 85)
	Primalac     ReceiverInfo `json:"Primalac"`     // ОБЯЗАТЕЛЬНО: Получатель
	Masa         int          `json:"Masa"`         // ОБЯЗАТЕЛЬНО: Вес в ГРАММАХ (integer!)

	// === COD и ценности (в PARA - 1 RSD = 100 para) ===
	Otkupnina *OtkupninaData `json:"Otkupnina,omitempty"` // COD с полными банковскими данными (структура!)
	Vrednost  int            `json:"Vrednost,omitempty"`  // Объявленная ценность в para (ОБЯЗАТЕЛЬНО для COD!)

	// === Дополнительные услуги (строка через запятую!) ===
	PosebneUsluge string `json:"PosebneUsluge,omitempty"` // "PNA,OTK,VD" - НЕ массив!

	// === Опциональные поля ===
	Sadrzaj          string `json:"Sadrzaj,omitempty"`          // Описание содержимого
	ReferencaBroj    string `json:"ReferencaBroj,omitempty"`    // Референсный номер
	Napomena         string `json:"Napomena,omitempty"`         // Примечание
	ParcelLockerCode string `json:"ParcelLockerCode,omitempty"` // Код паккетомата (для IdRukovanje=85)
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
	PrijemniBroj   string `json:"PrijemniBroj"`       // Трек-номер (реальный от Post Express)
	TrackingNumber string `json:"TrackingNumber"`     // Трек-номер (алиас для совместимости)
	Status         string `json:"Status"`             // Статус отправления
	Rezultat       int    `json:"Rezultat"`           // 0 - успех
	Poruka         string `json:"Poruka,omitempty"`   // Сообщение об ошибке
	LabelURL       string `json:"LabelUrl,omitempty"` // URL этикетки для печати
	Postarina      int    `json:"Postarina,omitempty"` // Стоимость доставки в para (1 RSD = 100 para)
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

// ========================================
// TX 3 - GetNaselje (Поиск населённых пунктов)
// ========================================

// GetSettlementsRequest - запрос поиска населённых пунктов (TX 3)
type GetSettlementsRequest struct {
	Naziv string `json:"Naziv"` // Название или часть названия населённого пункта
}

// Settlement - данные населённого пункта
type Settlement struct {
	Id             int    `json:"Id"`             // ID населённого пункта (ВАЖНО: API возвращает "Id", не "IdNaselje"!)
	IdNaselje      int    `json:"IdNaselje"`      // Алиас для совместимости
	Naziv          string `json:"Naziv"`          // Название
	PostanskiBroj  string `json:"PostanskiBroj"`  // Почтовый индекс
	IdOkrug        int    `json:"IdOkrug"`        // ID округа
	NazivOkruga    string `json:"NazivOkruga"`    // Название округа
}

// GetSettlementsResponse - ответ поиска населённых пунктов (TX 3)
type GetSettlementsResponse struct {
	Rezultat int          `json:"Rezultat"`       // 0 - успех, 1 - ошибка
	Poruka   string       `json:"Poruka,omitempty"` // Сообщение об ошибке
	Naselja  []Settlement `json:"Naselja"`        // Массив найденных населённых пунктов
}

// ========================================
// TX 4 - GetUlica (Поиск улиц)
// ========================================

// GetStreetsRequest - запрос поиска улиц в населённом пункте (TX 4)
type GetStreetsRequest struct {
	IdNaselje int    `json:"IdNaselje"` // ID населённого пункта (из TX 3)
	Naziv     string `json:"Naziv"`     // Название или часть названия улицы
}

// Street - данные улицы
type Street struct {
	Id        int    `json:"Id"`        // ID улицы (ВАЖНО: API возвращает "Id", не "IdUlica"!)
	IdUlica   int    `json:"IdUlica"`   // Алиас для совместимости
	Naziv     string `json:"Naziv"`     // Название улицы
	IdNaselje int    `json:"IdNaselje"` // ID населённого пункта
}

// GetStreetsResponse - ответ поиска улиц (TX 4)
type GetStreetsResponse struct {
	Rezultat int      `json:"Rezultat"`       // 0 - успех, 1 - ошибка
	Poruka   string   `json:"Poruka,omitempty"` // Сообщение об ошибке
	Ulice    []Street `json:"Ulice"`          // Массив найденных улиц
}

// ========================================
// TX 6 - ProveraAdrese (Валидация адреса)
// ========================================

// AddressValidationRequest - запрос валидации адреса (TX 6)
type AddressValidationRequest struct {
	TipAdrese      int    `json:"TipAdrese"`      // Тип адреса (0, 1, 2)
	IdRukovanje    int    `json:"IdRukovanje"`    // ID услуги доставки
	IdNaselje      int    `json:"IdNaselje"`      // ID населённого пункта
	IdUlica        int    `json:"IdUlica,omitempty"` // ID улицы (опционально)
	BrojPodbroj    string `json:"BrojPodbroj"`    // Номер дома (например, "2" или "2a")
	PostanskiBroj  string `json:"PostanskiBroj"`  // Почтовый индекс
	Datum          string `json:"Datum,omitempty"` // Дата (опционально, формат: YYYY-MM-DD)
}

// AddressValidationResponse - ответ валидации адреса (TX 6)
type AddressValidationResponse struct {
	Rezultat      int    `json:"Rezultat"`       // 0 - успех, 1 - ошибка
	Poruka        string `json:"Poruka,omitempty"` // Сообщение об ошибке
	PostojiAdresa bool   `json:"PostojiAdresa"`  // Существует ли адрес
	IdNaselje     int    `json:"IdNaselje"`      // ID населённого пункта
	NazivNaselja  string `json:"NazivNaselja"`   // Название населённого пункта
	IdUlica       int    `json:"IdUlica"`        // ID улицы
	NazivUlice    string `json:"NazivUlice"`     // Название улицы
	Broj          string `json:"Broj"`           // Номер дома
	PostanskiBroj string `json:"PostanskiBroj"`  // Почтовый индекс
	PAK           string `json:"PAK"`            // Postal Address Code
	IdPoste       int    `json:"IdPoste"`        // ID почтового отделения
	NazivPoste    string `json:"NazivPoste"`     // Название почтового отделения
}

// ========================================
// TX 9 - ProveraDostupnostiUsluge (Проверка доступности услуги)
// ========================================

// ServiceAvailabilityRequest - запрос проверки доступности услуги (TX 9)
type ServiceAvailabilityRequest struct {
	TipAdrese              int    `json:"TipAdrese"`              // Тип адреса (0, 1, 2)
	IdRukovanje            int    `json:"IdRukovanje"`            // ID услуги (29, 30, 55, 58, 59, 71, 85)
	IdNaseljeOdlaska       int    `json:"IdNaseljeOdlaska,omitempty"`       // ID населённого пункта отправления
	IdNaseljeDolaska       int    `json:"IdNaseljeDolaska,omitempty"`       // ID населённого пункта прибытия
	PostanskiBrojOdlaska   string `json:"PostanskiBrojOdlaska"`   // Почтовый индекс отправления
	PostanskiBrojDolaska   string `json:"PostanskiBrojDolaska"`   // Почтовый индекс прибытия
	Datum                  string `json:"Datum,omitempty"`        // Дата (опционально, формат: YYYY-MM-DD)
}

// ServiceAvailabilityResponse - ответ проверки доступности услуги (TX 9)
type ServiceAvailabilityResponse struct {
	Rezultat      int    `json:"Rezultat"`       // 0 - успех, 1 - ошибка
	Poruka        string `json:"Poruka,omitempty"` // Сообщение об ошибке
	Dostupna      bool   `json:"Dostupna"`       // Доступна ли услуга
	IdRukovanje   int    `json:"IdRukovanje"`    // ID услуги
	NazivUsluge   string `json:"NazivUsluge"`    // Название услуги
	OcekivanoDana int    `json:"OcekivanoDana"`  // Ожидаемое время доставки (дни)
	Napomena      string `json:"Napomena,omitempty"` // Примечание
}

// ========================================
// TX 11 - PostarinaPosiljke (Расчёт стоимости доставки)
// ========================================

// PostageCalculationRequest - запрос расчёта стоимости доставки (TX 11)
type PostageCalculationRequest struct {
	IdRukovanje            int    `json:"IdRukovanje"`            // ID услуги доставки
	IdZemlja               int    `json:"IdZemlja"`               // ID страны (0 = внутренние отправления)
	PostanskiBrojOdlaska   string `json:"PostanskiBrojOdlaska"`   // Почтовый индекс отправления
	PostanskiBrojDolaska   string `json:"PostanskiBrojDolaska"`   // Почтовый индекс прибытия
	Masa                   int    `json:"Masa"`                   // Вес в граммах
	Otkupnina              int    `json:"Otkupnina,omitempty"`    // COD в para (1 RSD = 100 para)
	Vrednost               int    `json:"Vrednost,omitempty"`     // Объявленная ценность в para
	PosebneUsluge          string `json:"PosebneUsluge,omitempty"` // Дополнительные услуги через запятую
}

// PostageCalculationResponse - ответ расчёта стоимости доставки (TX 11)
type PostageCalculationResponse struct {
	Rezultat             int    `json:"Rezultat"`       // 0 - успех, 1 - ошибка
	Poruka               string `json:"Poruka,omitempty"` // Сообщение об ошибке
	Postarina            int    `json:"Postarina"`      // Стоимость в para (295 RSD = 29500 para)
	IdRukovanje          int    `json:"IdRukovanje"`    // ID услуги
	NazivUsluge          string `json:"NazivUsluge"`    // Название услуги
	PostanskiBrojOdlaska string `json:"PostanskiBrojOdlaska"` // Почтовый индекс отправления
	PostanskiBrojDolaska string `json:"PostanskiBrojDolaska"` // Почтовый индекс прибытия
	Masa                 int    `json:"Masa"`           // Вес в граммах
	Otkupnina            int    `json:"Otkupnina"`      // COD в para
	Vrednost             int    `json:"Vrednost"`       // Объявленная ценность в para
	PosebneUsluge        string `json:"PosebneUsluge"`  // Дополнительные услуги
	Napomena             string `json:"Napomena,omitempty"` // Примечание
}
