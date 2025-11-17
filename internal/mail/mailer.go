package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/BlaccStacc/blaccend/internal/config"
)

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	AppURL   string
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPass,
		From:     cfg.SMTPFrom,
		AppURL:   cfg.AppURL,
	}
}

func (m *Mailer) Send(to string, subject string, htmlBody string) error {
	if m.Host == "" || m.Port == 0 {
		return fmt.Errorf("smtp not configured (host=%q port=%d)", m.Host, m.Port)
	}
	if m.From == "" {
		return fmt.Errorf("smtp from address is empty")
	}

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	var auth smtp.Auth
	if m.Username != "" {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}

	if err := smtp.SendMail(addr, auth, m.From, []string{to}, msg.Bytes()); err != nil {
		return fmt.Errorf("smtp sendmail failed: %w", err)
	}

	return nil
}

func SendVerificationEmail(to, token string) error {
	cfg := config.Load()
	mailer := NewMailer(cfg)

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", cfg.AppURL, token)

	const tpl = `
		<h2>Verify your email</h2>
		<p>Click the link below to confirm your account:</p>
		<p><a href="{{.URL}}">{{.URL}}</a></p>
	`

	t, err := template.New("verify").Parse(tpl)
	if err != nil {
		return fmt.Errorf("parse verify template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, map[string]string{"URL": verifyURL}); err != nil {
		return fmt.Errorf("execute verify template: %w", err)
	}

	if err := mailer.Send(to, "Verify your email", body.String()); err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}

	return nil
}
