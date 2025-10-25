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
  - Search & advanced filtering
  - Stock management
  - Image upload support
- **Shopping Cart**
  - Add/remove items
  - Update quantities
  - Real-time cart count
  - Glassmorphism cart drawer
- **Coupon/Promo System** 🆕
  - Create & manage discount coupons
  - Percentage or fixed amount discounts
  - Minimum purchase requirements
  - Max discount caps
  - Usage limits & tracking
  - Date-based validity
  - Apply coupons at checkout
  - Auto-calculate discounts
- **Order Management**
  - Checkout process with coupon support
  - Order history
  - Order status tracking
  - Order details view
  - Coupon discount tracking per order

### 👨‍💼 Admin Dashboard
- Product CRUD operations
- Category management
- Coupon management 🆕
- Order management
- Sales analytics
- Revenue tracking
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
│   ├── coupon_handler.go # Coupon management 🆕
│   ├── review_handler.go
│   ├── wishlist_handler.go
│   ├── profile_handler.go
│   ├── dashboard.go
│   └── user_handler.go
├── middleware/           # HTTP middlewares
│   ├── auth.go          # JWT middleware
│   ├── cors.go
│   └── logger.go
├── models/               # Data models & repositories
│   ├── user.go
│   ├── profile.go
│   ├── product.go
│   ├── category.go
│   ├── coupon.go        # Coupon model 🆕
│   ├── order.go
│   ├── order_item.go
│   ├── review.go
│   └── wishlist.go
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
│   │   │   ├── Checkout.jsx          # With coupon support 🆕
│   │   │   ├── Orders.jsx
│   │   │   ├── Profile.jsx
│   │   │   ├── Dashboard.jsx
│   │   │   ├── AdminDashboard.jsx
│   │   │   ├── AdminProducts.jsx
│   │   │   ├── AdminCategories.jsx
│   │   │   └── AdminCoupons.jsx      # Coupon management 🆕
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
GET    /api/v1/orders              - Get user orders (Protected)
GET    /api/v1/orders/:id          - Get order by ID (Protected)
POST   /api/v1/orders              - Create order with coupon (Protected) 🆕
PUT    /api/v1/orders/:id/status   - Update order status (Protected)
POST   /api/v1/orders/:id/cancel   - Cancel order (Protected)
```

### Coupons 🆕
```
# Public endpoints
POST   /api/v1/coupons/validate    - Validate coupon code
GET    /api/v1/coupons/active      - Get active coupons

# Admin endpoints (Protected)
POST   /api/v1/coupons             - Create coupon
GET    /api/v1/coupons             - Get all coupons
GET    /api/v1/coupons/:id         - Get coupon by ID
PUT    /api/v1/coupons/:id         - Update coupon
DELETE /api/v1/coupons/:id         - Delete coupon
POST   /api/v1/coupons/:id/toggle  - Toggle active status
```

### Reviews
```
GET    /api/v1/reviews/product/:id - Get product reviews
POST   /api/v1/reviews             - Create review (Protected)
GET    /api/v1/reviews/my          - Get my reviews (Protected)
PUT    /api/v1/reviews/:id         - Update review (Protected)
DELETE /api/v1/reviews/:id         - Delete review (Protected)
```

### Wishlist
```
POST   /api/v1/wishlist            - Add to wishlist (Protected)
GET    /api/v1/wishlist            - Get my wishlist (Protected)
DELETE /api/v1/wishlist            - Remove from wishlist (Protected)
POST   /api/v1/wishlist/check      - Check if in wishlist (Protected)
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
- **Apply discount coupons** 🆕
- Input shipping information
- **See discount breakdown** 🆕
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
- **Manage coupons** 🆕
  - Create coupons dengan percentage/fixed discount
  - Set validity period & usage limits
  - Track coupon usage statistics
  - Toggle active/inactive status
- View and update order status
- User management

## 🎟️ Coupon System Usage

### Admin - Create Coupon

1. **Login sebagai admin** dan navigate ke Admin → Manage Coupons
2. **Click "Create New Coupon"**
3. **Fill in coupon details:**
   ```
   Code: SAVE20
   Description: Save 20% on all products
   Discount Type: Percentage
   Discount Value: 20
   Min Purchase: 100000 (Rp 100,000)
   Max Discount: 50000 (Rp 50,000 - optional)
   Usage Limit: 100 (or 0 for unlimited)
   Start Date: 2025-01-01 00:00
   End Date: 2025-12-31 23:59
   Active: ✓
   ```
4. **Click "Create Coupon"**

### Customer - Apply Coupon at Checkout

1. **Add products to cart** (min Rp 100,000 untuk contoh di atas)
2. **Go to Checkout**
3. **In Order Summary section**, enter code: `SAVE20`
4. **Click "Apply"**
5. **See discount breakdown:**
   - Subtotal: Rp 150,000
   - Discount (SAVE20): -Rp 30,000
   - Total: Rp 120,000
6. **Complete checkout** - order akan tersimpan dengan informasi kupon

### Coupon Validation Rules

Sistem akan otomatis validate:
- ✅ Kode kupon valid & exists
- ✅ Kupon masih aktif (is_active = true)
- ✅ Tanggal masih dalam periode valid
- ✅ Order amount memenuhi minimum purchase
- ✅ Usage limit belum tercapai

### Example API Request

**Validate Coupon:**
```bash
curl -X POST http://localhost:8080/api/v1/coupons/validate \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SAVE20",
    "order_amount": 150000
  }'

# Response
{
  "valid": true,
  "coupon": {
    "id": 1,
    "code": "SAVE20",
    "discount_type": "percentage",
    "discount_value": 20,
    "min_purchase": 100000,
    "max_discount": 50000
  },
  "order_amount": 150000,
  "discount": 30000,
  "final_amount": 120000
}
```

**Create Order with Coupon:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address": "Jl. Contoh No. 123",
    "shipping_city": "Jakarta",
    "shipping_province": "DKI Jakarta",
    "shipping_postal_code": "12345",
    "payment_method": "bank_transfer",
    "coupon_code": "SAVE20",
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

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
\d coupons     # 🆕 New coupon table
\d orders
\d order_items
\d reviews
\d wishlists
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

- [x] Product reviews & ratings ✅
- [x] Wishlist functionality ✅
- [x] Advanced search & filters ✅
- [x] Admin analytics dashboard ✅
- [x] Coupon/Promo system ✅ 🆕
- [ ] Payment gateway integration (Midtrans/Stripe)
- [ ] Email notifications
- [ ] Product recommendations
- [ ] Multi-image upload for products
- [ ] Social login (Google, Facebook)
- [ ] Forgot password functionality
- [ ] Export orders to CSV/Excel
- [ ] Advanced coupon rules (first-time user, category-specific)
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
