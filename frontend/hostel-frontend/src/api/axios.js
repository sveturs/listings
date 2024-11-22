import axios from 'axios';

// Установите базовый URL для вашего API-сервера
const instance = axios.create({
    baseURL: "http://192.168.100.14:3000",
    headers: {
      "Content-Type": "application/json",
    },
  });
  

export default instance;
