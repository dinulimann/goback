import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productService } from '../services/api';
import { useCart } from '../context/CartContext';
import './ProductDetail.css';

function ProductDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const { addToCart, toggleCart } = useCart();
  const [addedToCart, setAddedToCart] = useState(false);

  useEffect(() => {
    fetchProduct();
  }, [id]);

  const fetchProduct = async () => {
    try {
      setLoading(true);
      const response = await productService.getProductById(id);
      setProduct(response.data);
      setError(null);
    } catch (err) {
      setError('Failed to load product details');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = () => {
    if (product && quantity > 0 && quantity <= product.stock) {
      addToCart(product, quantity);
      setAddedToCart(true);
      setTimeout(() => setAddedToCart(false), 2000);
    }
  };

  const handleBuyNow = () => {
    if (product && quantity > 0 && quantity <= product.stock) {
      addToCart(product, quantity);
      toggleCart();
    }
  };

  const handleQuantityChange = (newQuantity) => {
    if (newQuantity >= 1 && newQuantity <= product.stock) {
      setQuantity(newQuantity);
    }
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(price);
  };

  if (loading) {
    return (
      <div className="product-detail-container">
        <div className="loading">Loading product...</div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="product-detail-container">
        <div className="error">{error || 'Product not found'}</div>
        <button onClick={() => navigate('/products')} className="back-btn">
          Back to Products
        </button>
      </div>
    );
  }

  return (
    <div className="product-detail-container">
      <button onClick={() => navigate('/products')} className="back-btn">
        ← Back to Products
      </button>

      <div className="product-detail-content">
        <div className="product-image-section">
          {product.image_url ? (
            <img src={product.image_url} alt={product.name} className="product-detail-image" />
          ) : (
            <div className="product-detail-image-placeholder">
              <span>No Image Available</span>
            </div>
          )}
        </div>

        <div className="product-info-section">
          <div className="product-category-badge">{product.category?.name}</div>
          <h1 className="product-title">{product.name}</h1>

          {product.sku && (
            <p className="product-sku">SKU: {product.sku}</p>
          )}

          <div className="product-price-section">
            <p className="product-detail-price">{formatPrice(product.price)}</p>
            <div className="product-availability">
              {product.stock > 0 ? (
                <span className="in-stock">
                  ✓ In Stock ({product.stock} available)
                </span>
              ) : (
                <span className="out-of-stock">✗ Out of Stock</span>
              )}
            </div>
          </div>

          <div className="product-description">
            <h3>Description</h3>
            <p>{product.description || 'No description available.'}</p>
          </div>

          {product.stock > 0 && (
            <div className="product-actions">
              <div className="quantity-selector">
                <label>Quantity:</label>
                <div className="quantity-controls">
                  <button
                    onClick={() => handleQuantityChange(quantity - 1)}
                    disabled={quantity <= 1}
                  >
                    -
                  </button>
                  <input
                    type="number"
                    value={quantity}
                    onChange={(e) => handleQuantityChange(parseInt(e.target.value) || 1)}
                    min="1"
                    max={product.stock}
                  />
                  <button
                    onClick={() => handleQuantityChange(quantity + 1)}
                    disabled={quantity >= product.stock}
                  >
                    +
                  </button>
                </div>
              </div>

              <div className="action-buttons">
                <button
                  className="add-to-cart-btn"
                  onClick={handleAddToCart}
                  disabled={addedToCart}
                >
                  {addedToCart ? '✓ Added to Cart!' : 'Add to Cart'}
                </button>
                <button className="buy-now-btn" onClick={handleBuyNow}>
                  Buy Now
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ProductDetail;
