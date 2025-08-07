package main

// cd C:\dev\projects\test_go\server\auth-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/pkg/database"
	"auth-service/internal/repositories"
	"auth-service/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetPostgresConfig()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	favouriteItemsRepository := repositories.NewFavouriteItemsRepository(db)
	favouriteItemsService := services.NewFavouriteItemsService(favouriteItemsRepository)
	favouriteItemsHandler := handlers.NewFavouriteItemsHandler(favouriteItemsService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)
	
	router := gin.Default()
	router.POST("/user/:user_id/favourites", favouriteItemsHandler.AddToFavourites)
	router.GET("/user/:user_id/favourites", favouriteItemsHandler.GetFavouritesByUserID)
	router.DELETE("/user/:user_id/favourites/:product_id", favouriteItemsHandler.DeleteFromFavourites)

	router.POST("/user/check/:email_or_phone", authHandler.CheckUserExists)
	router.GET("/user/:user_id", authHandler.GetUser)
	router.GET("/user/:user_id/session", authHandler.GetUserSession)
	router.POST("user/logout/:user_id", authHandler.LogOut)

	router.POST("auth/login", authHandler.LogIn)
	router.POST("auth/register", authHandler.Register)
	router.POST("auth/code/send", authHandler.SendCode)
	router.POST("auth/code/verify",authHandler.VerifyCode)
	// router.POST("auth/refresh", handlers.RefreshToken)
	
	router.POST("auth/check", authHandler.CheckToken)

	router.Run("0.0.0.0:8084")
	// mailHandler.SendMailCode("muhamedmakusev@gmail.com")
}

// var categories = []Category{
// 	// {ID: 0, Name: "Аксессуары", Image: "imgs/accessories.png"},
// 	// {ID: 1, Name: "Книги", Image: "imgs/books.png"},
// 	{ID: 2, Name: "Категория 2"},
// 	{ID: 2, Name: "Категория 3"},
// 	// {ID: 3, Name: "Одежда", Image: "imgs/clothes.png"},
// 	{ID: 4, Name: "Категория 1", Image: "imgs/food.png"},
// 	// {ID: 5, Name: "Электроника", Image: "imgs/electronic.png"},
// 	// {ID: 6, Name: "Распродажа", Image: "imgs/sale.png"},
// }
