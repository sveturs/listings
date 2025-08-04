#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

console.log('📝 Добавление отсутствующих ключей переводов...\n');

// Загружаем отчёт
const reportPath = path.join(__dirname, 'missing-keys-report.json');
if (!fs.existsSync(reportPath)) {
  console.error('❌ Файл missing-keys-report.json не найден!');
  console.log('   Сначала запустите: node scripts/find-missing-translation-keys.js');
  process.exit(1);
}

const report = JSON.parse(fs.readFileSync(reportPath, 'utf8'));
const missingKeys = report.missingKeys;

// Локали для обновления
const locales = ['ru', 'en', 'sr'];

// Счётчики
let totalAdded = 0;
let totalSkipped = 0;

// Функция для добавления ключа в объект по пути
function setNestedValue(obj, path, value) {
  const parts = path.split('.');
  let current = obj;
  
  for (let i = 0; i < parts.length - 1; i++) {
    const part = parts[i];
    if (!current[part] || typeof current[part] !== 'object') {
      current[part] = {};
    }
    current = current[part];
  }
  
  const lastPart = parts[parts.length - 1];
  if (lastPart && current) {
    current[lastPart] = value;
  }
}

// Функция для генерации значения по умолчанию
function generateDefaultValue(key, locale) {
  // Очищаем ключ от префиксов модуля
  const cleanKey = key.split('.').pop();
  
  // Преобразуем camelCase в readable text
  const readable = cleanKey
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, str => str.toUpperCase())
    .trim();
  
  // Базовые переводы для общих ключей
  const commonTranslations = {
    ru: {
      'loading': 'Загрузка...',
      'save': 'Сохранить',
      'cancel': 'Отмена',
      'delete': 'Удалить',
      'edit': 'Редактировать',
      'close': 'Закрыть',
      'search': 'Поиск',
      'all': 'Все',
      'active': 'Активный',
      'inactive': 'Неактивный',
      'yes': 'Да',
      'no': 'Нет',
      'back': 'Назад',
      'next': 'Далее',
      'error': 'Ошибка',
      'success': 'Успешно',
      'noData': 'Нет данных',
    },
    en: {
      'loading': 'Loading...',
      'save': 'Save',
      'cancel': 'Cancel',
      'delete': 'Delete',
      'edit': 'Edit',
      'close': 'Close',
      'search': 'Search',
      'all': 'All',
      'active': 'Active',
      'inactive': 'Inactive',
      'yes': 'Yes',
      'no': 'No',
      'back': 'Back',
      'next': 'Next',
      'error': 'Error',
      'success': 'Success',
      'noData': 'No data',
    },
    sr: {
      'loading': 'Učitavanje...',
      'save': 'Sačuvaj',
      'cancel': 'Otkaži',
      'delete': 'Obriši',
      'edit': 'Izmeni',
      'close': 'Zatvori',
      'search': 'Pretraga',
      'all': 'Sve',
      'active': 'Aktivno',
      'inactive': 'Neaktivno',
      'yes': 'Da',
      'no': 'Ne',
      'back': 'Nazad',
      'next': 'Sledeće',
      'error': 'Greška',
      'success': 'Uspešno',
      'noData': 'Nema podataka',
    }
  };
  
  // Проверяем есть ли готовый перевод
  if (commonTranslations[locale] && commonTranslations[locale][cleanKey]) {
    return commonTranslations[locale][cleanKey];
  }
  
  // Иначе возвращаем сгенерированное значение
  if (locale === 'ru') {
    return `[RU] ${readable}`;
  } else if (locale === 'sr') {
    return `[SR] ${readable}`;
  }
  
  return readable;
}

// Обрабатываем каждый модуль
Object.entries(missingKeys).forEach(([module, keys]) => {
  console.log(`\n📦 Обработка модуля: ${module}`);
  
  if (keys.length === 0) {
    console.log('   ✅ Нет отсутствующих ключей');
    return;
  }
  
  // Фильтруем невалидные ключи
  const validKeys = keys.filter(key => {
    // Пропускаем пустые ключи и пути к файлам
    if (!key || 
        key.includes('/') || 
        key.includes('@/') || 
        key.trim() === ',' || 
        key.trim() === '_' ||
        key.trim() === 'a' ||
        key.trim() === 'T' ||
        key.includes('2d') ||
        key.includes('canvas') ||
        key.length > 100 || // Слишком длинные ключи обычно ошибочные
        /^[0-9]+$/.test(key) || // Числовые ключи
        key.includes('Facebook') ||
        key.includes('Instagram') ||
        key.includes('Геолокация') || // Длинные сообщения на русском
        key.includes('Не удалось')) {
      totalSkipped++;
      return false;
    }
    return true;
  });
  
  if (validKeys.length === 0) {
    console.log('   ⏭️  Все ключи невалидны, пропускаем');
    return;
  }
  
  // Обновляем каждую локаль
  locales.forEach(locale => {
    const filePath = path.join(__dirname, `../src/messages/${locale}/${module}.json`);
    
    // Проверяем существование файла
    if (!fs.existsSync(filePath)) {
      console.log(`   ⚠️  Файл не существует: ${locale}/${module}.json`);
      return;
    }
    
    // Загружаем существующие переводы
    let translations = {};
    try {
      translations = JSON.parse(fs.readFileSync(filePath, 'utf8'));
    } catch (e) {
      console.error(`   ❌ Ошибка чтения ${locale}/${module}.json:`, e.message);
      return;
    }
    
    // Добавляем отсутствующие ключи
    let addedCount = 0;
    validKeys.forEach(key => {
      // Проверяем, существует ли уже ключ
      const parts = key.split('.');
      let current = translations;
      let exists = true;
      
      for (const part of parts) {
        if (!current || !current[part]) {
          exists = false;
          break;
        }
        current = current[part];
      }
      
      if (!exists) {
        // Генерируем значение по умолчанию
        const defaultValue = generateDefaultValue(key, locale);
        setNestedValue(translations, key, defaultValue);
        addedCount++;
      }
    });
    
    if (addedCount > 0) {
      // Сохраняем обновлённый файл
      fs.writeFileSync(filePath, JSON.stringify(translations, null, 2) + '\n');
      console.log(`   ✅ ${locale}: добавлено ${addedCount} ключей`);
      totalAdded += addedCount;
    } else {
      console.log(`   ℹ️  ${locale}: все ключи уже существуют`);
    }
  });
});

console.log('\n📊 Итоги:');
console.log(`   Добавлено ключей: ${totalAdded}`);
console.log(`   Пропущено невалидных: ${totalSkipped}`);

console.log('\n⚠️  Важно:');
console.log('   1. Добавлены ключи с значениями-заглушками');
console.log('   2. Необходимо заменить их на правильные переводы');
console.log('   3. Перезапустите frontend для применения изменений');

console.log('\n✨ Готово!');