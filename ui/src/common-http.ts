import axios from 'axios';

const http = axios.create({
  baseURL: import.meta.env.VITE_APP_API_URL ? import.meta.env.VITE_APP_API_URL : window.location.origin,
  headers: {
    'Content-type': 'application/json'
  },
  withCredentials: true
});

http.interceptors.response.use(
  function (response) {
    return response;
  },
  async (error) => {
    if (error.response.status === 401) {
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default http;
