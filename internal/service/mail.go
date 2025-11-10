package service

import (
	"auth-server/internal/config"

	"gopkg.in/gomail.v2"
)

type MailService struct {
}

func NewMailService() *MailService {
	return &MailService{}
}

func (s *MailService) SendActivationMAil(link string, to string) {
	cfg := config.Load()

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Host_Mail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Активация аккаунта")
	m.SetBody("text/html", `
        <h1>Подтверждение email</h1>
        <p>Нажмите на ссылку для активации:</p>
        <a href="http://localhost:`+cfg.ServerPort+`/activate/ `+link+`">Активировать</a>
    `)

	// Добавление вложения (опционально)
	// m.Attach("/path/to/file.pdf")

	d := gomail.NewDialer(cfg.Host_Mail, 587, cfg.Host_Mail, cfg.Mail_Passwd)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
