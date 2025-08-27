package client

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"minio-client/internal/models"
	"minio-client/internal/repositories"
)

// MinioClient обёртка
type MinioClient struct {
	Client     *minio.Client
	BucketName string
	UrlsRepo   *repositories.UrlsRepository
}

// NewMinioClient создаёт клиента MinIO
func NewMinioClient(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool, urlsRepo *repositories.UrlsRepository) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли бакет, если нет — создаём
	ctx := context.Background()
	exists, errBucketExists := client.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		return nil, errBucketExists
	}
	if !exists {
		if err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinioClient{Client: client, BucketName: bucketName, UrlsRepo: urlsRepo}, nil
}

// UploadImage загружает файл в MinIO
func (mc *MinioClient) UploadImage(ctx context.Context, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	_, err := mc.Client.PutObject(ctx, mc.BucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// GetPresignedURL генерирует ссылку для скачивания изображения (GET)
func (mc *MinioClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	isExist, err := mc.UrlsRepo.IsUrlExist(ctx, objectName)
	if err != nil {
		return "", err
	}
	var presignedURL *models.Url
	if isExist {
		presignedURL, err = mc.UrlsRepo.GetUrlByObjectName(ctx, objectName)
		if err != nil {
			return "", err
		}
	} else {
		reqParams := make(url.Values)
		url, err := mc.Client.PresignedGetObject(ctx, mc.BucketName, objectName, expiry, reqParams)
		presignedURL = &models.Url{
			ObjectName: objectName,
			BucketName: mc.BucketName,
			Url:        url.String(),
			ExpiresAt:  time.Now().Add(expiry),
		}
		log.Println("presignedURL", presignedURL)
		if err != nil {
			return "", err
		}
		err = mc.UrlsRepo.AddUrl(ctx, presignedURL)
		if err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}
	return presignedURL.Url, nil
}

// DeleteImage удаляет изображение из MinIO
func (mc *MinioClient) DeleteImage(ctx context.Context, objectName string) error {
	log.Println("objectName", objectName, mc.BucketName)
	_, err := mc.Client.StatObject(context.Background(), mc.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			log.Println("Объект не найден!")
		} else {
			log.Printf("Ошибка при проверке объекта: %v", err)
		}
		// return err
	}
	err = mc.Client.RemoveObject(ctx, mc.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// UploadViaPresignedPut генерирует presigned PUT URL и загружает файл туда через HTTP PUT
func (mc *MinioClient) PutPresignedUrl(ctx context.Context, bucket, objectName string, data []byte, contentType string, expiry time.Duration) (string, error) {
	// 1. Генерируем presigned PUT URL
	presignedURL, err := mc.Client.PresignedPutObject(ctx, bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return presignedURL.String(), nil
}

func (mc *MinioClient) RenameMinIOObject(bucket, oldName, newName string) error {
	// 1. Копируем объект
	_, err := mc.Client.CopyObject(
		context.Background(),
		minio.CopyDestOptions{
			Bucket: bucket,
			Object: newName,
		},
		minio.CopySrcOptions{
			Bucket: bucket,
			Object: oldName,
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка копирования: %v", err)
	}

	// 2. Удаляем исходный объект
	err = mc.Client.RemoveObject(
		context.Background(),
		bucket,
		oldName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %v", err)
	}

	return nil
}

func (mc *MinioClient) RenumberImages(bucket, productID string) error {
	// 1. Получаем список объектов
	objectsCh := mc.Client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Prefix:    fmt.Sprintf("products/%s/", productID),
		Recursive: true,
	})

	var objectKeys []string
	for object := range objectsCh {
		if object.Err != nil {
			log.Fatal(object.Err)
		}
		objectKeys = append(objectKeys, object.Key)
	}
	log.Println("objectKeys", objectKeys)
	// Сортируем по алфавиту
	sort.Strings(objectKeys)

	// 2. Обрабатываем каждый файл
	for i, key := range objectKeys {

		// Парсим номер из имени (например, "image_3.jpg" → 3)
		parts := strings.Split(key, "_")
		if len(parts) < 2 {
			continue
		}

		numExt := strings.Split(parts[1], ".")
		oldNum, _ := strconv.Atoi(numExt[0])
		ext := numExt[1]

		// Если номер > удаленного — уменьшаем его
		if oldNum != i {
			newNum := i
			newKey := fmt.Sprintf("products/%s_%d.%s", productID, newNum, ext)

			// Копируем с новым именем
			err := mc.RenameMinIOObject(bucket, key, newKey)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}
