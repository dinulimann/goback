package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dinulimann/goback/config"
	"github.com/dinulimann/goback/models"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	paymentRepo models.PaymentRepository
	orderRepo   *models.OrderRepository
	config      *config.Config
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentRepo models.PaymentRepository, orderRepo *models.OrderRepository, cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		config:      cfg,
	}
}

// CreateSnapTokenRequest represents request to create snap token
type CreateSnapTokenRequest struct {
	OrderID uint `json:"order_id" binding:"required"`
}

// CreateSnapToken creates Midtrans Snap token for payment
func (h *PaymentHandler) CreateSnapToken(c *gin.Context) {
	var req CreateSnapTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get order
	order, err := h.orderRepo.GetByID(req.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
		return
	}

	// Verify order belongs to user
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this order"})
		return
	}

	// Check if payment already exists
	existingPayment, err := h.paymentRepo.GetByOrderID(req.OrderID)
	if err == nil && existingPayment.PaymentStatus == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already paid"})
		return
	}

	// Initialize Midtrans client
	var s snap.Client
	s.New(h.config.MidtransServerKey, midtrans.Sandbox)

	// Set to Production if environment is production
	if h.config.Environment == "production" {
		s.New(h.config.MidtransServerKey, midtrans.Production)
	}

	// Create transaction details
	transactionID := fmt.Sprintf("ORDER-%d-%d", order.ID, time.Now().Unix())

	// Add items
	var items []midtrans.ItemDetails
	for _, item := range order.OrderItems {
		productName := "Product"
		if item.Product.Name != "" {
			productName = item.Product.Name
		}
		items = append(items, midtrans.ItemDetails{
			ID:    fmt.Sprintf("ITEM-%d", item.ProductID),
			Name:  productName,
			Price: int64(item.Price),
			Qty:   int32(item.Quantity),
		})
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionID,
			GrossAmt: int64(order.FinalAmount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.User.Name,
			Email: order.User.Email,
		},
		Items:           &items,
		EnabledPayments: snap.AllSnapPaymentType,
	}

	// Create snap token
	snapResp, err := s.CreateTransaction(snapReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment token",
			"details": err.Error(),
		})
		return
	}

	// Save or update payment record
	payment := &models.Payment{
		OrderID:         order.ID,
		PaymentMethod:   "midtrans",
		TransactionID:   transactionID,
		SnapToken:       snapResp.Token,
		GrossAmount:     order.FinalAmount,
		PaymentStatus:   "pending",
		PaymentProvider: "midtrans",
	}

	// Set expiration (24 hours from now)
	expiredAt := time.Now().Add(24 * time.Hour)
	payment.ExpiredAt = &expiredAt

	if existingPayment != nil {
		payment.ID = existingPayment.ID
		if err := h.paymentRepo.Update(payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
			return
		}
	} else {
		if err := h.paymentRepo.Create(payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment token created successfully",
		"snap_token":     snapResp.Token,
		"redirect_url":   snapResp.RedirectURL,
		"transaction_id": transactionID,
	})
}

// MidtransNotification handles Midtrans payment notification callback
func (h *PaymentHandler) MidtransNotification(c *gin.Context) {
	var notificationPayload map[string]interface{}

	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get transaction details
	transactionID, ok := notificationPayload["order_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_id"})
		return
	}

	transactionStatus, _ := notificationPayload["transaction_status"].(string)
	fraudStatus, _ := notificationPayload["fraud_status"].(string)
	paymentType, _ := notificationPayload["payment_type"].(string)

	// Get payment record
	payment, err := h.paymentRepo.GetByTransactionID(transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Update payment based on transaction status
	payment.TransactionStatus = transactionStatus
	payment.FraudStatus = fraudStatus

	// Extract VA number or bank info if available
	if vaNumber, ok := notificationPayload["va_numbers"].([]interface{}); ok && len(vaNumber) > 0 {
		if vaMap, ok := vaNumber[0].(map[string]interface{}); ok {
			if va, ok := vaMap["va_number"].(string); ok {
				payment.VANumber = va
			}
			if bank, ok := vaMap["bank"].(string); ok {
				payment.Bank = bank
			}
		}
	}

	// Update payment status based on transaction status
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			payment.PaymentStatus = "paid"
			now := time.Now()
			payment.PaidAt = &now

			// Update order status
			order, _ := h.orderRepo.GetByID(payment.OrderID)
			if order != nil {
				order.Status = "processing"
				order.PaymentStatus = "paid"
				h.orderRepo.Update(order)
			}
		}
	case "settlement":
		payment.PaymentStatus = "paid"
		now := time.Now()
		payment.PaidAt = &now

		// Update order status
		order, _ := h.orderRepo.GetByID(payment.OrderID)
		if order != nil {
			order.Status = "processing"
			order.PaymentStatus = "paid"
			h.orderRepo.Update(order)
		}
	case "pending":
		payment.PaymentStatus = "pending"
	case "deny", "cancel", "expire":
		payment.PaymentStatus = "failed"
		payment.ErrorMessage = fmt.Sprintf("Payment %s", transactionStatus)

		// Update order status
		order, _ := h.orderRepo.GetByID(payment.OrderID)
		if order != nil {
			order.Status = "cancelled"
			order.PaymentStatus = "failed"
			h.orderRepo.Update(order)
		}
	}

	payment.PaymentMethod = paymentType

	// Save payment update
	if err := h.paymentRepo.Update(payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification processed successfully",
	})
}

// GetPaymentByOrderID retrieves payment details by order ID
func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get order to verify ownership
	order, err := h.orderRepo.GetByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Verify order belongs to user
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this payment"})
		return
	}

	// Get payment
	payment, err := h.paymentRepo.GetByOrderID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment retrieved successfully",
		"payment": payment,
	})
}

// GetPaymentStatus checks current payment status
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	transactionID := c.Param("transactionId")

	payment, err := h.paymentRepo.GetByTransactionID(transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify order belongs to user
	order, _ := h.orderRepo.GetByID(payment.OrderID)
	if order != nil && order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Payment status retrieved successfully",
		"payment_status":   payment.PaymentStatus,
		"transaction_id":   payment.TransactionID,
		"gross_amount":     payment.GrossAmount,
		"payment_method":   payment.PaymentMethod,
		"va_number":        payment.VANumber,
		"bank":             payment.Bank,
		"paid_at":          payment.PaidAt,
		"expired_at":       payment.ExpiredAt,
	})
}

// GetAllPayments retrieves all payments (admin only)
func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	offset := (page - 1) * limit

	var payments []models.Payment
	var total int64
	var err error

	if status != "" {
		payments, total, err = h.paymentRepo.GetByStatus(status, limit, offset)
	} else {
		payments, total, err = h.paymentRepo.GetAll(limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Payments retrieved successfully",
		"payments": payments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
