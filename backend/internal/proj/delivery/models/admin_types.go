package models

import (
	"time"
)

// ProviderConfig represents API configuration for a provider
type ProviderConfig struct {
	APIURL    string `json:"api_url"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	TestMode  bool   `json:"test_mode"`
}

// ProblemAssignment represents problem assignment request
type ProblemAssignment struct {
	AdminID int `json:"admin_id"`
}

// ProblemResolution represents problem resolution request
type ProblemResolution struct {
	Resolution     string `json:"resolution"`
	NotifyCustomer bool   `json:"notify_customer"`
}

// ProblemShipment represents a shipment with problems
type ProblemShipment struct {
	ID             int       `json:"id"`
	TrackingNumber string    `json:"tracking_number"`
	ProviderName   string    `json:"provider_name"`
	ProblemType    string    `json:"problem_type"`
	Priority       string    `json:"priority"`
	Status         string    `json:"status"`
	AssignedTo     *string   `json:"assigned_to,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	CustomerName   string    `json:"customer_name"`
	CustomerPhone  string    `json:"customer_phone"`
	Description    string    `json:"description"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TodayShipments  int             `json:"today_shipments"`
	TodayDelivered  int             `json:"today_delivered"`
	InTransit       int             `json:"in_transit"`
	Problems        int             `json:"problems"`
	AvgDeliveryTime string          `json:"avg_delivery_time"`
	SuccessRate     float64         `json:"success_rate"`
	ProviderStats   []ProviderStats `json:"provider_stats"`
	CostAnalysis    CostAnalysis    `json:"cost_analysis"`
}

// ProviderStats represents statistics per provider
type ProviderStats struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Shipments   int     `json:"shipments"`
	Delivered   int     `json:"delivered"`
	SuccessRate float64 `json:"success_rate"`
	AvgTime     string  `json:"avg_time"`
}

// CostAnalysis represents cost analysis data
type CostAnalysis struct {
	AvgCost   float64 `json:"avg_cost"`
	TotalCost float64 `json:"total_cost"`
	Savings   float64 `json:"savings"`
}

// AnalyticsData represents comprehensive analytics
type AnalyticsData struct {
	Period                 string             `json:"period"`
	TotalShipments         int                `json:"total_shipments"`
	DeliveryRate           float64            `json:"delivery_rate"`
	AvgDeliveryTime        string             `json:"avg_delivery_time"`
	CustomerSatisfaction   float64            `json:"customer_satisfaction"`
	CostPerShipment        float64            `json:"cost_per_shipment"`
	ProblemRate            float64            `json:"problem_rate"`
	TrendData              []TrendPoint       `json:"trend_data"`
	ProviderComparison     []ProviderStats    `json:"provider_comparison"`
	GeographicDistribution map[string]float64 `json:"geographic_distribution"`
}

// TrendPoint represents a point in trend data
type TrendPoint struct {
	Date      string `json:"date"`
	Shipments int    `json:"shipments"`
	Delivered int    `json:"delivered"`
	Problems  int    `json:"problems"`
}
