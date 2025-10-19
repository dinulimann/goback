import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { orderService } from '../services/api';
import './Orders.css';

function Orders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      setLoading(true);
      const response = await orderService.getMyOrders();
      setOrders(response.data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load orders');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCancelOrder = async (orderId) => {
    if (!window.confirm('Are you sure you want to cancel this order?')) {
      return;
    }

    try {
      await orderService.cancelOrder(orderId);
      alert('Order cancelled successfully');
      fetchOrders();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to cancel order');
    }
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(price);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusBadge = (status) => {
    const statusClasses = {
      pending: 'status-pending',
      processing: 'status-processing',
      shipped: 'status-shipped',
      delivered: 'status-delivered',
      cancelled: 'status-cancelled',
    };

    return <span className={`status-badge ${statusClasses[status]}`}>{status}</span>;
  };

  if (loading) {
    return (
      <div className="orders-container">
        <div className="loading">Loading orders...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="orders-container">
        <div className="error">{error}</div>
      </div>
    );
  }

  return (
    <div className="orders-container">
      <div className="orders-header">
        <h1>My Orders</h1>
        <Link to="/products" className="continue-shopping-btn">
          Continue Shopping
        </Link>
      </div>

      {orders.length === 0 ? (
        <div className="no-orders">
          <p>You haven't placed any orders yet</p>
          <Link to="/products" className="browse-btn">
            Browse Products
          </Link>
        </div>
      ) : (
        <div className="orders-list">
          {orders.map((order) => (
            <div key={order.id} className="order-card">
              <div className="order-header">
                <div className="order-info">
                  <h3>Order #{order.order_number}</h3>
                  <p className="order-date">{formatDate(order.created_at)}</p>
                </div>
                {getStatusBadge(order.status)}
              </div>

              <div className="order-items">
                {order.order_items?.map((item) => (
                  <div key={item.id} className="order-item">
                    <span>{item.product?.name}</span>
                    <span className="item-quantity">x{item.quantity}</span>
                    <span className="item-price">{formatPrice(item.subtotal)}</span>
                  </div>
                ))}
              </div>

              <div className="order-footer">
                <div className="order-total">
                  <span>Total:</span>
                  <span className="total-amount">{formatPrice(order.total_amount)}</span>
                </div>
                <div className="order-actions">
                  <Link to={`/orders/${order.id}`} className="view-details-btn">
                    View Details
                  </Link>
                  {(order.status === 'pending' || order.status === 'processing') && (
                    <button
                      onClick={() => handleCancelOrder(order.id)}
                      className="cancel-order-btn"
                    >
                      Cancel Order
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Orders;
