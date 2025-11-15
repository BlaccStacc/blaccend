package mail

import (
	"github.com/BlaccStacc/blaccend/internal/config"
)

// added SMTP env variables and possibly ruined config

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPass,
		From:     cfg.SMTPFrom,
	}
}

func (m *Mailer) Send(to string, subject string, htmlBody string) error {

	return nil
}

func SendVerificationEmail(to, token string) error {

}

func SendPasswordResetEmail(to, token string) error {

}
