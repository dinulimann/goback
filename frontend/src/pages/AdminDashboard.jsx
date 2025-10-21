import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import '../styles/AdminDashboard.css';

function AdminDashboard() {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [revenueData, setRevenueData] = useState([]);
  const [revenuePeriod, setRevenuePeriod] = useState('month');
  const navigate = useNavigate();

  useEffect(() => {
    fetchDashboardStats();
    fetchRevenueStats(revenuePeriod);
  }, []);

  const fetchDashboardStats = async () => {
    try {
      setLoading(true);
      const response = await api.get('/dashboard/stats');
      setStats(response.data.data);
      setError(null);
    } catch (err) {
      console.error('Error fetching dashboard stats:', err);
      setError('Failed to load dashboard statistics');
    } finally {
      setLoading(false);
    }
  };

  const fetchRevenueStats = async (period) => {
    try {
      const response = await api.get(`/dashboard/revenue?period=${period}`);
      setRevenueData(response.data.data);
    } catch (err) {
      console.error('Error fetching revenue stats:', err);
    }
  };

  const handlePeriodChange = (period) => {
    setRevenuePeriod(period);
    fetchRevenueStats(period);
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const getStatusBadgeClass = (status) => {
    const statusClasses = {
      pending: 'status-badge-pending',
      processing: 'status-badge-processing',
      shipped: 'status-badge-shipped',
      delivered: 'status-badge-delivered',
      cancelled: 'status-badge-cancelled',
    };
    return statusClasses[status] || 'status-badge-default';
  };

  if (loading) {
    return (
      <div className="admin-dashboard">
        <div className="loading">Loading dashboard...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="admin-dashboard">
        <div className="error">{error}</div>
      </div>
    );
  }

  return (
    <div className="admin-dashboard">
      <div className="dashboard-header">
        <h1>Admin Dashboard</h1>
        <p>Welcome to your analytics dashboard</p>
      </div>

      {/* Stats Cards */}
      <div className="stats-grid">
        <div className="stat-card stat-revenue">
          <div className="stat-icon">💰</div>
          <div className="stat-content">
            <h3>Total Revenue</h3>
            <p className="stat-value">{formatCurrency(stats.total_revenue)}</p>
            <span className="stat-label">All time</span>
          </div>
        </div>

        <div className="stat-card stat-orders">
          <div className="stat-icon">📦</div>
          <div className="stat-content">
            <h3>Total Orders</h3>
            <p className="stat-value">{stats.total_orders}</p>
            <span className="stat-label">{stats.pending_orders} pending</span>
          </div>
        </div>

        <div className="stat-card stat-products">
          <div className="stat-icon">🏷️</div>
          <div className="stat-content">
            <h3>Total Products</h3>
            <p className="stat-value">{stats.total_products}</p>
            <span className="stat-label">Active products</span>
          </div>
        </div>

        <div className="stat-card stat-users">
          <div className="stat-icon">👥</div>
          <div className="stat-content">
            <h3>Total Users</h3>
            <p className="stat-value">{stats.total_users}</p>
            <span className="stat-label">Registered users</span>
          </div>
        </div>
      </div>

      {/* Charts Section */}
      <div className="dashboard-grid">
        {/* Revenue Chart */}
        <div className="dashboard-card revenue-chart">
          <div className="card-header">
            <h2>Revenue Overview</h2>
            <div className="period-selector">
              <button
                className={revenuePeriod === 'day' ? 'active' : ''}
                onClick={() => handlePeriodChange('day')}
              >
                Daily
              </button>
              <button
                className={revenuePeriod === 'week' ? 'active' : ''}
                onClick={() => handlePeriodChange('week')}
              >
                Weekly
              </button>
              <button
                className={revenuePeriod === 'month' ? 'active' : ''}
                onClick={() => handlePeriodChange('month')}
              >
                Monthly
              </button>
              <button
                className={revenuePeriod === 'year' ? 'active' : ''}
                onClick={() => handlePeriodChange('year')}
              >
                Yearly
              </button>
            </div>
          </div>
          <div className="chart-container">
            {revenueData.length > 0 ? (
              <div className="simple-bar-chart">
                {revenueData.map((item, index) => (
                  <div key={index} className="chart-bar-wrapper">
                    <div
                      className="chart-bar"
                      style={{
                        height: `${(item.revenue / Math.max(...revenueData.map(d => d.revenue))) * 200}px`
                      }}
                      title={`${item.date}: ${formatCurrency(item.revenue)}`}
                    ></div>
                    <span className="chart-label">{item.date}</span>
                  </div>
                ))}
              </div>
            ) : (
              <p className="no-data">No revenue data available</p>
            )}
          </div>
        </div>

        {/* Top Products */}
        <div className="dashboard-card top-products">
          <div className="card-header">
            <h2>Top Selling Products</h2>
          </div>
          <div className="products-list">
            {stats.top_products && stats.top_products.length > 0 ? (
              stats.top_products.map((product, index) => (
                <div key={product.product_id} className="product-item">
                  <div className="product-rank">#{index + 1}</div>
                  <div className="product-info">
                    <h4>{product.product_name}</h4>
                    <p>{product.total_sold} units sold</p>
                  </div>
                  <div className="product-revenue">
                    {formatCurrency(product.total_revenue)}
                  </div>
                </div>
              ))
            ) : (
              <p className="no-data">No sales data available</p>
            )}
          </div>
        </div>
      </div>

      {/* Orders by Status */}
      <div className="dashboard-grid">
        <div className="dashboard-card orders-status">
          <div className="card-header">
            <h2>Orders by Status</h2>
          </div>
          <div className="status-grid">
            {stats.orders_by_status && stats.orders_by_status.length > 0 ? (
              stats.orders_by_status.map((status) => (
                <div key={status.status} className="status-item">
                  <span className={`status-badge ${getStatusBadgeClass(status.status)}`}>
                    {status.status}
                  </span>
                  <span className="status-count">{status.count}</span>
                </div>
              ))
            ) : (
              <p className="no-data">No order data available</p>
            )}
          </div>
        </div>

        {/* Monthly Revenue */}
        <div className="dashboard-card monthly-revenue">
          <div className="card-header">
            <h2>Monthly Revenue (Last 6 Months)</h2>
          </div>
          <div className="revenue-list">
            {stats.revenue_by_month && stats.revenue_by_month.length > 0 ? (
              stats.revenue_by_month.map((month, index) => (
                <div key={index} className="revenue-item">
                  <div className="month-info">
                    <span className="month-name">{month.month} {month.year}</span>
                    <span className="month-orders">{month.orders} orders</span>
                  </div>
                  <div className="month-revenue">
                    {formatCurrency(month.revenue)}
                  </div>
                </div>
              ))
            ) : (
              <p className="no-data">No revenue data available</p>
            )}
          </div>
        </div>
      </div>

      {/* Recent Orders */}
      <div className="dashboard-card recent-orders">
        <div className="card-header">
          <h2>Recent Orders</h2>
          <button
            className="view-all-btn"
            onClick={() => navigate('/admin/orders')}
          >
            View All Orders
          </button>
        </div>
        <div className="orders-table">
          <table>
            <thead>
              <tr>
                <th>Order #</th>
                <th>Customer</th>
                <th>Amount</th>
                <th>Status</th>
                <th>Date</th>
              </tr>
            </thead>
            <tbody>
              {stats.recent_orders && stats.recent_orders.length > 0 ? (
                stats.recent_orders.map((order) => (
                  <tr key={order.id}>
                    <td className="order-number">{order.order_number}</td>
                    <td>{order.user_name}</td>
                    <td className="order-amount">{formatCurrency(order.total_amount)}</td>
                    <td>
                      <span className={`status-badge ${getStatusBadgeClass(order.status)}`}>
                        {order.status}
                      </span>
                    </td>
                    <td>{formatDate(order.created_at)}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="5" className="no-data">No recent orders</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default AdminDashboard;
