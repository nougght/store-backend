package config

import (
	"fmt"
	"os"
)

type JwtConfig struct {
	SecretKey string
}

type YandexMapkitConfig struct {
	APIKey string
}
type MinioConfig struct {
	Endpoint    string
	AccessKeyID string
	SecretKey   string
	BucketName  string
	UseSSL      bool
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type Config struct {
	Jwt          *JwtConfig
	Postgres     *PostgresConfig
	Minio        *MinioConfig
	SMTP         *SMTPConfig
	YandexMapkit *YandexMapkitConfig
}

func LoadConfig() (*Config, error) {
	jwtSecretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is not set")
	}

	postgresHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return nil, fmt.Errorf("POSTGRES_HOST environment variable is not set")
	}
	postgresPort, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, fmt.Errorf("POSTGRES_PORT environment variable is not set")
	}
	postgresUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, fmt.Errorf("POSTGRES_USER environment variable is not set")
	}
	postgresPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("POSTGRES_PASSWORD environment variable is not set")
	}
	postgresDB, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return nil, fmt.Errorf("POSTGRES_DB environment variable is not set")
	}
	postgresSSLMode, ok := os.LookupEnv("POSTGRES_SSLMODE")
	if !ok {
		return nil, fmt.Errorf("POSTGRES_SSLMODE environment variable is not set")
	}

	minioEndpoint, ok := os.LookupEnv("MINIO_ENDPOINT")
	if !ok {
		return nil, fmt.Errorf("MINIO_ENDPOINT environment variable is not set")
	}
	minioAccessKeyID, ok := os.LookupEnv("MINIO_ACCESS_KEY_ID")
	if !ok {
		return nil, fmt.Errorf("MINIO_ACCESS_KEY_ID environment variable is not set")
	}
	minioSecretAccessKey, ok := os.LookupEnv("MINIO_SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("MINIO_SECRET_KEY environment variable is not set")
	}
	minioBucketName, ok := os.LookupEnv("MINIO_BUCKET_NAME")
	if !ok {
		return nil, fmt.Errorf("MINIO_BUCKET_NAME environment variable is not set")
	}
	minioUseSSL, _ := os.LookupEnv("MINIO_USE_SSL")

	smtpHost, ok := os.LookupEnv("SMTP_HOST")
	if !ok {
		return nil, fmt.Errorf("SMTP_HOST environment variable is not set")
	}
	smtpPort, ok := os.LookupEnv("SMTP_PORT")
	if !ok {
		return nil, fmt.Errorf("SMTP_PORT environment variable is not set")
	}
	smtpUsername, ok := os.LookupEnv("SMTP_USERNAME")
	if !ok {
		return nil, fmt.Errorf("SMTP_USERNAME environment variable is not set")
	}
	smtpPassword, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("SMTP_PASSWORD environment variable is not set")
	}
	smtpFrom, ok := os.LookupEnv("SMTP_FROM")
	if !ok {
		return nil, fmt.Errorf("SMTP_FROM environment variable is not set")
	}
	yandexMapkitAPIKey, ok := os.LookupEnv("YANDEX_MAPKIT_API_KEY")
	if !ok {
		return nil, fmt.Errorf("YANDEX_MAPKIT_API_KEY environment variable is not set")
	}

	return &Config{
		Jwt: &JwtConfig{
			SecretKey: jwtSecretKey,
		},
		YandexMapkit: &YandexMapkitConfig{
			APIKey: yandexMapkitAPIKey,
		},
		Postgres: &PostgresConfig{
			Host:     postgresHost,
			Port:     postgresPort,
			User:     postgresUser,
			Password: postgresPassword,
			DBName:   postgresDB,
			SSLMode:  postgresSSLMode,
		},
		Minio: &MinioConfig{
			Endpoint:    minioEndpoint,
			AccessKeyID: minioAccessKeyID,
			SecretKey:   minioSecretAccessKey,
			BucketName:  minioBucketName,
			UseSSL:      minioUseSSL == "true",
		},
		SMTP: &SMTPConfig{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUsername,
			Password: smtpPassword,
			From:     smtpFrom,
		},
	}, nil
}
