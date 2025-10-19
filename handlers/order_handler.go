package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	orderRepo     *models.OrderRepository
	orderItemRepo *models.OrderItemRepository
	productRepo   *models.ProductRepository
	db            *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{
		orderRepo:     models.NewOrderRepository(db),
		orderItemRepo: models.NewOrderItemRepository(db),
		productRepo:   models.NewProductRepository(db),
		db:            db,
	}
}

// generateOrderNumber generates a unique order number
func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

// GetAllOrders retrieves all orders (admin)
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}

// GetMyOrders retrieves orders for authenticated user
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	orders, err := h.orderRepo.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}

// GetOrder retrieves an order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := h.orderRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order",
		})
		return
	}

	// Check if user owns this order (unless admin)
	userID, exists := c.Get("user_id")
	if exists && order.UserID != userID.(uint) {
		// TODO: Add admin check here
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to view this order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var input models.OrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate total and create order items
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range input.Items {
		// Get product
		product, err := h.productRepo.GetByID(item.ProductID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Product with ID %d not found", item.ProductID),
			})
			return
		}

		// Check stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Insufficient stock for product %s", product.Name),
			})
			return
		}

		// Calculate subtotal
		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		})
	}

	// Create order
	order := &models.Order{
		UserID:             userID.(uint),
		OrderNumber:        generateOrderNumber(),
		Status:             models.OrderStatusPending,
		TotalAmount:        totalAmount,
		ShippingAddress:    input.ShippingAddress,
		ShippingCity:       input.ShippingCity,
		ShippingProvince:   input.ShippingProvince,
		ShippingPostalCode: input.ShippingPostalCode,
		PaymentMethod:      input.PaymentMethod,
		PaymentStatus:      "unpaid",
		Notes:              input.Notes,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	// Create order items and update stock
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create order items",
			})
			return
		}

		// Update product stock
		if err := tx.Model(&models.Product{}).Where("id = ?", orderItems[i].ProductID).
			Update("stock", gorm.Expr("stock - ?", orderItems[i].Quantity)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update product stock",
			})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit order",
		})
		return
	}

	// Load relationships
	order, _ = h.orderRepo.GetByID(order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    order,
	})
}

// UpdateOrderStatus updates order status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var input models.OrderUpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.orderRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order",
		})
		return
	}

	order.Status = input.Status
	if err := h.orderRepo.Update(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"data":    order,
	})
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	order, err := h.orderRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order",
		})
		return
	}

	// Check if user owns this order
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to cancel this order",
		})
		return
	}

	// Check if order can be cancelled
	if order.Status == models.OrderStatusDelivered || order.Status == models.OrderStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order cannot be cancelled",
		})
		return
	}

	// Start transaction to restore stock
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Restore product stock
	orderItems, _ := h.orderItemRepo.GetByOrderID(order.ID)
	for _, item := range orderItems {
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to restore product stock",
			})
			return
		}
	}

	// Update order status
	order.Status = models.OrderStatusCancelled
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cancel order",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit cancellation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"data":    order,
	})
}
