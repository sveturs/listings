// frontend/hostel-frontend/src/api/axios.js
import axios from 'axios';

const instance = axios.create({
  baseURL: process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000',
  withCredentials: true,
  headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
  },
  validateStatus: function (status) {
      return status >= 200 && status < 500;
  }
});
instance.interceptors.response.use(
  response => response,
  error => {
      if (error.response?.status === 401) {
          // Можно добавить редирект на страницу логина
          window.location.href = '/login';
      }
      return Promise.reject(error);
  }
);

export default instance;