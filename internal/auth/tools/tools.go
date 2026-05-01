package tools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GenerateAuthCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	code := fmt.Sprintf("%06d", n) // Добавляем ведущие нули
	return code
}

type JwtTools struct {
	secretKey string
}

func NewJwtTools(secretKey string) *JwtTools {
	return &JwtTools{secretKey: secretKey}
}

func (t *JwtTools) GenerateJWTToken(userId string) (string, error) {
	claims := CustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *JwtTools) ValidateJWTToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
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

type SmsApi struct {
}

func NewSmsApi() *SmsApi {
	return &SmsApi{}
}
