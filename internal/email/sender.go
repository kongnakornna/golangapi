package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

// SMTPConfig holds SMTP server settings.
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// SMTPSender implements EmailSender using SMTP.
type SMTPSender struct {
	config SMTPConfig
}

// NewSMTPSender creates a new SMTP sender from environment variables.
func NewSMTPSender() *SMTPSender {
	return &SMTPSender{
		config: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},
	}
}

// SendOTP sends a 6‑digit OTP to the given email address.
func (s *SMTPSender) SendOTP(email, otp string) error {
	// Build the email message
	subject := "รหัสยืนยันการรีเซ็ตรหัสผ่าน"
	body := fmt.Sprintf("รหัส OTP ของคุณคือ: %s\nรหัสนี้ใช้ได้ 15 นาที", otp)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// SMTP authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Server address
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// Optional: if using port 465 (SSL/TLS directly)
	if s.config.Port == "465" {
		return s.sendTLS(addr, auth, email, msg)
	}

	// Standard SMTP (port 25, 587)
	return smtp.SendMail(addr, auth, s.config.From, []string{email}, msg)
}

// sendTLS handles SMTP over explicit TLS (port 465).
func (s *SMTPSender) sendTLS(addr string, auth smtp.Auth, to string, msg []byte) error {
	// Connect to the server
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.config.Host})
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	// Set sender and recipient
	if err = client.Mail(s.config.From); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}