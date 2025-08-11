package main

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	productServiceURL := "http://localhost:8081"
	cartServiceURL := "http://localhost:8082"
	categoriesServiceURL := "http://localhost:8083"
	authServiceURL := "http://localhost:8084"

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Или конкретный URL фронтенда
		AllowMethods: []string{"GET", "POST", "DELETE", "PATCH"},
		AllowHeaders: []string{"Content-Type"},
	}))

	r.RedirectTrailingSlash = false // Отключаем авторедиректы

	r.POST("/auth/check", reverseProxy(authServiceURL))

	r.POST("/user/check/:email_or_phone", reverseProxy(authServiceURL))
	r.GET("/user/:user_id", reverseProxy(authServiceURL))
	r.DELETE("/user/:user_id", reverseProxy(authServiceURL))
	r.GET("/user/:user_id/session", reverseProxy(authServiceURL))
	r.POST("/user/logout/:user_id", reverseProxy(authServiceURL))

	r.POST("/auth/login", reverseProxy(authServiceURL))
	r.POST("/auth/register", reverseProxy(authServiceURL))
	r.POST("/auth/code/send", reverseProxy(authServiceURL))
	r.POST("/auth/code/verify", reverseProxy(authServiceURL))

	r.GET("/cart/:user_id", reverseProxy(cartServiceURL))
	r.POST("/cart/:user_id", reverseProxy(cartServiceURL))
	r.GET("/cart/items/:cart_id", reverseProxy(cartServiceURL))
	r.POST("/cart/items", reverseProxy(cartServiceURL))
	r.PATCH("/cart/items", reverseProxy(cartServiceURL))
	r.DELETE("/cart/items", reverseProxy(cartServiceURL))
	r.DELETE("/cart/items/:id", reverseProxy(cartServiceURL))

	// товары
	r.GET("/products", reverseProxy(productServiceURL))
	r.GET("/products:ids", reverseProxy(productServiceURL))
	r.POST("/products", reverseProxy(productServiceURL))
	r.PUT("/products/:id", reverseProxy(productServiceURL))
	r.DELETE("/products/:id", reverseProxy(productServiceURL))


	// категории
	r.GET("/categories", reverseProxy(categoriesServiceURL))
	r.POST("/categories", reverseProxy(categoriesServiceURL))
	r.PUT("/categories/:id", reverseProxy(categoriesServiceURL))
	r.DELETE("/categories/:id", reverseProxy(categoriesServiceURL))

	

	r.GET("/user/:user_id/favourites", reverseProxy(authServiceURL))
	r.POST("/user/:user_id/favourites", reverseProxy(authServiceURL))
	r.DELETE("/user/:user_id/favourites/:product_id", reverseProxy(authServiceURL))

	log.Println("API Gateway запущен на :8080")
	r.Run(":8080")
}

func reverseProxy(target string) gin.HandlerFunc {
	remote, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Блокируем редиректы от сервисов
	// proxy.ModifyResponse = func(r *http.Response) error {
	// 	if r.StatusCode >= 300 && r.StatusCode < 400 {
	// 		return fmt.Errorf("редиректы запрещены")
	// 	}
	// 	return nil
	// }

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
