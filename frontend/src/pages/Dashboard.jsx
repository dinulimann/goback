import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import UserForm from '../components/UserForm';
import UserList from '../components/UserList';
import '../App.css';

function Dashboard() {
  const [editingUser, setEditingUser] = useState(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const { user, logout } = useAuth();

  const handleSuccess = () => {
    setEditingUser(null);
    setRefreshTrigger(prev => prev + 1);
  };

  const handleEdit = (user) => {
    setEditingUser(user);
  };

  const handleCancel = () => {
    setEditingUser(null);
  };

  return (
    <div className="dashboard">
      <header className="app-header">
        <div className="header-content">
          <div>
            <h1>Admin Dashboard</h1>
            <p>Welcome back, {user?.name}!</p>
          </div>
          <button onClick={logout} className="btn-logout">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            Logout
          </button>
        </div>
      </header>

      {/* Admin Quick Links */}
      <div style={{ maxWidth: '1200px', margin: '20px auto', padding: '0 20px' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '20px', marginBottom: '30px' }}>
          <Link to="/admin/categories" style={{
            padding: '30px',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '12px',
            textAlign: 'center',
            fontWeight: '600',
            fontSize: '1.1rem',
            transition: 'transform 0.3s',
            boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
          }}>
            📁 Manage Categories
          </Link>
          <Link to="/admin/products" style={{
            padding: '30px',
            background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '12px',
            textAlign: 'center',
            fontWeight: '600',
            fontSize: '1.1rem',
            transition: 'transform 0.3s',
            boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
          }}>
            📦 Manage Products
          </Link>
          <Link to="/orders" style={{
            padding: '30px',
            background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '12px',
            textAlign: 'center',
            fontWeight: '600',
            fontSize: '1.1rem',
            transition: 'transform 0.3s',
            boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
          }}>
            🛒 View Orders
          </Link>
          <Link to="/products" style={{
            padding: '30px',
            background: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '12px',
            textAlign: 'center',
            fontWeight: '600',
            fontSize: '1.1rem',
            transition: 'transform 0.3s',
            boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
          }}>
            🏪 View Store
          </Link>
        </div>
      </div>

      <div className="app-container">
        <div className="form-section">
          <UserForm
            editingUser={editingUser}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        </div>

        <div className="list-section">
          <UserList
            onEdit={handleEdit}
            refreshTrigger={refreshTrigger}
          />
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
