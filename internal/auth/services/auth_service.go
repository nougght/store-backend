package services

import (
	"context"
	"fmt"
	"store-server/config"
	"store-server/internal/auth/models"
	"store-server/internal/auth/repositories"
	"store-server/internal/auth/smtp"
	"store-server/internal/auth/tools"
)

type AuthService struct {
	repo      *repositories.AuthRepository
	smtp      *smtp.SMTPClient
	smsSender *tools.SmsApi
	tools     *tools.JwtTools
	yandexCfg *config.YandexMapkitConfig
}

func NewAuthService(cfg *config.Config, repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo, smtp: smtp.NewSMTPClient(cfg.SMTP), smsSender: nil,
		tools: tools.NewJwtTools(cfg.Jwt.SecretKey), yandexCfg: cfg.YandexMapkit}
}

func (s *AuthService) GetYandexAPIKey(ctx context.Context) (string, error) {
	return s.yandexCfg.APIKey, nil
}

func (s *AuthService) CheckUserExists(ctx context.Context, emailOrPhone string) (string, error) {
	return s.repo.CheckUserExists(ctx, emailOrPhone)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *AuthService) DeleteUserByID(ctx context.Context, userID string) error {
	return s.repo.DeleteUserByID(ctx, userID)
}
func (s *AuthService) SendCode(ctx context.Context, recipient string, code *models.AuthCode) error {
	newCode := tools.GenerateAuthCode()
	code.Code = newCode
	err := s.repo.CreateAuthCode(ctx, code)
	if err != nil {
		return err
	}
	fmt.Println("code", code)
	fmt.Println("channel", code.Channel)
	if code.Channel == "email" {
		err := s.smtp.SendVerificationCode(recipient, code.Code)
		if err != nil {
			return err
		}
	} else if code.Channel == "sms" {
		fmt.Println("sms")
	}

	return nil
}

func (s *AuthService) VerifyCode(ctx context.Context, recipient string, code string) (bool, error) {
	isValid, err := s.repo.VerifyAuthCode(ctx, recipient, code)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, fmt.Errorf("invalid code")
	}

	return true, nil
}

func (s *AuthService) ValidateJWTToken(ctx context.Context, token string) (string, error) {
	return s.tools.ValidateJWTToken(token)
}

func (s *AuthService) LogIn(ctx context.Context, user *models.User, code *models.AuthCode) (*models.Session, error) {
	isValid, err := s.VerifyCode(ctx, code.Recipient, code.Code)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("invalid code")
	}
	token, err := s.tools.GenerateJWTToken(user.UserID)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		UserID: user.UserID,
		Token:  token,
	}
	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}
	code.Used = true
	code.UserID = user.UserID
	err = s.repo.UpdateAuthCode(ctx, code)
	if err != nil {
		fmt.Println("ошибка обновления кода", err)
	}
	return session, nil
}

func (s *AuthService) Register(ctx context.Context, user *models.User, code *models.AuthCode) (*models.Session, error) {
	isValid, err := s.VerifyCode(ctx, code.Recipient, code.Code)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("invalid code")
	}
	token, err := s.tools.GenerateJWTToken(user.UserID)
	if err != nil {
		return nil, err
	}
	err = s.repo.CreateUser(ctx, user)
	if user.UserID == "" {
		s.repo.DeleteUserByID(ctx, user.UserID)
		return nil, fmt.Errorf("user not created")
	}
	if err != nil {
		s.repo.DeleteUserByID(ctx, user.UserID)
		return nil, err
	}

	session := &models.Session{
		UserID: user.UserID,
		Token:  token,
	}
	fmt.Println("session", session, "user", user)
	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		s.repo.DeleteUserByID(ctx, user.UserID)
		return nil, err
	}
	code.Used = true
	code.UserID = user.UserID
	err = s.repo.UpdateAuthCode(ctx, code)
	if err != nil {
		fmt.Println("ошибка обновления кода", err)
	}
	return session, nil
}

func (s *AuthService) LogOut(ctx context.Context, userID string) error {
	session, err := s.repo.GetSessionByUserID(ctx, userID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteSessionByID(ctx, session.SessionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error) {
	return s.repo.GetSessionByUserID(ctx, userID)
}
