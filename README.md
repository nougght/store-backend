# Бэкенд для онлайн магазина на Golang

Клиентское приложение (Flutter) - <https://github.com/nougght/store-client>

БД: PostgreSQL \
Хранилище для файлов: MinIO

### Указать в .env:
```
# Конфигурация для PostgreSQL
POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_SSLMODE=

# Конфигурация для MinIO
MINIO_ENDPOINT=
MINIO_ACCESS_KEY_ID=
MINIO_SECRET_KEY=
MINIO_BUCKET_NAME=
MINIO_USE_SSL=

# Конфигурация для SMTP (отправка писем для авторизации)
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=

# JWT секретный ключ
JWT_SECRET_KEY=

# API ключ для Yandex MapKit (запрашивается клиентом)
YANDEX_MAPKIT_API_KEY=
```
