import { useState } from 'react'
import './App.css'
import UserForm from './components/UserForm'
import UserList from './components/UserList'

function App() {
  const [editingUser, setEditingUser] = useState(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Handle saat user berhasil dibuat/diupdate
  const handleSuccess = () => {
    setEditingUser(null);
    setRefreshTrigger(prev => prev + 1); // Trigger refresh list
  };

  // Handle saat tombol edit diklik
  const handleEdit = (user) => {
    setEditingUser(user);
  };

  // Handle saat tombol cancel diklik
  const handleCancel = () => {
    setEditingUser(null);
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>User Management System</h1>
        <p>GoBack + React Frontend</p>
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
  )
}

export default App
