import { useState, useEffect } from 'react';
import { userService } from '../services/api';

function UserForm({ editingUser, onSuccess, onCancel }) {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Set form data jika sedang edit user
  useEffect(() => {
    if (editingUser) {
      setFormData({
        name: editingUser.name,
        email: editingUser.email,
      });
    } else {
      setFormData({
        name: '',
        email: '',
      });
    }
  }, [editingUser]);

  // Handle input change
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // Handle form submit
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      if (editingUser) {
        // Update existing user
        await userService.updateUser(editingUser.id, formData);
      } else {
        // Create new user
        await userService.createUser(formData);
      }

      // Reset form
      setFormData({ name: '', email: '' });

      // Callback untuk refresh list
      onSuccess();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to save user');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="user-form">
      <h2>{editingUser ? 'Edit User' : 'Create New User'}</h2>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Name:</label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            placeholder="Enter user name"
          />
        </div>

        <div className="form-group">
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
            placeholder="Enter user email"
          />
        </div>

        <div className="form-actions">
          <button
            type="submit"
            disabled={loading}
            className="btn-submit"
          >
            {loading ? 'Saving...' : editingUser ? 'Update User' : 'Create User'}
          </button>

          {editingUser && (
            <button
              type="button"
              onClick={onCancel}
              className="btn-cancel"
            >
              Cancel
            </button>
          )}
        </div>
      </form>
    </div>
  );
}

export default UserForm;
