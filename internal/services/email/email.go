package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// Структура EmailService для хранения настроек SMTP
type EmailService struct {
	SMTPServer string
	SMTPPort   string
	Email      string
	Password   string
}

// Конструктор для создания нового экземпляра EmailService
func New(smtpServer, smtpPort, email, password string) *EmailService {
	return &EmailService{
		SMTPServer: smtpServer,
		SMTPPort:   smtpPort,
		Email:      email,
		Password:   password,
	}
}

// Функция для отправки письма
func (e *EmailService) SendEmail(ctx context.Context, toAddress, subject, body string) error {
	op := "EmailService.SendEmail"

	// Устанавливаем соединение с SMTP-сервером через TLS
	serverAddr := e.SMTPServer + ":" + e.SMTPPort
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Установите true, если самоподписанный сертификат (не рекомендуется)
		ServerName:         e.SMTPServer,
	}

	conn, err := tls.Dial("tcp", serverAddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("%s: ошибка TLS-соединения: %w", op, err)
	}
	defer conn.Close()

	// Создаём новый SMTP клиент поверх TLS-соединения
	client, err := smtp.NewClient(conn, e.SMTPServer)
	if err != nil {
		return fmt.Errorf("%s: ошибка создания SMTP клиента: %w", op, err)
	}
	defer client.Quit()

	// Аутентификация
	auth := smtp.PlainAuth("", e.Email, e.Password, e.SMTPServer)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("%s: ошибка аутентификации: %w", op, err)
	}

	// Устанавливаем адрес отправителя
	if err := client.Mail(e.Email); err != nil {
		return fmt.Errorf("%s: ошибка установки отправителя: %w", op, err)
	}

	// Указываем получателя
	if err := client.Rcpt(toAddress); err != nil {
		return fmt.Errorf("%s: ошибка установки получателя: %w", op, err)
	}

	// Записываем сообщение в SMTP-соединение
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("%s: ошибка открытия потока для данных: %w", op, err)
	}
	// Создаем сообщение с заголовком Content-Type для HTML
	message := fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", subject, body)
	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("%s: ошибка записи сообщения: %w", op, err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("%s: ошибка закрытия потока данных: %w", op, err)
	}

	return nil
}

func (e *EmailService) SendActivationLink(ctx context.Context, toAddress string, activationLink string) error {
	op := "EmailService.SendActivationLink"

	htmlBody := `
    <!DOCTYPE html>
    <html lang="ru">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Подтверждение регистрации</title>
    </head>
    <body>
        <h2>Спасибо за регистрацию!</h2>
        <p>Для подтверждения вашего аккаунта, пожалуйста, перейдите по следующей <a href="{{.ActivationLink}}">ссылке</a>.</p>
    </body>
    </html>
`
	htmlBody = strings.Replace(htmlBody, "{{.ActivationLink}}", activationLink, -1)

	err := e.SendEmail(ctx, toAddress, "Подтверждение регестрации", htmlBody)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
