import axios from 'axios';

// базовый URL для API-сервера
const instance = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});


export default instance;
