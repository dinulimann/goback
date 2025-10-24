import { useState, useEffect } from 'react';
import api from '../services/api';
import Navbar from '../components/Navbar';
import '../styles/AdminCoupons.css';

function AdminCoupons() {
  const [coupons, setCoupons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingCoupon, setEditingCoupon] = useState(null);
  const [formData, setFormData] = useState({
    code: '',
    description: '',
    discount_type: 'percentage',
    discount_value: '',
    min_purchase: '',
    max_discount: '',
    usage_limit: '',
    start_date: '',
    end_date: '',
    is_active: true,
  });

  useEffect(() => {
    fetchCoupons();
  }, []);

  const fetchCoupons = async () => {
    try {
      const response = await api.get('/coupons');
      setCoupons(response.data.coupons || []);
    } catch (error) {
      console.error('Failed to fetch coupons:', error);
      alert('Failed to load coupons');
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'checkbox' ? checked : value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Prepare data for API
    const couponData = {
      code: formData.code.toUpperCase(),
      description: formData.description,
      discount_type: formData.discount_type,
      discount_value: parseFloat(formData.discount_value),
      min_purchase: parseFloat(formData.min_purchase) || 0,
      max_discount: parseFloat(formData.max_discount) || 0,
      usage_limit: parseInt(formData.usage_limit) || 0,
      start_date: new Date(formData.start_date).toISOString(),
      end_date: new Date(formData.end_date).toISOString(),
      is_active: formData.is_active,
    };

    try {
      if (editingCoupon) {
        await api.put(`/coupons/${editingCoupon.id}`, couponData);
        alert('Coupon updated successfully!');
      } else {
        await api.post('/coupons', couponData);
        alert('Coupon created successfully!');
      }

      fetchCoupons();
      resetForm();
    } catch (error) {
      console.error('Failed to save coupon:', error);
      alert(error.response?.data?.error || 'Failed to save coupon');
    }
  };

  const handleEdit = (coupon) => {
    setEditingCoupon(coupon);
    setFormData({
      code: coupon.code,
      description: coupon.description,
      discount_type: coupon.discount_type,
      discount_value: coupon.discount_value,
      min_purchase: coupon.min_purchase,
      max_discount: coupon.max_discount,
      usage_limit: coupon.usage_limit,
      start_date: new Date(coupon.start_date).toISOString().slice(0, 16),
      end_date: new Date(coupon.end_date).toISOString().slice(0, 16),
      is_active: coupon.is_active,
    });
    setShowForm(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this coupon?')) {
      return;
    }

    try {
      await api.delete(`/coupons/${id}`);
      alert('Coupon deleted successfully!');
      fetchCoupons();
    } catch (error) {
      console.error('Failed to delete coupon:', error);
      alert('Failed to delete coupon');
    }
  };

  const handleToggleActive = async (coupon) => {
    try {
      await api.post(`/coupons/${coupon.id}/toggle`);
      fetchCoupons();
    } catch (error) {
      console.error('Failed to toggle coupon status:', error);
      alert('Failed to toggle coupon status');
    }
  };

  const resetForm = () => {
    setFormData({
      code: '',
      description: '',
      discount_type: 'percentage',
      discount_value: '',
      min_purchase: '',
      max_discount: '',
      usage_limit: '',
      start_date: '',
      end_date: '',
      is_active: true,
    });
    setEditingCoupon(null);
    setShowForm(false);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getDiscountDisplay = (coupon) => {
    if (coupon.discount_type === 'percentage') {
      return `${coupon.discount_value}%`;
    }
    return `Rp ${coupon.discount_value.toLocaleString()}`;
  };

  if (loading) {
    return (
      <div className="admin-coupons">
        <Navbar />
        <div className="loading-spinner">Loading coupons...</div>
      </div>
    );
  }

  return (
    <div className="admin-coupons">
      <Navbar />
      <div className="admin-coupons-container">
        <div className="admin-coupons-header">
          <h1>Coupon Management</h1>
          <button
            className="btn-create-coupon"
            onClick={() => setShowForm(!showForm)}
          >
            {showForm ? 'Cancel' : '+ Create New Coupon'}
          </button>
        </div>

        {showForm && (
          <div className="coupon-form-card">
            <h2>{editingCoupon ? 'Edit Coupon' : 'Create New Coupon'}</h2>
            <form onSubmit={handleSubmit}>
              <div className="form-grid">
                <div className="form-group">
                  <label>Coupon Code *</label>
                  <input
                    type="text"
                    name="code"
                    value={formData.code}
                    onChange={handleInputChange}
                    placeholder="e.g., SAVE20"
                    required
                    disabled={editingCoupon !== null}
                  />
                </div>

                <div className="form-group">
                  <label>Discount Type *</label>
                  <select
                    name="discount_type"
                    value={formData.discount_type}
                    onChange={handleInputChange}
                    required
                  >
                    <option value="percentage">Percentage</option>
                    <option value="fixed">Fixed Amount</option>
                  </select>
                </div>

                <div className="form-group">
                  <label>
                    Discount Value *
                    {formData.discount_type === 'percentage' ? ' (%)' : ' (Rp)'}
                  </label>
                  <input
                    type="number"
                    name="discount_value"
                    value={formData.discount_value}
                    onChange={handleInputChange}
                    placeholder={formData.discount_type === 'percentage' ? '20' : '50000'}
                    min="0"
                    step={formData.discount_type === 'percentage' ? '1' : '1000'}
                    max={formData.discount_type === 'percentage' ? '100' : undefined}
                    required
                  />
                </div>

                <div className="form-group">
                  <label>Minimum Purchase (Rp)</label>
                  <input
                    type="number"
                    name="min_purchase"
                    value={formData.min_purchase}
                    onChange={handleInputChange}
                    placeholder="0"
                    min="0"
                    step="1000"
                  />
                </div>

                <div className="form-group">
                  <label>
                    Max Discount (Rp)
                    {formData.discount_type === 'percentage' && ' - Optional'}
                  </label>
                  <input
                    type="number"
                    name="max_discount"
                    value={formData.max_discount}
                    onChange={handleInputChange}
                    placeholder="0 = No limit"
                    min="0"
                    step="1000"
                    disabled={formData.discount_type === 'fixed'}
                  />
                </div>

                <div className="form-group">
                  <label>Usage Limit</label>
                  <input
                    type="number"
                    name="usage_limit"
                    value={formData.usage_limit}
                    onChange={handleInputChange}
                    placeholder="0 = Unlimited"
                    min="0"
                  />
                </div>

                <div className="form-group">
                  <label>Start Date *</label>
                  <input
                    type="datetime-local"
                    name="start_date"
                    value={formData.start_date}
                    onChange={handleInputChange}
                    required
                  />
                </div>

                <div className="form-group">
                  <label>End Date *</label>
                  <input
                    type="datetime-local"
                    name="end_date"
                    value={formData.end_date}
                    onChange={handleInputChange}
                    required
                  />
                </div>
              </div>

              <div className="form-group full-width">
                <label>Description</label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleInputChange}
                  placeholder="Describe your coupon promotion..."
                  rows="3"
                />
              </div>

              <div className="form-group checkbox-group">
                <label>
                  <input
                    type="checkbox"
                    name="is_active"
                    checked={formData.is_active}
                    onChange={handleInputChange}
                  />
                  <span>Active</span>
                </label>
              </div>

              <div className="form-actions">
                <button type="submit" className="btn-submit">
                  {editingCoupon ? 'Update Coupon' : 'Create Coupon'}
                </button>
                <button type="button" onClick={resetForm} className="btn-cancel">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="coupons-list">
          <h2>All Coupons ({coupons.length})</h2>
          {coupons.length === 0 ? (
            <div className="empty-state">
              <p>No coupons yet. Create your first coupon to get started!</p>
            </div>
          ) : (
            <div className="coupons-grid">
              {coupons.map((coupon) => (
                <div key={coupon.id} className={`coupon-card ${!coupon.is_active ? 'inactive' : ''}`}>
                  <div className="coupon-header">
                    <div className="coupon-code-badge">{coupon.code}</div>
                    <div className={`coupon-status ${coupon.is_active ? 'active' : 'inactive'}`}>
                      {coupon.is_active ? 'Active' : 'Inactive'}
                    </div>
                  </div>

                  <div className="coupon-discount">
                    <span className="discount-value">{getDiscountDisplay(coupon)}</span>
                    <span className="discount-type">
                      {coupon.discount_type === 'percentage' ? 'OFF' : 'DISCOUNT'}
                    </span>
                  </div>

                  {coupon.description && (
                    <p className="coupon-description">{coupon.description}</p>
                  )}

                  <div className="coupon-details">
                    {coupon.min_purchase > 0 && (
                      <div className="detail-item">
                        <span className="detail-label">Min Purchase:</span>
                        <span className="detail-value">Rp {coupon.min_purchase.toLocaleString()}</span>
                      </div>
                    )}

                    {coupon.max_discount > 0 && coupon.discount_type === 'percentage' && (
                      <div className="detail-item">
                        <span className="detail-label">Max Discount:</span>
                        <span className="detail-value">Rp {coupon.max_discount.toLocaleString()}</span>
                      </div>
                    )}

                    <div className="detail-item">
                      <span className="detail-label">Usage:</span>
                      <span className="detail-value">
                        {coupon.used_count} / {coupon.usage_limit > 0 ? coupon.usage_limit : '∞'}
                      </span>
                    </div>

                    <div className="detail-item">
                      <span className="detail-label">Valid From:</span>
                      <span className="detail-value">{formatDate(coupon.start_date)}</span>
                    </div>

                    <div className="detail-item">
                      <span className="detail-label">Valid Until:</span>
                      <span className="detail-value">{formatDate(coupon.end_date)}</span>
                    </div>
                  </div>

                  <div className="coupon-actions">
                    <button
                      className="btn-toggle"
                      onClick={() => handleToggleActive(coupon)}
                    >
                      {coupon.is_active ? 'Deactivate' : 'Activate'}
                    </button>
                    <button
                      className="btn-edit"
                      onClick={() => handleEdit(coupon)}
                    >
                      Edit
                    </button>
                    <button
                      className="btn-delete"
                      onClick={() => handleDelete(coupon.id)}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default AdminCoupons;
