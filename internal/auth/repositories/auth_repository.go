package repositories

import (
	"context"
	"fmt"
	"store-server/internal/auth/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	fmt.Println("username", user.Username, "email", user.Email, "phone", user.Phone)
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO auth.users (username, email, phone)
		VALUES ($1, $2, $3) returning user_id
	`, user.Username, user.Email, user.Phone).Scan(&user.UserID)
	return err
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM auth.users WHERE user_id = $1", userID)
	return &user, err
}

func (r *AuthRepository) DeleteUserByID(ctx context.Context, userID string) error {
	query := `DELETE FROM auth.users WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *AuthRepository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM auth.users WHERE phone = $1", phone)
	return &user, err
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM auth.users WHERE email = $1", email)
	return &user, err
}

func (r *AuthRepository) CheckUserExists(ctx context.Context, emailOrPhone string) (string, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM auth.users WHERE email = $1 OR phone = $2)", emailOrPhone, emailOrPhone)
	if err != nil {
		return "", err
	} else if exists {
		var user_id string
		err := r.db.GetContext(ctx, &user_id, "SELECT user_id FROM auth.users WHERE email = $1 OR phone = $2", emailOrPhone, emailOrPhone)
		if err == nil {
			return user_id, nil
		}
	}
	return "", nil
}

func (r *AuthRepository) CreateAuthCode(ctx context.Context, code *models.AuthCode) error {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO auth.auth_codes (code, channel, expires_at, recipient)
		VALUES ($1, $2, $3, $4) returning code_id
	`, code.Code, code.Channel, code.ExpiresAt, code.Recipient).StructScan(code)
	return err
}

func (r *AuthRepository) UpdateAuthCode(ctx context.Context, code *models.AuthCode) error {
	d, err := r.db.ExecContext(ctx, `
		UPDATE auth.auth_codes
		SET user_id = $1, code = $2, channel = $3, expires_at = $4, used = $5
		WHERE code_id = $6 AND recipient = $7
	`, code.UserID, code.Code, code.Channel, code.ExpiresAt, code.Used, code.AuthCodeID, code.Recipient)
	if err != nil {
		return err
	}
	rows, err := d.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("updated", rows, "code", code.Code, "user_id", code.UserID, "channel", code.Channel, "expires_at", code.ExpiresAt, "used", code.Used, "code_id", code.AuthCodeID, "recipient", code.Recipient)
	return err
}

func (r *AuthRepository) VerifyAuthCode(ctx context.Context, recipient string, code string) (bool, error) {
	var isValid bool
	err := r.db.GetContext(ctx, &isValid, "SELECT EXISTS(SELECT 1 FROM auth.auth_codes WHERE code = $1 AND used = false AND recipient = $2)", code, recipient)
	return isValid, err
}

func (r *AuthRepository) GetAuthCodeByRecipient(ctx context.Context, recipient string) (*models.AuthCode, error) {
	var authCode models.AuthCode
	err := r.db.GetContext(ctx, &authCode, "SELECT * FROM auth.auth_codes WHERE recipient = $1 AND used = false", recipient)
	if authCode.AuthCodeID == "" {
		return nil, err
	}
	return &authCode, err
}

func (r *AuthRepository) GetAuthCodeByCodeAndUserID(ctx context.Context, code string, userID string) (*models.AuthCode, error) {
	var authCode models.AuthCode
	err := r.db.GetContext(ctx, &authCode, "SELECT * FROM auth.auth_codes WHERE code = $1 AND user_id = $2", code, userID)
	return &authCode, err
}

func (r *AuthRepository) GetAuthCodeByID(ctx context.Context, codeID string) (*models.AuthCode, error) {
	var authCode models.AuthCode
	err := r.db.GetContext(ctx, &authCode, "SELECT * FROM auth.auth_codes WHERE code_id = $1", codeID)
	return &authCode, err
}

func (r *AuthRepository) DeleteAuthCodeByID(ctx context.Context, codeID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM auth.auth_codes WHERE code_id = $1", codeID)
	return err
}

func (r *AuthRepository) CreateSession(ctx context.Context, session *models.Session) error {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO auth.sessions (token, user_id)
		VALUES ($1, $2) returning session_id
	`, session.Token, session.UserID).StructScan(session)
	return err
}

func (r *AuthRepository) GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error) {
	var session models.Session
	err := r.db.GetContext(ctx, &session, "SELECT * FROM auth.sessions WHERE user_id = $1", userID)
	return &session, err
}

func (r *AuthRepository) DeleteSessionByID(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM auth.sessions WHERE session_id = $1", sessionID)
	return err
}
