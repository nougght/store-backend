package main

import (
	"fmt"
	"store-server/config"
	"store-server/database"
	"store-server/internal/auth"
	"store-server/internal/cart"
	"store-server/internal/category"
	"store-server/internal/minio"
	"store-server/internal/od"
	"store-server/internal/product"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(".env file not found")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Configuration loading error: %v", err)
		return
	}

	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		fmt.Printf("Database initialization error: %v", err)
		return
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или "*" для всех
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.RedirectTrailingSlash = false

	authModule := auth.NewAuthModule(cfg, db)
	productModule := product.NewProductModule(db)
	cartModule := cart.NewCartModule(db)
	categoryModule := category.NewCategoryModule(db)
	odModule := od.NewODModule(db)
	minioModule, err := minio.NewMinioModule(cfg.Minio, db)
	if err != nil {
		fmt.Printf("Minio initialization error: %v", err)
		return
	}

	authModule.RegisterRoutes(r)
	productModule.RegisterRoutes(r)
	cartModule.RegisterRoutes(r)
	categoryModule.RegisterRoutes(r)
	minioModule.RegisterRoutes(r)
	odModule.RegisterRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	err = r.Run(":8080")
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
