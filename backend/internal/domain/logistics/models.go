package logistics

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// LogisticsMetrics представляет агрегированные метрики логистики
type LogisticsMetrics struct {
	ID   int       `json:"id" db:"id"`
	Date time.Time `json:"date" db:"date"`

	// Общие метрики
	TotalShipments int `json:"total_shipments" db:"total_shipments"`
	Delivered      int `json:"delivered" db:"delivered"`
	InTransit      int `json:"in_transit" db:"in_transit"`
	Problems       int `json:"problems" db:"problems"`
	Returns        int `json:"returns" db:"returns"`
	Canceled       int `json:"canceled" db:"canceled"`

	// Метрики по курьерским службам
	BexShipments         int `json:"bex_shipments" db:"bex_shipments"`
	BexDelivered         int `json:"bex_delivered" db:"bex_delivered"`
	PostExpressShipments int `json:"postexpress_shipments" db:"postexpress_shipments"`
	PostExpressDelivered int `json:"postexpress_delivered" db:"postexpress_delivered"`

	// Временные метрики
	AvgDeliveryTimeHours *float64 `json:"avg_delivery_time_hours" db:"avg_delivery_time_hours"`
	MinDeliveryTimeHours *float64 `json:"min_delivery_time_hours" db:"min_delivery_time_hours"`
	MaxDeliveryTimeHours *float64 `json:"max_delivery_time_hours" db:"max_delivery_time_hours"`

	// Финансовые метрики
	TotalShippingCost float64 `json:"total_shipping_cost" db:"total_shipping_cost"`
	TotalCODCollected float64 `json:"total_cod_collected" db:"total_cod_collected"`
	TotalReturnCost   float64 `json:"total_return_cost" db:"total_return_cost"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProblemType тип проблемы с отправлением
type ProblemType string

const (
	ProblemTypeDelayed         ProblemType = "delayed"
	ProblemTypeLost            ProblemType = "lost"
	ProblemTypeDamaged         ProblemType = "damaged"
	ProblemTypeReturnRequested ProblemType = "return_requested"
	ProblemTypeWrongAddress    ProblemType = "wrong_address"
	ProblemTypeComplaint       ProblemType = "complaint"
)

// ProblemSeverity уровень критичности проблемы
type ProblemSeverity string

const (
	SeverityLow      ProblemSeverity = "low"
	SeverityMedium   ProblemSeverity = "medium"
	SeverityHigh     ProblemSeverity = "high"
	SeverityCritical ProblemSeverity = "critical"
)

// ProblemStatus статус решения проблемы
type ProblemStatus string

const (
	StatusOpen            ProblemStatus = "open"
	StatusInvestigating   ProblemStatus = "investigating"
	StatusWaitingResponse ProblemStatus = "waiting_response"
	StatusResolved        ProblemStatus = "resolved"
	StatusClosed          ProblemStatus = "closed"
)

// ProblemShipment представляет проблемное отправление
type ProblemShipment struct {
	ID             int             `json:"id" db:"id"`
	ShipmentID     int             `json:"shipment_id" db:"shipment_id"`
	ShipmentType   string          `json:"shipment_type" db:"shipment_type"`
	TrackingNumber *string         `json:"tracking_number" db:"tracking_number"`
	ProblemType    ProblemType     `json:"problem_type" db:"problem_type"`
	Severity       ProblemSeverity `json:"severity" db:"severity"`
	Description    *string         `json:"description" db:"description"`
	Status         ProblemStatus   `json:"status" db:"status"`
	AssignedTo     *int            `json:"assigned_to" db:"assigned_to"`
	AssignedUser   *User           `json:"assigned_user,omitempty"`
	Resolution     *string         `json:"resolution" db:"resolution"`
	OrderID        *int            `json:"order_id" db:"order_id"`
	UserID         *int            `json:"user_id" db:"user_id"`
	User           *User           `json:"user,omitempty"`
	Metadata       JSONB           `json:"metadata" db:"metadata"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
	ResolvedAt     *time.Time      `json:"resolved_at" db:"resolved_at"`

	// Дополнительные поля
	Comments []ProblemComment       `json:"comments,omitempty"`
	History  []ProblemStatusHistory `json:"history,omitempty"`
}

