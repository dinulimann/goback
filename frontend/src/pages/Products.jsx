import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { productService, categoryService } from '../services/api';
import { useCart } from '../context/CartContext';
import './Products.css';

function Products() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { addToCart } = useCart();
  const [addedToCart, setAddedToCart] = useState({});

  useEffect(() => {
    fetchCategories();
    fetchProducts();
  }, []);

  useEffect(() => {
    if (selectedCategory) {
      fetchProductsByCategory(selectedCategory);
    } else {
      fetchProducts();
    }
  }, [selectedCategory]);

  const fetchCategories = async () => {
    try {
      const response = await categoryService.getAllCategories();
      setCategories(response.data || []);
    } catch (err) {
      console.error('Error fetching categories:', err);
    }
  };

  const fetchProducts = async () => {
    try {
      setLoading(true);
      const response = await productService.getAllProducts();
      setProducts(response.data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load products');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchProductsByCategory = async (categoryId) => {
    try {
      setLoading(true);
      const response = await productService.getProductsByCategory(categoryId);
      setProducts(response.data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load products');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = (product) => {
    addToCart(product, 1);
    setAddedToCart({ ...addedToCart, [product.id]: true });
    setTimeout(() => {
      setAddedToCart({ ...addedToCart, [product.id]: false });
    }, 2000);
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(price);
  };

  if (loading) {
    return (
      <div className="products-container">
        <div className="loading">Loading products...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="products-container">
        <div className="error">{error}</div>
      </div>
    );
  }

  return (
    <div className="products-container">
      <div className="products-header">
        <h1>Our Products</h1>
        <p>Browse our collection of amazing products</p>
      </div>

      <div className="products-content">
        <aside className="category-sidebar">
          <h3>Categories</h3>
          <ul className="category-list">
            <li
              className={selectedCategory === null ? 'active' : ''}
              onClick={() => setSelectedCategory(null)}
            >
              All Products
            </li>
            {categories.map((category) => (
              <li
                key={category.id}
                className={selectedCategory === category.id ? 'active' : ''}
                onClick={() => setSelectedCategory(category.id)}
              >
                {category.name}
              </li>
            ))}
          </ul>
        </aside>

        <div className="products-grid">
          {products.length === 0 ? (
            <div className="no-products">
              <p>No products found</p>
            </div>
          ) : (
            products.map((product) => (
              <div key={product.id} className="product-card">
                <Link to={`/products/${product.id}`} className="product-image-link">
                  {product.image_url ? (
                    <img
                      src={product.image_url}
                      alt={product.name}
                      className="product-image"
                    />
                  ) : (
                    <div className="product-image-placeholder">
                      <span>No Image</span>
                    </div>
                  )}
                </Link>
                <div className="product-info">
                  <Link to={`/products/${product.id}`} className="product-name">
                    {product.name}
                  </Link>
                  <p className="product-category">{product.category?.name}</p>
                  <p className="product-price">{formatPrice(product.price)}</p>
                  <div className="product-stock">
                    {product.stock > 0 ? (
                      <span className="in-stock">In Stock ({product.stock})</span>
                    ) : (
                      <span className="out-of-stock">Out of Stock</span>
                    )}
                  </div>
                  <button
                    className="add-to-cart-btn"
                    onClick={() => handleAddToCart(product)}
                    disabled={product.stock === 0 || addedToCart[product.id]}
                  >
                    {addedToCart[product.id]
                      ? '✓ Added!'
                      : product.stock > 0
                      ? 'Add to Cart'
                      : 'Out of Stock'}
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

export default Products;
