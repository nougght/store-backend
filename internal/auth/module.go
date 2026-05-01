package auth

import (
	"store-server/config"
	"store-server/internal/auth/handlers"
	"store-server/internal/auth/repositories"
	"store-server/internal/auth/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthModule struct {
	AuthService     *services.AuthService
	FavItemsService *services.FavouriteItemsService

	authRepo     *repositories.AuthRepository
	favItemsRepo *repositories.FavouriteItemsRepository
}

func NewAuthModule(cfg *config.Config, db *sqlx.DB) *AuthModule {
	authRepo := repositories.NewAuthRepository(db)
	favItemsRepo := repositories.NewFavouriteItemsRepository(db)

	return &AuthModule{
		authRepo:        authRepo,
		favItemsRepo:    favItemsRepo,
		AuthService:     services.NewAuthService(cfg, authRepo),
		FavItemsService: services.NewFavouriteItemsService(favItemsRepo),
	}
}

func (m *AuthModule) RegisterRoutes(r *gin.Engine) {
	authHandler := handlers.NewAuthHandler(m.AuthService)
	favItemsHandler := handlers.NewFavouriteItemsHandler(m.FavItemsService)

	r.GET("/user/:user_id/favourites", favItemsHandler.GetFavouritesByUserID)
	r.POST("/user/:user_id/favourites", favItemsHandler.AddToFavourites)
	r.DELETE("/user/:user_id/favourites/:product_id", favItemsHandler.DeleteFromFavourites)

	r.POST("/user/check/:email_or_phone", authHandler.CheckUserExists)
	r.DELETE("/user/:user_id", authHandler.DeleteUserByID)
	r.GET("/user/:user_id", authHandler.GetUser)
	r.GET("/user/:user_id/session", authHandler.GetUserSession)
	r.POST("/user/logout/:user_id", authHandler.LogOut)

	r.POST("/auth/login", authHandler.LogIn)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/code/send", authHandler.SendCode)
	r.POST("/auth/code/verify", authHandler.VerifyCode)
	// router.POST("auth/refresh", handlers.RefreshToken)
	r.POST("/auth/check", authHandler.CheckToken)
	r.GET("/yandex-map-key", authHandler.GetYandexAPIKey)
}
