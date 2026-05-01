package minio

import (
	"store-server/config"
	"store-server/internal/minio/client"
	"store-server/internal/minio/handlers"
	"store-server/internal/minio/repositories"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type MinioModule struct {
	minioClient *client.MinioClient
}

func NewMinioModule(cfg *config.MinioConfig, db *sqlx.DB) (*MinioModule, error) {
	urlsRepository := repositories.NewUrlsRepository(db)
	minioClient, err := client.NewMinioClient(cfg.Endpoint, cfg.AccessKeyID, cfg.SecretKey, cfg.BucketName, cfg.UseSSL, urlsRepository)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &MinioModule{minioClient: minioClient}, nil
}

func (m *MinioModule) RegisterRoutes(r *gin.Engine) {
	minioHandler := handlers.NewMinioHandler(m.minioClient)

	r.GET("/products/:id/images/:number/upload_url/:ext", minioHandler.GetProductUploadURL) // получить PUT presigned URL
	r.GET("/products/:id/images/:number", minioHandler.GetProductImageURL)                  // GET presigned URL для скачивания
	r.DELETE("/products/:id/images/:number/:ext", minioHandler.DeleteProductImage)
	r.GET("/products/:id/images", minioHandler.ListProductImages)

	// Категории
	r.GET("/categories/:id/image/upload_url/:ext", minioHandler.GetCategoryUploadURL)
	r.GET("/categories/:id/image", minioHandler.GetCategoryImageURL)
	r.DELETE("/categories/:id/image", minioHandler.DeleteCategoryImage)

}
