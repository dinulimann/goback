import { useState, useEffect } from 'react';
import { productService, categoryService } from '../services/api';
import './AdminProducts.css';

function AdminProducts() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
    stock: '',
    sku: '',
    category_id: '',
    image_url: '',
    is_active: true,
  });

  useEffect(() => {
    fetchProducts();
    fetchCategories();
  }, []);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      const response = await productService.getAllProducts();
      setProducts(response.data || []);
    } catch (err) {
      console.error(err);
      alert('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const response = await categoryService.getAllCategories();
      setCategories(response.data || []);
    } catch (err) {
      console.error(err);
    }
  };

  const handleOpenModal = (product = null) => {
    if (product) {
      setEditingProduct(product);
      setFormData({
        name: product.name,
        description: product.description || '',
        price: product.price,
        stock: product.stock,
        sku: product.sku || '',
        category_id: product.category_id,
        image_url: product.image_url || '',
        is_active: product.is_active,
      });
    } else {
      setEditingProduct(null);
      setFormData({
        name: '',
        description: '',
        price: '',
        stock: '',
        sku: '',
        category_id: '',
        image_url: '',
        is_active: true,
      });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingProduct(null);
  };

  const handleChange = (e) => {
    const value = e.target.type === 'checkbox' ? e.target.checked : e.target.value;
    setFormData({
      ...formData,
      [e.target.name]: value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const productData = {
      ...formData,
      price: parseFloat(formData.price),
      stock: parseInt(formData.stock),
      category_id: parseInt(formData.category_id),
    };

    try {
      if (editingProduct) {
        await productService.updateProduct(editingProduct.id, productData);
        alert('Product updated successfully!');
      } else {
        await productService.createProduct(productData);
        alert('Product created successfully!');
      }
      handleCloseModal();
      fetchProducts();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to save product');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this product?')) {
      return;
    }

    try {
      await productService.deleteProduct(id);
      alert('Product deleted successfully!');
      fetchProducts();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to delete product');
    }
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(price);
  };

  if (loading) {
    return <div className="admin-container">Loading...</div>;
  }

  return (
    <div className="admin-container">
      <div className="admin-header">
        <h1>Product Management</h1>
        <button onClick={() => handleOpenModal()} className="add-btn">
          + Add New Product
        </button>
      </div>

      <div className="admin-table-container">
        <table className="admin-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Image</th>
              <th>Name</th>
              <th>Category</th>
              <th>Price</th>
              <th>Stock</th>
              <th>SKU</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((product) => (
              <tr key={product.id}>
                <td>{product.id}</td>
                <td>
                  {product.image_url ? (
                    <img src={product.image_url} alt={product.name} className="product-thumb" />
                  ) : (
                    <div className="product-thumb-placeholder">No Image</div>
                  )}
                </td>
                <td>{product.name}</td>
                <td>{product.category?.name || '-'}</td>
                <td>{formatPrice(product.price)}</td>
                <td>{product.stock}</td>
                <td>{product.sku || '-'}</td>
                <td>
                  <span className={`status-badge ${product.is_active ? 'active' : 'inactive'}`}>
                    {product.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td>
                  <button onClick={() => handleOpenModal(product)} className="edit-btn-small">
                    Edit
                  </button>
                  <button onClick={() => handleDelete(product.id)} className="delete-btn-small">
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>{editingProduct ? 'Edit Product' : 'Add New Product'}</h2>
              <button onClick={handleCloseModal} className="close-btn">✕</button>
            </div>

            <form onSubmit={handleSubmit} className="modal-form">
              <div className="form-group">
                <label>Product Name *</label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                />
              </div>

              <div className="form-group">
                <label>Description</label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  rows="3"
                />
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>Price *</label>
                  <input
                    type="number"
                    name="price"
                    value={formData.price}
                    onChange={handleChange}
                    step="0.01"
                    min="0"
                    required
                  />
                </div>

                <div className="form-group">
                  <label>Stock *</label>
                  <input
                    type="number"
                    name="stock"
                    value={formData.stock}
                    onChange={handleChange}
                    min="0"
                    required
                  />
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>SKU</label>
                  <input
                    type="text"
                    name="sku"
                    value={formData.sku}
                    onChange={handleChange}
                  />
                </div>

                <div className="form-group">
                  <label>Category *</label>
                  <select
                    name="category_id"
                    value={formData.category_id}
                    onChange={handleChange}
                    required
                  >
                    <option value="">Select Category</option>
                    {categories.map((cat) => (
                      <option key={cat.id} value={cat.id}>
                        {cat.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="form-group">
                <label>Image URL</label>
                <input
                  type="url"
                  name="image_url"
                  value={formData.image_url}
                  onChange={handleChange}
                  placeholder="https://example.com/image.jpg"
                />
              </div>

              <div className="form-group">
                <label className="checkbox-label">
                  <input
                    type="checkbox"
                    name="is_active"
                    checked={formData.is_active}
                    onChange={handleChange}
                  />
                  <span>Active</span>
                </label>
              </div>

              <div className="form-actions">
                <button type="submit" className="submit-btn">
                  {editingProduct ? 'Update' : 'Create'} Product
                </button>
                <button type="button" onClick={handleCloseModal} className="cancel-btn">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default AdminProducts;
