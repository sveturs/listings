#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Файлы для обработки
const files = ['src/messages/ru/common.json', 'src/messages/sr/common.json'];

files.forEach((file) => {
  console.log(`Processing ${file}...`);

  const content = JSON.parse(fs.readFileSync(file, 'utf8'));

  // Если есть объект common внутри файла common.json
  if (content.common && typeof content.common === 'object') {
    console.log(
      `Found nested common object with ${Object.keys(content.common).length} keys`
    );

    // Переносим все ключи из common на верхний уровень
    Object.keys(content.common).forEach((key) => {
      // Если ключ уже существует на верхнем уровне
      if (content[key] !== undefined) {
        console.log(
          `  - Key "${key}" already exists: "${content[key]}" -> replacing with "${content.common[key]}"`
        );
      }
      // Перезаписываем значением из common
      content[key] = content.common[key];
    });

    // Удаляем вложенный объект common
    delete content.common;

    // Сохраняем обновленный файл
    fs.writeFileSync(file, JSON.stringify(content, null, 2) + '\n', 'utf8');
    console.log(`  ✅ Fixed ${file}`);
  } else {
    console.log(`  ℹ️  No nested common object found`);
  }
});

console.log('\nDone! Structure fixed.');
