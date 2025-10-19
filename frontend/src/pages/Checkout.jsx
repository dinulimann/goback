import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { orderService, profileService } from '../services/api';
import './Checkout.css';

function Checkout() {
  const navigate = useNavigate();
  const { cartItems, getCartTotal, clearCart } = useCart();
  const [loading, setLoading] = useState(false);
  const [profile, setProfile] = useState(null);
  const [formData, setFormData] = useState({
    shipping_address: '',
    shipping_city: '',
    shipping_province: '',
    shipping_postal_code: '',
    payment_method: 'bank_transfer',
    notes: '',
  });

  useEffect(() => {
    if (cartItems.length === 0) {
      navigate('/products');
    }
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const response = await profileService.getMyProfile();
      const profileData = response.data;
      setProfile(profileData);

      // Pre-fill form with profile data
      if (profileData) {
        setFormData({
          ...formData,
          shipping_address: profileData.address || '',
          shipping_city: profileData.city || '',
          shipping_province: profileData.province || '',
          shipping_postal_code: profileData.postal_code || '',
        });
      }
    } catch (err) {
      console.log('No profile found, user can enter manually');
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (cartItems.length === 0) {
      alert('Your cart is empty');
      return;
    }

    setLoading(true);

    try {
      const orderData = {
        ...formData,
        items: cartItems.map((item) => ({
          product_id: item.id,
          quantity: item.quantity,
        })),
      };

      const response = await orderService.createOrder(orderData);

      if (response.data) {
        clearCart();
        alert('Order placed successfully!');
        navigate(`/orders/${response.data.id}`);
      }
    } catch (err) {
      console.error('Error creating order:', err);
      alert(err.response?.data?.error || 'Failed to create order. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(price);
  };

  return (
    <div className="checkout-container">
      <div className="checkout-content">
        <div className="checkout-form-section">
          <h1>Checkout</h1>

          <form onSubmit={handleSubmit}>
            <div className="form-section">
              <h2>Shipping Address</h2>

              <div className="form-group">
                <label>Street Address *</label>
                <textarea
                  name="shipping_address"
                  value={formData.shipping_address}
                  onChange={handleChange}
                  required
                  rows="3"
                  placeholder="Enter your full address"
                />
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>City *</label>
                  <input
                    type="text"
                    name="shipping_city"
                    value={formData.shipping_city}
                    onChange={handleChange}
                    required
                    placeholder="City"
                  />
                </div>

                <div className="form-group">
                  <label>Province *</label>
                  <input
                    type="text"
                    name="shipping_province"
                    value={formData.shipping_province}
                    onChange={handleChange}
                    required
                    placeholder="Province"
                  />
                </div>
              </div>

              <div className="form-group">
                <label>Postal Code *</label>
                <input
                  type="text"
                  name="shipping_postal_code"
                  value={formData.shipping_postal_code}
                  onChange={handleChange}
                  required
                  placeholder="Postal Code"
                />
              </div>
            </div>

            <div className="form-section">
              <h2>Payment Method</h2>

              <div className="payment-methods">
                <label className="payment-option">
                  <input
                    type="radio"
                    name="payment_method"
                    value="bank_transfer"
                    checked={formData.payment_method === 'bank_transfer'}
                    onChange={handleChange}
                  />
                  <span>Bank Transfer</span>
                </label>

                <label className="payment-option">
                  <input
                    type="radio"
                    name="payment_method"
                    value="credit_card"
                    checked={formData.payment_method === 'credit_card'}
                    onChange={handleChange}
                  />
                  <span>Credit Card</span>
                </label>

                <label className="payment-option">
                  <input
                    type="radio"
                    name="payment_method"
                    value="e_wallet"
                    checked={formData.payment_method === 'e_wallet'}
                    onChange={handleChange}
                  />
                  <span>E-Wallet</span>
                </label>

                <label className="payment-option">
                  <input
                    type="radio"
                    name="payment_method"
                    value="cash_on_delivery"
                    checked={formData.payment_method === 'cash_on_delivery'}
                    onChange={handleChange}
                  />
                  <span>Cash on Delivery</span>
                </label>
              </div>
            </div>

            <div className="form-section">
              <h2>Order Notes (Optional)</h2>

              <div className="form-group">
                <textarea
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  rows="4"
                  placeholder="Any special instructions for your order?"
                />
              </div>
            </div>

            <button type="submit" className="place-order-btn" disabled={loading}>
              {loading ? 'Processing...' : 'Place Order'}
            </button>
          </form>
        </div>

        <div className="order-summary-section">
          <div className="order-summary">
            <h2>Order Summary</h2>

            <div className="summary-items">
              {cartItems.map((item) => (
                <div key={item.id} className="summary-item">
                  <div className="item-info">
                    <span className="item-name">{item.name}</span>
                    <span className="item-quantity">x{item.quantity}</span>
                  </div>
                  <span className="item-price">
                    {formatPrice(item.price * item.quantity)}
                  </span>
                </div>
              ))}
            </div>

            <div className="summary-divider" />

            <div className="summary-total">
              <span>Total</span>
              <span className="total-amount">{formatPrice(getCartTotal())}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Checkout;
