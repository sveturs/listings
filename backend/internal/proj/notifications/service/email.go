package service

import (
    "crypto/tls"
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
    log.Printf("Attempting to send email to %s with subject: %s", to, subject)
    
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

    // Формируем адрес SMTP-сервера
    addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)
    
    // Создаем незащищенное соединение
    c, err := smtp.Dial(addr)
    if err != nil {
        log.Printf("Error connecting to SMTP server: %v", err)
        return err
    }
    defer c.Quit()

    // Выполняем STARTTLS с отключенной проверкой сертификата
    tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         e.smtpHost,
    }
    
    if err = c.StartTLS(tlsConfig); err != nil {
        log.Printf("Error starting TLS: %v", err)
        
        // Если StartTLS не работает, попробуем без шифрования
        log.Printf("Trying without TLS...")
        c, err = smtp.Dial(addr)
        if err != nil {
            log.Printf("Error reconnecting to SMTP server: %v", err)
            return err
        }
        defer c.Quit()
    }

    // Аутентификация
    auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)
    if err = c.Auth(auth); err != nil {
        log.Printf("Authentication error: %v", err)
        return err
    }

    // Отправитель и получатель
    if err = c.Mail(e.senderEmail); err != nil {
        log.Printf("Error setting sender: %v", err)
        return err
    }
    
    if err = c.Rcpt(to); err != nil {
        log.Printf("Error setting recipient: %v", err)
        return err
    }

    // Отправка данных
    w, err := c.Data()
    if err != nil {
        log.Printf("Error getting data writer: %v", err)
        return err
    }
    
    _, err = w.Write([]byte(message))
    if err != nil {
        log.Printf("Error writing message: %v", err)
        return err
    }
    
    err = w.Close()
    if err != nil {
        log.Printf("Error closing writer: %v", err)
        return err
    }

    log.Printf("Email successfully sent to %s", to)
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