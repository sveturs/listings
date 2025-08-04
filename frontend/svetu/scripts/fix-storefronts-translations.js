#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Файл для обновления
const filePath = path.join(__dirname, '../src/app/[locale]/profile/storefronts/page.tsx');

// Читаем файл
let content = fs.readFileSync(filePath, 'utf8');

// Замены для основного namespace storefronts
const replacements = [
  // Простые замены
  [/t\('storefronts\.alwaysOpen'\)/g, "t('alwaysOpen')"],
  [/t\('storefronts\.closedToday'\)/g, "t('closedToday')"],
  [/t\('storefronts\.myStorefronts'\)/g, "t('myStorefronts')"],
  [/t\('storefronts\.manageDescription'\)/g, "t('manageDescription')"],
  [/t\('storefronts\.createNew'\)/g, "t('createNew')"],
  [/t\('storefronts\.totalStorefronts'\)/g, "t('totalStorefronts')"],
  [/t\('storefronts\.activeCount'\)/g, "t('activeCount')"],
  [/t\('storefronts\.totalViews'\)/g, "t('totalViews')"],
  [/t\('storefronts\.totalProducts'\)/g, "t('totalProducts')"],
  [/t\('storefronts\.acrossAllStorefronts'\)/g, "t('acrossAllStorefronts')"],
  [/t\('storefronts\.monthlyRevenue'\)/g, "t('monthlyRevenue')"],
  [/t\('storefronts\.thisMonth'\)/g, "t('thisMonth')"],
  [/t\('storefronts\.active'\)/g, "t('active')"],
  [/t\('storefronts\.inactive'\)/g, "t('inactive')"],
  [/t\('storefronts\.noStorefronts'\)/g, "t('noStorefronts')"],
  [/t\('storefronts\.noStorefrontsInCategory'\)/g, "t('noStorefrontsInCategory')"],
  [/t\('storefronts\.createFirstStorefront'\)/g, "t('createFirstStorefront')"],
  [/t\('storefronts\.createStorefront'\)/g, "t('createStorefront')"],
  [/t\('storefronts\.noDescription'\)/g, "t('noDescription')"],
  [/t\('storefronts\.views'\)/g, "t('views')"],
  [/t\('storefronts\.orders'\)/g, "t('orders')"],
  [/t\('storefronts\.dashboard'\)/g, "t('dashboard')"],
  [/t\('storefronts\.settings'\)/g, "t('settings')"],
  [/t\('storefronts\.viewPublicPage'\)/g, "t('viewPublicPage')"],
  [/t\('storefronts\.editDetails'\)/g, "t('editDetails')"],
  [/t\('storefronts\.manageProducts'\)/g, "t('manageProducts')"],
  [/t\('storefronts\.manageStaff'\)/g, "t('manageStaff')"],
  [/t\('storefronts\.manageReviews'\)/g, "t('manageReviews')"],
  [/t\('storefronts\.messages'\)/g, "t('messages')"],
  
  // Вложенные ключи
  [/t\('storefronts\.status\.active'\)/g, "t('status.active')"],
  [/t\('storefronts\.status\.inactive'\)/g, "t('status.inactive')"],
  [/t\('storefronts\.products\.title'\)/g, "t('products.title')"],
];

// Применяем замены
replacements.forEach(([pattern, replacement]) => {
  content = content.replace(pattern, replacement);
});

// Записываем обратно
fs.writeFileSync(filePath, content);

console.log('✅ Успешно обновлены переводы в MyStorefrontsPage');