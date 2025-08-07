package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/smtp"
)

type MailHandler struct {
}

func NewMailHandler() *MailHandler {
	return &MailHandler{}
}
func (h *MailHandler) SendMailCode(mailAdress string) {
	// Данные SMTP Yandex
	smtpHost := "smtp.yandex.ru"
	smtpPort := "587"
	smtpUsername := "nougght@yandex.ru" // Ваш логин Yandex
	smtpPassword := "oezxxqzwefhqdvnl"  // Пароль или app-specific пароль

	// Настройки письма
	from := "nougght@yandex.ru"
	to := []string{"muhamedmakusev@gmail.com"}
	// subject := "Код подтверждения"
	// body := fmt.Sprintf(`Здравствуйте! Ваш код подтверждения - %s. Спасибо за использование нашего сервиса.`, code)
	subject := "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte("Код подтверждения")) + "?="
	// Формируем MIME-письмо

	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	code := fmt.Sprintf("%06d", n) // Добавляем ведущие нули
	message := "From: Test <nougght@yandex.ru>\r\n" +
		"To: muhamedmakusev@gmail.com\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		"Ваш код: " + code

	// Аутентификация и отправка
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		to,
		[]byte(message),
	)

	if err != nil {
		fmt.Println("Ошибка отправки:", err)
	} else {
		fmt.Println("Письмо успешно отправлено! ")
	}
}
