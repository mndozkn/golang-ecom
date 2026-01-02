package mailer

import (
	"net/smtp"
	"os"
)

type Mailer interface {
	SendWelcomeEmail(to string) error
}

type mailtrapService struct {
	auth smtp.Auth
	addr string
	from string
}

func NewMailtrapService() Mailer {
	host := "live.smtp.mailtrap.io"
	port := "587"
	username := os.Getenv("MAILTRAP_USER")
	password := os.Getenv("MAILTRAP_PASS")

	return &mailtrapService{
		auth: smtp.PlainAuth("", username, password, host),
		addr: host + ":" + port,
		from: "hello@demomailtrap.co",
	}
}

func (m *mailtrapService) SendWelcomeEmail(to string) error {
	msg := []byte("Subject: Hoş Geldin!\r\n" +
		"To: " + to + "\r\n" +
		"\r\n" +
		"Uygulamamıza başarıyla kayıt oldunuz!")

	return smtp.SendMail(m.addr, m.auth, m.from, []string{to}, msg)
}
