# 🛒 GoBack - E-Commerce Platform

Fullstack E-Commerce aplikasi dengan Golang Backend dan React Frontend yang modern dan responsive.

## ✨ Fitur Utama

### 🔐 Authentication & Authorization
- **JWT (JSON Web Token)** untuk autentikasi yang aman
- Register dan Login dengan password hashing (bcrypt)
- Protected routes dengan JWT middleware
- Token expiration (7 hari)
- Profile management

### 🛍️ E-Commerce Features
- **Product Management**
  - Browse products dengan category filtering
  - Product detail page
  - Stock management
  - Image upload support
- **Shopping Cart**
  - Add/remove items
  - Update quantities
  - Real-time cart count
  - Glassmorphism cart drawer
- **Order Management**
  - Checkout process
  - Order history
  - Order status tracking
  - Order details view

### 👨‍💼 Admin Dashboard
- Product CRUD operations
- Category management
- Order management
- Sales analytics
- User management

### 🎨 Modern UI/UX
- **Glassmorphism design** dengan backdrop blur
- **Animated gradients** dan smooth transitions
- **Floating particles** background effects
- **Micro-interactions** untuk better UX
- **Responsive design** untuk semua device
- **Dark mode ready** color palette

## 🏗️ Struktur Project

```
goback/
├── cmd/api/              # Backend entry point
│   └── main.go
├── config/               # Configuration management
│   └── config.go
├── handlers/             # HTTP request handlers
│   ├── auth_handler.go   # Authentication
│   ├── product_handler.go
│   ├── category_handler.go
│   ├── order_handler.go
│   └── user_handler.go
├── middleware/           # HTTP middlewares
│   ├── auth.go          # JWT middleware
│   ├── cors.go
│   └── logger.go
├── models/               # Data models & repositories
│   ├── user.go
│   ├── product.go
│   ├── category.go
│   └── order.go
├── routes/               # Route definitions
│   └── routes.go
├── frontend/             # React Frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── Navbar.jsx
│   │   │   ├── Cart.jsx
│   │   │   └── ...
│   │   ├── pages/
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   ├── Products.jsx
│   │   │   ├── ProductDetail.jsx
│   │   │   ├── Checkout.jsx
│   │   │   ├── Orders.jsx
│   │   │   ├── Profile.jsx
│   │   │   ├── Dashboard.jsx
│   │   │   ├── AdminDashboard.jsx
│   │   │   ├── AdminProducts.jsx
│   │   │   └── AdminCategories.jsx
│   │   ├── context/
│   │   │   ├── AuthContext.jsx
│   │   │   └── CartContext.jsx
│   │   ├── services/
│   │   │   └── api.js
│   │   ├── styles/
│   │   ├── App.jsx
│   │   └── main.jsx
│   └── package.json
├── .env.example
├── go.mod
└── README.md
```

## 📋 Prerequisites

- **Go 1.21+** atau lebih baru
- **Node.js 18+** dan npm
- **PostgreSQL 14+**
- Git

## 🚀 Instalasi & Setup

### 1. Clone Repository

```bash
git clone https://github.com/dinulimann/goback.git
cd goback
```

### 2. Setup PostgreSQL Database

```bash
# Login ke PostgreSQL
psql postgres

# Buat database baru
CREATE DATABASE goback_db;

# Keluar
\q
```

### 3. Setup Backend (Golang)

```bash
# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Edit .env dengan konfigurasi Anda
nano .env
```

**Konfigurasi .env:**
```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=goback_db
DB_SSLMODE=disable

# JWT Secret (WAJIB GANTI DI PRODUCTION!)
JWT_SECRET=your-super-secret-jwt-key-change-this
```

```bash
# Jalankan backend
go run main.go
```

Output yang diharapkan:
```
✓ Database connection established successfully
✓ Database migration completed
✓ Server starting on port 8080
```

### 4. Setup Frontend (React)

Buka terminal baru:

```bash
# Masuk ke folder frontend
cd frontend

# Install dependencies
npm install

# Jalankan development server
npm run dev
```

Frontend akan berjalan di `http://localhost:5173`

### 5. Akses Aplikasi

Buka browser: **http://localhost:5173**

**Default Admin Account** (jika belum ada, buat manual):
- Email: `admin@example.com`
- Password: `admin123`

## 📡 API Endpoints

### Authentication
```
POST   /api/auth/register       - Register user baru
POST   /api/auth/login          - Login user
GET    /api/auth/profile        - Get user profile (Protected)
```

### Products
```
GET    /api/products            - Get all products
GET    /api/products/:id        - Get product by ID
GET    /api/products/category/:id - Get products by category
POST   /api/products            - Create product (Protected)
PUT    /api/products/:id        - Update product (Protected)
DELETE /api/products/:id        - Delete product (Protected)
```

### Categories
```
GET    /api/categories          - Get all categories
GET    /api/categories/:id      - Get category by ID
POST   /api/categories          - Create category (Protected)
PUT    /api/categories/:id      - Update category (Protected)
DELETE /api/categories/:id      - Delete category (Protected)
```

