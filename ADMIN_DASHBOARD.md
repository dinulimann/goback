# Admin Dashboard Feature

## Overview
Admin Dashboard adalah fitur analytics dan monitoring untuk admin e-commerce yang menampilkan statistik penting seperti revenue, orders, products, dan users.

## Features

### 1. Statistics Cards
- **Total Revenue**: Total pendapatan dari semua order dengan status "paid"
- **Total Orders**: Jumlah total order dengan informasi pending orders
- **Total Products**: Jumlah produk aktif di sistem
- **Total Users**: Jumlah user yang terdaftar

### 2. Revenue Overview Chart
- Visualisasi revenue dalam bentuk bar chart
- Filter berdasarkan periode:
  - Daily (30 hari terakhir)
  - Weekly (12 minggu terakhir)
  - Monthly (12 bulan terakhir)
  - Yearly (5 tahun terakhir)
- Interactive hover untuk detail

### 3. Top Selling Products
- Menampilkan 5 produk terlaris
- Informasi:
  - Ranking
  - Nama produk
  - Jumlah unit terjual
  - Total revenue dari produk

### 4. Orders by Status
- Breakdown jumlah order berdasarkan status
- Status yang tersedia:
  - Pending
  - Processing
  - Shipped
  - Delivered
  - Cancelled

### 5. Monthly Revenue
- Revenue 6 bulan terakhir
- Informasi per bulan:
  - Nama bulan dan tahun
  - Jumlah orders
  - Total revenue

### 6. Recent Orders Table
- 10 order terbaru
- Informasi yang ditampilkan:
  - Order Number
  - Customer Name
  - Total Amount
  - Order Status
  - Created Date
- Button untuk melihat semua orders

## API Endpoints

### Get Dashboard Statistics
```
GET /api/v1/dashboard/stats
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_revenue": 10000000,
    "total_orders": 150,
    "total_products": 50,
    "total_users": 100,
    "pending_orders": 10,
    "recent_orders": [...],
    "top_products": [...],
    "revenue_by_month": [...],
    "orders_by_status": [...]
  }
}
```

### Get Revenue Statistics
```
GET /api/v1/dashboard/revenue?period=month
Authorization: Bearer <token>
```

**Query Parameters:**
- `period`: day | week | month | year (default: month)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "date": "2025-01",
      "revenue": 5000000,
      "orders": 50
    }
  ]
}
```

## Frontend Routes

### Admin Dashboard Page
```
/admin/dashboard
```

**Requirements:**
- User harus sudah login
- Protected route (requires authentication)

## Access via Navigation

1. Login ke aplikasi
2. Klik menu **Admin** di navbar
3. Pilih **Admin Dashboard** dari dropdown menu
4. Dashboard akan ditampilkan dengan semua statistik

## Dropdown Menu Admin

Navbar sekarang memiliki dropdown menu **Admin** dengan 3 pilihan:
- **Admin Dashboard** - Analytics dan statistik
- **Manage Products** - Kelola produk
- **Manage Categories** - Kelola kategori

## Technical Details

### Backend Files
- `handlers/dashboard.go` - Handler untuk dashboard endpoints
- `routes/routes.go` - Route definitions (lines 140-145)

### Frontend Files
- `frontend/src/pages/AdminDashboard.jsx` - Dashboard page component
- `frontend/src/styles/AdminDashboard.css` - Dashboard styling
- `frontend/src/components/Navbar.jsx` - Navigation with admin dropdown
- `frontend/src/components/Navbar.css` - Navbar styling with dropdown
- `frontend/src/App.jsx` - Route configuration (lines 157-163)

### Database Queries
Dashboard menggunakan query yang optimal untuk performa:
- Aggregate functions untuk statistik
- JOIN untuk relasi data
- Preload untuk eager loading
- Indexing pada field yang sering di-query

### Styling Features
- Responsive design (mobile-friendly)
- Modern gradient colors
- Smooth animations
- Interactive hover effects
- Clean card-based layout
- Professional color scheme

## How to Run

### Backend
```bash
cd /Users/monang/Documents/Projek/belajar/goback
go run cmd/api/main.go
```

Server akan berjalan di: `http://localhost:8080`

### Frontend
```bash
cd /Users/monang/Documents/Projek/belajar/goback/frontend
npm run dev
```

Frontend akan berjalan di: `http://localhost:5173`

## Testing

1. Pastikan backend dan frontend berjalan
2. Login dengan user credentials
3. Akses `/admin/dashboard`
4. Verify semua statistik muncul dengan benar
5. Test filter revenue chart (Daily, Weekly, Monthly, Yearly)
6. Hover pada chart bars untuk melihat detail
7. Check responsive design di mobile view

## Future Enhancements

Fitur yang bisa ditambahkan di masa depan:
- Real-time updates dengan WebSocket
- Export data to PDF/Excel
- Date range picker untuk custom period
- More detailed charts (pie chart, line chart)
- User analytics (new users per day/month)
- Product inventory alerts
- Sales forecasting
- Revenue comparison (current vs previous period)
- Admin notifications
- Advanced filters

## Notes

- Dashboard hanya menampilkan revenue dari order dengan `payment_status = 'paid'`
- Top products diurutkan berdasarkan total quantity sold
- Recent orders diurutkan dari yang terbaru
- Semua endpoint memerlukan JWT authentication
- Currency format menggunakan Indonesian Rupiah (IDR)
- Date format menggunakan locale Indonesia

## Screenshot Locations

Dashboard terdiri dari:
1. Stats Cards (Top section - 4 cards)
2. Revenue Chart & Top Products (Middle section - 2 columns)
3. Orders by Status & Monthly Revenue (Below middle - 2 columns)
4. Recent Orders Table (Bottom section - full width)

Semua section responsive dan akan menyesuaikan layout di layar kecil.
