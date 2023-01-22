import axios from "axios";

const http = axios.create({
  baseURL: process.env.VUE_APP_API_URL ? process.env.VUE_APP_API_URL : window.location.origin,
  headers: {
    "Content-type": "application/json",
    'Authorization': sessionStorage.getItem('token') ? `Bearer ${sessionStorage.getItem('token')}` : '',
  },
  withCredentials: true,
});

http.interceptors.response.use(function (response) {
  return response;
}, async (error) => {
  const originalRequest = error.config;
  if (error.response.status === 401 && !originalRequest?.doRefresh) {
    // Try only one time to refresh the token
    originalRequest.doRefresh = true;
    await http.post('/api/v1/auth/refresh', { token: sessionStorage.getItem('refreshToken') }).then(response => {
      if (response.status === 200) {
        sessionStorage.setItem('token', response.data.token);
        originalRequest.headers['Authorization'] = `Bearer ${response.data.token}`;
      }
      else {
        location.assign('/login');
      }
    }).catch(() => {
      location.assign('/login');
    })
    return axios(originalRequest);
  }
  return Promise.reject(error);
});

export default http;

