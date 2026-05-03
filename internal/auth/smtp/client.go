package smtp

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"store-server/config"
)

type SMTPClient struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

func NewSMTPClient(cfg *config.SMTPConfig) *SMTPClient {
	return &SMTPClient{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
	}
}

func (e *SMTPClient) SendVerificationCode(email string, code string) error {
	subject := "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte("Код подтверждения")) + "?="
	message := "From: Test <" + e.Username + ">\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		"Ваш код: " + code

	// Аутентификация и отправка
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	err := smtp.SendMail(
		e.Host+":"+e.Port,
		auth,
		e.From,
		[]string{email},
		[]byte(message),
	)

	if err != nil {
		fmt.Println("Ошибка отправки:", err, message)
		return err
	} else {
		fmt.Println("Письмо успешно отправлено! ")
		return nil
	}
}
