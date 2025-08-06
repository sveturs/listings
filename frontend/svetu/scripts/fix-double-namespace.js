#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Функция для чтения файла
function readFile(filePath) {
  return fs.readFileSync(filePath, 'utf8');
}

// Функция для записи файла
function writeFile(filePath, content) {
  fs.writeFileSync(filePath, content);
}

// Находим все файлы TypeScript и TSX
const files = glob.sync('src/**/*.{ts,tsx}', {
  ignore: ['**/node_modules/**'],
});

let fixedCount = 0;

files.forEach((file) => {
  let content = readFile(file);
  let modified = false;

  // Ищем паттерн useTranslations('namespace')
  const translationMatches = content.match(/useTranslations\(['"](\w+)['"]\)/g);

  if (translationMatches) {
    translationMatches.forEach((match) => {
      const namespace = match.match(/useTranslations\(['"](\w+)['"]\)/)[1];

      // Ищем использование t('namespace.key') с тем же namespace
      const pattern = new RegExp(`t\\(['"]${namespace}\\.([^'"]+)['"]`, 'g');

      if (pattern.test(content)) {
        // Заменяем namespace.key на просто key
        content = content.replace(pattern, `t('$1'`);
        modified = true;
        console.log(`Fixed double namespace in ${file}: ${namespace}.* → *`);
      }
    });
  }

  if (modified) {
    writeFile(file, content);
    fixedCount++;
  }
});

console.log(`\n✅ Fixed ${fixedCount} files with double namespace usage`);
