import { useState, useEffect } from 'react';
import { userService } from '../services/api';

function UserList({ onEdit, onRefresh, refreshTrigger }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch users dari API
  const fetchUsers = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await userService.getAllUsers();
      setUsers(data || []);
    } catch (err) {
      setError('Failed to fetch users. Make sure backend is running.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Load users saat komponen mount dan saat refreshTrigger berubah
  useEffect(() => {
    fetchUsers();
  }, [refreshTrigger]);

  // Handle delete user
  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this user?')) {
      return;
    }

    try {
      await userService.deleteUser(id);
      // Refresh list setelah delete
      fetchUsers();
    } catch (err) {
      alert('Failed to delete user');
      console.error(err);
    }
  };

  if (loading) {
    return <div className="loading">Loading users...</div>;
  }

  if (error) {
    return (
      <div className="error">
        <p>{error}</p>
        <button onClick={fetchUsers}>Retry</button>
      </div>
    );
  }

  return (
    <div className="user-list">
      <div className="list-header">
        <h2>User List</h2>
        <button onClick={fetchUsers} className="btn-refresh">
          Refresh
        </button>
      </div>

      {users.length === 0 ? (
        <p className="no-users">No users found. Create one to get started!</p>
      ) : (
        <table className="user-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Email</th>
              <th>Created At</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>{user.id}</td>
                <td>{user.name}</td>
                <td>{user.email}</td>
                <td>{new Date(user.created_at).toLocaleString()}</td>
                <td className="actions">
                  <button
                    onClick={() => onEdit(user)}
                    className="btn-edit"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(user.id)}
                    className="btn-delete"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default UserList;
