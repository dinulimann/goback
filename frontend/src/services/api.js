import axios from 'axios';

// Base URL untuk backend API
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Create axios instance dengan konfigurasi default
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests automatically
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle 401 errors (unauthorized)
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, redirect to login
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// User API Services
export const userService = {
  // Get all users
  getAllUsers: async () => {
    try {
      const response = await apiClient.get('/users');
      return response.data;
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },

  // Get user by ID
  getUserById: async (id) => {
    try {
      const response = await apiClient.get(`/users/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching user ${id}:`, error);
      throw error;
    }
  },

  // Create new user
  createUser: async (userData) => {
    try {
      const response = await apiClient.post('/users', userData);
      return response.data;
    } catch (error) {
      console.error('Error creating user:', error);
      throw error;
    }
  },

  // Update user
  updateUser: async (id, userData) => {
    try {
      const response = await apiClient.put(`/users/${id}`, userData);
      return response.data;
    } catch (error) {
      console.error(`Error updating user ${id}:`, error);
      throw error;
    }
  },

  // Delete user
  deleteUser: async (id) => {
    try {
      const response = await apiClient.delete(`/users/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error deleting user ${id}:`, error);
      throw error;
    }
  },
};

// Category API Services
export const categoryService = {
  // Get all categories
  getAllCategories: async () => {
    try {
      const response = await apiClient.get('/categories');
      return response.data;
    } catch (error) {
      console.error('Error fetching categories:', error);
      throw error;
    }
  },

  // Get category by ID
  getCategoryById: async (id) => {
    try {
      const response = await apiClient.get(`/categories/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching category ${id}:`, error);
      throw error;
    }
  },

  // Get category by slug
  getCategoryBySlug: async (slug) => {
    try {
      const response = await apiClient.get(`/categories/slug/${slug}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching category ${slug}:`, error);
      throw error;
    }
  },

  // Create category (admin)
  createCategory: async (categoryData) => {
    try {
      const response = await apiClient.post('/categories', categoryData);
      return response.data;
    } catch (error) {
      console.error('Error creating category:', error);
      throw error;
    }
  },

  // Update category (admin)
  updateCategory: async (id, categoryData) => {
    try {
      const response = await apiClient.put(`/categories/${id}`, categoryData);
      return response.data;
    } catch (error) {
      console.error(`Error updating category ${id}:`, error);
      throw error;
    }
  },

  // Delete category (admin)
  deleteCategory: async (id) => {
    try {
      const response = await apiClient.delete(`/categories/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error deleting category ${id}:`, error);
      throw error;
    }
  },
};

// Product API Services
export const productService = {
  // Get all products with optional filters
  getAllProducts: async (filters = {}) => {
    try {
      const params = new URLSearchParams();
      if (filters.search) params.append('search', filters.search);
      if (filters.category_id) params.append('category_id', filters.category_id);
      if (filters.min_price) params.append('min_price', filters.min_price);
      if (filters.max_price) params.append('max_price', filters.max_price);

      const queryString = params.toString();
      const url = queryString ? `/products?${queryString}` : '/products';
      const response = await apiClient.get(url);
      return response.data;
    } catch (error) {
      console.error('Error fetching products:', error);
      throw error;
    }
  },

  // Get product by ID
  getProductById: async (id) => {
    try {
      const response = await apiClient.get(`/products/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching product ${id}:`, error);
      throw error;
    }
  },

  // Get product by SKU
  getProductBySKU: async (sku) => {
    try {
      const response = await apiClient.get(`/products/sku/${sku}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching product ${sku}:`, error);
      throw error;
    }
  },

  // Get products by category
  getProductsByCategory: async (categoryId) => {
    try {
      const response = await apiClient.get(`/products/category/${categoryId}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching products for category ${categoryId}:`, error);
      throw error;
    }
  },

  // Create product (admin)
  createProduct: async (productData) => {
    try {
      const response = await apiClient.post('/products', productData);
      return response.data;
    } catch (error) {
      console.error('Error creating product:', error);
      throw error;
    }
  },

  // Update product (admin)
  updateProduct: async (id, productData) => {
    try {
      const response = await apiClient.put(`/products/${id}`, productData);
      return response.data;
    } catch (error) {
      console.error(`Error updating product ${id}:`, error);
      throw error;
    }
  },

  // Delete product (admin)
  deleteProduct: async (id) => {
    try {
      const response = await apiClient.delete(`/products/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error deleting product ${id}:`, error);
      throw error;
    }
  },
};

// Profile API Services
export const profileService = {
  // Get my profile
  getMyProfile: async () => {
    try {
      const response = await apiClient.get('/profile');
      return response.data;
    } catch (error) {
      console.error('Error fetching profile:', error);
      throw error;
    }
  },

  // Create or update my profile
  updateMyProfile: async (profileData) => {
    try {
      const response = await apiClient.post('/profile', profileData);
      return response.data;
    } catch (error) {
      console.error('Error updating profile:', error);
      throw error;
    }
  },

  // Get profile by user ID (admin)
  getProfileByUserId: async (userId) => {
    try {
      const response = await apiClient.get(`/profile/user/${userId}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching profile for user ${userId}:`, error);
      throw error;
    }
  },
};

// Order API Services
export const orderService = {
  // Create order
  createOrder: async (orderData) => {
    try {
      const response = await apiClient.post('/orders', orderData);
      return response.data;
    } catch (error) {
      console.error('Error creating order:', error);
      throw error;
    }
  },

  // Get my orders
  getMyOrders: async () => {
    try {
      const response = await apiClient.get('/orders');
      return response.data;
    } catch (error) {
      console.error('Error fetching orders:', error);
      throw error;
    }
  },

  // Get all orders (admin)
  getAllOrders: async () => {
    try {
      const response = await apiClient.get('/orders/all');
      return response.data;
    } catch (error) {
      console.error('Error fetching all orders:', error);
      throw error;
    }
  },

  // Get order by ID
  getOrderById: async (id) => {
    try {
      const response = await apiClient.get(`/orders/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching order ${id}:`, error);
      throw error;
    }
  },

  // Update order status (admin)
  updateOrderStatus: async (id, status) => {
    try {
      const response = await apiClient.put(`/orders/${id}/status`, { status });
      return response.data;
    } catch (error) {
      console.error(`Error updating order ${id} status:`, error);
      throw error;
    }
  },

  // Cancel order
  cancelOrder: async (id) => {
    try {
      const response = await apiClient.post(`/orders/${id}/cancel`);
      return response.data;
    } catch (error) {
      console.error(`Error cancelling order ${id}:`, error);
      throw error;
    }
  },
};

// Review API Services
export const reviewService = {
  // Get reviews for a product
  getProductReviews: async (productId) => {
    try {
      const response = await apiClient.get(`/reviews/product/${productId}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching reviews for product ${productId}:`, error);
      throw error;
    }
  },

  // Get my reviews
  getMyReviews: async () => {
    try {
      const response = await apiClient.get('/reviews/my');
      return response.data;
    } catch (error) {
      console.error('Error fetching my reviews:', error);
      throw error;
    }
  },

  // Create review
  createReview: async (reviewData) => {
    try {
      const response = await apiClient.post('/reviews', reviewData);
      return response.data;
    } catch (error) {
      console.error('Error creating review:', error);
      throw error;
    }
  },

  // Update review
  updateReview: async (id, reviewData) => {
    try {
      const response = await apiClient.put(`/reviews/${id}`, reviewData);
      return response.data;
    } catch (error) {
      console.error(`Error updating review ${id}:`, error);
      throw error;
    }
  },

  // Delete review
  deleteReview: async (id) => {
    try {
      const response = await apiClient.delete(`/reviews/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error deleting review ${id}:`, error);
      throw error;
    }
  },
};

// Wishlist API Services
export const wishlistService = {
  // Get my wishlist
  getWishlist: async () => {
    try {
      const response = await apiClient.get('/wishlist');
      return response.data;
    } catch (error) {
      console.error('Error fetching wishlist:', error);
      throw error;
    }
  },

  // Add to wishlist
  addToWishlist: async (productId) => {
    try {
      const response = await apiClient.post('/wishlist', { product_id: productId });
      return response.data;
    } catch (error) {
      console.error('Error adding to wishlist:', error);
      throw error;
    }
  },

  // Remove from wishlist
  removeFromWishlist: async (productId) => {
    try {
      const response = await apiClient.delete('/wishlist', { data: { product_id: productId } });
      return response.data;
    } catch (error) {
      console.error('Error removing from wishlist:', error);
      throw error;
    }
  },

  // Check if product is in wishlist
  checkWishlist: async (productId) => {
    try {
      const response = await apiClient.post('/wishlist/check', { product_id: productId });
      return response.data;
    } catch (error) {
      console.error('Error checking wishlist:', error);
      throw error;
    }
  },
};

// Health check
export const healthCheck = async () => {
  try {
    const response = await axios.get('http://localhost:8080/health');
    return response.data;
  } catch (error) {
    console.error('Health check failed:', error);
    throw error;
  }
};
