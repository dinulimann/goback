package handlers

import (
	"net/http"
	"strconv"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	repo *models.ProductRepository
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		repo: models.NewProductRepository(db),
	}
}

// GetAllProducts retrieves all products with optional filters
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Check for query parameters
	search := c.Query("search")
	categoryIDStr := c.Query("category_id")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	// If no filters, return all products
	if search == "" && categoryIDStr == "" && minPriceStr == "" && maxPriceStr == "" {
		products, err := h.repo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve products",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": products,
		})
		return
	}

	// Parse optional parameters
	var categoryID *uint
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			catID := uint(id)
			categoryID = &catID
		}
	}

	var minPrice *float64
	if minPriceStr != "" {
		price, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			minPrice = &price
		}
	}

	var maxPrice *float64
	if maxPriceStr != "" {
		price, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			maxPrice = &price
		}
	}

	// Get filtered products
	products, err := h.repo.GetWithFilters(search, categoryID, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"filters": gin.H{
			"search":      search,
			"category_id": categoryID,
			"min_price":   minPrice,
			"max_price":   maxPrice,
		},
	})
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

// GetProductsBySKU retrieves a product by SKU
func (h *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku := c.Param("sku")

	product, err := h.repo.GetBySKU(sku)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

// GetProductsByCategory retrieves products by category ID
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	products, err := h.repo.GetByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product := &models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		SKU:         input.SKU,
		CategoryID:  input.CategoryID,
		ImageURL:    input.ImageURL,
	}

	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	} else {
		product.IsActive = true
	}

	if err := h.repo.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    product,
	})
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve product",
		})
		return
	}

	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.SKU = input.SKU
	product.CategoryID = input.CategoryID
	product.ImageURL = input.ImageURL
	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	}

	if err := h.repo.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    product,
	})
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
