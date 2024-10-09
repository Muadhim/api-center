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
	mailer.SetHeader("From", config.SmptUser)
	mailer.SetHeader("To", em.Email)
	mailer.SetAddressHeader("Cc", config.SmptUser, "Admin")
	mailer.SetHeader("Subject", em.Subject)
	mailer.SetBody("text/html", em.Message)

	port, err := strconv.Atoi(config.SmptPort)
	if err != nil {
		return err
	}
	dialer := gomail.NewDialer(config.SmptHost, port, config.SmptUser, config.SmptPass)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
