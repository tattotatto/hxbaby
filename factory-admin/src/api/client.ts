import axios from 'axios';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:9090/api/v1',
  timeout: 30000,
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('factory_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('factory_token');
      localStorage.removeItem('factory_customer');
      // Only redirect if not already on login/register page
      if (!window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
        window.location.href = '/login';
      }
    }
    // Network error
    if (!err.response) {
      console.error('Network error:', err.message);
    }
    return Promise.reject(err);
  }
);

export default client;
