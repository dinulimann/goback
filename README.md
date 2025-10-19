# GoBack - User Management System

Fullstack aplikasi User Management dengan Golang Backend dan React Frontend.

## Fitur

### Backend (Golang + Gin)
- RESTful API untuk manajemen User (CRUD)
- Middleware untuk CORS dan Logging
- Configuration management dengan environment variables
- Struktur project yang terorganisir
- Health check endpoint

### Frontend (React + Vite)
- User interface untuk manajemen user
- Create, Read, Update, Delete users
- Responsive design
- Real-time form validation
- Integration dengan backend API

## Struktur Project

```
goback/
├── cmd/api/           # Backend entry point
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
├── frontend/          # React Frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── UserForm.jsx
│   │   │   └── UserList.jsx
│   │   ├── services/
│   │   │   └── api.js
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── README.md
├── .env.example       # Environment variables example
├── go.mod             # Go module file
└── README.md
```

## Prerequisites

- Go 1.21.6 atau lebih baru
- Node.js 18+ dan npm
- PostgreSQL 14+
- Git

## Instalasi & Menjalankan Aplikasi

### 0. Setup PostgreSQL (WAJIB)

Aplikasi sekarang menggunakan PostgreSQL untuk data persistence.

**Lihat panduan lengkap:** [POSTGRESQL_SETUP.md](POSTGRESQL_SETUP.md)

**Quick Setup:**
```bash
# 1. Install PostgreSQL (sesuai OS Anda)

# 2. Buat database
psql postgres
CREATE DATABASE goback_db;
\q

# 3. Copy dan edit .env
cp .env.example .env
# Edit .env dengan password postgres Anda
```

### 1. Backend (Golang)

```bash
# Install dependencies
go mod tidy

# Setup environment variables
cp .env.example .env
# Edit .env dengan konfigurasi PostgreSQL Anda

# Jalankan backend server
go run cmd/api/main.go
```

Output yang diharapkan:
```
Database connection established successfully
Database migration completed
Server starting on port 8080
```

Backend akan berjalan di `http://localhost:8080`

### 2. Frontend (React)

Buka terminal baru:

```bash
# Masuk ke folder frontend
cd frontend

# Install dependencies (jika belum)
npm install

# Jalankan development server
npm run dev
```

Frontend akan berjalan di `http://localhost:5173`

### 3. Akses Aplikasi

Buka browser dan kunjungi: `http://localhost:5173`

**PENTING:** Pastikan backend sudah berjalan terlebih dahulu sebelum membuka frontend!

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

## Teknologi Stack

### Backend
- **Golang 1.21.6** - Bahasa pemrograman
- **Gin Framework** - Web framework
- **godotenv** - Environment variable management

### Frontend
- **React 18** - UI library
- **Vite** - Build tool & dev server
- **Axios** - HTTP client
- **CSS3** - Styling

## Development Notes

- Backend menggunakan **PostgreSQL database** untuk data persistence
- Data tersimpan **PERMANEN** (tidak hilang saat restart)
- GORM Auto Migration membuat/update table otomatis
- CORS sudah dikonfigurasi untuk development
- Hot reload aktif di frontend dan backend

## Cara Menggunakan Aplikasi

1. **Membuat User Baru**
   - Isi form "Name" dan "Email" di sebelah kiri
   - Klik tombol "Create User"
   - User akan muncul di table

2. **Melihat Daftar User**
   - Semua user ditampilkan di table sebelah kanan
   - Klik "Refresh" untuk reload data

3. **Edit User**
   - Klik tombol "Edit" pada user yang ingin diubah
   - Form akan terisi otomatis dengan data user
   - Ubah data, lalu klik "Update User"

4. **Hapus User**
   - Klik tombol "Delete" pada user
   - Konfirmasi penghapusan
   - User akan terhapus dari list

## Troubleshooting

### Backend error saat start
- Pastikan port 8080 tidak dipakai aplikasi lain
- Cek `go mod tidy` sudah dijalankan

### Frontend tidak bisa fetch data
- Pastikan backend berjalan di `http://localhost:8080`
- Buka browser console (F12) untuk lihat error
- Cek CORS middleware aktif di backend

### Port conflict
- Backend default: 8080 (ubah di .env → PORT)
- Frontend default: 5173 (ubah di vite.config.js)

## Next Steps

Fitur yang bisa ditambahkan:
- Database integration (PostgreSQL/MySQL)
- Authentication & Authorization (JWT)
- Pagination
- Search & Filter
- Validation yang lebih kompleks
- Unit testing
- Docker containerization
- CI/CD pipeline

## License

MIT