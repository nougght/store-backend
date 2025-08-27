package repositories

import (
	"context"
	"log"
	"minio-client/internal/models"

	"github.com/jmoiron/sqlx"
)

type UrlsRepository struct {
	db *sqlx.DB
}

func NewUrlsRepository(db *sqlx.DB) *UrlsRepository {
	return &UrlsRepository{db: db}
}

func (r *UrlsRepository) IsUrlExist(ctx context.Context, objectName string) (bool, error) {
	var exists bool
	//  AND NOW() < expires_at)
	query := "SELECT EXISTS(SELECT 1 FROM minio.urls WHERE object_name = $1)"
	err := r.db.GetContext(ctx, &exists, query, objectName)
	if err != nil {
		return false, err
	}
	log.Println("exists", exists)
	return exists, nil
}

func (r *UrlsRepository) GetUrlByObjectName(ctx context.Context, objectName string) (*models.Url, error) {
	var url models.Url
	query := "SELECT * FROM minio.urls WHERE object_name = $1"
	err := r.db.GetContext(ctx, &url, query, objectName)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *UrlsRepository) CreateUrl(ctx context.Context, url *models.Url) error {
	query := `INSERT INTO minio.urls (object_name, bucket_name, url, expires_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowxContext(ctx, query, url.ObjectName, url.BucketName, url.Url, url.ExpiresAt).StructScan(url)
}

func (r *UrlsRepository) UpdateUrl(ctx context.Context, url *models.Url) error {
	query := `UPDATE minio.urls SET bucket_name = $1, url = $2, expires_at = $3 WHERE object_name = $4`
	_, err := r.db.ExecContext(ctx, query, url.BucketName, url.Url, url.ExpiresAt, url.ObjectName)
	return err
}

func (r *UrlsRepository) AddUrl(ctx context.Context, url *models.Url) error {
	query := `INSERT INTO minio.urls (object_name, bucket_name, url, expires_at) VALUES ($1, $2, $3, $4)
				ON CONFLICT(object_name)
				DO UPDATE SET
					url = EXCLUDED.url,
					expires_at = EXCLUDED.expires_at,
					bucket_name = EXCLUDED.bucket_name,`
	_, err := r.db.ExecContext(ctx, query, url.BucketName, url.Url, url.ExpiresAt, url.ObjectName)
	return err
}

func (r *UrlsRepository) DeleteUrl(ctx context.Context, objectName string) error {
	query := `DELETE FROM minio.urls WHERE object_name = $1`
	_, err := r.db.ExecContext(ctx, query, objectName)
	return err
}
