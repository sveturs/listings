#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Имя временного файла для изолированной проверки
const TEMP_FILE = path.join(__dirname, '.temp-check.tsx');

// Функция для извлечения только кода компонента, без импортов библиотек
function extractComponentCode(filePath) {
  const content = fs.readFileSync(filePath, 'utf8');
  
  // Удаляем импорты i18next и других проблемных библиотек
  const modifiedContent = content
    .replace(/import.*from ['"]i18next['"].*;\n/g, '')
    .replace(/import.*from ['"]react-i18next['"].*;\n/g, '')
    .replace(/import.*from ['"].*\/LocationContext['"].*;\n/g, '')
    .replace(/import.*from ['"].*\/AuthContext['"].*;\n/g, '')
    .replace(/import.*from ['"].*\/axios['"].*;\n/g, '');

  return `
// Временный файл для проверки типов
import React from 'react';
import { Box, Typography } from '@mui/material';

// Заглушки для контекстов
const useLocation = () => ({
  userLocation: { city: '', country: '' },
  setCity: () => {},
  detectUserLocation: async () => {},
  isGeolocating: false
});

const useAuth = () => ({
  user: null,
  isAuthenticated: false
});

const useTranslation = () => ({
  t: (key, options) => key
});

// Заглушка для axios
const axios = {
  get: async () => ({ data: { data: [] } }),
  put: async () => ({})
};

${modifiedContent}
`;
}

// Проверяем, передан ли путь к файлу
if (process.argv.length < 3) {
  console.error('Необходимо указать путь к файлу для проверки');
  process.exit(1);
}

const filePath = process.argv[2];

try {
  // Создаем временный файл с изолированным компонентом
  const code = extractComponentCode(filePath);
  fs.writeFileSync(TEMP_FILE, code);
  
  // Запускаем проверку типов
  console.log(`Проверка изолированного компонента из ${filePath}...`);
  execSync(`export NODE_OPTIONS=--openssl-legacy-provider && npx tsc ${TEMP_FILE} --noEmit --jsx react --skipLibCheck`, { stdio: 'inherit' });
  
  console.log('✅ Компонент прошел проверку типов');
} catch (error) {
  console.error('❌ Ошибка при проверке типов:', error.message);
  process.exit(1);
} finally {
  // Удаляем временный файл
  if (fs.existsSync(TEMP_FILE)) {
    fs.unlinkSync(TEMP_FILE);
  }
}