package models

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"time"

	"encoding/base64"
	"net/smtp"

	"github.com/golang-jwt/jwt"

	"github.com/google/uuid"
)

type User struct {
	UserID     string     `db:"user_id" json:"user_id"`
	Email      string     `db:"email" json:"email"`
	Phone      string     `db:"phone" json:"phone"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	LastActive *time.Time `db:"last_active" json:"last_active"`
	Username   string     `db:"username" json:"username"`
}

type FavouriteItem struct {
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Session struct {
	SessionID  string     `db:"session_id" json:"session_id"`
	UserID     string     `db:"user_id" json:"user_id"`
	Token      string     `db:"token" json:"token"`
	DeviceInfo *string    `db:"device_info" json:"device_info,omitempty"`
	IPAddress  *net.IP    `db:"ip_address" json:"ip_address,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ExpiresAt  *time.Time `db:"expires_at" json:"expires_at"`
}

type AuthCode struct {
	AuthCodeID string    `db:"code_id" json:"code_id"`
	UserID     string    `db:"user_id" json:"user_id"`
	Recipient  string    `db:"recipient" json:"recipient"`
	Code       string    `db:"code" json:"code"`
	Channel    string    `db:"channel" json:"channel"`
	ExpiresAt  time.Time `db:"expires_at" json:"expires_at"`
	Used       bool      `db:"used" json:"used"`
	IPAddress  *net.IP   `db:"ip_address" json:"ip_address,omitempty"`
}

type EmailSMTP struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

func NewEmailSMTP() *EmailSMTP {
	return &EmailSMTP{
		Host:     "smtp.yandex.ru",
		Port:     587,
		Username: "nsrexu@yandex.ru",
		Password: "uzidfzwztlpecrup",
		From:     "nsrexu@yandex.ru",
	}
}

func (e *EmailSMTP) SendVerificationCode(email string, code string) error {
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
		e.Host+":"+fmt.Sprint(e.Port),
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

type SmsApi struct {
}

func NewSmsApi() *SmsApi {
	return &SmsApi{}
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
type Tools struct {
}

func (t *Tools) IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (t *Tools) GenerateAuthCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	code := fmt.Sprintf("%06d", n) // Добавляем ведущие нули
	return code
}

func (t *Tools) GenerateJWTToken(userId string) (string, error) {
	claims := CustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *Tools) ValidateJWTToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	// fmt.Println("token", token)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		fmt.Println(userID)
		if !ok {
			return "", fmt.Errorf("invalid token")
		}

		return userID, nil
	}
	return "", fmt.Errorf("invalid token")
}
