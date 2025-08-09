package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"fmt"
	"net/http"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
	tools   *models.Tools
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service, tools: &models.Tools{}}
}

func (s *AuthHandler) CheckUserExists(ctx *gin.Context) {
	emailOrPhone := ctx.Param("email_or_phone")
	if emailOrPhone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email_or_phone is required"})
		return
	}
	fmt.Println(emailOrPhone)
	userID, err := s.service.CheckUserExists(ctx, emailOrPhone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (s *AuthHandler) GetUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	user, err := s.service.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *AuthHandler) SendCode(ctx *gin.Context) {
	var code models.AuthCode
	err := ctx.ShouldBindJSON(&code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code.ExpiresAt = time.Now().Add(time.Minute * 5)
	err = s.service.SendCode(ctx, code.Recipient, &code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(code)
	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

func (s *AuthHandler) VerifyCode(ctx *gin.Context) {
	var input struct {
		Recipient string `json:"recipient"`
		Code      string `json:"code"`
	}

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipient := input.Recipient
	code := input.Code
	_, err = s.service.VerifyCode(ctx, recipient, code)
	fmt.Println(recipient, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *AuthHandler) GetUserSession(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	session, err := s.service.GetSessionByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"session": session})
}

func (s *AuthHandler) Register(ctx *gin.Context) {
	var input struct {
		User models.User     `json:"user"`
		Code models.AuthCode `json:"code"`
	}

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := input.User
	code := input.Code

	session, err := s.service.Register(ctx, &user, &code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"session": session, "user": user})
}

func (s *AuthHandler) LogIn(ctx *gin.Context) {
	var input struct {
		User models.User     `json:"user"`
		Code models.AuthCode `json:"code"`
	}

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := input.User
	code := input.Code

	session, err := s.service.LogIn(ctx, &user, &code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userP, err := s.service.GetUserByID(ctx, user.UserID)
	user = *userP
	if err != nil {
		fmt.Println("ошибка получения пользователя", err)
	} else {
		fmt.Println("user", user)
	}
	ctx.JSON(http.StatusOK, gin.H{"session": session, "user": user})
}

func (s *AuthHandler) LogOut(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	err := s.service.LogOut(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *AuthHandler) DeleteUserByID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	err := s.service.DeleteUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + userID})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *AuthHandler) CheckToken(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	fmt.Println(tokenString)
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	user_id, err := s.tools.ValidateJWTToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": user_id})
}
