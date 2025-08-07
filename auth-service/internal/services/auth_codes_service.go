package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"context"
	"fmt"
)

type AuthService struct {
	repo      *repositories.AuthRepository
	emailSMTP *models.EmailSMTP
	smsSender *models.SmsApi
	tools     *models.Tools
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo, emailSMTP: &models.EmailSMTP{}, smsSender: nil, tools: &models.Tools{}}
}

func (s *AuthService) CheckUserExists(ctx context.Context, emailOrPhone string) (bool, error) {
	return s.repo.CheckUserExists(ctx, emailOrPhone)
}

func (s *AuthService) SendCode(ctx context.Context, recipient string, code *models.AuthCode) error {
	newCode := s.tools.GenerateAuthCode()
	code.Code = newCode
	err := s.repo.CreateAuthCode(ctx, code)
	if err != nil {
		return err
	}

	if code.Channel == "email" {
		err := s.emailSMTP.SendVerificationCode(recipient, code.Code)
		if err != nil {
			return err
		}
	} else if code.Channel == "sms" {
		fmt.Println("sms")
	}

	return nil
}

func (s *AuthService) VerifyCode(ctx context.Context, recipient string, code string) error {
	isValid, err := s.repo.VerifyAuthCode(ctx, recipient, code)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("invalid code")
	}

	return nil
}

func (s *AuthService) LogIn(ctx context.Context, user *models.User, code *models.AuthCode) (string, error) {

	token, err := s.tools.GenerateJWTToken(user.UserID)
	if err != nil {
		return "", err
	}
	session := &models.Session{
		UserID: user.UserID,
		Token:  token,
	}
	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return "", err
	}

	s.repo.UpdateAuthCode(ctx, code)
	return token, nil
}

func (s *AuthService) Register(ctx context.Context, user *models.User, code *models.AuthCode) (string, error) {
	token, err := s.tools.GenerateJWTToken(user.UserID)
	if err != nil {
		return "", err
	}
	s.repo.CreateUser(ctx, user)

	session := &models.Session{
		UserID: user.UserID,
		Token:  token,
	}

	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return "", err
	}

	s.repo.UpdateAuthCode(ctx, code)
	return token, nil
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