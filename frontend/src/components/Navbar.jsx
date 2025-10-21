import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import './Navbar.css';

function Navbar() {
  const { isAuthenticated, logout } = useAuth();
  const { getCartItemCount, toggleCart } = useCart();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          🛒 GoBack Store
        </Link>

        <div className="navbar-menu">
          <Link to="/products" className="nav-link">
            Products
          </Link>

          {isAuthenticated ? (
            <>
              <Link to="/orders" className="nav-link">
                My Orders
              </Link>
              <Link to="/profile" className="nav-link">
                Profile
              </Link>
              <Link to="/dashboard" className="nav-link">
                Dashboard
              </Link>
              <div className="dropdown">
                <button className="dropdown-btn">
                  Admin
                </button>
                <div className="dropdown-content">
                  <Link to="/admin/dashboard">Admin Dashboard</Link>
                  <Link to="/admin/products">Manage Products</Link>
                  <Link to="/admin/categories">Manage Categories</Link>
                </div>
              </div>
              <button onClick={toggleCart} className="cart-btn">
                🛒 Cart
                {getCartItemCount() > 0 && (
                  <span className="cart-badge">{getCartItemCount()}</span>
                )}
              </button>
              <button onClick={handleLogout} className="logout-btn">
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">
                Login
              </Link>
              <Link to="/register" className="nav-link-primary">
                Register
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
