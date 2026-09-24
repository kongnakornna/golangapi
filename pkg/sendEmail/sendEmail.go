package sendEmail

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"icmongolang/config"

	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(ctx context.Context, from string, to string, subject string, bodyHtml string, bodyPlain string) error
}

type emailSender struct {
	cfg *config.Config
}

func NewEmailSender(cfg *config.Config) EmailSender {
	return &emailSender{
		cfg: cfg,
	}
}

func (es *emailSender) SendEmail(ctx context.Context, from string, to string, subject string, bodyHtml string, bodyPlain string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyHtml)
	m.AddAlternative("text/plain", bodyPlain)

	d := gomail.NewDialer(es.cfg.SmtpEmail.Host, es.cfg.SmtpEmail.Port, es.cfg.SmtpEmail.User, es.cfg.SmtpEmail.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	timeout := es.cfg.SmtpEmail.Timeout
	if timeout <= 0 {
		timeout = 30
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- d.DialAndSend(m)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("smtp timeout after %d seconds: %w", timeout, ctx.Err())
	}
}
