package models

import (
	"api-center/configs"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Email struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Subject string `json:"subject"`
}

func (em *Email) SendEmail() error {
	config := configs.LoadConfig()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.SmtpUser)
	mailer.SetHeader("To", em.Email)
	mailer.SetAddressHeader("Cc", config.SmtpUser, "Admin")
	mailer.SetHeader("Subject", em.Subject)
	mailer.SetBody("text/html", em.Message)

	port, err := strconv.Atoi(config.SmtpPort)
	if err != nil {
		return err
	}
	dialer := gomail.NewDialer(config.SmtpHost, port, config.SmtpUser, config.SmtpPass)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
