package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"backend/pkg/logger"
)

// AlertManager manages different alert channels
type AlertManager struct {
	config   *AlertConfig
	logger   *logger.Logger
	channels []AlertChannel
}

// AlertConfig holds configuration for alerts
type AlertConfig struct {
	// Email configuration
	EmailEnabled bool
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	AlertEmails  []string

	// Slack configuration
	SlackEnabled    bool
	SlackWebhookURL string
	SlackChannel    string

	// Telegram configuration
	TelegramEnabled  bool
	TelegramBotToken string
	TelegramChatID   string

	// General settings
	Environment string
	ServiceName string
}

// AlertChannel interface for different alert channels
type AlertChannel interface {
	Send(alert Alert) error
	Name() string
}

// Alert represents an alert to be sent
type Alert struct {
	Level       AlertLevel             `json:"level"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Environment string                 `json:"environment"`
	Service     string                 `json:"service"`
}

// AlertLevel represents the severity of an alert
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "INFO"
	AlertLevelWarning  AlertLevel = "WARNING"
	AlertLevelCritical AlertLevel = "CRITICAL"
	AlertLevelSecurity AlertLevel = "SECURITY"
)

// NewAlertManager creates a new alert manager
func NewAlertManager(config *AlertConfig, logger *logger.Logger) *AlertManager {
	am := &AlertManager{
		config:   config,
		logger:   logger,
		channels: make([]AlertChannel, 0),
	}

	// Initialize enabled channels
	if config.EmailEnabled {
		am.channels = append(am.channels, NewEmailChannel(config))
	}
	if config.SlackEnabled {
		am.channels = append(am.channels, NewSlackChannel(config))
	}
	if config.TelegramEnabled {
		am.channels = append(am.channels, NewTelegramChannel(config))
	}

	// Always add log channel
	am.channels = append(am.channels, NewLogChannel(logger))

	return am
}

// SendAlert sends an alert through all configured channels
func (am *AlertManager) SendAlert(ctx context.Context, alert Alert) {
	// Set default values
	alert.Timestamp = time.Now()
	alert.Environment = am.config.Environment
	alert.Service = am.config.ServiceName

	// Send through all channels
	for _, channel := range am.channels {
		go func(ch AlertChannel) {
			if err := ch.Send(alert); err != nil {
				am.logger.Error("Failed to send alert via %s: %v", ch.Name(), err)
			} else {
				am.logger.Info("Alert sent via %s: %s", ch.Name(), alert.Title)
			}
		}(channel)
	}
}

// SendSecurityAlert sends a security-specific alert
func (am *AlertManager) SendSecurityAlert(ctx context.Context, alertType, title, message string, details map[string]interface{}) {
	alert := Alert{
		Level:   AlertLevelSecurity,
		Type:    alertType,
		Title:   fmt.Sprintf("🚨 SECURITY ALERT: %s", title),
		Message: message,
		Details: details,
	}
	am.SendAlert(ctx, alert)
}

// SendPaymentFraudAlert sends an alert about potential payment fraud
func (am *AlertManager) SendPaymentFraudAlert(ctx context.Context, userID int, ip string, details map[string]interface{}) {
	alert := Alert{
		Level:   AlertLevelSecurity,
		Type:    "payment_fraud",
		Title:   "🚨 Potential Payment Fraud Detected",
		Message: fmt.Sprintf("Suspicious payment activity detected from User ID: %d, IP: %s", userID, ip),
		Details: details,
	}
	am.SendAlert(ctx, alert)
}

// SendRateLimitAlert sends an alert about rate limit violations
func (am *AlertManager) SendRateLimitAlert(ctx context.Context, endpoint string, userID int, ip string, violations int) {
	alert := Alert{
		Level:   AlertLevelWarning,
		Type:    "rate_limit_exceeded",
		Title:   "⚠️ Rate Limit Violations",
		Message: fmt.Sprintf("Multiple rate limit violations on %s", endpoint),
		Details: map[string]interface{}{
			"endpoint":   endpoint,
			"user_id":    userID,
			"ip":         ip,
			"violations": violations,
		},
	}
	am.SendAlert(ctx, alert)
}

// SendWebhookFailureAlert sends an alert about webhook failures
func (am *AlertManager) SendWebhookFailureAlert(ctx context.Context, webhookType string, failures int, lastError string) {
	alert := Alert{
		Level:   AlertLevelWarning,
		Type:    "webhook_failure",
		Title:   "⚠️ Webhook Processing Failures",
		Message: fmt.Sprintf("Multiple webhook failures for type: %s", webhookType),
		Details: map[string]interface{}{
			"webhook_type": webhookType,
			"failures":     failures,
			"last_error":   lastError,
		},
	}
	am.SendAlert(ctx, alert)
}

// EmailChannel sends alerts via email
type EmailChannel struct {
	config *AlertConfig
}

func NewEmailChannel(config *AlertConfig) *EmailChannel {
	return &EmailChannel{config: config}
}

func (e *EmailChannel) Name() string {
	return "Email"
}

func (e *EmailChannel) Send(alert Alert) error {
	if len(e.config.AlertEmails) == 0 {
		return fmt.Errorf("no alert emails configured")
	}

	// Format email body
	body := fmt.Sprintf(`
Alert Level: %s
Type: %s
Time: %s
Environment: %s
Service: %s

%s

Details:
%s
	`, alert.Level, alert.Type, alert.Timestamp.Format(time.RFC3339),
		alert.Environment, alert.Service, alert.Message,
		formatDetails(alert.Details))

	// Send email
	auth := smtp.PlainAuth("", e.config.SMTPUser, e.config.SMTPPassword, e.config.SMTPHost)
	to := e.config.AlertEmails
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s",
		strings.Join(to, ","), alert.Title, body))

	addr := fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort)
	return smtp.SendMail(addr, auth, e.config.SMTPUser, to, msg)
}

// SlackChannel sends alerts to Slack
type SlackChannel struct {
	config *AlertConfig
}

func NewSlackChannel(config *AlertConfig) *SlackChannel {
	return &SlackChannel{config: config}
}

func (s *SlackChannel) Name() string {
	return "Slack"
}

func (s *SlackChannel) Send(alert Alert) error {
	// Format Slack message
	color := "#17a2b8" // info
	switch alert.Level {
	case AlertLevelInfo:
		color = "#17a2b8"
	case AlertLevelWarning:
		color = "#ffc107"
	case AlertLevelCritical:
		color = "#dc3545"
	case AlertLevelSecurity:
		color = "#721c24"
	}

	payload := map[string]interface{}{
		"channel": s.config.SlackChannel,
		"attachments": []map[string]interface{}{
			{
				"color":     color,
				"title":     alert.Title,
				"text":      alert.Message,
				"fields":    formatSlackFields(alert.Details),
				"footer":    fmt.Sprintf("%s | %s", alert.Service, alert.Environment),
				"timestamp": alert.Timestamp.Unix(),
			},
		},
	}

	// Send to Slack webhook
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	resp, err := http.Post(s.config.SlackWebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// TelegramChannel sends alerts to Telegram
type TelegramChannel struct {
	config *AlertConfig
}

func NewTelegramChannel(config *AlertConfig) *TelegramChannel {
	return &TelegramChannel{config: config}
}

func (t *TelegramChannel) Name() string {
	return "Telegram"
}

func (t *TelegramChannel) Send(alert Alert) error {
	// Format Telegram message
	icon := "ℹ️"
	switch alert.Level {
	case AlertLevelInfo:
		icon = "ℹ️"
	case AlertLevelWarning:
		icon = "⚠️"
	case AlertLevelCritical:
		icon = "🔴"
	case AlertLevelSecurity:
		icon = "🚨"
	}

	message := fmt.Sprintf("%s *%s*\n\n%s\n\n_%s | %s_",
		icon, escapeMarkdown(alert.Title),
		escapeMarkdown(alert.Message),
		alert.Environment, alert.Service)

	// Add details if present
	if len(alert.Details) > 0 {
		message += "\n\n*Details:*\n"
		for k, v := range alert.Details {
			message += fmt.Sprintf("• %s: %v\n", k, v)
		}
	}

	// Send to Telegram
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.config.TelegramBotToken)
	payload := map[string]interface{}{
		"chat_id":    t.config.TelegramChatID,
		"text":       message,
		"parse_mode": "Markdown",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload)) // #nosec G107 - URL is Telegram API with validated token
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

// LogChannel logs alerts
type LogChannel struct {
	logger *logger.Logger
}

func NewLogChannel(logger *logger.Logger) *LogChannel {
	return &LogChannel{logger: logger}
}

func (l *LogChannel) Name() string {
	return "Log"
}

func (l *LogChannel) Send(alert Alert) error {
	switch alert.Level {
	case AlertLevelCritical, AlertLevelSecurity:
		l.logger.Error("ALERT [%s] %s: %s | Details: %+v",
			alert.Level, alert.Title, alert.Message, alert.Details)
	case AlertLevelInfo, AlertLevelWarning:
		l.logger.Info("ALERT [%s] %s: %s | Details: %+v",
			alert.Level, alert.Title, alert.Message, alert.Details)
	}
	return nil
}

// Helper functions

func formatDetails(details map[string]interface{}) string {
	if len(details) == 0 {
		return "No additional details"
	}

	var parts []string
	for k, v := range details {
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}
	return strings.Join(parts, "\n")
}

func formatSlackFields(details map[string]interface{}) []map[string]interface{} {
	fields := make([]map[string]interface{}, 0, len(details))
	for k, v := range details {
		fields = append(fields, map[string]interface{}{
			"title": k,
			"value": fmt.Sprintf("%v", v),
			"short": true,
		})
	}
	return fields
}

func escapeMarkdown(text string) string {
	// Escape special markdown characters for Telegram
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}
