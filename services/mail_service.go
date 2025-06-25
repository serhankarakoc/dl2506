package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"davet.link/configs/logconfig"
	"go.uber.org/zap"
)

// IMailService defines the interface for mail operations
type IMailService interface {
	SendMail(to, subject, body string) error
}

// MailService implements IMailService
type MailService struct {
	host     string
	port     string
	username string
	password string
}

// NewMailService creates a new MailService instance
func NewMailService() IMailService {
	return &MailService{
		host:     getEnvWithDefault("SMTP_HOST", "smtp.example.com"),
		port:     getEnvWithDefault("SMTP_PORT", "587"),
		username: getEnvWithDefault("SMTP_USERNAME", ""),
		password: getEnvWithDefault("SMTP_PASSWORD", ""),
	}
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SendMail sends an email with the given parameters
func (m *MailService) SendMail(to, subject, body string) error {
	if to == "" {
		return fmt.Errorf("alıcı e-posta adresi boş olamaz")
	}

	message, err := m.buildMessage(to, subject, body)
	if err != nil {
		return fmt.Errorf("e-posta mesajı oluşturulamadı: %w", err)
	}

	client, err := m.createSMTPClient()
	if err != nil {
		return fmt.Errorf("SMTP istemcisi oluşturulamadı: %w", err)
	}
	defer func() {
		if err := client.Quit(); err != nil {
			logconfig.Log.Warn("SMTP bağlantısı kapatılırken hata oluştu", zap.Error(err))
		}
	}()

	if err := m.sendMail(client, to, message); err != nil {
		return fmt.Errorf("e-posta gönderilemedi: %w", err)
	}

	return nil
}

// buildMessage constructs the email message
func (m *MailService) buildMessage(to, subject, body string) ([]byte, error) {
	if subject == "" {
		subject = "(Konu Belirtilmemiş)"
	}

	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n",
		m.username, to, subject)
	return []byte(header + body), nil
}

// createSMTPClient establishes a secure SMTP connection
func (m *MailService) createSMTPClient() (*smtp.Client, error) {
	address := fmt.Sprintf("%s:%s", m.host, m.port)

	// Establish TLS connection
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         m.host,
	})
	if err != nil {
		return nil, fmt.Errorf("TLS bağlantısı kurulamadı: %w", err)
	}

	// Create SMTP client
	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("SMTP istemcisi oluşturulamadı: %w", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	if err := client.Auth(auth); err != nil {
		client.Quit()
		return nil, fmt.Errorf("kimlik doğrulama başarısız: %w", err)
	}

	return client, nil
}

// sendMail performs the actual email sending
func (m *MailService) sendMail(client *smtp.Client, to string, message []byte) error {
	// Set sender
	if err := client.Mail(m.username); err != nil {
		return fmt.Errorf("gönderici ayarlanamadı: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("alıcı ayarlanamadı: %w", err)
	}

	// Send email data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("veri gönderimi başlatılamadı: %w", err)
	}

	if _, err := writer.Write(message); err != nil {
		writer.Close()
		return fmt.Errorf("mesaj yazılamadı: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("mesaj gönderimi tamamlanamadı: %w", err)
	}

	return nil
}

// Ensure MailService implements IMailService
var _ IMailService = (*MailService)(nil)
