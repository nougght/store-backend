package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"net/http"

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
	emailOrPhone := ctx.Request.FormValue("email_or_phone")
	exists, err := s.service.CheckUserExists(ctx, emailOrPhone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (s *AuthHandler) GetUser(ctx *gin.Context) {
	userID := ctx.Query("user_id")
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
	err = s.service.SendCode(ctx, code.recipient, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

func (s *AuthHandler) VerifyCode(ctx *gin.Context) {
	recipient := ctx.Request.FormValue("recipient")
	code := ctx.Request.FormValue("code")

	_, err := s.service.VerifyCode(ctx, recipient, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *AuthHandler) GetUserSession(ctx *gin.Context) {
	userID := ctx.Query("user_id")
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
	ctx.JSON(http.StatusOK, gin.H{"session": session, "user": user})
}

func (s *AuthHandler) LogOut(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	err := s.service.LogOut(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *AuthHandler) CheckToken(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id, err := s.tools.ValidateJWTToken(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": user_id})
}
