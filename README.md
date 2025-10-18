# GoBack - Simple Backend API

Backend sederhana menggunakan Golang dengan Gin Framework.

## Fitur

- RESTful API untuk manajemen User (CRUD)
- Middleware untuk CORS dan Logging
- Configuration management dengan environment variables
- Struktur project yang terorganisir
- Health check endpoint

## Struktur Project

```
goback/
├── cmd/api/           # Entry point aplikasi
│   └── main.go
├── config/            # Configuration management
│   └── config.go
├── handlers/          # HTTP request handlers
│   ├── health_handler.go
│   └── user_handler.go
├── middleware/        # HTTP middlewares
│   ├── cors.go
│   └── logger.go
├── models/            # Data models
│   └── user.go
├── routes/            # Route definitions
│   └── routes.go
├── utils/             # Utility functions
├── .env.example       # Environment variables example
├── go.mod             # Go module file
└── README.md
```

## Prerequisites

- Go 1.21.6 atau lebih baru
- Git

## Instalasi

1. Clone repository ini
2. Copy file `.env.example` menjadi `.env`:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Jalankan aplikasi:
   ```bash
   go run cmd/api/main.go
   ```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/health` - Mengecek status server

### Users
- **POST** `/api/v1/users` - Membuat user baru
- **GET** `/api/v1/users` - Mendapatkan semua users
- **GET** `/api/v1/users/:id` - Mendapatkan user berdasarkan ID
- **PUT** `/api/v1/users/:id` - Update user
- **DELETE** `/api/v1/users/:id` - Hapus user

## Contoh Request

### Membuat User Baru
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### Mendapatkan Semua Users
```bash
curl http://localhost:8080/api/v1/users
```

### Mendapatkan User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Hapus User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Environment Variables

Lihat file `.env.example` untuk konfigurasi yang tersedia:

- `PORT` - Port server (default: 8080)
- `ENVIRONMENT` - Environment mode (development/production)
- `DB_*` - Database configuration (untuk future development)

## Development

Saat ini aplikasi menggunakan in-memory storage. Untuk production, sebaiknya integrasikan dengan database seperti PostgreSQL, MySQL, atau MongoDB.

## License

MIT