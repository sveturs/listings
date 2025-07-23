package service

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	senderEmail  string
	senderName   string
	smtpUsername string
	smtpPassword string
}

func NewEmailService(smtpHost, smtpPort, senderEmail, senderName, smtpUsername, smtpPassword string) *EmailService {
	return &EmailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		senderEmail:  senderEmail,
		senderName:   senderName,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
	}
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	log.Printf("Sending email to %s with subject: %s", to, subject)

	// Формируем заголовки письма
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", e.senderName, e.senderEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Собираем сообщение из заголовков и тела
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Альтернативный способ отправки - через smtp.SendMail с полным обходом проверки TLS
	// Создаем кастомную аутентификацию
	auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)

	// Формируем адрес SMTP-сервера
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)

	// Пробуем разные варианты отправки
	var lastErr error

	// Вариант 1: Прямая отправка без TLS
	err := smtp.SendMail(addr, auth, e.senderEmail, []string{to}, []byte(message))
	if err == nil {
		log.Printf("Email sent successfully to %s using direct SendMail", to)
		return nil
	}
	lastErr = err
	log.Printf("Direct email sending failed: %v. Trying alternative methods...", err)

	// Вариант 2: Ручное подключение без TLS
	conn, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("Error connecting to mail server: %v", err)
		return fmt.Errorf("all email sending methods failed: %w", lastErr)
	}
	defer func() {
		if err := conn.Quit(); err != nil {
			log.Printf("Error closing SMTP connection: %v", err)
		}
	}()

	// Установка параметров отправки
	if err = conn.Mail(e.senderEmail); err != nil {
		log.Printf("Error in MAIL FROM command: %v", err)
		return fmt.Errorf("MAIL FROM error: %w", err)
	}

	if err = conn.Rcpt(to); err != nil {
		log.Printf("Error in RCPT TO command: %v", err)
		return fmt.Errorf("RCPT TO error: %w", err)
	}

	w, err := conn.Data()
	if err != nil {
		log.Printf("Error getting data writer: %v", err)
		return fmt.Errorf("DATA command error: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("Error writing message body: %v", err)
		return fmt.Errorf("message writing error: %v", err)
	}

	if err = w.Close(); err != nil {
		log.Printf("Error closing message writer: %v", err)
		return fmt.Errorf("error finalizing message: %v", err)
	}

	log.Printf("Email sent successfully to %s using manual connection", to)
	return nil
}

// Метод для форматирования HTML-шаблона уведомления
func (e *EmailService) FormatNotificationEmail(title, message, listingID string) string {
	var listingLink string
	if listingID != "" {
		listingLink = fmt.Sprintf("<p><a href='https://SveTu.rs/marketplace/listings/%s' style='display: inline-block; background-color: #4CAF50; color: white; padding: 10px 15px; text-decoration: none; border-radius: 4px;'>Перейти к объявлению</a></p>", listingID)
	}

	template := `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>%s</title>
        <style>
            body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
            .container { max-width: 600px; margin: 0 auto; padding: 20px; }
            .header { background-color: #f8f9fa; padding: 15px; border-bottom: 1px solid #e9ecef; }
            .content { padding: 20px 0; }
            .footer { font-size: 12px; color: #6c757d; border-top: 1px solid #e9ecef; padding-top: 15px; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h2>%s</h2>
            </div>
            <div class="content">
                <p>%s</p>
                %s
            </div>
            <div class="footer">
                <p>Это автоматическое уведомление от SveTu.rs. Пожалуйста, не отвечайте на это письмо.</p>
                <p>Вы можете изменить настройки уведомлений в <a href="https://SveTu.rs/notifications/settings">личном кабинете</a>.</p>
            </div>
        </div>
    </body>
    </html>
    `

	return fmt.Sprintf(template, title, title, strings.ReplaceAll(message, "\n", "<br>"), listingLink)
}
