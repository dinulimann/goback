package handlers

import (
	"net/http"
	"time"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

type DashboardStats struct {
	TotalRevenue    float64              `json:"total_revenue"`
	TotalOrders     int64                `json:"total_orders"`
	TotalProducts   int64                `json:"total_products"`
	TotalUsers      int64                `json:"total_users"`
	PendingOrders   int64                `json:"pending_orders"`
	RecentOrders    []OrderSummary       `json:"recent_orders"`
	TopProducts     []TopProduct         `json:"top_products"`
	RevenueByMonth  []MonthlyRevenue     `json:"revenue_by_month"`
	OrdersByStatus  []OrderStatusCount   `json:"orders_by_status"`
}

type OrderSummary struct {
	ID           uint                `json:"id"`
	OrderNumber  string              `json:"order_number"`
	UserName     string              `json:"user_name"`
	TotalAmount  float64             `json:"total_amount"`
	Status       models.OrderStatus  `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
}

type TopProduct struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	TotalSold    int64   `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Year    int     `json:"year"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type OrderStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// Get total revenue
	h.db.Model(&models.Order{}).
		Where("payment_status = ?", "paid").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	// Get total orders
	h.db.Model(&models.Order{}).Count(&stats.TotalOrders)

	// Get total products
	h.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&stats.TotalProducts)

	// Get total users
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)

	// Get pending orders
	h.db.Model(&models.Order{}).Where("status = ?", "pending").Count(&stats.PendingOrders)

	// Get recent orders (last 10)
	var orders []models.Order
	h.db.Preload("User").
		Order("created_at DESC").
		Limit(10).
		Find(&orders)

	stats.RecentOrders = make([]OrderSummary, len(orders))
	for i, order := range orders {
		stats.RecentOrders[i] = OrderSummary{
			ID:          order.ID,
			OrderNumber: order.OrderNumber,
			UserName:    order.User.Name,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
		}
	}

	// Get top selling products
	type ProductSales struct {
		ProductID    uint
		ProductName  string
		TotalSold    int64
		TotalRevenue float64
	}
	var topProducts []ProductSales
	h.db.Table("order_items").
		Select("products.id as product_id, products.name as product_name, SUM(order_items.quantity) as total_sold, SUM(order_items.subtotal) as total_revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.id, products.name").
		Order("total_sold DESC").
		Limit(5).
		Scan(&topProducts)

	stats.TopProducts = make([]TopProduct, len(topProducts))
	for i, product := range topProducts {
		stats.TopProducts[i] = TopProduct{
			ProductID:    product.ProductID,
			ProductName:  product.ProductName,
			TotalSold:    product.TotalSold,
			TotalRevenue: product.TotalRevenue,
		}
	}

	// Get revenue by month (last 6 months)
	type MonthlyData struct {
		Month   int
		Year    int
		Revenue float64
		Orders  int64
	}
	var monthlyData []MonthlyData
	h.db.Table("orders").
		Select("EXTRACT(MONTH FROM created_at)::integer as month, EXTRACT(YEAR FROM created_at)::integer as year, SUM(total_amount) as revenue, COUNT(*) as orders").
		Where("created_at >= ?", time.Now().AddDate(0, -6, 0)).
		Where("payment_status = ?", "paid").
		Group("EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)").
		Order("year DESC, month DESC").
		Scan(&monthlyData)

	monthNames := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	stats.RevenueByMonth = make([]MonthlyRevenue, len(monthlyData))
	for i, data := range monthlyData {
		stats.RevenueByMonth[i] = MonthlyRevenue{
			Month:   monthNames[data.Month],
			Year:    data.Year,
			Revenue: data.Revenue,
			Orders:  data.Orders,
		}
	}

	// Get orders by status
	var statusCounts []OrderStatusCount
	h.db.Table("orders").
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts)

	stats.OrdersByStatus = statusCounts

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

func (h *DashboardHandler) GetRevenueStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month") // day, week, month, year

	type RevenueData struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
		Orders  int64   `json:"orders"`
	}

	var revenueData []RevenueData
	var dateFormat string
	var timeFilter time.Time

	switch period {
	case "day":
		dateFormat = "TO_CHAR(created_at, 'YYYY-MM-DD')"
		timeFilter = time.Now().AddDate(0, 0, -30) // Last 30 days
	case "week":
		dateFormat = "TO_CHAR(created_at, 'IYYY-IW')"
		timeFilter = time.Now().AddDate(0, 0, -84) // Last 12 weeks
	case "year":
		dateFormat = "TO_CHAR(created_at, 'YYYY')"
		timeFilter = time.Now().AddDate(-5, 0, 0) // Last 5 years
	default: // month
		dateFormat = "TO_CHAR(created_at, 'YYYY-MM')"
		timeFilter = time.Now().AddDate(0, -12, 0) // Last 12 months
	}

	h.db.Table("orders").
		Select(dateFormat+" as date, SUM(total_amount) as revenue, COUNT(*) as orders").
		Where("created_at >= ?", timeFilter).
		Where("payment_status = ?", "paid").
		Group("date").
		Order("date ASC").
		Scan(&revenueData)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    revenueData,
	})
}