// ProblemComment представляет комментарий к проблеме
type ProblemComment struct {
	ID          int       `json:"id" db:"id"`
	ProblemID   int       `json:"problem_id" db:"problem_id"`
	AdminID     int       `json:"admin_id" db:"admin_id"`
	Comment     string    `json:"comment" db:"comment"`
	CommentType string    `json:"comment_type" db:"comment_type"` // 'comment', 'status_change', 'assignment', 'resolution'
	Metadata    JSONB     `json:"metadata" db:"metadata"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// JOIN поля
	Admin *User `json:"admin,omitempty"`
}

// ProblemStatusHistory представляет историю изменений статуса проблемы
type ProblemStatusHistory struct {
	ID            int       `json:"id" db:"id"`
	ProblemID     int       `json:"problem_id" db:"problem_id"`
	AdminID       *int      `json:"admin_id" db:"admin_id"`
	OldStatus     *string   `json:"old_status" db:"old_status"`
	NewStatus     string    `json:"new_status" db:"new_status"`
	OldAssignedTo *int      `json:"old_assigned_to" db:"old_assigned_to"`
	NewAssignedTo *int      `json:"new_assigned_to" db:"new_assigned_to"`
	Comment       *string   `json:"comment" db:"comment"`
	Metadata      JSONB     `json:"metadata" db:"metadata"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`

	// JOIN поля
	Admin       *User `json:"admin,omitempty"`
	OldAssigned *User `json:"old_assigned,omitempty"`
	NewAssigned *User `json:"new_assigned,omitempty"`
}

// AdminLog представляет лог действий администратора
type AdminLog struct {
	ID         int       `json:"id" db:"id"`
	AdminID    int       `json:"admin_id" db:"admin_id"`
	AdminEmail string    `json:"admin_email" db:"admin_email"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   *int      `json:"entity_id" db:"entity_id"`
	Action     string    `json:"action" db:"action"`
	Details    JSONB     `json:"details" db:"details"`
	IPAddress  *string   `json:"ip_address" db:"ip_address"`
	UserAgent  *string   `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// MonitoringSettings настройки системы мониторинга
type MonitoringSettings struct {
	ID                  int       `json:"id" db:"id"`
	DelayThresholdHours int       `json:"delay_threshold_hours" db:"delay_threshold_hours"`
	CriticalDelayHours  int       `json:"critical_delay_hours" db:"critical_delay_hours"`
	NotifyOnDelays      bool      `json:"notify_on_delays" db:"notify_on_delays"`
	NotifyOnReturns     bool      `json:"notify_on_returns" db:"notify_on_returns"`
	NotifyOnProblems    bool      `json:"notify_on_problems" db:"notify_on_problems"`
	NotificationEmails  []string  `json:"notification_emails" db:"notification_emails"`
	AutoCreateProblems  bool      `json:"auto_create_problems" db:"auto_create_problems"`
	AutoAssignProblems  bool      `json:"auto_assign_problems" db:"auto_assign_problems"`
	DailyReportEnabled  bool      `json:"daily_report_enabled" db:"daily_report_enabled"`
	WeeklyReportEnabled bool      `json:"weekly_report_enabled" db:"weekly_report_enabled"`
	ReportRecipients    []string  `json:"report_recipients" db:"report_recipients"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy           *int      `json:"updated_by" db:"updated_by"`
}

// DashboardCache кеш данных для dashboard
type DashboardCache struct {
	ID        int       `json:"id" db:"id"`
	CacheKey  string    `json:"cache_key" db:"cache_key"`
	CacheData JSONB     `json:"cache_data" db:"cache_data"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// User минимальная информация о пользователе
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CourierStats статистика по курьерской службе
type CourierStats struct {
	Name        string  `json:"name"`
	Shipments   int     `json:"shipments"`
	Delivered   int     `json:"delivered"`
	AvgTime     float64 `json:"avg_time"`
	SuccessRate float64 `json:"success_rate"`
}

// DailyStats статистика за день
type DailyStats struct {
	Date      string `json:"date"`
	Shipments int    `json:"shipments"`
	Delivered int    `json:"delivered"`
	InTransit int    `json:"in_transit"`
	Problems  int    `json:"problems"`
}

// DashboardStats статистика для dashboard
type DashboardStats struct {
	TodayShipments      int            `json:"today_shipments"`
	TodayDelivered      int            `json:"today_delivered"`
	ActiveShipments     int            `json:"active_shipments"`
	ProblemShipments    int            `json:"problem_shipments"`
	AvgDeliveryTime     float64        `json:"avg_delivery_time"`
	DeliverySuccessRate float64        `json:"delivery_success_rate"`
	WeeklyDeliveries    []DailyStats   `json:"weekly_deliveries"`
	StatusDistribution  map[string]int `json:"status_distribution"`
	CourierPerformance  []CourierStats `json:"courier_performance"`
}

// JSONB тип для работы с JSONB полями PostgreSQL
type JSONB map[string]interface{}

// Value реализует интерфейс driver.Valuer для JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Scan реализует интерфейс sql.Scanner для JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// ShipmentDetails детальная информация об отправлении для админки
type ShipmentDetails struct {
	ID             int        `json:"id"`
	Type           string     `json:"type"` // bex, postexpress, marketplace
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`

	// Информация о товаре
	OrderInfo OrderInfo `json:"order_info"`

	// Адреса
	SenderAddress   Address `json:"sender_address"`
	ReceiverAddress Address `json:"receiver_address"`

	// История статусов
	StatusHistory []StatusHistoryEntry `json:"status_history"`

	// Проблемы
	Problems []ProblemShipment `json:"problems"`

	// Документы
	Documents []Document `json:"documents"`
}

// OrderInfo информация о заказе
type OrderInfo struct {
	OrderID      int     `json:"order_id"`
	ProductName  string  `json:"product_name"`
	ProductImage *string `json:"product_image"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
}

// Address адрес доставки
type Address struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// StatusHistoryEntry запись в истории статусов
type StatusHistoryEntry struct {
	Status     string    `json:"status"`
	StatusText string    `json:"status_text"`
	Timestamp  time.Time `json:"timestamp"`
	Location   *string   `json:"location"`
}

// Document документ отправления
type Document struct {
	Type      string    `json:"type"` // label, invoice, customs
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// ShipmentsFilter фильтры для списка отправлений
type ShipmentsFilter struct {
	Status         *string    `json:"status"`
	CourierService *string    `json:"courier_service"`
	DateFrom       *time.Time `json:"date_from"`
	DateTo         *time.Time `json:"date_to"`
	City           *string    `json:"city"`
	TrackingNumber *string    `json:"tracking_number"`
	HasProblems    *bool      `json:"has_problems"`
	Page           int        `json:"page"`
	Limit          int        `json:"limit"`
	SortBy         string     `json:"sort_by"`
	SortOrder      string     `json:"sort_order"`
}
