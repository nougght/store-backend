package main

// cd C:\dev\projects\test_go\server\auth-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (
	"auth-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user/:user_id/favourites", func(c *gin.Context) {
		userId := c.Param("user_id")
		handlers.AddToFavourites(c, userId)

	})
	router.GET("/user/:user_id/favourites", func(c *gin.Context) {
		userId := c.Param("user_id")
		handlers.GetFavouritesByUserID(c, userId)
	})
	router.DELETE("/user/:user_id/favourites/:product_id", func(c *gin.Context) {
		userId := c.Param("user_id")
		productId := c.Param("product_id")
		handlers.DeleteFromFavourites(c, userId, productId)
	})

	router.Run("0.0.0.0:8084")
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
