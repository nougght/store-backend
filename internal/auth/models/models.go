package models

import (
	"net"
	"time"
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
