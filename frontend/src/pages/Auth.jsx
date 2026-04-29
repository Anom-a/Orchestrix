import React, { useState } from 'react';
import { api } from '../api';
import { Database, Lock, Mail, User } from 'lucide-react';
import './Auth.css'; // Let's put specific auth CSS here or just use inline/index.css

export default function Auth({ setIsAuthenticated }) {
  const [isLogin, setIsLogin] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: ''
  });

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      if (isLogin) {
        const res = await api.login(formData.email, formData.password);
        localStorage.setItem('token', res.token);
        setIsAuthenticated(true);
      } else {
        await api.register(formData.name, formData.email, formData.password);
        // Auto login after register
        const res = await api.login(formData.email, formData.password);
        localStorage.setItem('token', res.token);
        setIsAuthenticated(true);
      }
    } catch (err) {
      setError(err.message || 'Authentication failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-brand-section">
        <div className="auth-brand-content fade-in">
          <div className="brand-icon-wrapper">
            <Database size={48} className="brand-icon" />
          </div>
          <h1>Orchestrix</h1>
          <p className="brand-tagline">
            Your personal AI knowledge base.<br />
            Upload documents. Ask anything.
          </p>
          
          <div className="feature-list">
            <div className="feature-item">
              <span className="feature-dot"></span>
              <span>RAG-powered contextual answers</span>
            </div>
            <div className="feature-item">
              <span className="feature-dot"></span>
              <span>Multi-format document support</span>
            </div>
            <div className="feature-item">
              <span className="feature-dot"></span>
              <span>Instant source-cited responses</span>
            </div>
          </div>
        </div>
      </div>

      <div className="auth-form-section">
        <div className="glass-panel auth-card fade-in" style={{ animationDelay: '0.1s' }}>
          <div className="auth-tabs">
            <button 
              className={`auth-tab ${isLogin ? 'active' : ''}`}
              onClick={() => setIsLogin(true)}
              type="button"
            >
              Sign In
            </button>
            <button 
              className={`auth-tab ${!isLogin ? 'active' : ''}`}
              onClick={() => setIsLogin(false)}
              type="button"
            >
              Create Account
            </button>
          </div>

          <form onSubmit={handleSubmit} className="auth-form">
            {!isLogin && (
              <div className="form-group">
                <label>Full Name</label>
                <div className="input-wrapper">
                  <User size={18} className="input-icon" />
                  <input
                    type="text"
                    name="name"
                    className="input-field with-icon"
                    placeholder="John Doe"
                    value={formData.name}
                    onChange={handleChange}
                    required={!isLogin}
                  />
                </div>
              </div>
            )}

            <div className="form-group">
              <label>Email Address</label>
              <div className="input-wrapper">
                <Mail size={18} className="input-icon" />
                <input
                  type="email"
                  name="email"
                  className="input-field with-icon"
                  placeholder="you@example.com"
                  value={formData.email}
                  onChange={handleChange}
                  required
                />
              </div>
            </div>

            <div className="form-group">
              <label>Password</label>
              <div className="input-wrapper">
                <Lock size={18} className="input-icon" />
                <input
                  type="password"
                  name="password"
                  className="input-field with-icon"
                  placeholder={isLogin ? '••••••••' : 'Min 6 characters'}
                  value={formData.password}
                  onChange={handleChange}
                  required
                  minLength={6}
                />
              </div>
            </div>

            {error && <div className="auth-error">{error}</div>}

            <button type="submit" className="btn btn-primary auth-submit" disabled={isLoading}>
              {isLoading ? (
                <div className="spinner"></div>
              ) : (
                isLogin ? 'Sign In' : 'Create Account'
              )}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
