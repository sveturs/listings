import axios from 'axios';

// Установите базовый URL для вашего API-сервера
const instance = axios.create({
  baseURL: 'http://localhost:3000', // Замените на URL вашего API
});

export default instance;
