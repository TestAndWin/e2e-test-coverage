import axios from 'axios';

const http = axios.create({
  baseURL: import.meta.env.VITE_APP_API_URL ? import.meta.env.VITE_APP_API_URL : window.location.origin,
  headers: {
    'Content-type': 'application/json'
  },
  withCredentials: true
});

// Add request interceptor
http.interceptors.request.use(
  function (config) {
    return config;
  },
  function (error) {
    return Promise.reject(error);
  }
);

let isRedirecting = false;
// Add response interceptor for handling errors
http.interceptors.response.use(
  function (response) {
    return response;
  },
  async (error) => {
    if (error.response && error.response.status === 401) {
      const path = window.location.pathname;
      const requestUrl = error.config?.url || '';
      if (!isRedirecting && !path.startsWith('/login') && !requestUrl.startsWith('/api/v1/auth/me')) {
        isRedirecting = true;
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export default http;
