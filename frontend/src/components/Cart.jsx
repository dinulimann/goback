import { useCart } from '../context/CartContext';
import { useNavigate } from 'react-router-dom';
import './Cart.css';

function Cart() {
  const {
    cartItems,
    removeFromCart,
    updateQuantity,
    getCartTotal,
    isCartOpen,
    setIsCartOpen,
  } = useCart();
  const navigate = useNavigate();

  const formatPrice = (price) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(price);
  };

  const handleCheckout = () => {
    setIsCartOpen(false);
    navigate('/checkout');
  };

  if (!isCartOpen) return null;

  return (
    <>
      <div className="cart-overlay" onClick={() => setIsCartOpen(false)} />
      <div className="cart-drawer">
        <div className="cart-header">
          <h2>Shopping Cart</h2>
          <button className="close-btn" onClick={() => setIsCartOpen(false)}>
            ✕
          </button>
        </div>

        <div className="cart-content">
          {cartItems.length === 0 ? (
            <div className="empty-cart">
              <p>Your cart is empty</p>
              <button
                onClick={() => {
                  setIsCartOpen(false);
                  navigate('/products');
                }}
                className="browse-btn"
              >
                Browse Products
              </button>
            </div>
          ) : (
            <>
              <div className="cart-items">
                {cartItems.map((item) => (
                  <div key={item.id} className="cart-item">
                    <div className="cart-item-image">
                      {item.image_url ? (
                        <img src={item.image_url} alt={item.name} />
                      ) : (
                        <div className="cart-item-placeholder">No Image</div>
                      )}
                    </div>
                    <div className="cart-item-info">
                      <h4>{item.name}</h4>
                      <p className="cart-item-price">{formatPrice(item.price)}</p>
                      <div className="cart-item-quantity">
                        <button
                          onClick={() => updateQuantity(item.id, item.quantity - 1)}
                          disabled={item.quantity <= 1}
                        >
                          -
                        </button>
                        <span>{item.quantity}</span>
                        <button
                          onClick={() => updateQuantity(item.id, item.quantity + 1)}
                          disabled={item.quantity >= item.stock}
                        >
                          +
                        </button>
                      </div>
                    </div>
                    <button
                      className="remove-btn"
                      onClick={() => removeFromCart(item.id)}
                      title="Remove item"
                    >
                      🗑️
                    </button>
                  </div>
                ))}
              </div>

              <div className="cart-footer">
                <div className="cart-total">
                  <span>Total:</span>
                  <span className="total-price">{formatPrice(getCartTotal())}</span>
                </div>
                <button className="checkout-btn" onClick={handleCheckout}>
                  Proceed to Checkout
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </>
  );
}

export default Cart;
