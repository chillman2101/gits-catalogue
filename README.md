# Gits Catalogue API

## Tech Stack

- **Language**: Go 1.23
- **Framework**: Gin
- **ORM**: GORM + PostgreSQL
- **Auth**: JWT (access token 15 menit + refresh token 7 hari)
- **Cache**: Redis (opsional — app bisa berjalan tanpa redis)
- **Docs**: Swagger (swaggo)
- **Container**: Docker + Docker Compose

## Project Structure

```
.
├── cmd/
│   ├── app/        # Main entrypoint
│   ├── migrate/    # Database migration
│   └── seed/       # Database seeder
├── internal/
│   ├── config/     # App & DB configuration
│   ├── dto/        # Request & response DTOs
│   ├── handler/    # HTTP handlers
│   ├── middleware/  # JWT auth middleware
│   ├── mocks/      # Testify mocks
│   ├── model/      # GORM models
│   ├── query/      # Pagination & sorting
│   ├── repository/ # Database queries
│   ├── response/   # Response helpers & error types
│   ├── router/     # Route definitions
│   └── service/    # Business logic
└── pkg/
    └── redis/      # Redis client & cache helpers
```

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL 16+
- Redis 7+ (opsional)
- Docker & Docker Compose (opsional)

### Environment Variables

Salin `.env.example` lalu sesuaikan nilainya:

```bash
cp .env.example .env
```

| Variable           | Description                        | Default       |
|--------------------|------------------------------------|---------------|
| APP_PORT           | Port server                        | 8080          |
| GIN_MODE           | `debug` atau `release`             | release       |
| JWT_SECRET         | Secret key untuk access token      | -             |
| JWT_REFRESH_SECRET | Secret key untuk refresh token     | -             |
| DB_HOST            | PostgreSQL host                    | localhost     |
| DB_PORT            | PostgreSQL port                    | 5432          |
| DB_USER            | PostgreSQL user                    | postgres      |
| DB_PASSWORD        | PostgreSQL password                | -             |
| DB_NAME            | Nama database                      | catalogue     |
| DB_SSLMODE         | SSL mode (`disable`/`require`)     | disable       |
| REDIS_ADDR         | Redis address `host:port`          | 127.0.0.1:6379|
| REDIS_PASSWORD     | Redis password (kosongkan jika tidak ada) | -      |
| REDIS_DB           | Redis database index               | 0             |

### Run with Docker Compose

```bash
# Start app (includes auto-migrate)
docker compose up --build

# Start app with seed data
docker compose --profile seed up --build

# Stop
docker compose down
```

### Run Locally

```bash
# 1. Clone repo
git clone https://github.com/chillman2101/gits-catalogue.git
cd gits-catalogue

# 2. Copy dan isi environment variables
cp .env.example .env

# 3. Install dependencies
make download

# 4. Jalankan migration (auto-create database + tabel)
make migrate

# 5. (Opsional) Seed data awal
make seed

# 6. Start server
make server
```

Server berjalan di `http://localhost:8080`

## API Endpoints

Base URL: `http://localhost:8080`

### Auth

| Method | Endpoint         | Description              | Auth |
|--------|-----------------|--------------------------|------|
| POST   | /auth/register  | Register new user        | -    |
| POST   | /auth/login     | Login & get token        | -    |
| POST   | /auth/refresh   | Refresh access token     | JWT  |
| POST   | /auth/logout    | Logout & invalidate token| JWT  |

### Authors

| Method | Endpoint              | Description       | Auth |
|--------|-----------------------|-------------------|------|
| GET    | /api/v1/authors       | List all authors  | JWT  |
| GET    | /api/v1/authors/:id   | Get author by ID  | JWT  |
| POST   | /api/v1/authors       | Create author     | JWT  |
| PUT    | /api/v1/authors/:id   | Update author     | JWT  |
| DELETE | /api/v1/authors/:id   | Delete author     | JWT  |

### Publishers

| Method | Endpoint                 | Description          | Auth |
|--------|--------------------------|----------------------|------|
| GET    | /api/v1/publishers       | List all publishers  | JWT  |
| GET    | /api/v1/publishers/:id   | Get publisher by ID  | JWT  |
| POST   | /api/v1/publishers       | Create publisher     | JWT  |
| PUT    | /api/v1/publishers/:id   | Update publisher     | JWT  |
| DELETE | /api/v1/publishers/:id   | Delete publisher     | JWT  |

### Books

| Method | Endpoint            | Description     | Auth |
|--------|---------------------|-----------------|------|
| GET    | /api/v1/books       | List all books  | JWT  |
| GET    | /api/v1/books/:id   | Get book by ID  | JWT  |
| POST   | /api/v1/books       | Create book     | JWT  |
| PUT    | /api/v1/books/:id   | Update book     | JWT  |
| DELETE | /api/v1/books/:id   | Delete book     | JWT  |

### Query Parameters (list endpoints)

| Parameter | Type   | Default | Description                    |
|-----------|--------|---------|--------------------------------|
| page      | int    | 1       | Page number                    |
| limit     | int    | 10      | Items per page (max 100)       |
| search    | string | -       | Search by title, name, or ISBN |
| sort      | string | id      | Sort field                     |
| order     | string | ASC     | Sort order: ASC or DESC        |

### Response Format

All endpoints return a consistent JSON structure:

```json
{
  "success": true,
  "message": "books retrieved successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 42,
    "total_pages": 5
  }
}
```

Validation errors return:

```json
{
  "success": false,
  "message": "validation failed",
  "data": ["email must be a valid email address", "password is required"]
}
```

## API Documentation

### Swagger

Setelah server berjalan, buka di browser:

```
http://localhost:8080/swagger/index.html
```

Swagger menyediakan dokumentasi interaktif untuk semua endpoint — bisa langsung coba request dari browser.

Untuk regenerate docs setelah ada perubahan kode:

```bash
make swag
```

### Postman

Import file Swagger ke Postman:
1. Buka Postman → **Import**
2. Pilih **URL** → masukkan `http://localhost:8080/swagger/doc.json`
3. Postman akan generate semua collection secara otomatis

Atau download file JSON-nya langsung: `http://localhost:8080/swagger/doc.json`

## Testing

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage
make test-cover
```

## Entity Relationships

- **Author** has many **Books**
- **Book** belongs to one **Publisher**
- **Book** belongs to one **Author**
