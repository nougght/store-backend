package main

// cd C:\dev\projects\test_go\server\product-service
// ngrok http 8080 --subdomain=tesvi
// Зе1е pqsl

import (

	"product-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/products", handlers.GetProducts)
	router.GET("/products:ids", handlers.GetProductsByIDs)
	router.Run("0.0.0.0:8081")
}
