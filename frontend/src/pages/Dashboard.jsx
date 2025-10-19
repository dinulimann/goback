import { useState } from 'react';
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
            <h1>User Management System</h1>
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
