# Бэкенд для онлайн магазина на Golang

_Клиентское приложение (Flutter) - <https://github.com/nougght/store-client>_

## Технологии

- Gin для API
- PostgreSQL для хранения данных
- MinIO для хранения файлов (изображений)
- SMTP для авторизации по email

## Rest API

#### Auth

- POST /auth/login
- POST /auth/register
- POST /auth/code/send
- POST /auth/code/verify
- POST /auth/check

#### User

- GET /user/:user_id
- DELETE /user/:user_id
- GET /user/:user_id/session
- POST /user/logout/:user_id
- POST /user/check/:email_or_phone
- GET /user/:user_id/favourites
- POST /user/:user_id/favourites
- DELETE /user/:user_id/favourites/:product_id

#### Products

- GET /products
- POST /products
- PUT /products/:id
- DELETE /products/:id
- GET /products/:id/images
- GET /products/:id/images/:number/upload_url/:ext
- GET /products/:id/images/:number
- DELETE /products/:id/images/:number/:ext

#### Categories

- GET /categories
- POST /categories
- PUT /categories/:id
- DELETE /categories/:id
- GET /categories/:id/image/upload_url/:ext
- GET /categories/:id/image
- DELETE /categories/:id/image

#### Cart

- GET /cart/:user_id
- POST /cart/:user_id
- GET /cart/items/:cart_id
- POST /cart/items
- PATCH /cart/items
- DELETE /cart/items
- DELETE /cart/items/:id

#### Orders

- GET /orders
- POST /order
- GET /order/:id
- PUT /order/:id
- DELETE /order/:id
- GET /users/:user_id/orders
- GET /order/:id/items
- POST /order/items
- POST /order/:id/items
- GET /order/items/:id
- PUT /order/items/:id
- DELETE /order/items/:id

#### Delivery

- POST /delivery
- GET /delivery/:id
- GET /order/:id/delivery
- PUT /delivery/:id
- DELETE /delivery/:id

#### Other

- GET /yandex-map-key
- GET /health

## Указать в .env

``` .env
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
