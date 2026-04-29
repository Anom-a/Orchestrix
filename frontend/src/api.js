const API_URL = 'http://localhost:8080';

export const api = {
  getHeaders(isFormData = false) {
    const token = localStorage.getItem('token');
    const headers = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    if (!isFormData) {
      headers['Content-Type'] = 'application/json';
    }
    return headers;
  },

  async handleResponse(response) {
    if (!response.ok) {
      let errorMessage = 'An error occurred';
      try {
        const errorData = await response.json();
        errorMessage = errorData.error || errorMessage;
      } catch (e) {
        // Not JSON
      }
      throw new Error(errorMessage);
    }
    return response.json();
  },

  // Auth
  async login(email, password) {
    const res = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ email, password }),
    });
    return this.handleResponse(res);
  },

  async register(name, email, password) {
    const res = await fetch(`${API_URL}/auth/register`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ name, email, password }),
    });
    return this.handleResponse(res);
  },

  // Documents
  async getDocuments() {
    const res = await fetch(`${API_URL}/api/documents/`, {
      headers: this.getHeaders(),
    });
    return this.handleResponse(res);
  },

  async getDocument(id) {
    const res = await fetch(`${API_URL}/api/documents/${id}`, {
      headers: this.getHeaders(),
    });
    return this.handleResponse(res);
  },

  async uploadDocument(file) {
    const formData = new FormData();
    formData.append('file', file);
    
    const res = await fetch(`${API_URL}/api/documents/upload`, {
      method: 'POST',
      headers: this.getHeaders(true),
      body: formData,
    });
    return this.handleResponse(res);
  },

  async processDocument(id) {
    const res = await fetch(`${API_URL}/api/documents/${id}/process`, {
      method: 'POST',
      headers: this.getHeaders(),
    });
    return this.handleResponse(res);
  },

  async queryDocument(id, question) {
    const res = await fetch(`${API_URL}/api/documents/${id}/query`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ question }),
    });
    return this.handleResponse(res);
  }
};
