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

Salin dan sesuaikan file `.env`:

```env
APP_PORT=8080
GIN_MODE=release

JWT_SECRET=your-super-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=catalogue
DB_SSLMODE=disable

REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
```

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
# Install dependencies
make download

# Run migration
make migrate

# (Optional) Seed data
make seed

# Start server
make server
```

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

## Swagger Docs

After starting the server, open:

```
http://localhost:8080/swagger/index.html
```

To regenerate docs after code changes:

```bash
make swag
```

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
