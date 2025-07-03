// Package handler
// backend/internal/proj/payments/handler/responses.go
package handler

// WebhookResponse represents webhook processing response
// Commented out due to conflict with order_webhook_handler.go
// type WebhookResponse struct {
// 	Status  string `json:"status" example:"success"`
// 	Message string `json:"message" example:"payments.webhook.processed"`
// }

// StripeWebhookRequest represents the expected structure from Stripe webhook
// This is a simplified version - actual Stripe webhooks have more fields
type StripeWebhookRequest struct {
	ID      string                 `json:"id" example:"evt_1234567890"`
	Object  string                 `json:"object" example:"event"`
	Type    string                 `json:"type" example:"checkout.session.completed"`
	Created int64                  `json:"created" example:"1234567890"`
	Data    map[string]interface{} `json:"data"`
}
