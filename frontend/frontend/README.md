# Frontend - User Management System

Frontend React untuk aplikasi User Management yang terhubung dengan backend GoBack.

## Teknologi

- React 18
- Vite (Build tool)
- Axios (HTTP client)
- CSS3 (Styling)

## Struktur Folder

```
frontend/
├── src/
│   ├── components/
│   │   ├── UserForm.jsx      # Form create/edit user
│   │   └── UserList.jsx      # Tampilan daftar user
│   ├── services/
│   │   └── api.js            # API service layer
│   ├── App.jsx               # Main component
│   ├── App.css               # Styling aplikasi
│   ├── index.css             # Global styles
│   └── main.jsx              # Entry point
├── package.json
└── vite.config.js
```

## Instalasi

Dependencies sudah terinstall. Jika belum:

```bash
cd frontend
npm install
```

## Cara Menjalankan

### 1. Pastikan Backend Berjalan

Backend harus berjalan terlebih dahulu di `http://localhost:8080`

```bash
# Di root folder goback
go run cmd/api/main.go
```

### 2. Jalankan Frontend

```bash
cd frontend
npm run dev
```

Frontend akan berjalan di: `http://localhost:5173`

## Fitur

### 1. Create User
- Form di sebelah kiri untuk membuat user baru
- Validasi required untuk Name dan Email
- Email harus format valid

### 2. Read Users
- Table di sebelah kanan menampilkan semua user
- Auto refresh setelah create/update/delete
- Button Refresh manual

### 3. Update User
- Klik tombol "Edit" pada user
- Form akan terisi dengan data user
- Update dan kembali ke mode create

### 4. Delete User
- Klik tombol "Delete"
- Konfirmasi sebelum menghapus
- Auto refresh setelah delete

## API Endpoints yang Digunakan

Frontend berkomunikasi dengan backend melalui endpoints berikut:

- `GET /api/v1/users` - Ambil semua user
- `GET /api/v1/users/:id` - Ambil user by ID
- `POST /api/v1/users` - Buat user baru
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Hapus user

## Konfigurasi API

Base URL API dikonfigurasi di `src/services/api.js`:

```javascript
const API_BASE_URL = 'http://localhost:8080/api/v1';
```

Jika backend berjalan di port lain, ubah di file tersebut.

## Build untuk Production

```bash
npm run build
```

File hasil build akan ada di folder `dist/`
