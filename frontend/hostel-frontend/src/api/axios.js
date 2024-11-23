import axios from 'axios';

// базовый URL для API-сервера
const instance = axios.create({
    baseURL: "http://192.168.100.14:3000",
    headers: {
      "Content-Type": "application/json",
    },
  });
  

export default instance;
