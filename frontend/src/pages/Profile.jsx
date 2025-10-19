import { useState, useEffect } from 'react';
import { profileService } from '../services/api';
import './Profile.css';

function Profile() {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [formData, setFormData] = useState({
    phone_number: '',
    address: '',
    city: '',
    province: '',
    postal_code: '',
    country: 'Indonesia',
    bio: '',
    gender: '',
  });

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      setLoading(true);
      const response = await profileService.getMyProfile();
      setProfile(response.data);
      setFormData({
        phone_number: response.data.phone_number || '',
        address: response.data.address || '',
        city: response.data.city || '',
        province: response.data.province || '',
        postal_code: response.data.postal_code || '',
        country: response.data.country || 'Indonesia',
        bio: response.data.bio || '',
        gender: response.data.gender || '',
      });
    } catch (err) {
      console.log('No profile found');
      setEditing(true);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      await profileService.updateMyProfile(formData);
      alert('Profile updated successfully!');
      setEditing(false);
      fetchProfile();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to update profile');
    }
  };

  if (loading) {
    return (
      <div className="profile-container">
        <div className="loading">Loading profile...</div>
      </div>
    );
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <h1>My Profile</h1>
        {!editing && (
          <button onClick={() => setEditing(true)} className="edit-btn">
            Edit Profile
          </button>
        )}
      </div>

      {editing ? (
        <form onSubmit={handleSubmit} className="profile-form">
          <div className="form-section">
            <h2>Contact Information</h2>

            <div className="form-group">
              <label>Phone Number</label>
              <input
                type="tel"
                name="phone_number"
                value={formData.phone_number}
                onChange={handleChange}
                placeholder="Enter your phone number"
              />
            </div>

            <div className="form-group">
              <label>Gender</label>
              <select name="gender" value={formData.gender} onChange={handleChange}>
                <option value="">Select gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
              </select>
            </div>
          </div>

          <div className="form-section">
            <h2>Address</h2>

            <div className="form-group">
              <label>Street Address</label>
              <textarea
                name="address"
                value={formData.address}
                onChange={handleChange}
                rows="3"
                placeholder="Enter your address"
              />
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>City</label>
                <input
                  type="text"
                  name="city"
                  value={formData.city}
                  onChange={handleChange}
                  placeholder="City"
                />
              </div>

              <div className="form-group">
                <label>Province</label>
                <input
                  type="text"
                  name="province"
                  value={formData.province}
                  onChange={handleChange}
                  placeholder="Province"
                />
              </div>
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Postal Code</label>
                <input
                  type="text"
                  name="postal_code"
                  value={formData.postal_code}
                  onChange={handleChange}
                  placeholder="Postal Code"
                />
              </div>

              <div className="form-group">
                <label>Country</label>
                <input
                  type="text"
                  name="country"
                  value={formData.country}
                  onChange={handleChange}
                  placeholder="Country"
                />
              </div>
            </div>
          </div>

          <div className="form-section">
            <h2>About Me</h2>

            <div className="form-group">
              <label>Bio</label>
              <textarea
                name="bio"
                value={formData.bio}
                onChange={handleChange}
                rows="4"
                placeholder="Tell us about yourself"
              />
            </div>
          </div>

          <div className="form-actions">
            <button type="submit" className="save-btn">
              Save Changes
            </button>
            <button
              type="button"
              onClick={() => {
                setEditing(false);
                fetchProfile();
              }}
              className="cancel-btn"
            >
              Cancel
            </button>
          </div>
        </form>
      ) : (
        <div className="profile-view">
          <div className="profile-section">
            <h2>Contact Information</h2>
            <div className="profile-field">
              <span className="field-label">Phone:</span>
              <span>{profile?.phone_number || '-'}</span>
            </div>
            <div className="profile-field">
              <span className="field-label">Gender:</span>
              <span className="capitalize">{profile?.gender || '-'}</span>
            </div>
          </div>

          <div className="profile-section">
            <h2>Address</h2>
            <div className="profile-field">
              <span className="field-label">Street:</span>
              <span>{profile?.address || '-'}</span>
            </div>
            <div className="profile-field">
              <span className="field-label">City:</span>
              <span>{profile?.city || '-'}</span>
            </div>
            <div className="profile-field">
              <span className="field-label">Province:</span>
              <span>{profile?.province || '-'}</span>
            </div>
            <div className="profile-field">
              <span className="field-label">Postal Code:</span>
              <span>{profile?.postal_code || '-'}</span>
            </div>
            <div className="profile-field">
              <span className="field-label">Country:</span>
              <span>{profile?.country || '-'}</span>
            </div>
          </div>

          {profile?.bio && (
            <div className="profile-section">
              <h2>About Me</h2>
              <p>{profile.bio}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default Profile;