### Orders
```
GET    /api/orders              - Get user orders (Protected)
GET    /api/orders/:id          - Get order by ID (Protected)
POST   /api/orders              - Create order (Protected)
PUT    /api/orders/:id/status   - Update order status (Protected)
```

### Users
```
GET    /api/users               - Get all users (Protected)
GET    /api/users/:id           - Get user by ID (Protected)
PUT    /api/users/:id           - Update user (Protected)
DELETE /api/users/:id           - Delete user (Protected)
```

## 🔒 Authentication Flow

### Register/Login
```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Response
{
  "message": "Registration successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### Protected Request
```bash
# Gunakan token di header Authorization
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## 🛠️ Technology Stack

### Backend
- **Go 1.24** - Programming language
- **Gin** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **JWT** - Authentication (golang-jwt/jwt/v5)
- **bcrypt** - Password hashing
- **godotenv** - Environment variables

### Frontend
- **React 19** - UI library
- **React Router DOM** - Routing
- **Vite** - Build tool & dev server
- **Axios** - HTTP client
- **Context API** - State management
- **CSS3** - Modern styling dengan glassmorphism

## 🎨 UI Features

### Design System
- **Glassmorphism** dengan backdrop-filter blur
- **Animated gradients** yang smooth
- **Floating particles** untuk background depth
- **Smooth transitions** dengan cubic-bezier
- **Micro-interactions** pada hover states
- **Color palette** yang modern dan consistent

### Animations
- Card hover effects dengan scale & translate
- Gradient text animations
- Loading spinners dengan multi-color
- Button ripple effects
- Floating icon animations
- Border glow animations
- Slide-in/fade-in transitions

## 📱 Features Walkthrough

### 1. Register & Login
- Buat akun baru atau login dengan akun existing
- JWT token disimpan di localStorage
- Auto redirect ke dashboard setelah login

### 2. Browse Products
- Lihat semua products di halaman Products
- Filter by category menggunakan sidebar
- View product details dengan klik product card
- Stock indicator (In Stock / Out of Stock)

### 3. Shopping Cart
- Add products ke cart dari product page
- View cart dengan klik cart button di navbar
- Update quantities atau remove items
- Real-time cart count badge

### 4. Checkout
- Review cart items
- Input shipping information
- Place order
- Redirect ke order confirmation

### 5. Order History
- View semua orders di halaman Orders
- See order status (Pending, Processing, Shipped, Delivered)
- View order details dengan item list

### 6. Admin Dashboard
- Analytics dashboard dengan sales metrics
- Manage products (CRUD operations)
- Manage categories
- View and update order status
- User management

## 🔧 Development Tips

### Hot Reload
- Backend: Gunakan `air` untuk hot reload
  ```bash
  go install github.com/cosmtrek/air@latest
  air
  ```
- Frontend: Vite sudah include hot reload otomatis

### Database Migrations
GORM auto migration akan membuat/update tables otomatis saat aplikasi start.

Manual migration jika diperlukan:
```bash
# Masuk ke psql
psql goback_db

# Lihat tables
\dt

# Lihat struktur table
\d users
\d products
\d categories
\d orders
```

### CORS Configuration
CORS sudah dikonfigurasi untuk development. Untuk production, update di `middleware/cors.go`

## 🐛 Troubleshooting

### Backend tidak start
- ✅ Pastikan PostgreSQL running
- ✅ Check credentials di .env benar
- ✅ Port 8080 tidak digunakan aplikasi lain
- ✅ Run `go mod tidy`

### Frontend tidak bisa connect ke backend
- ✅ Backend harus running di port 8080
- ✅ Check browser console (F12) untuk error
- ✅ Verify API base URL di `frontend/src/services/api.js`

### JWT Token expired
- Token valid 7 hari
- Logout dan login ulang untuk refresh token

### Database migration error
- Drop database dan recreate:
  ```sql
  DROP DATABASE goback_db;
  CREATE DATABASE goback_db;
  ```

## 🚀 Production Deployment

### Backend
1. Build binary:
   ```bash
   go build -o goback-server main.go
   ```

2. Set environment variables:
   ```bash
   export ENVIRONMENT=production
   export JWT_SECRET=your-production-secret
   export DB_HOST=production-db-host
   ```

3. Run:
   ```bash
   ./goback-server
   ```

### Frontend
1. Build untuk production:
   ```bash
   cd frontend
   npm run build
   ```

2. Serve static files dengan nginx/apache atau deploy ke:
   - Vercel
   - Netlify
   - AWS S3 + CloudFront
   - Firebase Hosting

### Docker (Coming Soon)
```bash
docker-compose up -d
```

## 🔜 Roadmap

- [ ] Payment gateway integration (Midtrans/Stripe)
- [ ] Product reviews & ratings
- [ ] Wishlist functionality
- [ ] Advanced search & filters
- [ ] Email notifications
- [ ] Product recommendations
- [ ] Admin analytics dashboard
- [ ] Multi-image upload for products
- [ ] Social login (Google, Facebook)
- [ ] Forgot password functionality
- [ ] Unit & integration testing
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] API documentation (Swagger)

## 📝 License

MIT License - feel free to use this project for learning or commercial purposes.

## 👨‍💻 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

For questions or suggestions, please open an issue or contact via email.

---

**Made with ❤️ using Go & React**
