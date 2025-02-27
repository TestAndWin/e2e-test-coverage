import axios from 'axios';

const http = axios.create({
  baseURL: import.meta.env.VITE_APP_API_URL ? import.meta.env.VITE_APP_API_URL : window.location.origin,
  headers: {
    'Content-type': 'application/json'
  },
  withCredentials: true
});

// Add request interceptor for debugging
http.interceptors.request.use(
  function (config) {
    console.log('Making request to:', config.url);
    return config;
  },
  function (error) {
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor with more detailed error handling
http.interceptors.response.use(
  function (response) {
    console.log('Response from', response.config.url, ':', response.status);

    // Debug login response
    if (response.config.url === '/api/v1/auth/login') {
      console.log('Login response data:', response.data);
      if (response.data.roles) {
        console.log('Setting roles in session storage:', response.data.roles);
      } else {
        console.warn('No roles found in login response');
      }
    }

    return response;
  },
  async (error) => {
    if (error.response) {
      console.error('Response error:', error.response.status, error.response.config.url);
      if (error.response.status === 401) {
        window.location.href = '/login';
      }
    } else if (error.request) {
      console.error('Request error - no response received');
    } else {
      console.error('Error setting up request:', error.message);
    }
    return Promise.reject(error);
  }
);

export default http;
