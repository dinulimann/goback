# Setup PostgreSQL untuk GoBack

## Instalasi PostgreSQL

### MacOS
```bash
# Install PostgreSQL menggunakan Homebrew
brew install postgresql@15

# Start PostgreSQL service
brew services start postgresql@15

# Atau start manual
pg_ctl -D /opt/homebrew/var/postgresql@15 start
```

### Linux (Ubuntu/Debian)
```bash
# Update package list
sudo apt update

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### Windows
1. Download PostgreSQL installer dari https://www.postgresql.org/download/windows/
2. Jalankan installer
3. Ikuti wizard instalasi
4. Catat password yang Anda buat untuk user postgres

## Setup Database

### 1. Akses PostgreSQL
```bash
# Login sebagai user postgres
psql postgres

# Atau di Linux
sudo -u postgres psql
```

### 2. Buat Database
```sql
-- Buat database baru
CREATE DATABASE goback_db;

-- Buat user baru (opsional)
CREATE USER goback_user WITH PASSWORD 'your_secure_password';

-- Berikan akses ke database
GRANT ALL PRIVILEGES ON DATABASE goback_db TO goback_user;

-- Keluar dari psql
\q
```

### 3. Konfigurasi .env

Edit file `.env` di root project:

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# PostgreSQL Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres              # atau goback_user jika buat user baru
DB_PASSWORD=your_password     # password postgres atau user baru
DB_NAME=goback_db
```

## Menjalankan Aplikasi

```bash
# Jalankan backend
go run cmd/api/main.go
```

Output yang diharapkan:
```
2024/01/01 10:00:00 No .env file found, using system environment variables
2024/01/01 10:00:00 Database connection established successfully
2024/01/01 10:00:00 Database migration completed
2024/01/01 10:00:00 Server starting on port 8080
```

## Verifikasi Database

Cek apakah table sudah dibuat:

```bash
# Login ke database
psql -U postgres -d goback_db

# Atau jika pakai user baru
psql -U goback_user -d goback_db
```

```sql
-- Lihat semua tables
\dt

-- Lihat struktur table users
\d users

-- Lihat data users
SELECT * FROM users;

-- Keluar
\q
```

## Troubleshooting

### Error: "password authentication failed"
- Pastikan password di `.env` sesuai dengan password postgres
- Reset password postgres:
  ```bash
  sudo -u postgres psql
  ALTER USER postgres PASSWORD 'new_password';
  ```

### Error: "database does not exist"
- Pastikan sudah create database `goback_db`
- Cek nama database di `.env` sudah benar

### Error: "connection refused"
- Pastikan PostgreSQL service sudah running:
  ```bash
  # MacOS
  brew services list
  
  # Linux
  sudo systemctl status postgresql
  ```

### Port 5432 sudah dipakai
- Ubah port PostgreSQL atau ubah DB_PORT di `.env`

## Perintah PostgreSQL Berguna

```bash
# Cek status PostgreSQL
brew services list                    # MacOS
sudo systemctl status postgresql      # Linux

# Start PostgreSQL
brew services start postgresql@15     # MacOS
sudo systemctl start postgresql       # Linux

# Stop PostgreSQL
brew services stop postgresql@15      # MacOS
sudo systemctl stop postgresql        # Linux

# Restart PostgreSQL
brew services restart postgresql@15   # MacOS
sudo systemctl restart postgresql     # Linux
```

## Migration

Aplikasi menggunakan GORM Auto Migration, yang akan:
- Membuat table `users` jika belum ada
- Update schema jika ada perubahan pada model
- **TIDAK akan menghapus data existing**

Schema yang dibuat:
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

## Data Persistence

Sekarang data akan tersimpan PERMANEN di PostgreSQL:
- ✅ Data tetap ada setelah restart server
- ✅ Data tetap ada setelah reboot sistem
- ✅ Bisa di-backup dan di-restore
- ✅ Production-ready
