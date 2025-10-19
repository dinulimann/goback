package handlers

import (
	"net/http"
	"strconv"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	repo *models.CategoryRepository
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{
		repo: models.NewCategoryRepository(db),
	}
}

// GetAllCategories retrieves all categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// GetCategory retrieves a category by ID
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	category, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// GetCategoryBySlug retrieves a category by slug
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.repo.GetBySlug(slug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category := &models.Category{
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
	}

	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	} else {
		category.IsActive = true
	}

	if err := h.repo.Create(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}

// UpdateCategory updates a category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	category, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve category",
		})
		return
	}

	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category.Name = input.Name
	category.Description = input.Description
	category.Slug = input.Slug
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	if err := h.repo.Update(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    category,
	})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
