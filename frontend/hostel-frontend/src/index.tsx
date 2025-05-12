// frontend/hostel-frontend/src/index.tsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import './i18n/config';

import App from './App';
import reportWebVitals from './reportWebVitals';
import { ThemeProvider, createTheme } from "@mui/material/styles";
import { CssBaseline } from '@mui/material';

// Примечание: Глобальный тип window.ENV уже определен в src/types/global.d.ts

console.log('Environment check:', window.ENV);

// Создание темы Material UI
const theme = createTheme({
  // Кастомные настройки темы
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

// Получаем корневой элемент и рендерим приложение
const rootElement = document.getElementById('root');
if (!rootElement) throw new Error('Failed to find the root element');

const root = ReactDOM.createRoot(rootElement);
root.render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </React.StrictMode>
);

// Измерение производительности веб-приложения
reportWebVitals();

// Примечание: reportWebVitals() вызывается дважды в исходном файле, но это, вероятно, ошибка.
// Предполагаем, что один вызов был предназначен для примера, поэтому оставляем только один вызов.