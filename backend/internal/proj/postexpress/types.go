package postexpress

import "time"

// ManifestRequest - запрос создания манифеста (списка отправлений)
type ManifestRequest struct {
	ExtIDManifest string         `json:"ExtIdManifest"` // Уникальный ID манифеста
	IDTipPosiljke int            `json:"IdTipPosiljke"` // Тип пошљке: 1-обычная, 2-возврат
	Porudzbine    []OrderRequest `json:"Porudzbine"`    // Список заказов
}

// OrderRequest - данные заказа в манифесте
type OrderRequest struct {
	BrojPorudzbine string            `json:"BrojPorudzbine"` // Номер заказа
	Posiljke       []ShipmentRequest `json:"Posiljke"`       // Список отправлений в заказе
}

// ShipmentRequest - данные отправления
type ShipmentRequest struct {
	// Основные данные
	BrojPosiljke string  `json:"BrojPosiljke"`         // Номер отправления (уникальный)
	Tezina       float64 `json:"Tezina"`               // Вес в кг
	VrednostRSD  float64 `json:"VrednostRSD"`          // Объявленная ценность в RSD
	Otkupnina    float64 `json:"Otkupnina"`            // Наложенный платеж (COD) в RSD
	NacinPlacanj string  `json:"NacinPlacanjaDostave"` // Способ оплаты доставки

	// Получатель
	PrijemnoLice       string `json:"PrijemnoLice"`                // ФИО получателя
	PrijemnoLiceAdresa string `json:"PrijemnoLiceAdresa"`          // Адрес получателя
	PrijemnoLiceGrad   string `json:"PrijemnoLiceGrad"`            // Город получателя
	PrijemnoLicePosbr  string `json:"PrijemnoLicePosbr"`           // Почтовый индекс
	PrijemnoLiceTel    string `json:"PrijemnoLiceTel"`             // Телефон получателя
	PrijemnoLiceEmail  string `json:"PrijemnoLiceEmail,omitempty"` // Email получателя

	// Отправитель
	PosaljalacNaziv  string `json:"PosaljalacNaziv"`           // Название отправителя
	PosaljalacAdresa string `json:"PosaljalacAdresa"`          // Адрес отправителя
	PosaljalacGrad   string `json:"PosaljalacGrad"`            // Город отправителя
	PosaljalacPosbr  string `json:"PosaljalacPosbr"`           // Почтовый индекс отправителя
	PosaljalacTel    string `json:"PosaljalacTel"`             // Телефон отправителя
	PosaljalacEmail  string `json:"PosaljalacEmail,omitempty"` // Email отправителя

	// Возврат/забор
	MestoPreuzimanja string `json:"MestoPreuzimanja,omitempty"` // Место забора
	VracanjePosiljke string `json:"VracanjePosiljke,omitempty"` // Адрес возврата

	// Дополнительные услуги
	Usluge []ServiceRequest `json:"Usluge,omitempty"` // Доп. услуги (SMS и т.д.)

	// Опции
	GrupnaDostava bool   `json:"GrupnaDostava,omitempty"` // Групповая доставка
	Napomena      string `json:"Napomena,omitempty"`      // Примечание
}

// ServiceRequest - дополнительная услуга
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
